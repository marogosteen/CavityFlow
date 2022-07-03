package controller

import (
	"math"

	"github.com/marogosteen/cavityflow/volume"
)

type SimulationController struct {
	Dvs float64
	Dh  float64
	Rho float64
	Eps float64

	HVeloCV volume.HVeloCV
	VVeloCV volume.VVeloCV
	PressCV volume.PressCV
}

func (s *SimulationController) NextVelocity() {
	// TODO index違うかも，最後に確認．magic number
	newHorVeloCVMap := volume.NewCVMap(s.HVeloCV.ColSize, s.HVeloCV.RowSize, 0)
	newVerVeloCVMap := volume.NewCVMap(s.VVeloCV.ColSize, s.VVeloCV.RowSize, 0)
	// TODO magic nuimber
	for row := 2; row < 251; row++ {
		for col := 1; col < 250; col++ {
			newHorVeloCVMap[volume.Coodinate{Col: col, Row: row}] = s.horNS(row, col)
			newVerVeloCVMap[volume.Coodinate{Col: col, Row: row}] = s.verNS(row, col)
		}
	}

	s.HVeloCV.CVMap = newHorVeloCVMap
	s.VVeloCV.CVMap = newVerVeloCVMap
	s.HVeloCV.BoundaryCondition()
	s.VVeloCV.BoundaryCondition()
}

func (s *SimulationController) BoundaryCondition() {
	s.HVeloCV.BoundaryCondition()
	s.VVeloCV.BoundaryCondition()
	s.PressCV.BoundaryCondition()
}

// ver velo Navier–Stokes equationsによるflowの計算
func (s *SimulationController) horNS(row int, col int) float64 {
	u := s.HVeloCV.Get(col, row)
	lu := s.HVeloCV.Get(col-1, row)
	ru := s.HVeloCV.Get(col+1, row)
	ou := s.HVeloCV.Get(col, row-1)
	uu := s.HVeloCV.Get(col, row+1)
	v4 := s.VVeloCV.SurroundingVelo(col, row)
	p := s.PressCV.Get(col, row)
	lp := s.PressCV.Get(col-1, row)

	var uat float64 = 0
	uat += (u + math.Abs(u)) / 2 * (u - lu) / s.Dh
	uat += (u - math.Abs(u)) / 2 * (ru - u) / s.Dh
	uat += (v4 + math.Abs(v4)) / 2 * (u - uu) / s.Dh
	uat += (v4 - math.Abs(v4)) / 2 * (ou - u) / s.Dh

	pt := (p - lp) / (s.Dh * s.Rho)
	dt := (ru + lu + uu + ou - 4*u/math.Pow(s.Dh, 2)*s.Dvs)
	newUVelo := u - dt*(uat+pt-dt)
	return newUVelo
}

func (s *SimulationController) verNS(row int, col int) float64 {
	v := s.VVeloCV.Get(col, row)
	lv := s.VVeloCV.Get(col-1, row)
	rv := s.VVeloCV.Get(col+1, row)
	uv := s.VVeloCV.Get(col, row+1)
	ov := s.VVeloCV.Get(col, row-1)
	u4 := s.HVeloCV.SurroundingVelo(col, row)
	p := s.PressCV.Get(col, row)
	up := s.PressCV.Get(col, row+1)

	var vat float64 = 0
	vat += (u4 + math.Abs(u4)) / 2 * (v - uv) / s.Dh
	vat += (u4 - math.Abs(u4)) / 2 * (ov - v) / s.Dh
	vat += (v + math.Abs(v)) / 2 * (v - lv) / s.Dh
	vat += (v - math.Abs(v)) / 2 * (rv - v) / s.Dh

	pt := (p - up) / (s.Dh * s.Rho)
	dt := (rv + lv + uv + ov - 4*v/math.Pow(s.Dh, 2)*s.Dvs)
	newVVelo := v - dt*(vat+pt-dt)
	return newVVelo
}

func (s *SimulationController) NextPress(phi volume.CVMap) {
	newPressCVMap := volume.NewCVMap(s.HVeloCV.ColSize, s.HVeloCV.RowSize, 0)
	for loss := 0.; loss > s.Eps; {
		loss = 0.

		for row := 1; row <= 251; row++ {
			for col := 1; col <= 250; col++ {
				p := s.PressCV.Get(col, row)
				np := s.Poisson(col, row, phi[volume.Coodinate{Col: col, Row: row}])
				dp := np - p
				loss += dp * dp
				newPressCVMap[volume.Coodinate{Col: col, Row: row}] = p
			}
		}
	}
}

func (s *SimulationController) Poisson(col int, row int, phi float64) float64 {
	var newP float64 = 0
	newP += s.PressCV.Get(col+1, row)
	newP += s.PressCV.Get(col-1, row)
	newP += s.PressCV.Get(col, row+1)
	newP += s.PressCV.Get(col, row-1)
	newP /= 4
	newP -= phi / 4 * s.Dh * s.Dh
	return newP
}

func (s *SimulationController) Phi(colSize int, rowSize int) volume.CVMap {
	phi := volume.NewCVMap(colSize, rowSize, 0.0)
	for row := 1; row <= rowSize; row++ {
		for col := 1; col <= colSize; col++ {
			u := s.HVeloCV.Get(col, row)
			ou := s.HVeloCV.Get(col, row-1)
			uu := s.HVeloCV.Get(col, row+1)
			ru := s.HVeloCV.Get(col+1, row)
			r2u := s.HVeloCV.Get(col+2, row)
			lu := s.HVeloCV.Get(col-1, row)
			oru := s.HVeloCV.Get(col+1, row-1)
			uru := s.HVeloCV.Get(col+1, row+1)
			u4 := s.HVeloCV.SurroundingVelo(col, row)
			ou4 := s.HVeloCV.SurroundingVelo(col, row-1)

			v := s.VVeloCV.Get(col, row)
			ov := s.VVeloCV.Get(col, row-1)
			o2v := s.VVeloCV.Get(col, row-2)
			uv := s.VVeloCV.Get(col, row+1)
			rv := s.VVeloCV.Get(col+1, row)
			lv := s.VVeloCV.Get(col-1, row)
			orv := s.VVeloCV.Get(col+1, row-1)
			olv := s.VVeloCV.Get(col-1, row-1)
			v4 := s.VVeloCV.SurroundingVelo(col, row)
			rv4 := s.VVeloCV.SurroundingVelo(col+1, row)

			var (
				pp1 float64 = 0
				pp2 float64 = 0
				pp3 float64 = 0
			)

			pp1 += ((ru-u)/s.Dh - (ov-v)/s.Dh) / s.Dh

			pp2 += ru*(r2u-u)/(2*s.Dh) + rv4*(oru-uru)/(2*s.Dh)
			pp2 -= u*(ru-lu)/(2*s.Dh) + v4*(ou-uu)/(2*s.Dh)
			pp2 *= -1 / s.Dh

			pp3 += ou4*(orv-olv)/(2*s.Dh) + ov*(o2v-v)/(2*s.Dh)
			pp3 -= u4*(rv-lv)/(2*s.Dh) + v*(ov-uv)/(2*s.Dh)
			pp3 *= -1 / s.Dh
			phi[volume.Coodinate{Col: col, Row: row}] = s.Rho * (pp1 + pp2 + pp3)
		}
	}
	return phi
}
