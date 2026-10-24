package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	// Sesuaikan nama module import dengan yang ada di go.mod Anda
	"tapkasir-pos/internal/delivery/http/handler"
	"tapkasir-pos/internal/model"
	"tapkasir-pos/internal/repository"
	"tapkasir-pos/internal/usecase"
)

func main() {
	// 1. Konfigurasi Koneksi Database (Default Laragon biasanya user: root, password: kosong)
	dsn := "root:@tcp(127.0.0.1:3306)/tapkasir_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	// Rename existing plural tables before migration to preserve their data.
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

	// Auto-Migrate creates or updates the singularly named tables.
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

	// 3. Inisialisasi Gin (Web Framework)
	r := gin.Default()

	// 4. Wiring Dependency Injection
	// Memasukkan koneksi DB ke Repository
	menuRepo := repository.NewMenuRepository(db)
	// Memasukkan Repository ke Usecase
	menuUsecase := usecase.NewMenuUsecase(menuRepo)
	// Memasukkan Usecase ke Handler dan mendaftarkan route ke Gin (r)
	handler.NewMenuHandler(r, menuUsecase)

	// 5. Jalankan Server di port 8080
	log.Println("Server Tapkasir berjalan di http://localhost:8080")
	r.Run(":8080")
}
