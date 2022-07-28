/* TODOs
計算struct導入する
*/

package main

import (
	"fmt"

	"github.com/marogosteen/cavityflow/core/controller"
	"github.com/marogosteen/cavityflow/core/volume"
	"github.com/marogosteen/cavityflow/services"
)

const (
	cavityGridHeight int = 66
	cavityGridWidth  int = 66

	xl        float64 = 0.02
	dh        float64 = xl / float64(cavityGridHeight-2)
	dt        float64 = 0.001
	dvs       float64 = 1e-6
	rho       float64 = 1000
	eps       float64 = 1e-6
	omega     float64 = 1.1
	initVelo  float64 = 0.
	initPress float64 = 0.4
	mainFlow  float64 = 0.02
	logdir    string  = "../log/"
)

func main() {
	logService := services.LogService{}
	logService.InitializeLogDir(logdir)

	horVelo := volume.NewVolume(cavityGridWidth+1, cavityGridHeight, initVelo)
	verVelo := volume.NewVolume(cavityGridWidth, cavityGridHeight+1, initVelo)
	press := volume.NewVolume(cavityGridWidth, cavityGridHeight, initPress)
	bc := controller.NewBoudaryCondition(mainFlow)
	sc := controller.SimulationController{
		Dt:         dt,
		Dh:         dh,
		Eps:        eps,
		Omega:      omega,
		Conditions: bc,
		HorVelo:    &horVelo,
		VerVelo:    &verVelo,
		Press:      &press,
	}
	sc.SetConditions()
	exec_count := sc.NextPress()
	fmt.Printf("poisson: %d\n", exec_count)

	fp := fmt.Sprintf(logdir+"log%d.csv", 0)
	logService.WriteLog(fp, sc.HorVelo.Grid, sc.VerVelo.Grid, sc.Press.Grid)

	for epoch := 1; ; epoch++ {
		sc.CalcVelocity()
		exec_count := sc.NextPress()
		fp := fmt.Sprintf(logdir+"log%d.csv", epoch)
		if epoch % 10 == 0 {
			logService.WriteLog(fp, sc.HorVelo.Grid, sc.VerVelo.Grid, sc.Press.Grid)
			fmt.Printf("output: %s, velocity loss: %.5f, poisson: %d\n", fp, sc.VelocityLoss, exec_count)
		}

		if sc.VelocityLoss == 0. {
			break
		}
	}

	fmt.Println("done")
}
