package lib

type Symbol rune
type State map[Symbol]int

const (
    Start Symbol = '-'
    Completed Symbol = '/'
    WontComplete Symbol = 'x'
    Rename Symbol = 'r'
    Continue Symbol = '\\'
)

var states = State {
  Start: 1,
  Completed: 2,
  WontComplete: 3,
  Rename: 4,
  Continue: 5,
}

func unknownState(state []byte) (bool) {
	return (len(state) > 0 && states[Symbol(state[0])] == 0)
}
