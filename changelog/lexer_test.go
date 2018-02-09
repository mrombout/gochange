package changelog

import (
	"bufio"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestLexTokensIndividually(t *testing.T) {
	testCases := []struct {
		line           string
		expectedTokens []token
	}{
		{"# A Header 1", []token{header1Title{Content: "A Header 1"}}},
		{"\r\n", []token{emptyLine{}}},
		{"## [Unreleased]", []token{releaseTitle{Content: "Unreleased"}}},
		{"## [0.0.1] - 2018-12-06", []token{releaseTitle{Content: "0.0.1", Date: "2018-12-06"}}},
		{"### Added", []token{sectionTitle{Content: "Added"}}},
		{"- A massive bug", []token{changeEntry{Content: "A massive bug"}}},
		{"[1.0.0]: https://github.com/olivierlacan/keep-a-changelog/compare/v0.3.0...v1.0.0", []token{releaseCompareLink{
			Title:      "1.0.0",
			URL:        "https://github.com/olivierlacan/keep-a-changelog/compare/",
			FromTarget: "v0.3.0",
			ToTarget:   "v1.0.0",
		}}},
		{"This is a regular text line!", []token{textLine{Content: "This is a regular text line!"}}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.line, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(testCase.line))
			tokens, err := Lex(scanner)

			if err != nil {
				t.Fatalf("expected result to be nil, but got %t", err)
			}
			if len(tokens) != len(testCase.expectedTokens) {
				t.Fatalf("expected to have lexed exactly %d tokens, but was %d (%t)", len(testCase.expectedTokens), len(tokens), tokens)
			}
			for i := 0; i < len(testCase.expectedTokens); i++ {
				actualToken := tokens[i]
				expectedToken := testCase.expectedTokens[i]

				if !reflect.DeepEqual(expectedToken, actualToken) {
					t.Errorf("expected token %d to equal %t, but was %t", i, expectedToken, actualToken)
				}
			}
		})
	}
}

type ErroringScanner struct {
}

func (ErroringScanner) Err() error {
	return errors.New("this error is expected")
}

func (ErroringScanner) Scan() bool { return false }

func (ErroringScanner) Text() string { return "" }

func TestLexErrorDuringLexing(t *testing.T) {
	scanner := ErroringScanner{}

	_, err := Lex(scanner)

	if scanner.Err().Error() != err.Error() {
		t.Errorf("expected Lex to return error %t , but was %t", scanner.Err(), err)
	}
}

func TestIsHeader1Title(t *testing.T) {
	testCases := []struct {
		line           string
		expectedResult bool
	}{
		{"# A Header 1", true},
		{"## A Header 2", false},
		{"[A Link](http://golang.org/)", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.line, func(t *testing.T) {
			result := isHeader1Title(testCase.line)

			if result != testCase.expectedResult {
				t.Errorf("expected result to be %t, but got %t", testCase.expectedResult, result)
			}
		})
	}
}

func TestLexHeader1Title(t *testing.T) {
	// arrange
	line := "# A Header 1"

	// act
	result := lexHeader1Title(line)

	// assert
	if result.Content != "A Header 1" {
		t.Errorf("expected result to be %s, but got %s", "A Header 1", result.Content)
	}
}

func TestIsEmptyLine(t *testing.T) {
	testCases := []struct {
		line           string
		expectedResult bool
	}{
		{"This line is not empty.", false},
		{"", true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.line, func(t *testing.T) {
			result := isEmptyLine(testCase.line)

			if result != testCase.expectedResult {
				t.Errorf("expected result to be %t, but got %t", testCase.expectedResult, result)
			}
		})
	}
}

func TestLexEmptyLine(t *testing.T) {
	// arrange
	line := ""
	expectedResult := emptyLine{}

	// act
	result := lexEmptyLine(line)

	// assert
	if result != expectedResult {
		t.Errorf("expected result to be %s, but got %s", emptyLine{}, result)
	}
}

func TestIsUnreleasedTitle(t *testing.T) {
	testCases := []struct {
		line           string
		expectedResult bool
	}{
		{"## [Unreleased]", true},
		{"## [1.0.1] - 2018-12-24", false},
		{"## Header 2", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.line, func(t *testing.T) {
			result := isUnreleasedTitle(testCase.line)

			if result != testCase.expectedResult {
				t.Errorf("expected result to be %t, but got %t", testCase.expectedResult, result)
			}
		})
	}
}

func TestLexUnreleasedTitle(t *testing.T) {
	// arrange
	line := "## [Unreleased]"

	// act
	result := lexUnreleasedTitle(line)

	// assert
	if result.Content != "Unreleased" {
		t.Errorf("expected result to be %s, but got %s", "Unreleased", result)
	}
}

func TestIsReleaseTitle(t *testing.T) {
	testCases := []struct {
		line           string
		expectedResult bool
	}{
		// {"## [Unreleased]", false},
		{"## [1.0.1] - 2018-12-24", true},
		{"## Header 2", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.line, func(t *testing.T) {
			result := isReleaseTitle(testCase.line)

			if result != testCase.expectedResult {
				t.Errorf("expected result to be %t, but got %t", testCase.expectedResult, result)
			}
		})
	}
}

func TestLexReleaseTitle(t *testing.T) {
	// arrange
	line := "## [v1.0.0] - 2018-12-24"

	// act
	result := lexReleaseTitle(line)

	// assert
	if result.Content != "v1.0.0" {
		t.Errorf("expected result to be %s, but got %s", "v1.0.0", result)
	}
	if result.Date != "2018-12-24" {
		t.Errorf("expected result to be %s, but got %s", "2018-12-24", result)
	}
}

func TestIsSectionTitle(t *testing.T) {
	testCases := []struct {
		line           string
		expectedResult bool
	}{
		// {"## [Unreleased]", false},
		{"### Added", true},
		{"### Changed", true},
		{"### Deprecated", true},
		{"### Removed", true},
		{"### Fixed", true},
		{"### Security", true},
		{"## [v1.0.0] - 2018-12-28", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.line, func(t *testing.T) {
			result := isSectionTitle(testCase.line)

			if result != testCase.expectedResult {
				t.Errorf("expected result to be %t, but got %t", testCase.expectedResult, result)
			}
		})
	}
}

func TestLexSectionTitle(t *testing.T) {
	testCases := []struct {
		line            string
		expectedContent string
	}{
		{"### Added", "Added"},
		{"### Changed", "Changed"},
		{"### Deprecated", "Deprecated"},
		{"### Removed", "Removed"},
		{"### Fixed", "Fixed"},
		{"### Security", "Security"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.line, func(t *testing.T) {
			result := lexSectionTitle(testCase.line)

			if result.Content != testCase.expectedContent {
				t.Errorf("expected result to be %s, but got %s", testCase.expectedContent, result)
			}
		})
	}
}

func TestIsChangeEntry(t *testing.T) {
	testCases := []struct {
		line           string
		expectedResult bool
	}{
		{"* Not a change entry", false},
		{"- A change entry", true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.line, func(t *testing.T) {
			result := isChangeEntry(testCase.line)

			if result != testCase.expectedResult {
				t.Errorf("expected result to be %t, but got %t", testCase.expectedResult, result)
			}
		})
	}
}

func TestLexChangeEntry(t *testing.T) {
	// arrange
	line := "- Added some stuff"

	// act
	result := lexChangeEntry(line)

	// assert
	if result.Content != "Added some stuff" {
		t.Errorf("expected result to be %s, but got %s", "Added some stuff", result)
	}
}

func TestIsReleaseCompareLink(t *testing.T) {
	testCases := []struct {
		line           string
		expectedResult bool
	}{
		{"[Unreleased]: https://golang.org/", true},
		{"[v1.0.0]: https://golang.org/", true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.line, func(t *testing.T) {
			result := isReleaseCompareLink(testCase.line)

			if result != testCase.expectedResult {
				t.Errorf("expected result to be %t, but got %t", testCase.expectedResult, result)
			}
		})
	}
}

func TestLexReleaseCompareLink(t *testing.T) {
	// arrange
	line := "[v1.0.0]: https://golang.org/v1.0.0...HEAD"

	// act
	result := lexReleaseCompareLink(line)

	// assert
	if result.Title != "v1.0.0" {
		t.Errorf("expected result to be %s, but got %s", "v1.0.0", result)
	}
	if result.URL != "https://golang.org/" {
		t.Errorf("expected result to be %s, but got %s", "https://golang.org/", result)
	}
	if result.FromTarget != "v1.0.0" {
		t.Errorf("expected result to be %s, but got %s", "v1.0.0", result)
	}
	if result.ToTarget != "HEAD" {
		t.Errorf("expected result to be %s, but got %s", "HEAD", result)
	}
}

func TestLexTextLine(t *testing.T) {
	// arrange
	line := "A line of text."

	// act
	result := lexTextLine(line)

	// assert
	if result.Content != line {
		t.Errorf("expected result to be %s, but got %s", line, result)
	}
}
