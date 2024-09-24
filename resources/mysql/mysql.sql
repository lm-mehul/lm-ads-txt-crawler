
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
    bundle_id varchar(412) NOT NULL,
    `category` varchar(255) NOT NULL,
    domain VARCHAR(512) NOT NULL,
    demand_line VARCHAR(255) NOT NULL,
    `ads_page_url` varchar(512) DEFAULT NULL,
    `page_type` varchar(32) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE `ads_txt_demand_lines` (
  `demand_line` varchar(512) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`demand_line`)
);


INSERT IGNORE INTO ads_txt_demand_lines (demand_line) VALUES
('pubmatic.com, 156494, RESELLER'),
('onetag.com, 8d4bbebb8fa701e, DIRECT'),
('freewheel.tv, 1585652, RESELLER'),
('freewheel.tv, 1585681, RESELLER'),
('video.unrulymedia.com, 8167205979129043832, RESELLER'),
('growzee.io, 031ae447-d523-4804-b4c7-be282ad874c3, DIRECT'),
('instal.com, 031ae447-d523-4804-b4c7-be282ad874c3, DIRECT'),
('boldscreen.com, 765167, DIRECT'),
('boldscreen.com, 46220, DIRECT'),
('unifd.la, 765167, DIRECT'),
('unifd.la, 46220, DIRECT'),
('advangelists.com, 57b9c5d42a8fe216162345a05b6e6afb, RESELLER'),
('magicmedia.ae, 6f828808, DIRECT'),
('e-planning.net, 608359bb987625b2, DIRECT, c1ba615865ed87b2'),
('connatix.com, 1742365173855029, DIRECT'),
('appnexus.com, 2007, RESELLER'),
('affinity.com, 965, RESELLER'),
('smartadserver.com, 3564, RESELLER'),
('Contextweb.com, 563169, RESELLER'),
('Contextweb.com, 562930, RESELLER'),
('thegermanemedia.com, 546, DIRECT, 6716225838534735'),
('smartadserver.com, 4938, RESELLER, 060d053dcf45cbf3'),
('pubmatic.com, 164928, RESELLER, 5d62403b186f2ace'),
('freewheel.tv, 1607251, RESELLER'),
('lijit.com, 509119, DIRECT, fafdf38b16bf6b2b'),
('Vidoomy.com, 8700140, DIRECT');



-- 5. Penetration insights Table :

-- Inventory Type  
-- Ios inventory int 
-- android inventory int
-- Web   inventory int
-- CTV   inventory int

-- example : 
-- Lemma Direct | 75 | 45 | 32 | 23
-- Lemma Reseller | 56 | 45 | 12 | 12
