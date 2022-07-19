package word

import (
	"strings"
	"unicode"
)

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

// 下划线 => 大写驼峰
func UnderscoreToUpperCamelCase(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.Title(strings.ToLower(s))
	return strings.ReplaceAll(s, " ", "")
}

// 下划线 => 小写驼峰
func UnderscoreToLowerCamelCase(s string) string {
	s = UnderscoreToUpperCamelCase(s)
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}

// 驼峰 => 下划线
func CamelCaseToUnderscore(s string) string {
	var output []rune
	for k, v := range s {
		if k == 0 {
			output = append(output, unicode.ToLower(v))
			continue
		}

		if unicode.IsUpper(v) {
			output = append(output, '_')
		}

		output = append(output, unicode.ToLower(v))

	}

	return string(output)
}
