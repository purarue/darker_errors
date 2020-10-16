package main

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
)

type PageInfo struct {
	Title         template.HTML
	Heading       template.HTML
	Message       template.HTML
	HeadHtml      template.HTML
	BeforeHeading template.HTML
	AfterHeading  template.HTML
	AfterMessage  template.HTML
}

// replaces any of the STATUS_CODE/STATUS_MSG tokens on a string,
// and converts it to template.HTML so the values are interpolated
// into the HTML directly instead of being escaped
func renderField(buffer string, httpCode int) template.HTML {
	return template.HTML(
		strings.ReplaceAll(
			strings.ReplaceAll(buffer, "STATUS_CODE", fmt.Sprintf("%d", httpCode)),
			"STATUS_MSG", StatusCodeMap[httpCode],
		),
	)
}

func MergeWithDefaults(httpCode int) *PageInfo {
	title := "STATUS_CODE - STATUS_MSG"
	heading := "<h2>STATUS_CODE</h2>"
	msg := "<p>STATUS_MSG</p>"
	headHtml := ""
	beforeHeading := ""
	afterHeading := ""
	afterMessage := ""
	return &PageInfo{
		Title:         renderField(title, httpCode),
		Heading:       renderField(heading, httpCode),
		Message:       renderField(msg, httpCode),
		HeadHtml:      renderField(headHtml, httpCode),
		BeforeHeading: renderField(beforeHeading, httpCode),
		AfterHeading:  renderField(afterHeading, httpCode),
		AfterMessage:  renderField(afterMessage, httpCode),
	}
}

/// Render the template to a file like object
func RenderErrorBuffer(tmpl *template.Template, info *PageInfo, fo io.WriteCloser) error {
	return tmpl.Execute(fo, info)
}

/// Render the template and write the result to a filepath
func RenderErrorFile(tmpl *template.Template, info *PageInfo, filepath string) error {
	// open file object
	fo, err := os.Create(filepath)
	if err != nil {
		return err
	}
	err = RenderErrorBuffer(tmpl, info, fo)
	if err != nil {
		return err
	}
	return fo.Close()
}
