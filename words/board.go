package words

type value interface {
	Matches(v value) bool
	String() string
}

type stringValue struct {
	value string
}

func (s stringValue) Matches(v value) bool {
	other, ok := v.(stringValue)
	if !ok {
		return false
	}
	return s.value == other.value
}

func (s stringValue) String() string {
	return s.value
}

// newStringValue returns the string wrapped in a value
func newStringValue(v string) value {
	return stringValue{v}
}

type board interface {
	Set(x, y int, v value)
	Get(x, y int) value
	ToStringArray() [][]string
}

type arrayBoard struct {
	values [][]value
}

func newArrayBoard(x, y int) board {
	a := arrayBoard{}
	a.values = make([][]value, y)
	for i := range a.values {
		a.values[i] = make([]value, x)
	}
	return a
}

func (a arrayBoard) Set(x, y int, v value) {
	a.values[x][y] = v
}

func (a arrayBoard) Get(x, y int) value {
	return a.values[x][y]
}

func (a arrayBoard) ToStringArray() [][]string {
	arr := make([][]string, len(a.values))
	for x := 0; x < len(a.values); x++ {
		arr[x] = make([]string, len(a.values[0]))
		for y := 0; y < len(a.values[0]); y++ {
			arr[x][y] = a.Get(x, y).String()
		}
	}
	return arr
}
