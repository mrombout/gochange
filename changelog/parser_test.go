package changelog

import (
	"errors"
	"testing"
)

func TestTokenStackPeekEmptyTokenListReturnsNil(t *testing.T) {
	tokenStack := tokenStack{
		tokens: []token{},
	}

	result := tokenStack.peek()

	if result != nil {
		t.Errorf("expected .peek() to return nil , but was %t", *result)
	}
}

func TestTokenStackPeekReturnsFirstToken(t *testing.T) {
	expectedToken := textLine{Content: "Lorum ipsum"}
	tokenStack := tokenStack{
		tokens: []token{expectedToken},
	}

	result := tokenStack.peek()

	if _, ok := (*result).(textLine); !ok {
		t.Errorf("expected .peek() to return textLine, but it wasn't")
	}
}

func TestTokenStackPopEmptyTokenListReturnsNil(t *testing.T) {
	tokenStack := tokenStack{
		tokens: []token{},
	}

	result := tokenStack.pop()

	if result != nil {
		t.Errorf("expected .pop() to return nil , but was %t", *result)
	}
}

func TestTokenStackPopRemovesAndReturnsFirstToken(t *testing.T) {
	expectedFirstToken := textLine{Content: "Line1"}
	expectedSecondToken := textLine{Content: "Line2"}
	tokenStack := tokenStack{
		tokens: []token{expectedFirstToken, expectedSecondToken},
	}

	result1 := tokenStack.pop()
	result2 := tokenStack.pop()

	if (*result1).(textLine).Content != expectedFirstToken.Content {
		t.Errorf("expected the first .pop() to return %v, but was %v", expectedFirstToken, *result1)
	}
	if (*result2).(textLine).Content != expectedSecondToken.Content {
		t.Errorf("expected the second .pop() to return %v, but was %v", expectedSecondToken, *result2)
	}
	if len(tokenStack.tokens) != 0 {
		t.Errorf("expected the token stack to be empty, but wasn't")
	}
}

func TestParseHeaderWhenValidNoError(t *testing.T) {
	tokenStack := tokenStack{
		tokens: []token{header1Title{Content: "Changelog"}, emptyLine{}},
	}

	err := parseHeader(&tokenStack)

	if err != nil {
		t.Errorf("expected error to be nil, but was %v", err)
	}
}

func TestParseHeaderWhenInvalidReturnsError(t *testing.T) {
	testCases := []struct {
		name          string
		tokens        []token
		expectedError error
	}{
		{"no header 1 title", []token{textLine{}, emptyLine{}}, errors.New("unexpected token")},
		{"no empty line after header 1", []token{header1Title{}, textLine{}}, errors.New("unexpected token")},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tokenStack := tokenStack{
				tokens: testCase.tokens,
			}

			err := parseHeader(&tokenStack)

			if err.Error() != testCase.expectedError.Error() {
				t.Errorf("expected error to be '%v', but was '%v'", testCase.expectedError.Error(), err.Error())
			}
		})
	}
}

func TestParseDescription(t *testing.T) {
	expectedDescription := "Lorum ipsum."

	changelog := Changelog{}
	tokenStack := tokenStack{
		tokens: []token{textLine{Content: expectedDescription}, emptyLine{}, releaseTitle{}},
	}

	parseDescription(&tokenStack, &changelog)

	if changelog.Description != expectedDescription {
		t.Errorf("expected parser to have populated .Description with '%v', but was '%v'", expectedDescription, changelog.Description)
	}
}

func TestParseReleaseSections(t *testing.T) {
	testCases := []struct {
		name       string
		tokenStack []token
	}{
		{name: "added without entries", tokenStack: []token{sectionTitle{Content: "Added"}, emptyLine{}}},
		{name: "added with entries", tokenStack: []token{sectionTitle{Content: "Added"}, emptyLine{}, changeEntry{Content: "A couple of cool new bugs."}, emptyLine{}}},
		{name: "changed without entries", tokenStack: []token{sectionTitle{Content: "Changed"}, emptyLine{}}},
		{name: "changed with entries", tokenStack: []token{sectionTitle{Content: "Changed"}, emptyLine{}, changeEntry{Content: "A little bit too much."}, emptyLine{}}},
		{name: "deprecated without entries", tokenStack: []token{sectionTitle{Content: "Deprecated"}, emptyLine{}}},
		{name: "deprecated with entries", tokenStack: []token{sectionTitle{Content: "Deprecated"}, emptyLine{}, changeEntry{Content: "Everything."}, emptyLine{}}},
		{name: "removed without entries", tokenStack: []token{sectionTitle{Content: "Removed"}, emptyLine{}}},
		{name: "removed with entries", tokenStack: []token{sectionTitle{Content: "Removed"}, emptyLine{}, changeEntry{Content: "A useful feature."}, emptyLine{}}},
		{name: "fixed without entries", tokenStack: []token{sectionTitle{Content: "Fixed"}, emptyLine{}}},
		{name: "fixed with entries", tokenStack: []token{sectionTitle{Content: "Fixed"}, emptyLine{}, changeEntry{Content: "Something that isn't broken."}, emptyLine{}}},
		{name: "security without entries", tokenStack: []token{sectionTitle{Content: "Security"}, emptyLine{}}},
		{name: "security with entries", tokenStack: []token{sectionTitle{Content: "Security"}, emptyLine{}, changeEntry{Content: "There's a big security hole in the back!"}, emptyLine{}}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			changelog := Changelog{}
			tokenStack := tokenStack{
				tokens: testCase.tokenStack,
			}
			release := Release{}

			err := parseReleaseSections(&tokenStack, &changelog, &release)

			if err != nil {
				t.Errorf("expected error to be nil, but was '%v'", err.Error())
			}
		})
	}
}

func TestParseReleaseSectionsMissingEmptyLineAfterActionTitleReturnsError(t *testing.T) {
	expectedError := errors.New("unexpected token")
	testCases := []struct {
		name       string
		tokenStack []token
	}{
		{name: "added without empty line", tokenStack: []token{sectionTitle{Content: "Added"}}},
		{name: "changed without empty line", tokenStack: []token{sectionTitle{Content: "Changed"}}},
		{name: "deprecated without empty line", tokenStack: []token{sectionTitle{Content: "Deprecated"}}},
		{name: "removed without empty line", tokenStack: []token{sectionTitle{Content: "Removed"}}},
		{name: "fixed without empty line", tokenStack: []token{sectionTitle{Content: "Fixed"}}},
		{name: "security without empty line", tokenStack: []token{sectionTitle{Content: "Security"}}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			changelog := Changelog{}
			tokenStack := tokenStack{
				tokens: testCase.tokenStack,
			}
			release := Release{}

			err := parseReleaseSections(&tokenStack, &changelog, &release)

			if err.Error() != expectedError.Error() {
				t.Errorf("expected error to be '%v', but was '%v'", expectedError.Error(), err.Error())
			}
		})
	}
}

func TestParseReleaseSectionsNotAChangeSectionAfterSectionTitleReturnsError(t *testing.T) {
	expectedError := errors.New("unexpected token")
	testCases := []struct {
		name       string
		tokenStack []token
	}{
		{name: "added", tokenStack: []token{sectionTitle{Content: "Added"}, emptyLine{}, textLine{}}},
		{name: "changed", tokenStack: []token{sectionTitle{Content: "Changed"}, emptyLine{}, textLine{}}},
		{name: "deprecated", tokenStack: []token{sectionTitle{Content: "Deprecated"}, emptyLine{}, textLine{}}},
		{name: "removed", tokenStack: []token{sectionTitle{Content: "Removed"}, emptyLine{}, textLine{}}},
		{name: "fixed", tokenStack: []token{sectionTitle{Content: "Fixed"}, emptyLine{}, textLine{}}},
		{name: "security", tokenStack: []token{sectionTitle{Content: "Security"}, emptyLine{}, textLine{}}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			changelog := Changelog{}
			tokenStack := tokenStack{
				tokens: testCase.tokenStack,
			}
			release := Release{}

			err := parseReleaseSections(&tokenStack, &changelog, &release)

			if err.Error() != expectedError.Error() {
				t.Errorf("expected error to be '%v', but was '%v'", expectedError.Error(), err.Error())
			}
		})
	}
}

func TestConnectAllRelease(t *testing.T) {
	changelog := newChangelog()
	changelog.Releases = append(changelog.Releases, Release{
		Name: "v1.0.0",
	})
	changelog.Releases = append(changelog.Releases, Release{
		Name: "v0.8.0",
	})
	changelog.Releases = append(changelog.Releases, Release{
		Name: "v0.7.0",
	})
	changelog.Releases = append(changelog.Releases, Release{
		Name: "v0.0.1",
	})

	connectAllReleases(&changelog)
}

func TestFindAndSetLatestRelease(t *testing.T) {
	changelog := newChangelog()
	changelog.Releases = append(changelog.Releases, Release{
		Name: "v1.0.0",
	})

	findAndSetLatestRelease(&changelog)
}

func TestAcceptTokenReturnsErrorWhenTokenIsUnexpected(t *testing.T) {
	tokenStack := tokenStack{
		tokens: []token{textLine{Content: "Lorum ipsum."}},
	}

	acceptToken(&tokenStack, emptyLine{})
}

func TestParseUnreleasedWhenDoesNotStartWithReleaseTitleThrowsError(t *testing.T) {
	expectedError := errors.New("unexpected token")
	tokenStack := tokenStack{
		tokens: []token{
			emptyLine{},
		},
	}
	changelog := Changelog{}

	err := parseUnreleased(&tokenStack, &changelog)

	if err.Error() != expectedError.Error() {
		t.Errorf("expected error to be '%v', but was '%v'", expectedError.Error(), err.Error())
	}
}

func TestParseUnreleasedPopulatesUnreleasedName(t *testing.T) {
	expectedUnreleasedName := "Unreleased"
	testCases := []struct {
		name   string
		tokens []token
	}{
		{"with emptyline", []token{releaseTitle{Content: expectedUnreleasedName}, emptyLine{}}},
		{"without emptyline", []token{releaseTitle{Content: expectedUnreleasedName}}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tokenStack := tokenStack{
				tokens: testCase.tokens,
			}
			changelog := Changelog{
				Unreleased: Release{
					Name: "Wrong Release Name",
				},
			}

			err := parseUnreleased(&tokenStack, &changelog)

			if err != nil {
				t.Errorf("expected error to be nil, but was '%v'", err.Error())
			}
			if changelog.Unreleased.Name != expectedUnreleasedName {
				t.Errorf("expected .Unreleased.Name to be '%v', but was '%v'", expectedUnreleasedName, changelog.Unreleased.Name)
			}
		})
	}
}

func TestParseUnreleasedPopulatesUnreleased(t *testing.T) {
	tokenStack := tokenStack{
		tokens: []token{
			releaseTitle{Content: "Unreleased"},
			emptyLine{},
		},
	}
	changelog := Changelog{}

	parseUnreleased(&tokenStack, &changelog)
}

func TestParseReleasesNoEmptyLineAfterTitleReturnsError(t *testing.T) {
	expectedError := errors.New("unexpected token")
	tokenStack := tokenStack{
		tokens: []token{
			releaseTitle{Content: "v1.0.0"},
			sectionTitle{Content: "Added"},
		},
	}
	changelog := Changelog{}

	err := parseReleases(&tokenStack, &changelog)

	if err == nil {
		t.Fatalf("expected error to be '%v', but was nil", expectedError.Error())
	}
	if err.Error() != expectedError.Error() {
		t.Errorf("expected error to be '%v', but was '%v'", expectedError.Error(), err.Error())
	}
}

func TestParseReleases(t *testing.T) {
	tokenStack := tokenStack{
		tokens: []token{
			releaseTitle{Content: "v1.0.0"},
			emptyLine{},
		},
	}
	changelog := Changelog{}

	parseReleases(&tokenStack, &changelog)

	// TODO: Assertions
}

func TestParse(t *testing.T) {
	tokenStack := []token{
		header1Title{},
		emptyLine{},
		textLine{},
		emptyLine{},
		releaseTitle{},
		sectionTitle{Content: "Added"},
		emptyLine{},
		changeEntry{Content: "Added a bug."},
		emptyLine{},
		releaseCompareLink{},
	}

	Parse(tokenStack)

	// TODO: Assertions
}

func TestParseInvalidHeaderReturnsError(t *testing.T) {
	expectedError := errors.New("unexpected token")
	tokenStack := []token{
		textLine{Content: "Changelog"},
		emptyLine{},
		textLine{},
		emptyLine{},
		releaseTitle{},
		sectionTitle{Content: "Added"},
		emptyLine{},
		changeEntry{Content: "Added a bug."},
		emptyLine{},
		releaseCompareLink{},
	}

	_, err := Parse(tokenStack)

	if err == nil {
		t.Fatalf("expected error to be '%v', but was nil", expectedError.Error())
	}
	if err.Error() != expectedError.Error() {
		t.Errorf("expected error to be '%v', but was '%v'", expectedError.Error(), err.Error())
	}
}

func TestParseIncorrectUnreleasedError(t *testing.T) {
	expectedError := errors.New("unexpected token")
	tokenStack := []token{
		header1Title{Content: "Changelog"},
		emptyLine{},
		textLine{},
		emptyLine{},
		releaseTitle{Content: "Unreleased"},
		sectionTitle{Content: "Added"},
		changeEntry{Content: "Added a bug."},
		emptyLine{},
		releaseCompareLink{},
	}

	_, err := Parse(tokenStack)

	if err == nil {
		t.Fatalf("expected error to be '%v', but was nil", expectedError.Error())
	}
	if err.Error() != expectedError.Error() {
		t.Errorf("expected error to be '%v', but was '%v'", expectedError.Error(), err.Error())
	}
}

func TestParseIncorrectReleaseError(t *testing.T) {
	expectedError := errors.New("unexpected token")
	tokenStack := []token{
		header1Title{Content: "Changelog"},
		emptyLine{},
		textLine{},
		emptyLine{},
		releaseTitle{Content: "Unreleased"},
		sectionTitle{Content: "Added"},
		emptyLine{},
		changeEntry{Content: "Added a bug."},
		emptyLine{},
		releaseTitle{Content: "v1.0.0"},
		sectionTitle{Content: "Added"},
		changeEntry{Content: "Added an easter egg."},
		releaseCompareLink{},
	}

	_, err := Parse(tokenStack)

	if err == nil {
		t.Fatalf("expected error to be '%v', but was nil", expectedError.Error())
	}
	if err.Error() != expectedError.Error() {
		t.Errorf("expected error to be '%v', but was '%v'", expectedError.Error(), err.Error())
	}
}
