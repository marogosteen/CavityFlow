package volume

type CVMap map[Coodinate]float64

func NewCVMap(width int, height int, InitValue float64) CVMap {
	cvMap := make(CVMap)
	for y := 1; y <= height; y++ {
		for x := 1; x <= width; x++ {
			cvMap[Coodinate{X: x, Y: y}] = InitValue
		}
	}
	return cvMap
}
