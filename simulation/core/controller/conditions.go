package controller

import "github.com/marogosteen/cavityflow/core/volume"

type BoundaryCondition struct {
	mainFlow float64
}

func NewBoudaryCondition(mf float64) BoundaryCondition {
	return BoundaryCondition{
		mainFlow: mf,
	}
}

func (bf *BoundaryCondition) HorVelo(vol *volume.Volume) {
	// 最左2列の流速は0
	leftmost := 0
	for y := 0; y < vol.MaxHeight; y++ {
		vol.Set(leftmost, y, 0.)
		vol.Set(leftmost+1, y, 0.)
	}

	// 最右2列の流速は0
	rightmost := vol.MaxWidth - 1
	for y := 0; y < vol.MaxHeight; y++ {
		vol.Set(rightmost, y, 0.)
		vol.Set(rightmost-1, y, 0.)
	}

	// 最下2段の流速は上段と同じ
	bottom := 0
	for x := 0; x < vol.MaxWidth; x++ {
		v := vol.Get(x, bottom+1)
		vol.Set(x, bottom, v)
	}

	// 最上段の流速は0，上から二段目は主流
	top := vol.MaxHeight - 1
	for x := 0; x < vol.MaxWidth; x++ {
		vol.Set(x, top, bf.mainFlow)
	}
}

func (bf *BoundaryCondition) VerVelo(vol *volume.Volume) {
	// 最左2列の流速は右列と同じ
	leftmost := 0
	for y := 0; y < vol.MaxHeight; y++ {
		v := vol.Get(leftmost+1, y)
		vol.Set(leftmost, y, v)
	}

	// 最右2列の流速は左列と同じ
	rightmost := vol.MaxWidth - 1
	for y := 0; y < vol.MaxHeight; y++ {
		v := vol.Get(rightmost-1, y)
		vol.Set(rightmost, y, v)
	}

	// 最下2段の流速は0
	bottom := 0
	for x := 0; x < vol.MaxWidth; x++ {
		vol.Set(x, bottom, 0.)
		vol.Set(x, bottom+1, 0.)
	}

	// 最上2段の流速は0
	top := vol.MaxHeight - 1
	for x := 0; x < vol.MaxWidth; x++ {
		vol.Set(x, top, 0.)
		vol.Set(x, top-1, 0.)
	}
}

func (bf *BoundaryCondition) Press(vol *volume.Volume) {
	// 最左2列の圧力は右列と同じ
	leftmost := 0
	for y := 0; y < vol.MaxHeight; y++ {
		v := vol.Get(leftmost+1, y)
		vol.Set(leftmost, y, v)
	}

	// 最右2列の圧力は左列と同じ
	rightmost := vol.MaxWidth - 1
	for y := 0; y < vol.MaxHeight; y++ {
		v := vol.Get(rightmost-1, y)
		vol.Set(rightmost, y, v)
	}

	// 最下2段の圧力は上段と同じ
	bottom := 0
	for x := 0; x < vol.MaxWidth; x++ {
		v := vol.Get(x, bottom+1)
		vol.Set(x, bottom, v)
	}

	// 最上2段の圧力は下段と同じ
	top := vol.MaxHeight - 1
	for x := 0; x < vol.MaxWidth; x++ {
		v := vol.Get(x, top-1)
		vol.Set(x, top, v)
	}
}
