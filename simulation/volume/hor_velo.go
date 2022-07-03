package volume

type HVeloCV struct {
	CVMap    map[Coodinate]float64
	ColSize  int
	RowSize  int
	MainFlow float64
}

func NewHVeloCV(csize int, rsize int, initVelo float64, mainFlow float64) (cv HVeloCV) {
	cv.CVMap = NewCVMap(csize, rsize, initVelo)
	cv.ColSize = csize
	cv.RowSize = rsize
	cv.MainFlow = mainFlow
	return
}

// func (cv HVeloCV) Get(col int, row int) float64 {
// 	v, b := cv.CVMap[Coodinate{Col: col, Row: row}]
// 	if !b {
// 		log.Fatalln("ver velocity CVMap value is nil")
// 	}
// 	return v
// }

func (cv HVeloCV) Get(col int, row int) float64 {
	v := cv.CVMap[Coodinate{Col: col, Row: row}]
	return v
}

func (cv HVeloCV) BoundaryCondition() {
	// 最上段の流速はMainVelo
	top := 0
	for col := 1; col < cv.ColSize-1; col++ {
		cv.CVMap[Coodinate{Col: col, Row: top}] = cv.MainFlow
	}

	// 最下段の流速は上段と同じ
	bottom := cv.RowSize
	for col := 1; col < cv.ColSize-1; col++ {
		cv.CVMap[Coodinate{Col: col, Row: bottom}] = cv.CVMap[Coodinate{Col: col, Row: bottom - 1}]
	}

	// 最左列の流速は0
	leftmost := 0
	for row := 1; row < cv.RowSize-1; row++ {
		cv.CVMap[Coodinate{Col: leftmost, Row: row}] = 0
	}

	// 最右列の流速は0
	rightmost := cv.ColSize
	for row := 1; row < cv.RowSize-1; row++ {
		cv.CVMap[Coodinate{Col: rightmost, Row: row}] = 0
	}
}

// Vの定義点における周囲４点のU
func (cv HVeloCV) SurroundingVelo(col int, row int) float64 {
	var velo float64 = 0
	velo += cv.CVMap[Coodinate{Col: col - 1, Row: row}]
	velo += cv.CVMap[Coodinate{Col: col - 1, Row: row + 1}]
	velo += cv.CVMap[Coodinate{Col: col, Row: row}]
	velo += cv.CVMap[Coodinate{Col: col, Row: row + 1}]
	velo *= 0.25
	return velo
}
