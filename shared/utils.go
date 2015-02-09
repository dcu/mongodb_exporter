package shared

import(
    "regexp"
    "strings"
)

var(
    snakeRegexp = regexp.MustCompile("\\B[A-Z]+[^_$]")
    parameterizeRegexp = regexp.MustCompile("[^A-Za-z0-9_]+")
)

func SnakeCase(text string) string {
    result := snakeRegexp.ReplaceAllStringFunc(text, func(match string) string {
        return "_"+match
    })

    return ParameterizeString(result)
}

func ParameterizeString(text string) string {
    result := parameterizeRegexp.ReplaceAllString(text, "_")
    return strings.ToLower(result)
}


