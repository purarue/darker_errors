package darker_errors

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
func renderField(buffer string, httpCode HttpCode) template.HTML {
	return template.HTML(
		strings.ReplaceAll(
			strings.ReplaceAll(buffer, "STATUS_CODE", fmt.Sprintf("%d", httpCode)),
			"STATUS_MSG", StatusCodeMap[httpCode],
		),
	)
}

// Render the template to a file like object
func RenderErrorBuffer(tmpl *template.Template, info *PageInfo, fo io.WriteCloser) error {
	return tmpl.Execute(fo, info)
}

// Render the template and write the result to a filepath
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

// get information from the directives the user gave for this particular
// HTTP status, else use the default
// replace the STATUS_CODE/STATUS_MSG strings and return
// a struct compatible with the template
func GetPageInfo(dMap *DirectiveMap, httpCode HttpCode) *PageInfo {
	return &PageInfo{
		Title:         renderField(dMap.Match(TITLE, httpCode), httpCode),
		Heading:       renderField(dMap.Match(HEADING, httpCode), httpCode),
		Message:       renderField(dMap.Match(MESSAGE, httpCode), httpCode),
		HeadHtml:      renderField(dMap.Match(HEADHTML, httpCode), httpCode),
		BeforeHeading: renderField(dMap.Match(BEFOREHEADING, httpCode), httpCode),
		AfterHeading:  renderField(dMap.Match(AFTERHEADING, httpCode), httpCode),
		AfterMessage:  renderField(dMap.Match(AFTERMESSAGE, httpCode), httpCode),
	}
}
