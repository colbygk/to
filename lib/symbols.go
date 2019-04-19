package lib

type Symbol rune
type State map[Symbol]int

const (
  Start        Symbol = '-'
  Completed    Symbol = '/'
  WontComplete Symbol = 'x'
  Rename       Symbol = 'r'
  Note         Symbol = 'n'
  High         Symbol = '!'
  Continue     Symbol = '\\'
)

var states = State{
  Start:        1,
  Completed:    2,
  WontComplete: 3,
  Rename:       4,
  Note:         5,
  High:         6,
  Continue:     7,
}

func UnknownState(state []byte) bool {
  return (state != nil && len(state) > 0 && states[Symbol(state[0])] == 0)
}
