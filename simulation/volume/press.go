package volume

import "fmt"

type PressCV struct {
	CVMap     map[Coodinate]float64
	MaxWidth  int
	MaxHeight int
}

func NewPressCV(width int, height int, initPress float64) (cv PressCV) {
	cv.CVMap = NewCVMap(width, height, initPress)
	cv.MaxWidth = width
	cv.MaxHeight = height
	return
}

func (cv *PressCV) Get(x int, y int) float64 {
	v, b := cv.CVMap[Coodinate{X: x, Y: y}]
	if !b {
		msg := "pressure CVMap value is nil\n"
		msg += fmt.Sprintf("coodinate x: %d, y: %d", x, y)
		panic(msg)
	}
	return v
}

func (cv *PressCV) Set(x int, y int, v float64) {
	c := Coodinate{X: x, Y: y}
	cv.CVMap[c] = v
}

func (cv *PressCV) Clone() PressCV {
	cloneMap := make(CVMap)
	for key, v := range(cv.CVMap){
		cloneMap[key] = v
	}

	clone := PressCV{
		CVMap: cloneMap,
		MaxHeight: cv.MaxHeight,
		MaxWidth:  cv.MaxWidth,
	}
	return clone
}
