package volume

type PressCV struct {
	CVMap   map[Coodinate]float64
	ColSize int
	RowSize int
}

func NewPressCV(csize int, rsize int, initPress float64) (cv PressCV) {
	cv.CVMap = NewCVMap(csize, rsize, initPress)
	cv.ColSize = csize
	cv.RowSize = rsize
	return
}

func (cv PressCV) Get(col int, row int) float64 {
	return cv.CVMap[Coodinate{Col: col, Row: row}]
}

func (cv PressCV) BoundaryCondition() {
	// 最上段の圧力は下段と同じ
	top := 0
	for col := 1; col < cv.ColSize-1; col++ {
		cv.CVMap[Coodinate{Col: col, Row: top}] = cv.CVMap[Coodinate{Col: col, Row: top + 1}]
	}

	// 最下段の圧力は上段と同じ
	bottom := cv.RowSize
	for col := 1; col < cv.ColSize-1; col++ {
		cv.CVMap[Coodinate{Col: col, Row: bottom}] = cv.CVMap[Coodinate{Col: col, Row: bottom - 1}]
	}

	// 最左列の圧力は右列と同じ
	leftmost := 0
	for row := 1; row < cv.RowSize-1; row++ {
		cv.CVMap[Coodinate{Col: leftmost, Row: row}] = cv.CVMap[Coodinate{Col: leftmost + 1, Row: row}]
	}

	// 最右列の圧力は左列と同じ
	rightmost := cv.ColSize
	for row := 1; row < cv.RowSize-1; row++ {
		cv.CVMap[Coodinate{Col: rightmost, Row: row}] = cv.CVMap[Coodinate{Col: rightmost - 1, Row: row}]
	}
}
