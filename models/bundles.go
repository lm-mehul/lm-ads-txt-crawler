package models

import "time"

type BundleInfo struct {
	Id            int64
	Bundle        string
	Category      string
	Website       string
	Domain        string
	AdsTxtURL     string
	AppAdsTxtURL  string
	AdsTxtHash    string
	AppAdsTxtHash string
	CreatedAt     time.Time
}

// LemmaEntry represents a record in the lemma_entries table.
type LemmaEntry struct {
	Bundle        string    `json:"bundle"`
	Category      string    `json:"category"`
	AdsPageURL    string    `json:"ads_page_url"`
	PageType      string    `json:"page_type"`
	LemmaDirect   string    `json:"lemma_direct"`
	LemmaReseller string    `json:"lemma_reseller"`
	CreatedAt     time.Time `json:"creation_time"`
}

type DemandLinesEntry struct {
	Bundle     string    `json:"bundle"`
	Category   string    `json:"category"`
	AdsPageURL string    `json:"ads_page_url"`
	PageType   string    `json:"page_type"`
	DemandLine string    `json:"demand_line"`
	CreatedAt  time.Time `json:"creation_time"`
}

type AdsTxtRecord struct {
	AdSystemDomain string `json:"ad_system_domain"`
	PublisherID    string `json:"publisher_id"`
	Relationship   string `json:"relationship"`
	CertAuthID     string `json:"cert_auth_id"`
}
