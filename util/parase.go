package echoapp_util

import "regexp"

func FetchImageUrls(text string) []string {
	reg := regexp.MustCompile(`\<img .*?title=\"查看大图\" .*?data-src=\"(\S*)\" .*?\>`)
	imgTags := reg.FindAllStringSubmatch(text, -1)
	urls := make([]string, 0)
	for _, tag := range imgTags {
		url := tag[1]
		urls = append(urls, url)
	}
	return urls
}
