package controller

import (
	"math"

	"github.com/marogosteen/cavityflow/volume"
)

type SimulationController struct {
	Dt       float64
	Dvs      float64
	Dh       float64
	Rho      float64
	Eps      float64
	MainFlow float64

	HorVeloCV *volume.VeloCV
	VerVeloCV *volume.VeloCV
	PressCV   *volume.PressCV
}

func (s *SimulationController) BoundaryCondition() {
	s.horVeloBoundaryCondition()
	s.verVeloBoundaryCondition()
	s.pressBoundaryCondition()
}

// Navier–Stokes equationsによるflowの計算
func (s *SimulationController) NextVelocity() {
	nextHorVeloCV := s.HorVeloCV.Clone()
	nextVerVeloCV := s.VerVeloCV.Clone()
	// TODO magic nuimber
	for y := 3; y <= 253; y++ {
		for x := 3; x <= 252; x++ {
			nextVerVeloCV.Set(x, y, s.verNS(x, y))
			// 主流部分の水平流速は計算しない
			if y <= 252 {
				nextHorVeloCV.Set(x, y, s.horNS(x, y))
			}
		}
	}

	s.HorVeloCV = &nextHorVeloCV
	s.VerVeloCV = &nextVerVeloCV
	// new cvmapはcavity内の計算のみで，境界条件を適応させていない．
	s.BoundaryCondition()
}

func (s *SimulationController) NextPress(phi volume.CVMap) int {
	for count := 1; ; count++ {
		nextPressCV := s.PressCV.Clone()
		loss := 0.
		// TODO magic numebr
		for y := 3; y <= 252; y++ {
			for x := 3; x <= 252; x++ {
				np := s.Poisson(x, y, phi[volume.Coodinate{X: x, Y: y}])
				nextPressCV.Set(x, y, np)

				p := s.PressCV.Get(x, y)
				dp := np - p
				loss += math.Pow(dp, 2)
			}
		}
		s.PressCV = &nextPressCV
		if loss < s.Eps {
			return count
		}
	}
}

// Vの定義点における周囲４点のU
func (s *SimulationController) SurroundingHorVelo(x int, y int) float64 {
	var velo float64 = 0.
	velo += s.HorVeloCV.Get(x+1, y)
	velo += s.HorVeloCV.Get(x+1, y-1)
	velo += s.HorVeloCV.Get(x, y)
	velo += s.HorVeloCV.Get(x, y-1)
	velo *= 0.25
	return velo
}

// Uの定義点における周囲４点のV
func (s *SimulationController) SurroundingVerVelo(x int, y int) float64 {
	var velo float64 = 0.
	velo += s.VerVeloCV.Get(x-1, y)
	velo += s.VerVeloCV.Get(x-1, y+1)
	velo += s.VerVeloCV.Get(x, y)
	velo += s.VerVeloCV.Get(x, y+1)
	velo *= 0.25
	return velo
}
