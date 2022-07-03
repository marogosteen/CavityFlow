package volume

type PressCV struct {
	CVMap     map[Coodinate]float64
	MaxWidth  int
	MaxHeight int
	MinWidth  int
	MinHeight int
}

func NewPressCV(width int, height int, initPress float64) (cv PressCV) {
	cv.CVMap = NewCVMap(width, height, initPress)
	cv.MaxWidth = width
	cv.MaxHeight = height
	cv.MinWidth = 1
	cv.MinHeight = 1
	return
}

func (cv *PressCV) Get(x int, y int) float64 {
	v, b := cv.CVMap[Coodinate{X: x, Y: y}]
	if !b {
		panic("pressure CVMap value is nil")
	}
	return v
}

func (cv *PressCV) Set(x int, y int, v float64) {
	c := Coodinate{X: x, Y: y}
	cv.CVMap[c] = v
}
