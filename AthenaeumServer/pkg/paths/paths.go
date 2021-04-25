package paths

import "strings"

var charactersToBeRemoved = []string{
	"../",
	"<!--",
	"-->",
	"<",
	">",
	"'",
	"\"",
	":",
	";", "?", "%20", "%22",
	"%3c",   // <
	"%253c", // <
	"%3e",   // >
	"",      // > -- fill in with % 0 e - without spaces in between
	"%28",   // (
	"%29",   // )
	"%2528", // (
	"%26",   // &
	"%24",   // $
	"%3f",   // ?
	"%3b",   // ;
	"%3d",   // =
	"/",
}

var charactersToBeReplacedByPlus = []string{
	"&",
	"$",
	"#",
	"=",
}

var charactersToBeReplacedByLeftParenthisis = []string{
	"{",
	"[",
}

var charactersToBeReplacedByRightParenthisis = []string{
	"}",
	"]",
}

// Cleaning a string for use in a filepath. Note that this cannot clean a full path, as the "/" character is cleaned.
func CleanNameOnly(name string) string {

	if name == "" {
		return name
	}

	trimmed := strings.TrimSpace(name)

	// replace bad characters from filename
	trimmed = replaceBadCharacters(trimmed, charactersToBeRemoved, "")
	trimmed = replaceBadCharacters(trimmed, charactersToBeReplacedByPlus, "+")
	trimmed = replaceBadCharacters(trimmed, charactersToBeReplacedByLeftParenthisis, "(")
	trimmed = replaceBadCharacters(trimmed, charactersToBeReplacedByRightParenthisis, ")")

	stripped := strings.Replace(trimmed, "\\", "", -1)

	return stripped
}

func replaceBadCharacters(input string, dictionary []string, replacementCharacter string) string {

	temp := input
	for _, badChar := range dictionary {
		temp = strings.Replace(temp, badChar, replacementCharacter, -1)
	}
	return temp
}
