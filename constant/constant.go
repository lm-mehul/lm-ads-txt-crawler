package constant

import (
	"os"

	"github.com/lemmamedia/ads-txt-crawler/logger"
)

const (
	BUNDLE_MOBILE_ANDROID = "Mobile App Android"
	BUNDLE_CTV            = "CTV"
	BUNDLE_WEB            = "Web"
	BUNDLE_MOBILE_IOS     = "Mobile App IOS"

	BATCH_SIZE = 800
)

const (
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"
)

var (
	CurrentDirectory string
	Headers          = map[string]string{
		"User-Agent": UserAgent,
	}

	IABdata = map[string]string{
		"arts & entertainment":    "IAB1",
		"automotive":              "IAB2",
		"business":                "IAB3",
		"careers":                 "IAB4",
		"education":               "IAB5",
		"family & parenting":      "IAB6",
		"health & fitness":        "IAB7",
		"food & drink":            "IAB8",
		"hobbies & interests":     "IAB9",
		"home & garden":           "IAB10",
		"law, gov't & politics":   "IAB11",
		"news":                    "IAB12",
		"personal finance":        "IAB13",
		"society":                 "IAB14",
		"science":                 "IAB15",
		"pets":                    "IAB16",
		"sports":                  "IAB17",
		"style & fashion":         "IAB18",
		"technology & computing":  "IAB19",
		"travel":                  "IAB20",
		"real estate":             "IAB21",
		"shopping":                "IAB22",
		"religion & spirituality": "IAB23",
		"uncategorized":           "IAB24",
		"non-standard content":    "IAB25",
		"illegal content":         "IAB26",
	}
)

func InitConstants() {
	var err error
	CurrentDirectory, err = os.Getwd()
	if err != nil {
		logger.Error("Failed to fetch current directory with error : %v\n", err)
	}
}
