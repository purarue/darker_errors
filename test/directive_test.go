package darker_errors_test

import (
	dark "github.com/seanbreckenridge/darker_errors/src"
	"testing"
)

func TestParseDirectiveId(t *testing.T) {
	input := "ERROR_HEAD"
	id, err := dark.ParseDirectiveId(input)
	if err != nil {
		t.Errorf("Expected '%s' to parse properly\n", input)
	}
	if id != dark.HEADHTML {
		t.Errorf("Expected '%s' to match HEADHTML const\n", input)
	}

	input = "ERROR_AFTER_HEADING"
	id, err = dark.ParseDirectiveId(input)
	if err != nil {
		t.Errorf("Expected '%s' to parse properly\n", input)
	}
	if id != dark.AFTERHEADING {
		t.Errorf("Expected '%s' to match AFTERHEADING const\n", input)
	}

	input = "DoesntMatch"
	_, err = dark.ParseDirectiveId(input)
	if err == nil {
		t.Errorf("Expected '%s' to error\n", input)
	}
}

func TestParseDirective(t *testing.T) {
	// test erroneous input
	input := "Something"
	_, err := dark.ParseDirective(input)
	if err == nil {
		t.Errorf("Expected error for input '%s'\n", input)
	}

	// check basic input
	input = "ERROR_MSG:<h1>STATUS_MSG</h1>"
	directive, err := dark.ParseDirective(input)
	if err != nil {
		t.Errorf("Expected '%s' to parse successfully\n", input)
	}
	if directive.HttpCode != dark.NoHttpCode {
		t.Errorf("Expected NoHttpCode for '%s'\n", input)
	}
	if directive.DirectiveId != dark.MESSAGE {
		t.Errorf("Expected MESSAGE const for '%s'\n", input)
	}
	if directive.Replacement != "<h1>STATUS_MSG</h1>" {
		t.Errorf("Replacement text doesn't match input argument for '%s'\n", input)
	}

	// check HTTP code specific input
	input = `502:ERROR_HEAD:<meta http-equiv="refresh" content="2">`
	directive, err = dark.ParseDirective(input)
	if err != nil {
		t.Errorf("Expected '%s' to parse successfully\n", input)
	}
	if directive.HttpCode != 502 {
		t.Errorf("Expected 502 for HTTP code from '%s'\n", input)
	}
	if directive.DirectiveId != dark.HEADHTML {
		t.Errorf("Expected HEADHTML const for '%s'\n", input)
	}
	if directive.Replacement != `<meta http-equiv="refresh" content="2">` {
		t.Errorf("Replacement text doesn't match input argument for '%s'\n", input)
	}

}
