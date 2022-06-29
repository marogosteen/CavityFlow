package volume

type VVeloCV struct {
	CVMap   map[Coodinate]float64
	ColSize int
	RowSize int
}

func NewVVeloCV(csize int, rsize int, initVelo float64) (cv VVeloCV) {
	cv.CVMap = NewCVMap(csize, rsize, initVelo)
	cv.ColSize = csize
	cv.RowSize = rsize
	return
}

func (cv VVeloCV) Get(col int, row int) float64 {
	return cv.CVMap[Coodinate{Col: col, Row: row}]
}

func (cv VVeloCV) BoundaryCondition() {
	// 最上段の流速は0
	top := 0
	for col := 1; col < cv.ColSize-1; col++ {
		cv.CVMap[Coodinate{Col: col, Row: top}] = 0
	}

	// 最下段の流速は0
	bottom := cv.RowSize
	for col := 1; col < cv.ColSize-1; col++ {
		cv.CVMap[Coodinate{Col: col, Row: bottom}] = 0
	}

	// 最左列の流速は右列と同じ
	leftmost := 0
	for row := 1; row < cv.RowSize-1; row++ {
		cv.CVMap[Coodinate{Col: leftmost, Row: row}] = cv.CVMap[Coodinate{Col: leftmost + 1, Row: row}]
	}

	// 最右列の流速は左列と同じ
	rightmost := cv.ColSize
	for row := 1; row < cv.RowSize-1; row++ {
		cv.CVMap[Coodinate{Col: rightmost, Row: row}] = cv.CVMap[Coodinate{Col: rightmost - 1, Row: row}]
	}
}

// Uの定義点における周囲４点のV
func (cv VVeloCV) SurroundingVelo(col int, row int) float64 {
	var velo float64 = 0
	velo += cv.CVMap[Coodinate{Col: col, Row: row}]
	velo += cv.CVMap[Coodinate{Col: col, Row: row - 1}]
	velo += cv.CVMap[Coodinate{Col: col + 1, Row: row}]
	velo += cv.CVMap[Coodinate{Col: col + 1, Row: row - 1}]
	velo *= 0.25
	return velo
}
