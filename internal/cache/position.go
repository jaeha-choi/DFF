package cache

type Position int

const (
	None Position = iota - 1
	Top
	Jungle
	Mid
	Adc
	Support
)

var PositionList = [5]Position{Top, Jungle, Mid, Adc, Support}

func (p Position) String() string {
	switch p {
	case Top:
		return "Top"
	case Jungle:
		return "Jungle"
	case Mid:
		return "Mid"
	case Adc:
		return "Adc"
	case Support:
		return "Support"
	default:
		return ""
	}
}
