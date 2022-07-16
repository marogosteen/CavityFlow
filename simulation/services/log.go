package services

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

type LogService struct {
}

func (s *LogService) InitializeLogDir(dir string) {
	if _, err := os.Stat(dir); err != nil {
		os.Mkdir(dir, os.ModePerm)
	}

	pattern := dir + "*.csv"
	files, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}

func (s *LogService) WriteLog(fp string, hvv [][]float64, vvv [][]float64, pv [][]float64) {
	f, err := os.Create(fp)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	bw := bufio.NewWriter(f)
	bw.WriteString("y,x,u,v,p\n")
	fmt.Println("conplete: " + fp)

	// TODO magic number
	for y := 0; y < 252; y++ {
		for x := 0; x < 252; x++ {
			u := (hvv[y][x] + hvv[y][x+1]) / 2
			v := (vvv[y][x] + vvv[y+1][x]) / 2
			p := pv[y][x]
			s := fmt.Sprintf("%d,%d,%f,%f,%f\n", y, x, u, v, p)
			_, err = bw.WriteString(s)
			if err != nil {
				panic(err)
			}
		}
	}
	bw.Flush()
}
