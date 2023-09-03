package pieces

type Color int8

const (
	COLOR_red  Color = 0
	COLOR_blue Color = 1
)

func (c Color) String() string {
	switch c {
	case COLOR_red:
		return "Red"
	case COLOR_blue:
		return "Blue"
	}
	return "Unknown"
}
