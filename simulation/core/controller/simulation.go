package controller

import (
	"math"

	"github.com/marogosteen/cavityflow/core/volume"
)

type SimulationController struct {
	VelocityLoss float64

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
	for y := 2; y < 65; y++ {
		for x := 1; x < 65; x++ {
			nextVerVelo.Grid[y][x] = s.verNS(x, y)
		}
	}
	for y := 1; y < 65; y++ {
		for x := 2; x < 65; x++ {
			nextHorVelo.Grid[y][x] = s.horNS(x, y)
		}
	}

	s.calcVelocityLoss(nextHorVelo, nextVerVelo)

	s.VerVelo = &nextVerVelo
	s.HorVelo = &nextHorVelo
	// cavity内の計算のみで，境界条件を適応させる必要がある．
	s.Conditions.VerVelo(s.VerVelo)
	s.Conditions.HorVelo(s.HorVelo)
}

func (s *SimulationController) NextPress() int {
	phi := s.newPhi()
	for count := 1; ; count++ {
		next := s.Press.Clone()
		loss := 0.
		// TODO magic numebr
		for y := 1; y < 65; y++ {
			for x := 1; x < 65; x++ {
				np := s.Poisson(x, y, phi[y][x])
				// np = p*(1-s.Omega) + s.Omega*np
				next.Grid[y][x] = np

				dp := np - s.Press.Grid[y][x]
				loss += math.Pow(dp, 2)
			}
		}
		s.Press = &next
		// cavity内の計算のみで，境界条件を適応させる必要がある．
		s.Conditions.Press(s.Press)

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

func (s *SimulationController) newPhi() [][]float64 {
	var phi [][]float64
	for y := 0; y < 65; y++ {
		phi = append(phi, make([]float64, 65))
	}
	// TODO magic number これはNextPressとのつながりがあるはず．
	for y := 1; y < 65; y++ {
		for x := 1; x < 65; x++ {
			phi[y][x] = s.Phi(x, y)
		}
	}
	return phi
}

func (s *SimulationController) calcVelocityLoss(hvv volume.Volume, vvv volume.Volume) float64 {
	s.VelocityLoss = 0.
	for y := 1; y < s.HorVelo.Height-1; y++ {
		for x := 2; x < s.HorVelo.Width-2; x++ {
			s.VelocityLoss += math.Abs(s.HorVelo.Grid[y][x] - hvv.Grid[y][x])
		}
	}
	for y := 2; y < s.VerVelo.Height-2; y++ {
		for x := 1; x < s.VerVelo.Width-1; x++ {
			s.VelocityLoss += math.Abs(s.VerVelo.Grid[y][x] - vvv.Grid[y][x])
		}
	}
	return s.VelocityLoss
}
