package controller

import (
	"math"

	"github.com/marogosteen/cavityflow/core/volume"
)

type SimulationController struct {
	Dt    float64
	Dh    float64
	Eps   float64
	Omega float64

	Conditions BoundaryCondition

	HorVelo *volume.Volume
	VerVelo *volume.Volume
	Press   *volume.Volume
}

func (s *SimulationController) SetConditions() {
	s.Conditions.HorVelo(s.HorVelo)
	s.Conditions.VerVelo(s.VerVelo)
	s.Conditions.Press(s.Press)
}

// Navier–Stokes equationsによるflowの計算
func (s *SimulationController) CalcVelocity() {
	nextVerVelo := s.VerVelo.Clone()
	nextHorVelo := s.HorVelo.Clone()

	// TODO magic nuimber
	// // 主流部分の垂直流速の計算は縦に1Grid多い
	for y := 2; y < 251; y++ {
		for x := 1; x < 251; x++ {
			nextVerVelo.Grid[y][x] = s.verNS(x, y)
		}
	}
	for y := 1; y < 251; y++ {
		for x := 2; x < 251; x++ {
			nextVerVelo.Grid[y][x] = s.horNS(x, y)
		}
	}
	s.VerVelo = &nextVerVelo
	s.HorVelo = &nextHorVelo
	// cavity内の計算のみで，境界条件を適応させる必要がある．
	s.SetConditions()
}

func (s *SimulationController) NextPress(phi [][]float64) int {
	for count := 1; ; count++ {
		nextPressCV := s.Press.Clone()
		loss := 0.
		// TODO magic numebr
		for y := 1; y < 251; y++ {
			for x := 1; x < 251; x++ {
				p := s.Press.Grid[y][x]
				np := s.Poisson(x, y, phi[y][x])
				// np = p*(1-s.Omega) + s.Omega*np
				nextPressCV.Grid[y][x] = np

				dp := np - p
				loss += math.Pow(dp, 2)
			}
		}
		s.Press = &nextPressCV
		if loss < s.Eps {
			return count
		}
	}
}

// Vの定義点における周囲４点のU
func (s *SimulationController) SurroundingHorVelo(x int, y int) float64 {
	var velo float64 = 0.
	velo += s.HorVelo.Grid[y][x+1]
	velo += s.HorVelo.Grid[y-1][x+1]
	velo += s.HorVelo.Grid[y][x] 
	velo += s.HorVelo.Grid[y-1][x]
	velo *= 0.25
	return velo
}

// Uの定義点における周囲４点のV
func (s *SimulationController) SurroundingVerVelo(x int, y int) float64 {
	var velo float64 = 0.
	velo += s.VerVelo.Grid[y][x-1]
	velo += s.VerVelo.Grid[y+1][x-1]
	velo += s.VerVelo.Grid[y][x]
	velo += s.VerVelo.Grid[y+1][x]
	velo *= 0.25
	return velo
}

func (s *SimulationController) NewPhi() [][]float64 {
	var phi [][]float64
	for y := 0; y < 252; y++ {
			phi = append(phi, make([]float64, 252))
	}
	// TODO magic number これはNextPressとのつながりがあるはず．
	for y := 1; y < 251; y++ {
		for x := 1; x < 251; x++ {
			phi[y][x] = s.Phi(x, y)
		}
	}
	return phi
}
