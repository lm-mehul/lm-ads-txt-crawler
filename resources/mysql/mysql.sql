
CREATE DATABASE IF NOT EXISTS `lm_teda_crawler`;

-- 0 . Bundles :

CREATE TABLE bundles (
    id INT NOT NULL AUTO_INCREMENT,
    bundle VARCHAR(512) NOT NULL,
    category VARCHAR(255) NOT NULL,
    creation_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY unique_bundle_category (bundle, category)
);


-- 1 . Failed Bundles :

CREATE TABLE failed_bundles (
    id INT NOT NULL AUTO_INCREMENT,
    bundle VARCHAR(512) NOT NULL,
    category VARCHAR(255) NOT NULL,
    creation_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY unique_bundle_category (bundle, category)
);


-- 2. Crawled Bundles :

CREATE TABLE crawled_bundles (
    id INT NOT NULL AUTO_INCREMENT,
    bundle VARCHAR(512) NOT NULL,
    category VARCHAR(255) NOT NULL,
    website VARCHAR(512) NOT NULL,
    domain VARCHAR(512) NOT NULL,
    ads_txt_URL VARCHAR(512) NOT NULL,
    app_ads_txt_URL VARCHAR(512) NOT NULL,
    ads_txt_Hash VARCHAR(512),
    app_ads_txt_Hash VARCHAR(512),
    creation_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY unique_bundle_category (bundle, category)
);


-- 3. Lemma Entries :

CREATE TABLE `lemma_entries` (
  `id` int NOT NULL AUTO_INCREMENT,
  `bundle` varchar(512) NOT NULL,
  `category` varchar(255) NOT NULL,
  `Lemma_Direct` text,
  `Lemma_Reseller` text,
  `ads_page_url` varchar(512) DEFAULT NULL,
  `page_type` varchar(32) DEFAULT NULL,
  `creation_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_bundle_category_lemma` (`bundle`(191),`category`(191),`Lemma_Direct`(100),`Lemma_Reseller`(100))
);

-- 4. Demand Table

CREATE TABLE bundle_demand_lines (
    id INT AUTO_INCREMENT PRIMARY KEY,
    bundle_id INT NOT NULL,
    `category` varchar(255) NOT NULL,
    domain VARCHAR(512) NOT NULL,
    demand_line VARCHAR(255) NOT NULL,
    `ads_page_url` varchar(512) DEFAULT NULL,
    `page_type` varchar(32) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



-- 5. Penetration insights Table :

-- Inventory Type  
-- Ios inventory int 
-- android inventory int
-- Web   inventory int
-- CTV   inventory int

-- example : 
-- Lemma Direct | 75 | 45 | 32 | 23
-- Lemma Reseller | 56 | 45 | 12 | 12
