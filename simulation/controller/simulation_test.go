package controller

import (
	"testing"

	"github.com/marogosteen/cavityflow/volume"
)

const (
	cavityGridHeight int = 254
	cavityGridWidth  int = 254

	xl        float64 = 0.02
	dh        float64 = xl / float64(250)
	dt        float64 = 0.001
	dvs       float64 = 1e-6
	rho       float64 = 1000
	eps       float64 = 1e-6
	initVelo  float64 = 0.
	initPress float64 = 0.4
	mainFlow  float64 = 0.02
	logdir    string  = "../log/"
)

func Benchmark_CalclatePressure(b *testing.B) {
	horVeloCV := volume.NewVeloCV(cavityGridWidth, cavityGridHeight, initVelo)
	verVeloCV := volume.NewVeloCV(cavityGridWidth, cavityGridHeight, initVelo)
	pressCV := volume.NewPressCV(cavityGridWidth, cavityGridHeight, initPress)

	sc := SimulationController{
		Dt:        dt,
		Dvs:       dvs,
		Dh:        dh,
		Rho:       rho,
		Eps:       eps,
		MainFlow:  mainFlow,
		HorVeloCV: &horVeloCV,
		VerVeloCV: &verVeloCV,
		PressCV:   &pressCV,
	}
	sc.BoundaryCondition()

	phi := volume.NewCVMap(cavityGridWidth, cavityGridHeight, 0.0)
	// TODO magic number これはNextPressとのつながりがあるはず．
	for y := 3; y <= 252; y++ {
		for x := 3; x <= 252; x++ {
			phi[volume.Coodinate{X: x, Y: y}] = sc.Phi(x, y)
		}
	}

	b.ResetTimer()
	sc.NextPress(phi)
}
