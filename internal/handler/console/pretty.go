package console

import (
	"encoding/json"
	"io"
	"strings"
)

type PrettyJSONWriter struct {
	writer io.Writer
}

func NewPrettyJSONWriter(w io.Writer) *PrettyJSONWriter {
	return &PrettyJSONWriter{writer: w}
}

func (w *PrettyJSONWriter) Write(p []byte) (n int, err error) {
	var jsonData any
	if err := json.Unmarshal(p, &jsonData); err != nil {
		return w.writer.Write(p)
	}

	prettyBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return w.writer.Write(p)
	}

	prettyBytes = append(prettyBytes, '\n')

	return w.writer.Write(prettyBytes)
}

// Apply basic JSON syntax highlighting
func applyJSONSyntaxHighlighting(jsonStr string) string {
	result := ""
	inString := false
	inKey := false
	escaped := false

	for i, char := range jsonStr {
		switch char {
		case '"':
			if !escaped {
				if !inString {
					inString = true
					// Determine if this is a key (followed by :) or value
					colonIndex := findNextNonWhitespace(jsonStr, i+1)
					if colonIndex != -1 && jsonStr[colonIndex] == ':' {
						inKey = true
						result += Cyan + Bold + string(char) // Key color
					} else {
						result += Green + string(char) // String value color
					}
				} else {
					if inKey {
						result += string(char) + Reset
						inKey = false
					} else {
						result += string(char) + Reset
					}
					inString = false
				}
			} else {
				result += string(char)
			}
			escaped = false
		case '\\':
			result += string(char)
			escaped = !escaped
		case ':':
			if !inString {
				result += Yellow + string(char) + Reset
			} else {
				result += string(char)
			}
			escaped = false
		case '{', '}', '[', ']':
			if !inString {
				result += Magenta + Bold + string(char) + Reset
			} else {
				result += string(char)
			}
			escaped = false
		case ',':
			if !inString {
				result += Gray + string(char) + Reset
			} else {
				result += string(char)
			}
			escaped = false
		default:
			// Numbers
			if !inString && char >= '0' && char <= '9' {
				result += Cyan + string(char) + Reset
			} else if !inString && (char == 't' || char == 'f' || char == 'n') {
				// Boolean and null values
				if strings.HasPrefix(jsonStr[i:], "true") || strings.HasPrefix(jsonStr[i:], "false") {
					result += Yellow + Bold + string(char) + Reset
				} else if strings.HasPrefix(jsonStr[i:], "null") {
					result += Red + string(char) + Reset
				} else {
					result += string(char)
				}
			} else {
				result += string(char)
			}
			escaped = false
		}
	}

	return result
}

// Helper to find next non-whitespace character
func findNextNonWhitespace(s string, start int) int {
	for i := start; i < len(s); i++ {
		if s[i] != ' ' && s[i] != '\t' && s[i] != '\n' && s[i] != '\r' {
			return i
		}
	}
	return -1
}
