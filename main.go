package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

type DarkerConfig struct {
	outputDir     string
	nginxConf     bool
	rawDirectives []string
}

func parseFlags() *DarkerConfig {
	// flag definitions
	output_dir := flag.String("output-dir", "error_html", "output directory for *.html files")
	nginx_conf := flag.Bool("nginx-conf", false, "generate nginx configuration for mapping static html files")
	// parse flags
	flag.Parse()
	// make sure path is valid
	fileInfo, err := os.Stat(*output_dir)
	if !os.IsNotExist(err) {
		if !fileInfo.IsDir() {
			log.Fatalf("Error: Path '%s' is not a directory\n", *output_dir)
		}
	}
	// create directory
	os.Mkdir(*output_dir, os.FileMode(int(0755)))
	return &DarkerConfig{
		outputDir:     *output_dir,
		nginxConf:     *nginx_conf,
		rawDirectives: flag.Args(),
	}
}

func main() {
	config := parseFlags()
	if config.nginxConf {
		PrintNginxConf(config)
	} else {
		tmpl := SetupTemplate()
		for httpCode := range StatusCodeMap {
			// write to file
			info := MergeWithDefaults(httpCode)
			filepath := path.Join(config.outputDir, fmt.Sprintf("%d.html", httpCode))
			err := RenderFile(tmpl, filepath, info)
			if err != nil {
				log.Fatalf("Error rendering %s: %s\n", filepath, err)
			}
		}
	}
}
