package controller

import (
	"math"
)

const (
	rho       float64 = 1000
	dvs       float64 = 1e-6
)

func (s *SimulationController) horNS(x int, y int) float64 {
	u := s.HorVelo.Get(x, y)
	lu := s.HorVelo.Get(x-1, y)
	ru := s.HorVelo.Get(x+1, y)
	ou := s.HorVelo.Get(x, y+1)
	uu := s.HorVelo.Get(x, y-1)
	v4 := s.SurroundingVerVelo(x, y)
	p := s.Press.Get(x, y)
	lp := s.Press.Get(x-1, y)

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
	v := s.VerVelo.Get(x, y)
	lv := s.VerVelo.Get(x-1, y)
	rv := s.VerVelo.Get(x+1, y)
	ov := s.VerVelo.Get(x, y+1)
	uv := s.VerVelo.Get(x, y-1)
	u4 := s.SurroundingHorVelo(x, y)
	p := s.Press.Get(x, y)
	up := s.Press.Get(x, y-1)

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
	p += s.Press.Get(x+1, y)
	p += s.Press.Get(x-1, y)
	p += s.Press.Get(x, y+1)
	p += s.Press.Get(x, y-1)
	p -= phi * math.Pow(s.Dh, 2)
	p /= 4
	return p
}

func (s *SimulationController) Phi(x int, y int) float64 {
	u := s.HorVelo.Get(x, y)
	ou := s.HorVelo.Get(x, y+1)
	uu := s.HorVelo.Get(x, y-1)
	ru := s.HorVelo.Get(x+1, y)
	r2u := s.HorVelo.Get(x+2, y)
	lu := s.HorVelo.Get(x-1, y)
	oru := s.HorVelo.Get(x+1, y+1)
	uru := s.HorVelo.Get(x+1, y-1)
	u4 := s.SurroundingHorVelo(x, y)
	ou4 := s.SurroundingHorVelo(x, y+1)

	v := s.VerVelo.Get(x, y)
	ov := s.VerVelo.Get(x, y+1)
	o2v := s.VerVelo.Get(x, y+2)
	uv := s.VerVelo.Get(x, y-1)
	rv := s.VerVelo.Get(x+1, y)
	lv := s.VerVelo.Get(x-1, y)
	orv := s.VerVelo.Get(x+1, y+1)
	olv := s.VerVelo.Get(x-1, y+1)
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
