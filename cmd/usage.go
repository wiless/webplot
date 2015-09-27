package main

import (
	"log"
	"time"

	"github.com/wiless/vlib"
	"github.com/wiless/webplot"
)

func main() {

	log.Println("Reading after init")

	go func() {
		s := wm.NewSession("HETNET")

		for i := 0; i < 10; i++ {

			s.Plot(vlib.RandUFVec(10), "holdon", "title=CDF Plot of received signal", "LineWidth=2")
			time.Sleep(5 * time.Second)
		}
	}()
	time.Sleep(4 * time.Second)
	s := wm.NewSession("Single Cell")
	nsamples := 50
	x := vlib.NewVectorF(nsamples)
	for i := 0; i < nsamples; i++ {
		x[i] = float64(i) * 10
	}
	NPLOTS := 3
	for i := 0; i < NPLOTS; i++ {

		// s.Plot(vlib.RandUFVec(10), "handle=4", "holdon", "title=CDF Plot of received signal", "style=+", "LineWidth=2")
		if i < 2 {

			if i == 1 {

				s.PlotXY(x, vlib.RandUFVec(nsamples), "handle=4", "holdon", "title=CDF Plot of received signal", "LineWidth=2")
			} else {
				y := x.Add(5.5)
				s.PlotXY(y, vlib.RandUFVec(nsamples), "handle=4", "holdff", "title=CDF Plot of received signal", "LineWidth=2")
			}

		} else {
			s.Plot(vlib.RandUFVec(nsamples), "handle=4", "holdon", "title=CDF Plot of received signal", "style=+", "LineWidth=2")
		}
		time.Sleep(4 * time.Second)
	}

	// wait if someone closes

}
