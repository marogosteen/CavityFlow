package controller

import (
	"github.com/marogosteen/cavityflow/core/volume"
)

type BoundaryCondition struct {
	mainFlow float64
}

func NewBoudaryCondition(mf float64) BoundaryCondition {
	return BoundaryCondition{
		mainFlow: mf,
	}
}

func (bf *BoundaryCondition) HorVelo(vol *volume.Volume) {
	// 最左列の水平流速は0
	leftmost := 0
	for y := 0; y < vol.Height; y++ {
		vol.Grid[y][leftmost] = 0.
		vol.Grid[y][leftmost+1] = 0.
	}

	// 最右列の水平流速は0
	rightmost := vol.Width - 1
	for y := 0; y < vol.Height; y++ {
		vol.Grid[y][rightmost] = 0.
		vol.Grid[y][rightmost-1] = 0.
	}

	// 最下段の水平流速は上段と同じ
	bottom := 0
	for x := 0; x < vol.Width; x++ {
		vol.Grid[bottom][x] = -vol.Grid[bottom+1][x]
	}

	// 最上段の水平流速は主流速
	top := vol.Height - 1
	for x := 0; x < vol.Width; x++ {
		vol.Grid[top][x] = bf.mainFlow
	}
}

func (bf *BoundaryCondition) VerVelo(vol *volume.Volume) {
	// 最左列の垂直流速は右列と同じ
	leftmost := 0
	for y := 0; y < vol.Height; y++ {
		vol.Grid[y][leftmost] = -vol.Grid[y][leftmost+1]
	}

	// 最右列の垂直流速は左列と同じ
	rightmost := vol.Width - 1
	for y := 0; y < vol.Height; y++ {
		vol.Grid[y][rightmost] = -vol.Grid[y][rightmost-1]
	}

	// 最下2段の垂直流速は0
	bottom := 0
	for x := 0; x < vol.Width; x++ {
		vol.Grid[bottom][x] = 0.
		vol.Grid[bottom+1][x] = 0.
	}

	// 最上2段の垂直流速は0
	top := vol.Height - 1
	for x := 0; x < vol.Width; x++ {
		vol.Grid[top][x] = 0.
		vol.Grid[top-1][x] = 0.
	}
}

func (bf *BoundaryCondition) Press(vol *volume.Volume) {
	// 最左列は右列と同じ
	leftmost := 0
	for y := 0; y < vol.Height; y++ {
		vol.Grid[y][leftmost] = vol.Grid[y][leftmost+1]
	}

	// 最右列は左列と同じ
	rightmost := vol.Width - 1
	for y := 0; y < vol.Height; y++ {
		vol.Grid[y][rightmost] = vol.Grid[y][rightmost-1]
	}

	// 最下段の圧力は上段と同じ
	bottom := 0
	for x := 0; x < vol.Width; x++ {
		vol.Grid[bottom][x] = vol.Grid[bottom+1][x]
	}

	// 最上段の圧力は下段と同じ
	top := vol.Height - 1
	for x := 0; x < vol.Width; x++ {
		vol.Grid[top][x] = vol.Grid[top-1][x]
	}
}
