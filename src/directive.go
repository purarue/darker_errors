package darker_errors

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	TITLE         = 1
	HEADING       = 2
	MESSAGE       = 3
	HEADHTML      = 4
	BEFOREHEADING = 5
	AFTERHEADING  = 6
	AFTERMESSAGE  = 7
)

type DirectiveId = uint

/// a parsed directive from the user
type Directive struct {
	/// the HTTP code to apply this directive to
	/// NoHttpCode is a sentinel value (-1) which
	/// means this should apply to all items
	HttpCode int
	/// which part of the page the replacement should be inserted into
	DirectiveId DirectiveId
	/// an HTML string to replace into the template
	Replacement string
}

func ParseDirectiveId(directiveName string) (DirectiveId, error) {
	switch directiveName {
	case "ERROR_TITLE":
		return TITLE, nil
	case "ERROR_HEADING":
		return HEADING, nil
	case "ERROR_MSG":
		return MESSAGE, nil
	case "ERROR_HEAD":
		return HEADHTML, nil
	case "ERROR_BEFORE_HEADING":
		return BEFOREHEADING, nil
	case "ERROR_AFTER_HEADING":
		return AFTERHEADING, nil
	case "ERROR_AFTER_MSG":
		return AFTERMESSAGE, nil
	}
	return 0, errors.New(fmt.Sprintf("Could not parse %s into a valid directive name", directiveName))
}

func ParseDirective(directive string) (*Directive, error) {
	// two possible formats
	// one where this includes the code, like:
	// 502:ERROR_MSG:<p>The site is down</p>
	// and another where it should apply to all pages:
	// ERROR_TITLE:My Site name - STATUS_CODE
	idx := strings.Index(directive, ":")
	if idx == -1 {
		return nil, errors.New(fmt.Sprintf("Could not find a ':' in %s, no way to parse fields", directive))
	}
	// by default, this should apply to all pages
	httpCode := NoHttpCode
	// try to parse into HTTP code
	httpParsed, err := strconv.Atoi(directive[0:idx])
	if err == nil {
		// http code was parsed
		// save value and modify string so that it now starts with the directive name
		httpCode = httpParsed
		directive = directive[idx+1:]
		idx = strings.Index(directive, ":")
	}
	// try to parse directive
	if idx == -1 {
		return nil, errors.New(fmt.Sprintf("Could not find a ':' in %s, no way to parse fields", directive))
	}
	dirId, err := ParseDirectiveId(directive[0:idx])
	if err != nil {
		return nil, err
	}
	return &Directive{
		HttpCode:    httpCode,
		DirectiveId: dirId,
		Replacement: directive[idx+1:],
	}, nil
}

var DefaultDirectiveValues = map[DirectiveId]string{
	TITLE:         "STATUS_CODE - STATUS_MSG",
	HEADING:       "<h2>STATUS_CODE</h2>",
	MESSAGE:       "<p>STATUS_MSG</p>",
	HEADHTML:      "",
	BEFOREHEADING: "",
	AFTERHEADING:  "",
	AFTERMESSAGE:  "",
}

// maps the directive id to a list of directives (which may
// have different http statuses/no http status)
// if none exist, the default should be used.
type DirectiveMap struct {
	forDirective map[DirectiveId][]Directive
}

/// Create a new directive map
func NewDirectiveMap(directives []Directive) *DirectiveMap {
	// TODO: use 2 maps here? could improve speed a bit; but there aren't that many HTTP codes anyways
	// number of custom directives is probably low, so this is fine
	fDirective := make(map[DirectiveId][]Directive)
	for _, d := range directives {
		fDirective[d.DirectiveId] = append(fDirective[d.DirectiveId], d)
	}
	return &DirectiveMap{
		forDirective: fDirective,
	}
}

// return the Replacement if one matches this code/directiveId,
// else the default
func (dm *DirectiveMap) Match(id DirectiveId, httpCode HttpCode) string {
	var genericMatch *Directive = nil
	for _, dr := range dm.forDirective[id] {
		// a generic directive, like
		// ERROR_TITLE:My Site name - STATUS_CODE
		// if a http code - specific directive wasn't specified,
		// return the generic directive that was provided by the user
		if dr.HttpCode == NoHttpCode {
			genericMatch = &dr
		}
		// this specific httpCode matches this directive
		if httpCode != NoHttpCode && httpCode == dr.HttpCode {
			return dr.Replacement
		}
	}
	// if we didnt match a specific code, but a generic directive existed, use that
	if genericMatch != nil {
		return genericMatch.Replacement
	}
	// if nothing matched, return the default value
	return DefaultDirectiveValues[id]
}
