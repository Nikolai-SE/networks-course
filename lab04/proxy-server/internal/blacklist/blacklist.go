package blacklist

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Forbiden struct {
	Prefix *string `json:"prefix,omitempty"`
	Infix  *string `json:"infix,omitempty"`
}

func GetBlackList(blacklist string) *[]Forbiden {
	file, err := os.OpenFile(blacklist, os.O_RDONLY, 0)
	if err != nil {
		log.Printf("Blacklist was not found: '%s'.\n", blacklist)
		return nil
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Blacklist read error: %v.\n", err)
		return nil
	}

	var forbidenList []Forbiden
	json.Unmarshal(data, &forbidenList)
	return &forbidenList
}

func IsInBlackList(r *http.Request, forbidenList []Forbiden) bool {
	url := r.URL.String()
	for _, f := range forbidenList {
		if f.Prefix != nil && strings.HasPrefix(url, *f.Prefix) {
			return true
		}
		if f.Infix != nil && strings.Contains(url, *f.Infix) {
			return true
		}
	}

	return false
}
