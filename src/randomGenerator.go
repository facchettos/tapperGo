package gotapper

import (
	"strings"
)

func replaceFromUrl(url string, randomFields map[string]string) string {
	finalUrl := url
	for k, v := range randomFields {
		finalUrl = strings.Replace(finalUrl, k, v, -1)
	}
	return finalUrl
}
