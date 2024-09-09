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
	Bundle        string `json:"bundle"`
	Category      string `json:"category"`
	AdsPageURL    string
	PageType      string
	LemmaDirect   string    `json:"lemma_direct"`
	LemmaReseller string    `json:"lemma_reseller"`
	CreatedAt     time.Time `json:"creation_time"`
}

type DemandLinesEntry struct {
	Bundle        string `json:"bundle"`
	Category      string `json:"category"`
	AdsPageURL    string
	PageType      string
	LemmaDirect   string    `json:"lemma_direct"`
	LemmaReseller string    `json:"lemma_reseller"`
	CreatedAt     time.Time `json:"creation_time"`
}
