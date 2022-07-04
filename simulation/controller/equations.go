package controller

import (
	"math"
)

func (s *SimulationController) horNS(x int, y int) float64 {
	u := s.HorVeloCV.Get(x, y)
	lu := s.HorVeloCV.Get(x-1, y)
	ru := s.HorVeloCV.Get(x+1, y)
	ou := s.HorVeloCV.Get(x, y+1)
	uu := s.HorVeloCV.Get(x, y-1)
	v4 := s.SurroundingVerVelo(x, y)
	p := s.PressCV.Get(x, y)
	lp := s.PressCV.Get(x-1, y)

	var uat float64 = 0
	uat += (u + math.Abs(u)) / 2 * (u - lu) / s.Dh
	uat += (u - math.Abs(u)) / 2 * (ru - u) / s.Dh
	uat += (v4 + math.Abs(v4)) / 2 * (u - uu) / s.Dh
	uat += (v4 - math.Abs(v4)) / 2 * (ou - u) / s.Dh

	pt := (p - lp) / (s.Dh * s.Rho)
	dt := (ru + lu + uu + ou - 4*u/math.Pow(s.Dh, 2)*s.Dvs)
	newUVelo := u - s.Dt*(uat+pt-dt)
	return newUVelo
}

func (s *SimulationController) verNS(x int, y int) float64 {
	v := s.VerVeloCV.Get(x, y)
	lv := s.VerVeloCV.Get(x-1, y)
	rv := s.VerVeloCV.Get(x+1, y)
	ov := s.VerVeloCV.Get(x, y+1)
	uv := s.VerVeloCV.Get(x, y-1)
	u4 := s.SurroundingHorVelo(x, y)
	p := s.PressCV.Get(x, y)
	up := s.PressCV.Get(x, y-1)

	var vat float64 = 0
	vat += (u4 + math.Abs(u4)) / 2 * (v - uv) / s.Dh
	vat += (u4 - math.Abs(u4)) / 2 * (ov - v) / s.Dh
	vat += (v + math.Abs(v)) / 2 * (v - lv) / s.Dh
	vat += (v - math.Abs(v)) / 2 * (rv - v) / s.Dh

	pt := (p - up) / (s.Dh * s.Rho)
	dt := (rv + lv + uv + ov - 4*v/math.Pow(s.Dh, 2)*s.Dvs)
	newVVelo := v - s.Dt*(vat+pt-dt)
	return newVVelo
}

func (s *SimulationController) Poisson(x int, y int, phi float64) float64 {
	var newP float64 = 0
	newP += s.PressCV.Get(x+1, y)
	newP += s.PressCV.Get(x-1, y)
	newP += s.PressCV.Get(x, y+1)
	newP += s.PressCV.Get(x, y-1)
	newP -= phi * s.Dh * s.Dh
	newP /= 4
	return newP
}

func (s *SimulationController) Phi(x int, y int) float64 {
	u := s.HorVeloCV.Get(x, y)
	ou := s.HorVeloCV.Get(x, y+1)
	uu := s.HorVeloCV.Get(x, y-1)
	ru := s.HorVeloCV.Get(x+1, y)
	r2u := s.HorVeloCV.Get(x+2, y)
	lu := s.HorVeloCV.Get(x-1, y)
	oru := s.HorVeloCV.Get(x+1, y+1)
	uru := s.HorVeloCV.Get(x+1, y-1)
	u4 := s.SurroundingHorVelo(x, y)
	ou4 := s.SurroundingHorVelo(x, y+1)

	v := s.VerVeloCV.Get(x, y)
	ov := s.VerVeloCV.Get(x, y+1)
	o2v := s.VerVeloCV.Get(x, y+2)
	uv := s.VerVeloCV.Get(x, y-1)
	rv := s.VerVeloCV.Get(x+1, y)
	lv := s.VerVeloCV.Get(x-1, y)
	orv := s.VerVeloCV.Get(x+1, y+1)
	olv := s.VerVeloCV.Get(x-1, y+1)
	v4 := s.SurroundingVerVelo(x, y)
	rv4 := s.SurroundingVerVelo(x+1, y)

	var (
		pp1 float64 = 0
		pp2 float64 = 0
		pp3 float64 = 0
	)

	pp1 += ((ru-u)/s.Dh + (ov-v)/s.Dh) / s.Dt

	pp2 += ru*(r2u-u)/(2*s.Dh) + rv4*(oru-uru)/(2*s.Dh)
	pp2 -= u*(ru-lu)/(2*s.Dh) + v4*(ou-uu)/(2*s.Dh)
	pp2 *= -1 / s.Dh

	pp3 += ou4*(orv-olv)/(2*s.Dh) + ov*(o2v-v)/(2*s.Dh)
	pp3 -= u4*(rv-lv)/(2*s.Dh) + v*(ov-uv)/(2*s.Dh)
	pp3 *= -1 / s.Dh

	return s.Rho * (pp1 + pp2 + pp3)
}
