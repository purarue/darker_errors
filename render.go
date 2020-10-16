package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
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

/// Render tpml to a buffer
func RenderBuffer(tmpl *template.Template, buf *bytes.Buffer, info *PageInfo) error {
	err := tmpl.Execute(buf, info)
	return err
}

/// Render a template and write the result to a filepath
func RenderFile(tmpl *template.Template, filepath string, info *PageInfo) error {
	// create temporary bytes buffer to execute template
	buf := new(bytes.Buffer)
	err := RenderBuffer(tmpl, buf, info)
	if err != nil {
		return err
	}
	// open file object
	fo, err := os.Create(filepath)
	if err != nil {
		return err
	}
	// write html contents to file
	err = ioutil.WriteFile(filepath, buf.Bytes(), os.FileMode(int(0755)))
	if err != nil {
		return err
	}
	err = fo.Close()
	return err
}
