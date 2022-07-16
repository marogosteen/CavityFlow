package volume

type Volume struct {
	Grid      [][]float64
	MaxWidth  int
	MaxHeight int
}

func NewVolume(w int, h int, val float64) (cv Volume) {
	var s [][]float64
	for y := 0; y < h; y++ {
		s = append(s, make([]float64, w))
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			s[y][x] = val
		}
	}
	cv.Grid = s
	cv.MaxWidth = w
	cv.MaxHeight = h
	return
}

func (vol *Volume) Get(x int, y int) float64 {
	return vol.Grid[y][x]
}

func (vol *Volume) Set(x int, y int, v float64) {
	vol.Grid[y][x] = v
}

func (vol *Volume) Clone() Volume {
	var s [][]float64
	for y := 0; y < vol.MaxHeight; y++ {
		s = append(s, make([]float64, vol.MaxWidth))
	}
	copy(s, vol.Grid)
	return Volume{
		Grid:      s,
		MaxHeight: vol.MaxHeight,
		MaxWidth:  vol.MaxWidth,
	}
}
