package main

/*TODOs
座標をIndexから，XYに変える．
境界条件
時間変化しなかったバグとる
境界条件はシミュレーションが持つべきでは？
param struct必要かな？
equations調整（高度情報が上下逆になったので）
SurroundingVelo
*/

import (
	"bufio"
	"fmt"
	"os"

	"github.com/marogosteen/cavityflow/controller"
	"github.com/marogosteen/cavityflow/volume"
)

const (
	cavityGridHeight int = 254
	cavityGridWidth  int = 254

	xl        float64 = 0.02
	dh        float64 = xl / float64(cavityGridWidth)
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

	for y := 1; y <= cavityGridHeight; y++ {
		for x := 1; x <= cavityGridWidth; x++ {
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

	horVeloCV := volume.NewVeloCV(cavityGridWidth, cavityGridHeight, initVelo)
	verVeloCV := volume.NewVeloCV(cavityGridWidth, cavityGridHeight, initVelo)
	pressCV := volume.NewPressCV(cavityGridWidth, cavityGridHeight, initPress)

	sc := controller.SimulationController{
		Dvs:       dvs,
		Dh:        dh,
		Rho:       rho,
		Eps:       eps,
		HorVeloCV: &horVeloCV,
		VerVeloCV: &verVeloCV,
		PressCV:   &pressCV,
	}
	sc.BoundaryCondition()
	writeLog(0, sc)

	for epoch := 1; epoch < 5; epoch++ {
		sc.NextVelocity()

		phi := sc.Phi(cavityGridWidth, cavityGridHeight)
		sc.NextPress(phi)
		sc.BoundaryCondition()

		writeLog(epoch, sc)
	}

	fmt.Println("done")
}
