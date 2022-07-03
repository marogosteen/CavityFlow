package volume

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
		panic("velocity cvMap value is nil")
	}
	return v
}

func (cv *VeloCV) Set(x int, y int, v float64) {
	c := Coodinate{X: x, Y: y}
	cv.CVMap[c] = v
}

// Vの定義点における周囲４点のU
// func (cv VeloCV) SurroundingVelo(x int, y int) float64 {
// 	var velo float64 = 0
// 	velo += cv.cvMap[Coodinate{X: x - 1, Y: y}]
// 	velo += cv.cvMap[Coodinate{X: x - 1, Y: y + 1}]
// 	velo += cv.cvMap[Coodinate{X: x, Y: y}]
// 	velo += cv.cvMap[Coodinate{X: x, Y: y + 1}]
// 	velo *= 0.25
// 	return velo
// }

// Uの定義点における周囲４点のV
// func (cv VVeloCV) SurroundingVelo(x int, y int) float64 {
// 	var velo float64 = 0
// 	velo += cv.CVMap[Coodinate{X: x, Y: y}]
// 	velo += cv.CVMap[Coodinate{X: x, Y: y - 1}]
// 	velo += cv.CVMap[Coodinate{X: x + 1, Y: y}]
// 	velo += cv.CVMap[Coodinate{X: x + 1, Y: y - 1}]
// 	velo *= 0.25
// 	return velo
// }
