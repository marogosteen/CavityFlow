package volume

type CVMap map[Coodinate]float64

func NewCVMap(csize int, rsize int, InitValue float64) CVMap {
	cvMap := make(CVMap)
	for row := 0; row < rsize; row++ {
		for col := 0; col < csize; col++ {
			cvMap[Coodinate{Col: col, Row: row}] = InitValue
		}
	}
	return cvMap
}
