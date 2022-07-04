package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/marogosteen/cavityflow/controller"
	"github.com/marogosteen/cavityflow/volume"
)

const (
	cavityGridHeight int = 254
	cavityGridWidth  int = 254

	xl        float64 = 0.02
	dh        float64 = xl / float64(cavityGridWidth)
	dt        float64 = 0.001
	dvs       float64 = 1e-6
	rho       float64 = 1000
	eps       float64 = 1e-6
	initVelo  float64 = 0.
	initPress float64 = 0.4
	mainFlow  float64 = 0.02
	logdir    string  = "../log/"
)

func writeLog(epoch int, sc controller.SimulationController) {
	fp := fmt.Sprintf(logdir+"log%d.csv", epoch)
	f, err := os.Create(fp)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	bw := bufio.NewWriter(f)
	bw.WriteString("y,x,u,v,p\n")
	fmt.Println("conplete: " + fp)

	// TODO magic number
	for y := 3; y <= cavityGridHeight-1; y++ {
		for x := 3; x <= cavityGridWidth-2; x++ {
			u := (sc.HorVeloCV.Get(x, y) + sc.HorVeloCV.Get(x+1, y)) / 2
			v := (sc.VerVeloCV.Get(x, y) + sc.VerVeloCV.Get(x, y+1)) / 2
			p := sc.PressCV.Get(x, y)
			s := fmt.Sprintf("%d,%d,%f,%f,%f\n", y, x, u, v, p)
			_, err = bw.WriteString(s)
			if err != nil {
				panic(err)
			}
		}
	}
	bw.Flush()
}

func main() {
	if _, err := os.Stat(logdir); err != nil {
		os.Mkdir(logdir, os.ModePerm)
	}

	pattern := logdir + "*.csv"
	files, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}

	horVeloCV := volume.NewVeloCV(cavityGridWidth, cavityGridHeight, initVelo)
	verVeloCV := volume.NewVeloCV(cavityGridWidth, cavityGridHeight, initVelo)
	pressCV := volume.NewPressCV(cavityGridWidth, cavityGridHeight, initPress)

	sc := controller.SimulationController{
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
	writeLog(0, sc)

	for epoch := 1; epoch < 20; epoch++ {
		sc.NextVelocity()
		phi := volume.NewCVMap(cavityGridWidth, cavityGridHeight, 0.0)
		// TODO magic number これはNextPressとのつながりがあるはず．
		for y := 3; y <= 252; y++ {
			for x := 3; x <= 252; x++ {
				phi[volume.Coodinate{X: x, Y: y}] = sc.Phi(x, y)
			}
		}

		sc.NextPress(phi)
		sc.BoundaryCondition()

		writeLog(epoch, sc)
	}

	fmt.Println("done")
}
