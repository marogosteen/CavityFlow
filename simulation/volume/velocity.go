package volume

import "fmt"

type VeloCV struct {
	CVMap     map[Coodinate]float64
	MaxWidth  int
	MaxHeight int
	MinWidth  int
	MinHeight int
}

func NewVeloCV(width int, height int, initVelo float64) (cv VeloCV) {
	cv.CVMap = NewCVMap(width, height, initVelo)
	cv.MaxWidth = width
	cv.MaxHeight = height
	cv.MinWidth = 1
	cv.MinHeight = 1
	return
}

func (cv *VeloCV) Get(x int, y int) float64 {
	v, b := cv.CVMap[Coodinate{X: x, Y: y}]
	if !b {
		msg := "pressure CVMap value is nil\n"
		msg += fmt.Sprintf("coodinate x: %d, y: %d", x, y)
		panic(msg)
	}
	return v
}

func (cv *VeloCV) Set(x int, y int, v float64) {
	c := Coodinate{X: x, Y: y}
	cv.CVMap[c] = v
}

func (cv *VeloCV) Clone() VeloCV {
	cloneMap := make(CVMap)
	for key, v := range(cv.CVMap){
		cloneMap[key] = v
	}

	clone := VeloCV{
		CVMap: cloneMap,
		MaxHeight: cv.MaxHeight,
		MaxWidth:  cv.MaxWidth,
		MinHeight: cv.MinHeight,
		MinWidth:  cv.MinWidth,
	}
	return clone
}