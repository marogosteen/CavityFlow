package controller

import (
	"fmt"

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
	newHorVeloCVMap := volume.NewCVMap(s.HorVeloCV.MaxWidth, s.HorVeloCV.MaxHeight, 0.)
	newVerVeloCVMap := volume.NewCVMap(s.VerVeloCV.MaxWidth, s.VerVeloCV.MaxHeight, 0.)
	// TODO magic nuimber
	for y := 3; y <= 252; y++ {
		for x := 3; x <= 252; x++ {
			newHorVeloCVMap[volume.Coodinate{X: x, Y: y}] = s.horNS(x, y)
			newVerVeloCVMap[volume.Coodinate{X: x, Y: y}] = s.verNS(x, y)
		}
	}

	s.HorVeloCV.CVMap = newHorVeloCVMap
	s.VerVeloCV.CVMap = newVerVeloCVMap
	// new cvmapはcavity内の計算のみで，境界条件を適応させていない．
	s.BoundaryCondition()
}

func (s *SimulationController) NextPress(phi volume.CVMap) {
	for count := 1; ; count++ {
		nextPressCVMap := volume.NewCVMap(s.HorVeloCV.MaxWidth, s.HorVeloCV.MaxHeight, 0)
		loss := 0.
		// TODO magic numebr
		for y := 3; y <= 252; y++ {
			for x := 3; x <= 252; x++ {
				p := s.PressCV.Get(x, y)
				np := s.Poisson(x, y, phi[volume.Coodinate{X: x, Y: y}])
				dp := np - p
				loss += dp * dp
				nextPressCVMap[volume.Coodinate{X: x, Y: y}] = np
			}
		}
		s.PressCV.CVMap = nextPressCVMap
		if loss < s.Eps {
			fmt.Println(count)
			break
		}
	}
}

// Vの定義点における周囲４点のU
func (s *SimulationController) SurroundingHorVelo(x int, y int) float64 {
	var velo float64 = 0.
	velo += s.HorVeloCV.Get(x-1, y)
	velo += s.HorVeloCV.Get(x-1, y+1)
	velo += s.HorVeloCV.Get(x, y)
	velo += s.HorVeloCV.Get(x-1, y+1)
	velo *= 0.25
	return velo
}

// Uの定義点における周囲４点のV
func (s *SimulationController) SurroundingVerVelo(x int, y int) float64 {
	var velo float64 = 0.
	velo += s.VerVeloCV.Get(x, y)
	velo += s.VerVeloCV.Get(x, y-1)
	velo += s.VerVeloCV.Get(x+1, y)
	velo += s.VerVeloCV.Get(x+1, y-1)
	velo *= 0.25
	return velo
}
