CREATE TABLE ads_txt_error_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    domain_name VARCHAR(255) NOT NULL,
    error TEXT,
    status_code INT
);

CREATE TABLE failed_bundles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    bundle VARCHAR(255),
    category VARCHAR(255),
    creation_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updation_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    is_deleted TINYINT DEFAULT 0,
    UNIQUE KEY bundle_category (bundle, category)
);

CREATE TABLE `crawled_bundles` (
  `id` int NOT NULL AUTO_INCREMENT,
  `bundle` varchar(255) DEFAULT NULL,
  `category` varchar(255) DEFAULT NULL,
  `website` varchar(512) DEFAULT NULL,
  `domain` varchar(512) DEFAULT NULL,
  `txt_file_url` varchar(512) DEFAULT NULL,
  `creation_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updation_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `is_deleted` tinyint DEFAULT '0',
  `ads_txt_page_hash` varchar(512) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `bundle_category` (`bundle`,`category`)
) ENGINE=InnoDB AUTO_INCREMENT=5510 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci



ALTER TABLE crawled_bundles CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;


