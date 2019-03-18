package lib

import (
	"unicode"
)

func ReadItem(line *[]byte) (depth int, leading []byte, content []byte) {
    var cursor int
	cursor = 0
    leading = nil

	for cursor < len(*line) && unicode.IsSpace(rune((*line)[cursor])) {
		cursor += 1
	}
    depth = cursor

    cursor += 1
    for cursor < len(*line) && unicode.IsSpace(rune((*line)[cursor]))  {
        cursor += 1
    }
    cursor = min(cursor, len(*line))
    leading = (*line)[depth:cursor]
    content = (*line)[cursor:]

	return
}
