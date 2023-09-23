package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	dark "github.com/seanbreckenridge/darker_errors/src"
)

type DarkerConfig struct {
	outputDir  string
	nginxConf  bool
	directives *dark.DirectiveMap
}

func parseFlags() *DarkerConfig {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `A dark-themed HTTP error page generator
Additional positional arguments are interpreted as replacement directives
For more information see https://github.com/seanbreckenridge/darker_errors
`)
		flag.PrintDefaults()
	}
	// flag definitions
	output_dir := flag.String("output-dir", "error_html", "output directory for *.html files")
	nginx_conf := flag.Bool("nginx-conf", false, "generate nginx configuration for mapping static html files")
	// parse flags
	flag.Parse()
	// make sure path is valid
	fileInfo, err := os.Stat(*output_dir)
	if !os.IsNotExist(err) {
		if !fileInfo.IsDir() {
			fmt.Fprintf(os.Stderr, "Error: Path '%s' is not a directory\n", *output_dir)
			os.Exit(1)
		}
	}
	// use DirectiveMap in main to generate a PageInfo after editing MergeWithDefaults
	var directives []dark.Directive
	for _, rawDirective := range flag.Args() {
		parsedDir, err := dark.ParseDirective(rawDirective)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing directive: %s\n", err)
			os.Exit(1)
		}
		directives = append(directives, *parsedDir)
	}
	return &DarkerConfig{
		outputDir:  *output_dir,
		nginxConf:  *nginx_conf,
		directives: dark.NewDirectiveMap(directives),
	}
}

func main() {
	config := parseFlags()
	if config.nginxConf {
		dark.PrintNginxConf(config.outputDir)
	} else {
		// create directory
		os.Mkdir(config.outputDir, os.FileMode(int(0755)))
		tmpl := dark.DarkTheme()
		for httpCode := range dark.StatusCodeMap {
			// the values to interpolate into the HTML template
			// if the user provided values those are used, else uses
			// the defaults
			pageInfo := dark.GetPageInfo(config.directives, httpCode)
			filepath := path.Join(config.outputDir, fmt.Sprintf("%d.html", httpCode))
			err := dark.RenderErrorFile(tmpl, pageInfo, filepath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error rendering %s: %s\n", filepath, err)
				os.Exit(1)
			}
		}
	}
}
