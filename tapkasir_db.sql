-- --------------------------------------------------------
-- Host:                         127.0.0.1
-- Server version:               8.4.3 - MySQL Community Server - GPL
-- Server OS:                    Win64
-- HeidiSQL Version:             12.8.0.6908
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


-- Dumping database structure for tapkasir_db
CREATE DATABASE IF NOT EXISTS `tapkasir_db` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `tapkasir_db`;

-- Dumping structure for table tapkasir_db.bahan_baku
CREATE TABLE IF NOT EXISTS `bahan_baku` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `nama_bahan` varchar(100) NOT NULL,
  `harga_beli` double NOT NULL,
  `stok_sistem` double NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `satuan` varchar(20) DEFAULT 'gram',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Dumping data for table tapkasir_db.bahan_baku: ~1 rows (approximately)
INSERT INTO `bahan_baku` (`id`, `nama_bahan`, `harga_beli`, `stok_sistem`, `created_at`, `updated_at`, `satuan`) VALUES
	(1, 'Biji Kopi Arabika', 250, 6782, '2026-11-17 00:27:12.132', '2026-11-19 12:28:25.714', 'gram'),
	(2, 'Susu UHT', 20, 9350, '2026-11-19 11:11:40.469', '2026-11-19 12:28:25.719', 'ml'),
	(3, 'Gula Aren', 50, 2980, '2026-11-19 11:11:53.378', '2026-11-19 12:28:25.724', 'gram');

-- Dumping structure for table tapkasir_db.detail_pesanan
CREATE TABLE IF NOT EXISTS `detail_pesanan` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `pesanan_id` bigint unsigned NOT NULL,
  `menu_id` bigint unsigned NOT NULL,
  `kuantitas` bigint NOT NULL,
  `sub_total` double NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_detail_pesanan_pesanan_id` (`pesanan_id`),
  KEY `idx_detail_pesanan_menu_id` (`menu_id`),
  CONSTRAINT `fk_detail_pesanan_menu` FOREIGN KEY (`menu_id`) REFERENCES `menu` (`id`),
  CONSTRAINT `fk_pesanan_detail_pesanan` FOREIGN KEY (`pesanan_id`) REFERENCES `pesanan` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Dumping data for table tapkasir_db.detail_pesanan: ~2 rows (approximately)
INSERT INTO `detail_pesanan` (`id`, `pesanan_id`, `menu_id`, `kuantitas`, `sub_total`) VALUES
	(1, 1, 2, 2, 36000),
	(2, 2, 2, 1, 18000),
	(3, 3, 3, 2, 64000),
	(4, 4, 2, 1, 18000),
	(5, 5, 3, 1, 32000),
	(6, 6, 3, 1, 32000),
	(7, 7, 2, 1, 18000);

-- Dumping structure for table tapkasir_db.log_barang_rusak
CREATE TABLE IF NOT EXISTS `log_barang_rusak` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `bahan_baku_id` bigint unsigned NOT NULL,
  `jumlah_terbuang` double NOT NULL,
  `alasan` text,
  `waktu_pencatatan` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_log_barang_rusak_bahan_baku_id` (`bahan_baku_id`),
  CONSTRAINT `fk_log_barang_rusak_bahan_baku` FOREIGN KEY (`bahan_baku_id`) REFERENCES `bahan_baku` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Dumping data for table tapkasir_db.log_barang_rusak: ~0 rows (approximately)
INSERT INTO `log_barang_rusak` (`id`, `bahan_baku_id`, `jumlah_terbuang`, `alasan`, `waktu_pencatatan`) VALUES
	(1, 1, 150, 'Bungkus sobek terkena tumpahan air', '2026-11-17 01:14:16.983'),
	(2, 2, 500, 'Susu kedaluwarsa', '2026-11-19 11:26:19.025');

-- Dumping structure for table tapkasir_db.log_stok_masuk
CREATE TABLE IF NOT EXISTS `log_stok_masuk` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `bahan_baku_id` bigint unsigned NOT NULL,
  `jumlah_masuk` double NOT NULL,
  `pemasok` varchar(150) DEFAULT NULL,
  `waktu_masuk` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_log_stok_masuk_bahan_baku_id` (`bahan_baku_id`),
  CONSTRAINT `fk_log_stok_masuk_bahan_baku` FOREIGN KEY (`bahan_baku_id`) REFERENCES `bahan_baku` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Dumping data for table tapkasir_db.log_stok_masuk: ~0 rows (approximately)
INSERT INTO `log_stok_masuk` (`id`, `bahan_baku_id`, `jumlah_masuk`, `pemasok`, `waktu_masuk`) VALUES
	(1, 1, 2000, 'PT Kopi Mandiri Nusantara', '2026-11-17 01:13:06.388'),
	(2, 1, 2000, 'PT Kopi Mandiri', '2026-11-19 11:25:25.293');

-- Dumping structure for table tapkasir_db.log_void
CREATE TABLE IF NOT EXISTS `log_void` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `pesanan_id` bigint unsigned NOT NULL,
  `alasan_batal` text NOT NULL,
  `waktu_batal` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_log_void_pesanan_id` (`pesanan_id`),
  CONSTRAINT `fk_pesanan_log_void` FOREIGN KEY (`pesanan_id`) REFERENCES `pesanan` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Dumping data for table tapkasir_db.log_void: ~0 rows (approximately)
INSERT INTO `log_void` (`id`, `pesanan_id`, `alasan_batal`, `waktu_batal`) VALUES
	(1, 4, 'Berubah pikiran', '2026-11-19 11:37:39.003'),
	(2, 6, 'Pelanggan salah pesan menu, sudah konfirmasi', '2026-11-19 11:48:37.515');

-- Dumping structure for table tapkasir_db.menu
CREATE TABLE IF NOT EXISTS `menu` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `nama_menu` longtext,
  `harga_jual` double DEFAULT NULL,
  `status_ketersediaan` longtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Dumping data for table tapkasir_db.menu: ~1 rows (approximately)
INSERT INTO `menu` (`id`, `nama_menu`, `harga_jual`, `status_ketersediaan`) VALUES
	(2, 'Kopi Susu Gula Aren', 18000, 'tersedia'),
	(3, 'Matcha Latte', 32000, 'tersedia');

-- Dumping structure for table tapkasir_db.pesanan
CREATE TABLE IF NOT EXISTS `pesanan` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `nomor_meja` bigint NOT NULL,
  `nama_pelanggan` varchar(100) DEFAULT NULL,
  `nomor_telepon` varchar(20) DEFAULT NULL,
  `metode_pembayaran` varchar(20) NOT NULL,
  `status_pembayaran` varchar(30) NOT NULL DEFAULT 'unpaid',
  `status_pesanan` varchar(30) NOT NULL DEFAULT 'menunggu',
  `total_harga` double NOT NULL,
  `waktu_pesanan` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Dumping data for table tapkasir_db.pesanan: ~2 rows (approximately)
INSERT INTO `pesanan` (`id`, `nomor_meja`, `nama_pelanggan`, `nomor_telepon`, `metode_pembayaran`, `status_pembayaran`, `status_pesanan`, `total_harga`, `waktu_pesanan`) VALUES
	(1, 5, 'Andi Pratama', '081234567890', 'tunai', 'paid', 'selesai', 36000, '2026-11-17 01:29:01.849'),
	(2, 3, 'Siti Nurhaliza', '089876543210', 'qris', 'paid', 'selesai', 18000, '2026-11-17 01:29:42.954'),
	(3, 10, 'Budi Santoso', '08123456789', 'qris', 'paid', 'dimasak', 64000, '2026-11-19 11:29:47.355'),
	(4, 2, 'Sari Dewi', '08567890123', 'tunai', 'unpaid', 'dibatalkan', 18000, '2026-11-19 11:37:11.902'),
	(5, 7, 'Andi Wijaya', '08211234567', 'tunai', 'paid', 'selesai', 32000, '2026-11-19 11:38:06.611'),
	(6, 15, 'Void Test', '08900000001', 'qris', 'unpaid', 'dibatalkan', 32000, '2026-11-19 11:45:40.752'),
	(7, 1, 'bagas', '08987364873', 'qris', 'paid', 'dimasak', 18000, '2026-11-19 12:27:10.847');

-- Dumping structure for table tapkasir_db.resep
CREATE TABLE IF NOT EXISTS `resep` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `menu_id` bigint unsigned NOT NULL,
  `bahan_baku_id` bigint unsigned NOT NULL,
  `takaran` double NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_resep_menu_id` (`menu_id`),
  KEY `idx_resep_bahan_baku_id` (`bahan_baku_id`),
  CONSTRAINT `fk_resep_bahan_baku` FOREIGN KEY (`bahan_baku_id`) REFERENCES `bahan_baku` (`id`),
  CONSTRAINT `fk_resep_menu` FOREIGN KEY (`menu_id`) REFERENCES `menu` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Dumping data for table tapkasir_db.resep: ~0 rows (approximately)
INSERT INTO `resep` (`id`, `menu_id`, `bahan_baku_id`, `takaran`) VALUES
	(2, 2, 1, 18),
	(3, 2, 2, 150),
	(4, 2, 3, 20);

-- Dumping structure for table tapkasir_db.stock_opname
CREATE TABLE IF NOT EXISTS `stock_opname` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `bahan_baku_id` bigint unsigned NOT NULL,
  `stok_sistem` double NOT NULL,
  `stok_fisik` double NOT NULL,
  `selisih` double NOT NULL,
  `waktu_opname` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_stock_opname_bahan_baku_id` (`bahan_baku_id`),
  CONSTRAINT `fk_stock_opname_bahan_baku` FOREIGN KEY (`bahan_baku_id`) REFERENCES `bahan_baku` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Dumping data for table tapkasir_db.stock_opname: ~0 rows (approximately)
INSERT INTO `stock_opname` (`id`, `bahan_baku_id`, `stok_sistem`, `stok_fisik`, `selisih`, `waktu_opname`) VALUES
	(1, 1, 2850, 2800, -50, '2026-11-17 01:14:57.651'),
	(2, 1, 4746, 6800, 2054, '2026-11-19 11:28:02.594');

-- Dumping structure for table tapkasir_db.user
CREATE TABLE IF NOT EXISTS `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `role` enum('owner','manager','cashier') NOT NULL,
  `pin` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Dumping data for table tapkasir_db.user: ~2 rows (approximately)
INSERT INTO `user` (`id`, `role`, `pin`, `created_at`, `updated_at`) VALUES
	(1, 'owner', '$2a$10$Q1N7KhnS773P8Hog756S2.9AHy.iHxUUcrzD9VEbw65wpMTX2bERG', '2026-11-16 23:53:39.706', '2026-11-16 23:53:39.706'),
	(2, 'manager', '$2a$10$3MYDJ.C8s32NxHQNeOyiaOqLyaYb5Jtq5aBErQUOiOt/QQiT10tYK', '2026-11-16 23:53:39.791', '2026-11-16 23:53:39.791'),
	(3, 'cashier', '$2a$10$l8FFeVU/RNRPq7xX1WQgNuWwdBxSq.3cdSkG.8vE13drUWbcFw9AW', '2026-11-16 23:53:39.883', '2026-11-16 23:53:39.883');

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
