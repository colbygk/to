package lib

import (
  "regexp"
)

var toDoItem = regexp.MustCompile(`(\s+)?([-/xrn!\\])?(\s+)?(.+)$`)

func ReadItem(line *[]byte) (depth int, leading []byte, content []byte) {
  leading = []byte{'\\'}
  content = nil
  depth = 0

  subitems := toDoItem.FindSubmatch(*line)
  if len(subitems) > 1 {
    depth = len(subitems[1])
    if len(subitems) > 2 {
      leading = subitems[2]
      if len(subitems) > 4 {
        content = subitems[4]
      } else {
        content = subitems[3]
      }
    }
  }

  return
}
