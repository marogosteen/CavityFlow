package controller

import (
	"math"
)

const (
	rho float64 = 1000
	dvs float64 = 1e-6
)

func (s *SimulationController) horNS(x int, y int) float64 {
	u := s.HorVelo.Grid[y][x]
	lu := s.HorVelo.Grid[y][x-1]
	ru := s.HorVelo.Grid[y][x+1]
	ou := s.HorVelo.Grid[y+1][x]
	uu := s.HorVelo.Grid[y-1][x]
	v4 := s.SurroundingVerVelo(x, y)
	p := s.Press.Grid[y][x]
	lp := s.Press.Grid[y][x-1]

	var at float64 = 0
	at += (u + math.Abs(u)) / 2 * (u - lu) / s.Dh
	at += (u - math.Abs(u)) / 2 * (ru - u) / s.Dh
	at += (v4 + math.Abs(v4)) / 2 * (u - uu) / s.Dh
	at += (v4 - math.Abs(v4)) / 2 * (ou - u) / s.Dh

	pt := (p - lp) / (s.Dh * rho)
	dt := (ru + lu + uu + ou - 4*u) / math.Pow(s.Dh, 2) * dvs
	newUVelo := u - s.Dt*(at+pt-dt)
	return newUVelo
}

func (s *SimulationController) verNS(x int, y int) float64 {
	v := s.VerVelo.Grid[y][x]
	lv := s.VerVelo.Grid[y][x-1]
	rv := s.VerVelo.Grid[y][x+1]
	ov := s.VerVelo.Grid[y+1][x]
	uv := s.VerVelo.Grid[y-1][x]
	u4 := s.SurroundingHorVelo(x, y)
	p := s.Press.Grid[y][x]
	up := s.Press.Grid[y-1][x]

	var at float64 = 0
	at += (u4 + math.Abs(u4)) / 2 * (v - lv) / s.Dh
	at += (u4 - math.Abs(u4)) / 2 * (rv - v) / s.Dh
	at += (v + math.Abs(v)) / 2 * (v - uv) / s.Dh
	at += (v - math.Abs(v)) / 2 * (ov - v) / s.Dh

	pt := (p - up) / (s.Dh * rho)
	dt := (rv + lv + uv + ov - 4*v) / math.Pow(s.Dh, 2) * dvs
	newVVelo := v - s.Dt*(at+pt-dt)
	return newVVelo
}

func (s *SimulationController) Poisson(x int, y int, phi float64) float64 {
	var p float64 = 0
	p += s.Press.Grid[y][x+1]
	p += s.Press.Grid[y][x-1]
	p += s.Press.Grid[y+1][x]
	p += s.Press.Grid[y-1][x]
	p -= phi * math.Pow(s.Dh, 2)
	p /= 4
	return p
}

func (s *SimulationController) Phi(x int, y int) float64 {
	u := s.HorVelo.Grid[y][x]
	ou := s.HorVelo.Grid[y+1][x]
	uu := s.HorVelo.Grid[y-1][x]
	ru := s.HorVelo.Grid[y][x+1]
	r2u := s.HorVelo.Grid[y][x+2]
	lu := s.HorVelo.Grid[y][x-1]
	oru := s.HorVelo.Grid[y+1][x+1]
	uru := s.HorVelo.Grid[y-1][x+1]
	u4 := s.SurroundingHorVelo(x, y)
	ou4 := s.SurroundingHorVelo(x, y+1)

	v := s.VerVelo.Grid[y][x]
	ov := s.VerVelo.Grid[y+1][x]
	o2v := s.VerVelo.Grid[y+2][x]
	uv := s.VerVelo.Grid[y-1][x]
	rv := s.VerVelo.Grid[y][x+1]
	lv := s.VerVelo.Grid[y][x-1]
	orv := s.VerVelo.Grid[y+1][x+1]
	olv := s.VerVelo.Grid[y+1][x-1]
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

	return rho * (pp1 + pp2 + pp3)
}
