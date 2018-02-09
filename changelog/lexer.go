package changelog

import (
	"fmt"
	"regexp"
	"strings"
)

// Scanner provides a convenient interface for data such as a changelog file. It is a subset of bufio.Scanner.
type Scanner interface {
	Err() error
	Scan() bool
	Text() string
}

type token interface {
}

type header1Title struct {
	Content string
}

type textLine struct {
	Content string
}

type emptyLine struct {
}

type releaseTitle struct {
	Content string
	Date    string
}

type sectionTitle struct {
	Content string
}

type changeEntry struct {
	Content string
}

type releaseCompareLink struct {
	Title      string
	URL        string
	FromTarget string
	ToTarget   string
}

// Lex lexes a changelog into logical tokens that makes parsing easier.
//
// - H1_TITLE
// - TEXT_LINE
// - EMPTY_LINE
// - RELEASE_TITLE
// - SECTION_TITLE
// - CHANGE_ENTRY
// - RELEASE_COMPARE_LINK
func Lex(scanner Scanner) ([]token, error) {
	tokens := []token{}

	for scanner.Scan() {
		line := scanner.Text()

		var currentToken token
		switch {
		case isHeader1Title(line):
			currentToken = lexHeader1Title(line)
		case isEmptyLine(line):
			currentToken = lexEmptyLine(line)
		case isUnreleasedTitle(line):
			currentToken = lexUnreleasedTitle(line)
		case isReleaseTitle(line):
			currentToken = lexReleaseTitle(line)
		case isSectionTitle(line):
			currentToken = lexSectionTitle(line)
		case isChangeEntry(line):
			currentToken = lexChangeEntry(line)
		case isReleaseCompareLink(line):
			currentToken = lexReleaseCompareLink(line)
		default:
			currentToken = lexTextLine(line)
		}

		if currentToken != nil {
			tokens = append(tokens, currentToken)
		}
	}

	if err := scanner.Err(); err != nil {
		return tokens, err
	}

	return tokens, nil
}

func isHeader1Title(line string) bool {
	return strings.HasPrefix(line, "# ")
}

func lexHeader1Title(line string) header1Title {
	return header1Title{
		Content: line[2:],
	}
}

func isEmptyLine(line string) bool {
	return line == ""
}

func lexEmptyLine(line string) emptyLine {
	return emptyLine{}
}

func isUnreleasedTitle(line string) bool {
	return strings.HasPrefix(line, "## [") && strings.HasSuffix(line, "]")
}

func lexUnreleasedTitle(line string) releaseTitle {
	return releaseTitle{
		Content: line[4 : len(line)-1],
	}
}

func isReleaseTitle(line string) bool {
	return strings.HasPrefix(line, "## [")
}

func lexReleaseTitle(line string) releaseTitle {
	title := ""
	date := ""

	fmt.Sscanf(line, "## %s - %s", &title, &date)

	return releaseTitle{
		Content: title[1 : len(title)-1],
		Date:    date,
	}
}

func isSectionTitle(line string) bool {
	return strings.HasPrefix(line, "### ")
}

func lexSectionTitle(line string) sectionTitle {
	return sectionTitle{
		Content: line[4:],
	}
}

func isChangeEntry(line string) bool {
	return strings.HasPrefix(line, "- ")
}

func lexChangeEntry(line string) changeEntry {
	return changeEntry{
		Content: line[2:],
	}
}

func isReleaseCompareLink(line string) bool {
	return strings.HasPrefix(line, "[")
}

func lexReleaseCompareLink(line string) releaseCompareLink {
	match := compareLinkRegex.FindStringSubmatch(line)

	return releaseCompareLink{
		Title:      match[1],
		URL:        match[2],
		FromTarget: match[3],
		ToTarget:   match[4],
	}
}

func lexTextLine(line string) textLine {
	return textLine{
		Content: line,
	}
}

var compareLinkRegex = regexp.MustCompile(`\[(.*)\]: (https?:\/\/.*\/)(.*)\.\.\.(.*)`)
