package main

import (
	"fmt"
	"path"
	"sort"
	"strings"
)

func PrintNginxConf(conf *DarkerConfig) {
	folderName := strings.TrimSpace(strings.Trim(conf.outputDir, "./"))
	// To store the keys in slice in sorted order
	statusKeys := make([]int, len(StatusCodeMap))
	i := 0
	for k := range StatusCodeMap {
		statusKeys[i] = k
		i++
	}
	sort.Ints(statusKeys)
	for _, statusCode := range statusKeys {
		folderPath := path.Join(folderName, fmt.Sprintf("%d", statusCode))
		fmt.Printf("error %d = /%s.html\n", statusCode, folderPath)
	}
}
