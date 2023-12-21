-- --------------------------------------------------------
-- Sunucu:                       127.0.0.1
-- Sunucu sürümü:                10.11.2-MariaDB - mariadb.org binary distribution
-- Sunucu İşletim Sistemi:       Win64
-- HeidiSQL Sürüm:               11.3.0.6295
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


-- movder için veritabanı yapısı dökülüyor
CREATE DATABASE IF NOT EXISTS `movder` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */;
USE `movder`;

-- tablo yapısı dökülüyor movder.ayar
CREATE TABLE IF NOT EXISTS `ayar` (
  `ayar_id` int(11) NOT NULL AUTO_INCREMENT,
  `ayar_ad` varchar(50) DEFAULT NULL,
  `ayar_deger` text DEFAULT NULL,
  PRIMARY KEY (`ayar_id`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

-- Veri çıktısı seçilmemişti

-- tablo yapısı dökülüyor movder.blog
CREATE TABLE IF NOT EXISTS `blog` (
  `blog_id` int(11) NOT NULL AUTO_INCREMENT,
  `blog_url` varchar(255) DEFAULT NULL,
  `blog_title` varchar(255) DEFAULT NULL,
  `blog_overview` text DEFAULT NULL,
  `blog_content` text DEFAULT NULL,
  `blog_owner` int(11) DEFAULT NULL,
  `blog_date` datetime DEFAULT current_timestamp(),
  `blog_popular` int(11) DEFAULT NULL,
  `blog_image` varchar(50) DEFAULT NULL,
  `blog_tags` text DEFAULT NULL,
  `blog_category` int(11) DEFAULT NULL,
  `blog_comment` int(11) DEFAULT NULL,
  PRIMARY KEY (`blog_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Veri çıktısı seçilmemişti

-- tablo yapısı dökülüyor movder.category_blog
CREATE TABLE IF NOT EXISTS `category_blog` (
  `category_id` int(11) DEFAULT NULL,
  `category_url` varchar(50) DEFAULT NULL,
  `category_name` varchar(50) DEFAULT NULL,
  `category_color` varchar(8) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Veri çıktısı seçilmemişti

-- tablo yapısı dökülüyor movder.client_chat
CREATE TABLE IF NOT EXISTS `client_chat` (
  `client_id` int(11) NOT NULL AUTO_INCREMENT,
  `client_key` varchar(50) DEFAULT NULL,
  `client_owner` int(11) DEFAULT NULL,
  `client_target` int(11) DEFAULT NULL,
  `client_view` int(11) DEFAULT 0,
  `client_didntviewmsg` int(11) DEFAULT 0,
  PRIMARY KEY (`client_id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Veri çıktısı seçilmemişti

-- tablo yapısı dökülüyor movder.comment
CREATE TABLE IF NOT EXISTS `comment` (
  `comment_id` int(11) NOT NULL AUTO_INCREMENT,
  `comment_film` int(11) DEFAULT NULL,
  `comment_owner` int(11) DEFAULT NULL,
  `comment_content` text DEFAULT NULL,
  `comment_date` datetime DEFAULT current_timestamp(),
  `comment_rate` int(11) DEFAULT NULL,
  `comment_like` int(11) DEFAULT NULL,
  `comment_dislike` int(11) DEFAULT NULL,
  `comment_filmtype` varchar(15) DEFAULT NULL,
  PRIMARY KEY (`comment_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

-- Veri çıktısı seçilmemişti

-- tablo yapısı dökülüyor movder.comment_blog
CREATE TABLE IF NOT EXISTS `comment_blog` (
  `comment_id` int(11) NOT NULL AUTO_INCREMENT,
  `comment_blog` int(11) DEFAULT NULL,
  `comment_owner` int(11) DEFAULT NULL,
  `comment_content` text DEFAULT NULL,
  `comment_date` datetime DEFAULT current_timestamp(),
  `comment_rate` int(11) DEFAULT NULL,
  `comment_like` int(11) DEFAULT NULL,
  `comment_dislike` int(11) DEFAULT NULL,
  PRIMARY KEY (`comment_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

-- Veri çıktısı seçilmemişti

-- tablo yapısı dökülüyor movder.favorite
CREATE TABLE IF NOT EXISTS `favorite` (
  `favorite_id` int(11) NOT NULL AUTO_INCREMENT,
  `favorite_film` int(11) DEFAULT NULL,
  `favorite_owner` int(11) DEFAULT NULL,
  `favorite_date` datetime DEFAULT current_timestamp(),
  `favorite_filmtype` varchar(15) DEFAULT NULL,
  PRIMARY KEY (`favorite_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

-- Veri çıktısı seçilmemişti

-- tablo yapısı dökülüyor movder.follow
CREATE TABLE IF NOT EXISTS `follow` (
  `follow_id` int(11) NOT NULL AUTO_INCREMENT,
  `follow_user` int(11) DEFAULT NULL,
  `follow_owner` int(11) DEFAULT NULL,
  PRIMARY KEY (`follow_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Veri çıktısı seçilmemişti

-- tablo yapısı dökülüyor movder.iletisim
CREATE TABLE IF NOT EXISTS `iletisim` (
  `iletisim_id` int(11) NOT NULL AUTO_INCREMENT,
  `iletisim_isim` varchar(50) DEFAULT NULL,
  `iletisim_eposta` varchar(50) DEFAULT NULL,
  `iletisim_konu` varchar(100) DEFAULT NULL,
  `iletisim_mesaj` text DEFAULT NULL,
  `iletisim_Tarih` datetime DEFAULT NULL,
  PRIMARY KEY (`iletisim_id`)
) ENGINE=InnoDB AUTO_INCREMENT=101 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

-- Veri çıktısı seçilmemişti

-- tablo yapısı dökülüyor movder.kullanici
CREATE TABLE IF NOT EXISTS `kullanici` (
  `kullanici_id` int(11) NOT NULL AUTO_INCREMENT,
  `kullanici_nick` varchar(50) DEFAULT NULL,
  `kullanici_isim` varchar(50) DEFAULT NULL,
  `kullanici_posta` varchar(50) DEFAULT NULL,
  `kullanici_sifre` varchar(50) DEFAULT NULL,
  `kullanici_yetki` varchar(50) DEFAULT NULL,
  `kullanici_url` varchar(50) DEFAULT NULL,
  `kullanici_avatar` text DEFAULT NULL,
  `kullanici_hakkinda` text DEFAULT NULL,
  `kullanici_yazi` int(11) DEFAULT NULL,
  `kullanici_instagram` varchar(50) DEFAULT NULL,
  `kullanici_twitter` varchar(50) DEFAULT NULL,
  `kullanici_facebook` varchar(50) DEFAULT NULL,
  `kullanici_website` varchar(50) DEFAULT NULL,
  `kullanici_location` varchar(255) DEFAULT NULL,
  `kullanici_online` varchar(50) DEFAULT NULL,
  `kullanici_telephone` varchar(50) DEFAULT NULL,
  `kullanici_reviewer` int(11) DEFAULT 0,
  PRIMARY KEY (`kullanici_id`)
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

-- Veri çıktısı seçilmemişti

-- tablo yapısı dökülüyor movder.list
CREATE TABLE IF NOT EXISTS `list` (
  `list_id` int(11) NOT NULL AUTO_INCREMENT,
  `list_owner` int(11) DEFAULT NULL,
  `list_title` varchar(50) DEFAULT NULL,
  `list_url` varchar(50) DEFAULT NULL,
  `list_count` int(11) DEFAULT NULL,
  `list_favorite` int(11) DEFAULT 0,
  `list_type` varchar(50) DEFAULT NULL,
  `list_date` datetime DEFAULT current_timestamp(),
  `list_image` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`list_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

-- Veri çıktısı seçilmemişti

-- tablo yapısı dökülüyor movder.list_content
CREATE TABLE IF NOT EXISTS `list_content` (
  `listcontent_id` int(11) NOT NULL AUTO_INCREMENT,
  `listcontent_owner` int(11) DEFAULT NULL,
  `listcontent_type` varchar(50) DEFAULT NULL,
  `listcontent_film` int(11) DEFAULT NULL,
  PRIMARY KEY (`listcontent_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Veri çıktısı seçilmemişti

-- tablo yapısı dökülüyor movder.message
CREATE TABLE IF NOT EXISTS `message` (
  `message_id` int(11) NOT NULL AUTO_INCREMENT,
  `message_owner` int(11) DEFAULT NULL,
  `message_target` int(11) DEFAULT NULL,
  `message_content` text DEFAULT NULL,
  `message_view` int(11) DEFAULT NULL,
  `message_date` datetime NOT NULL DEFAULT current_timestamp(),
  `message_delete` int(11) DEFAULT NULL,
  `message_owner_delete` int(11) DEFAULT NULL,
  `message_target_delete` int(11) DEFAULT NULL,
  PRIMARY KEY (`message_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=614 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Veri çıktısı seçilmemişti

-- tablo yapısı dökülüyor movder.message_box
CREATE TABLE IF NOT EXISTS `message_box` (
  `box_id` int(11) NOT NULL AUTO_INCREMENT,
  `box_owner` int(11) DEFAULT NULL,
  `box_target` int(11) DEFAULT NULL,
  `box_count` int(11) DEFAULT NULL,
  `box_date` datetime NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`box_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

-- Veri çıktısı seçilmemişti

-- tablo yapısı dökülüyor movder.notification
CREATE TABLE IF NOT EXISTS `notification` (
  `notification_id` int(11) NOT NULL AUTO_INCREMENT,
  `notification_owner` int(11) DEFAULT NULL,
  `notification_target` int(11) DEFAULT NULL,
  `notification_content` text DEFAULT NULL,
  `notification_view` int(11) DEFAULT NULL,
  PRIMARY KEY (`notification_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Veri çıktısı seçilmemişti

-- tablo yapısı dökülüyor movder.rate
CREATE TABLE IF NOT EXISTS `rate` (
  `rate_id` int(11) NOT NULL AUTO_INCREMENT,
  `rate_owner` int(11) DEFAULT NULL,
  `rate_film` int(11) DEFAULT NULL,
  `rate_filmtype` varchar(15) DEFAULT NULL,
  `rate_value` int(11) DEFAULT 0,
  `rate_date` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`rate_id`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Veri çıktısı seçilmemişti

-- tablo yapısı dökülüyor movder.table
CREATE TABLE IF NOT EXISTS `table` (
  `id` int(11) NOT NULL,
  `title` varchar(50) DEFAULT NULL,
  `content` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

-- Veri çıktısı seçilmemişti

-- tablo yapısı dökülüyor movder.watchedlist
CREATE TABLE IF NOT EXISTS `watchedlist` (
  `watchedlist_id` int(11) NOT NULL AUTO_INCREMENT,
  `watchedlist_film` int(11) DEFAULT NULL,
  `watchedlist_owner` int(11) DEFAULT NULL,
  `watchedlist_date` datetime DEFAULT current_timestamp(),
  `watchedlist_diary` int(11) DEFAULT 0,
  `watchedlist_filmtype` varchar(14) DEFAULT NULL,
  PRIMARY KEY (`watchedlist_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=160 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;

-- Veri çıktısı seçilmemişti

-- tablo yapısı dökülüyor movder.watchlist
CREATE TABLE IF NOT EXISTS `watchlist` (
  `watchlist_id` int(11) NOT NULL AUTO_INCREMENT,
  `watchlist_film` int(11) DEFAULT NULL,
  `watchlist_filmtype` varchar(15) DEFAULT NULL,
  `watchlist_owner` int(11) DEFAULT NULL,
  `watchlist_date` datetime DEFAULT current_timestamp(),
  PRIMARY KEY (`watchlist_id`)
) ENGINE=InnoDB AUTO_INCREMENT=96 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Veri çıktısı seçilmemişti

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
