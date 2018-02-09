package changelog

import (
	"errors"
	"reflect"
)

type tokenStack struct {
	tokens []token
}

func (t *tokenStack) peek() *token {
	if len(t.tokens) <= 0 {
		return nil
	}
	return &t.tokens[0]
}

func (t *tokenStack) pop() *token {
	if len(t.tokens) <= 0 {
		return nil
	}

	token := t.tokens[0]
	t.tokens = t.tokens[1:]

	return &token
}

// Parse parses a list of tokens as returned by `changelog.Lex`.
func Parse(tokens []token) (Changelog, error) {
	changelog := newChangelog()

	stack := tokenStack{
		tokens: tokens,
	}

	if err := parseHeader(&stack); err != nil {
		return changelog, err
	}
	parseDescription(&stack, &changelog)
	if err := parseUnreleased(&stack, &changelog); err != nil {
		return changelog, err
	}
	if err := parseReleases(&stack, &changelog); err != nil {
		return changelog, err
	}

	connectAllReleases(&changelog)
	findAndSetLatestRelease(&changelog)

	return changelog, nil
}

func parseHeader(stack *tokenStack) error {
	if _, err := acceptToken(stack, header1Title{}); err != nil {
		return err
	}
	if _, err := acceptToken(stack, emptyLine{}); err != nil {
		return err
	}

	return nil
}

func parseDescription(stack *tokenStack, changelog *Changelog) {
	for token := stack.peek(); !isToken(stack, releaseTitle{}); token = stack.peek() {
		if val, ok := (*token).(textLine); ok {
			changelog.Description += val.Content + "\n"
		} else if _, ok := (*token).(emptyLine); ok {

		}

		stack.pop()
	}
	changelog.Description = changelog.Description[:len(changelog.Description)-1]
}

func parseUnreleased(stack *tokenStack, changelog *Changelog) error {
	if _, err := acceptToken(stack, releaseTitle{}); err != nil {
		return err
	}
	if isToken(stack, emptyLine{}) {
		acceptToken(stack, emptyLine{})
	}

	currentRelease := Release{
		Name: "Unreleased",
	}

	parseReleaseSections(stack, changelog, &currentRelease)

	changelog.Unreleased = currentRelease

	return nil
}

func parseReleases(stack *tokenStack, changelog *Changelog) error {
	for isToken(stack, releaseTitle{}) {
		currentReleaseTitle, _ := acceptToken(stack, releaseTitle{})
		if _, err := acceptToken(stack, emptyLine{}); err != nil {
			return err
		}

		if currentReleaseTitle, ok := (*currentReleaseTitle).(releaseTitle); ok {
			currentRelease := Release{
				Name: currentReleaseTitle.Content,
				Date: currentReleaseTitle.Date,
			}

			parseReleaseSections(stack, changelog, &currentRelease)

			changelog.Releases = append(changelog.Releases, currentRelease)
		}
	}

	return nil
}

func parseReleaseSections(stack *tokenStack, changelog *Changelog, release *Release) error {
	for isToken(stack, sectionTitle{}) {
		sectionToken, _ := acceptToken(stack, sectionTitle{})
		if _, err := acceptToken(stack, emptyLine{}); err != nil {
			return err
		}

		var list *[]Entry
		if currentSectionToken, ok := (*sectionToken).(sectionTitle); ok {
			switch currentSectionToken.Content {
			case "Added":
				list = &release.Added
			case "Changed":
				list = &release.Changed
			case "Deprecated":
				list = &release.Deprecated
			case "Removed":
				list = &release.Removed
			case "Fixed":
				list = &release.Fixed
			case "Security":
				list = &release.Security
			}

			for !isToken(stack, emptyLine{}) && len(stack.tokens) > 0 {
				changeEntryToken, err := acceptToken(stack, changeEntry{})
				if err != nil {
					return err
				}
				if val, ok := (*changeEntryToken).(changeEntry); ok {
					*list = append(*list, Entry{
						Description: val.Content,
					})
				}
			}

			if len(stack.tokens) > 0 {
				acceptToken(stack, emptyLine{})
			}
		}
	}

	return nil
}

func connectAllReleases(changelog *Changelog) {
	for index := range changelog.Releases {
		if index < len(changelog.Releases)-1 {
			changelog.Releases[index].PreviousRelease = &changelog.Releases[index+1]
		}
	}
}

func findAndSetLatestRelease(changelog *Changelog) {
	if len(changelog.Releases) > 0 {
		changelog.LatestRelease = changelog.Releases[len(changelog.Releases)-1]
	}
}

func acceptToken(stack *tokenStack, tokenType interface{}) (*token, error) {
	if !isToken(stack, tokenType) {
		return nil, errors.New("unexpected token")
	}

	return stack.pop(), nil
}

func isToken(stack *tokenStack, tokenType interface{}) bool {
	if len(stack.tokens) <= 0 {
		return false
	}
	token := stack.peek()

	typeA := reflect.TypeOf(*token)
	typeB := reflect.TypeOf(tokenType)

	return typeA == typeB
}
