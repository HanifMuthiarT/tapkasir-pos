package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"tapkasir-pos/internal/delivery/http/handler"
	"tapkasir-pos/internal/delivery/http/middleware"
	"tapkasir-pos/internal/model"
	"tapkasir-pos/internal/repository"
	"tapkasir-pos/internal/usecase"
)

func main() {
	// 1. Konfigurasi Koneksi Database (Default Laragon: user root, password kosong)
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:@tcp(127.0.0.1:3306)/tapkasir_db?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	// Rename existing plural tables before migration to preserve their data
	tableNames := []string{
		"Menu",
		"User",
		"BahanBaku",
		"LogStokMasuk",
		"LogBarangRusak",
		"StockOpname",
		"Resep",
		"Pesanan",
		"DetailPesanan",
		"LogVoid",
	}
	pluralNaming := schema.NamingStrategy{}
	singularNaming := schema.NamingStrategy{SingularTable: true}
	for _, name := range tableNames {
		oldTable := pluralNaming.TableName(name)
		newTable := singularNaming.TableName(name)
		if oldTable == newTable || !db.Migrator().HasTable(oldTable) || db.Migrator().HasTable(newTable) {
			continue
		}
		if err := db.Migrator().RenameTable(oldTable, newTable); err != nil {
			log.Fatalf("Gagal mengubah nama tabel %s menjadi %s: %v", oldTable, newTable, err)
		}
	}

	// Auto-Migrate tabel database
	err = db.AutoMigrate(
		&model.Menu{},
		&model.User{},
		&model.BahanBaku{},
		&model.LogStokMasuk{},
		&model.LogBarangRusak{},
		&model.StockOpname{},
		&model.Resep{},
		&model.Pesanan{},
		&model.DetailPesanan{},
		&model.LogVoid{},
	)
	if err != nil {
		log.Fatal("Gagal melakukan migrasi database:", err)
	}
	log.Println("Database terhubung dan migrasi tabel sukses!")

	// Pastikan kolom status di tabel pesanan mendukung seluruh transisi status
	_ = db.Exec("ALTER TABLE `pesanan` MODIFY COLUMN `status_pesanan` VARCHAR(30) NOT NULL DEFAULT 'menunggu'").Error
	_ = db.Exec("ALTER TABLE `pesanan` MODIFY COLUMN `status_pembayaran` VARCHAR(30) NOT NULL DEFAULT 'unpaid'").Error

	// 2. Secret Key untuk JWT Token
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "tapkasir-super-secret-key-2026"
	}

	// 3. Wiring Dependency Injection: User & PIN Authentication (Tahap 1)
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, jwtSecret)

	// Inisialisasi Akun/PIN Default jika database user masih kosong
	if err := userUsecase.SeedDefaultUsers(); err != nil {
		log.Println("Catatan: Terjadi kendala saat pengecekan user seeder:", err)
	}

	// 4. Wiring Dependency Injection: Menu (CRUD Menu)
	menuRepo := repository.NewMenuRepository(db)
	menuUsecase := usecase.NewMenuUsecase(menuRepo)

	// 5. Wiring Dependency Injection: Bahan Baku (Tahap 2)
	bahanBakuRepo := repository.NewBahanBakuRepository(db)
	bahanBakuUsecase := usecase.NewBahanBakuUsecase(bahanBakuRepo)

	// 6. Wiring Dependency Injection: Resep & Kalkulasi HPP (Tahap 3)
	resepRepo := repository.NewResepRepository(db)
	resepUsecase := usecase.NewResepUsecase(resepRepo, menuRepo, bahanBakuRepo)

	// 7. Wiring Dependency Injection: Operasional Inventori & Stok Gudang (Tahap 4)
	invRepo := repository.NewInventoryRepository(db)
	invUsecase := usecase.NewInventoryUsecase(invRepo, bahanBakuRepo)

	// 8. Wiring Dependency Injection: Pemesanan Pelanggan (Tahap 5)
	pesananRepo := repository.NewPesananRepository(db)
	customerOrderUsecase := usecase.NewCustomerOrderUsecase(pesananRepo, menuRepo, resepRepo, bahanBakuRepo)

	// 9. Wiring Dependency Injection: Operasional Kasir (Tahap 6)
	cashierOrderUsecase := usecase.NewCashierOrderUsecase(pesananRepo, resepRepo, bahanBakuRepo)

	// 10. Wiring Dependency Injection: Void / Pembatalan Pesanan (Tahap 7)
	voidUsecase := usecase.NewVoidUsecase(pesananRepo, resepRepo, bahanBakuRepo)

	// 11. Wiring Dependency Injection: Laporan & Export Excel (Tahap 8)
	reportRepo := repository.NewReportRepository(db)
	reportUsecase := usecase.NewReportUsecase(reportRepo, pesananRepo, invRepo)

	// 12. Inisialisasi Gin Web Framework & Daftarkan Handler Route
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	handler.NewUserHandler(r, userUsecase)
	handler.NewMenuHandler(r, menuUsecase, jwtSecret)
	handler.NewBahanBakuHandler(r, bahanBakuUsecase, jwtSecret)
	handler.NewResepHandler(r, resepUsecase, jwtSecret)
	handler.NewInventoryHandler(r, invUsecase, jwtSecret)
	handler.NewCustomerOrderHandler(r, customerOrderUsecase)
	handler.NewCashierHandler(r, cashierOrderUsecase)
	handler.NewVoidHandler(r, voidUsecase, jwtSecret)
	handler.NewReportHandler(r, reportUsecase, jwtSecret)

	// 13. Jalankan Server di port 8080
	log.Println("Server Tapkasir berjalan di http://localhost:8080")
	r.Run(":8080")
}
