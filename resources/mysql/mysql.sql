0

CREATE TABLE bundles (
    id INT NOT NULL AUTO_INCREMENT,
    bundle VARCHAR(512) NOT NULL,
    category VARCHAR(255) NOT NULL,
    creation_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY unique_bundle_category (bundle, category)
);


1 . Failed Bundles :

CREATE TABLE failed_bundles (
    id INT NOT NULL AUTO_INCREMENT,
    bundle VARCHAR(512) NOT NULL,
    category VARCHAR(255) NOT NULL,
    creation_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY unique_bundle_category (bundle, category)
);


2. Crawled Bundles :

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


3. Lemma Entries :

CREATE TABLE lemma_entries (
    id INT NOT NULL AUTO_INCREMENT,
    bundle VARCHAR(512) NOT NULL,
    category VARCHAR(255) NOT NULL,
    Lemma_Direct TEXT, 
    Lemma_Reseller TEXT, 
    creation_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);


4. Demand Table

id int not null auto increment
bundle    varchar(512) not null
Category  varchar(255) not null
demand columns - dynamic columns
creation_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,


5. Penetration insights Table :

Inventory Type  
Ios inventory int 
android inventory int
Web   inventory int
CTV   inventory int

example : 
Lemma Direct | 75 | 45 | 32 | 23
Lemma Reseller | 56 | 45 | 12 | 12

