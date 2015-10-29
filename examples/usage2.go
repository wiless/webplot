package main

import (
	"log"
	"time"

	"github.com/wiless/vlib"
	"github.com/wiless/webplot"
)

func main() {

	// log.Println("Reading after init")
	time.Sleep(4 * time.Second)
	// go func() {
	s := wm.NewSession("SSK Demo")
	// time.Sleep(4 * time.Second)
	log.Println("sending plot")
	for i := 0; i < 5; i++ {
		// fmt.Printf("Start plot")
		s.Plot(vlib.RandUFVec(10), "handle=1", "holdon", "title=CDF Plot of received signal", "LineWidth=2")
		log.Println("Sent Plot ", i)
		time.Sleep(1 * time.Second)
	}

	s.Plot(vlib.RandUFVec(30), "handle=3", "title=SINR CDF", "LineWidth=2")
	s.Plot(vlib.RandUFVec(30), "handle=6", "title=DIP CDF", "LineWidth=2")
	time.Sleep(10 * time.Second)

}
