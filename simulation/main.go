package main

/*TODOs
	Indexを扱って座標を示すのを止める．
*/
import (
	"bufio"
	"fmt"
	"os"

	"github.com/marogosteen/cavityflow/controller"
	"github.com/marogosteen/cavityflow/volume"
)

const (
	horCavityGridSize int = 250
	verCavityGridSize int = 251

	xl        float64 = 0.02
	dh        float64 = xl / float64(horCavityGridSize)
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

	for row := 0; row < 252; row++ {
		for col := 0; col < 251; col++ {
			h := (sc.HVeloCV.Get(col, row) + sc.HVeloCV.Get(col+1, row)) / 2
			v := (sc.VVeloCV.Get(col, row) + sc.VVeloCV.Get(col, row+1)) / 2
			p := sc.PressCV.Get(col, row)
			s := fmt.Sprintf("%d,%d,%f,%f,%f\n", row, col, h, v, p)
			_, err = bw.WriteString(s)
			if err != nil {
				panic(err)
			}
		}
	}
	bw.Flush()
}

func main() {
	hVeloCV := volume.NewHVeloCV(horCavityGridSize+1, verCavityGridSize+1, initVelo, mainFlow)
	vVeloCV := volume.NewVVeloCV(horCavityGridSize+2, verCavityGridSize+1, initVelo)
	pressCV := volume.NewPressCV(horCavityGridSize+2, verCavityGridSize+2, initPress)

	sc := controller.SimulationController{
		Dvs:     dvs,
		Dh:      dh,
		Rho:     rho,
		Eps:     eps,
		HVeloCV: hVeloCV,
		VVeloCV: vVeloCV,
		PressCV: pressCV,
	}

	if _, err := os.Stat(logdir); err != nil {
		os.Mkdir(logdir, os.ModePerm)
	}
	sc.BoundaryCondition()
	writeLog(0, sc)

	for epoch := 1; epoch < 5; epoch++ {
		sc.BoundaryCondition()
		sc.NextVelocity()
		phi := sc.Phi(horCavityGridSize, verCavityGridSize)
		sc.NextPress(phi)
		writeLog(epoch, sc)
	}

	fmt.Println("done")
}
