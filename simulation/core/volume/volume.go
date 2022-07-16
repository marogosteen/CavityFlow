package volume

type Volume struct {
	Grid   [][]float64
	Width  int
	Height int
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
	cv.Width = w
	cv.Height = h
	return
}

func (vol *Volume) Clone() Volume {
	var s [][]float64
	for y := 0; y < vol.Height; y++ {
		line := make([]float64, vol.Width)
		copy(line, vol.Grid[y])
		s = append(s, line)
	}
	return Volume{
		Grid:   s,
		Height: vol.Height,
		Width:  vol.Width,
	}
}
