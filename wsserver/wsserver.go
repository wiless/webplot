package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/wiless/vlib"

	"golang.org/x/net/websocket"
	// "websocket"

	// "github.com/gorilla/websocket"
)

var ch chan bool
var activePlotter *websocket.Conn

func init() {

	http.Handle("/", websocket.Handler(socketListener))
	// http.HandleFunc("/series", FetchSeries)

	// go func() {
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		fmt.Println("Error ", err)
	}
	// 	ch <- true
	// }()

}

/// A ZMQ event  watcher which listens to push messages from the simulator
func EventWatcher() {
	// router, err := goczmq.NewRouter("tcp://*:5555")
	// if err != nil {
	// 	log.Println("Unable to START router")
	// 	// log.Fatal(err)
	// 	return
	// }
	// defer router.Destroy()

	// // type MatlabCommand struct {
	// // 	ID     int
	// // 	Name   string
	// // 	Values map[string]float64
	// // }
	// log.Println("My Identity ", router.Identity())
	for {
		// // log.Println("Waiting for zeromq messages")
		// request, err := router.RecvMessage()

		// if err != nil {
		// 	log.Println("Some Error in Receiving msg from router ", err)
		// } else {

		// }
		// if len(request) == 2 {
		// 	// data := make(map[string]interface{})
		// 	// str := string(request[2])
		// 	dataframe := request[1]
		// 	if len(dataframe) == 0 {
		// 		continue
		// 	}

		// 	log.Printf("SRC Identity %d", request[0])
		// 	log.Printf("Received %s", dataframe)

		// 	jsonData, _ := json.Marshal(plotcmd)
		// 	n, _ := activePlotter.Write(jsonData)
		// 	fmt.Printf("\nSent %d bytes is %s", n, jsonData)

	}
	// }
}

func main() {

	/// Subscription test
	// go EventWatcher()
	// <-ch
}

type PlotOption struct {
	Marker    rune
	LineWidth int
	Color     string
	LineType  int
	Title     string
}

type PlotInfo struct {
	Type    string
	X, Y    vlib.VectorF
	Handle  int
	HoldOn  bool
	Options PlotOption
}

func socketListener(ws *websocket.Conn) {

	// var fa FeedArray
	// SubscriberLists[ws] = fa
	for {

		log.Printf("Connection Opened from %v", ws.RemoteAddr())
		activePlotter = ws
		cnt := 0
		for {

			in := bufio.NewReader(os.Stdin)
			fmt.Printf("Enter Message : ")
			str, _ := in.ReadString('\n')
			var plotcmd PlotInfo
			log.Println("Input Command ", str)
			str = strings.TrimSpace(str)

			if strings.Contains(str, "SIN") {
				log.Printf("YOU ARE IN ", str)
				plotcmd.Type = "plot"
				plotcmd.HoldOn = (str == "SINH")
				plotcmd.Handle = 120
				plotcmd.Y = vlib.RandNFVec(100, 1)
				plotcmd.Options.Color = "red"
				plotcmd.Options.Title = fmt.Sprintf("Figure %d - (%d)", plotcmd.Handle, cnt)
				cnt++
				log.Println(plotcmd.Y)
				data, err := json.Marshal(plotcmd)

				if err != nil {
					log.Println("Err is ", err)
					break
				}
				activePlotter.Write(data)
			}

			if str == "COS" {
				plotcmd.Type = "plot"
				plotcmd.HoldOn = true
				plotcmd.Handle = rand.Intn(100)
				plotcmd.Y = vlib.RandNFVec(100, 1)
				plotcmd.Options.Color = "red"
				plotcmd.Options.Title = fmt.Sprintf("Figure %d", cnt)
				cnt++
				log.Println(plotcmd.Y)
				data, err := json.Marshal(plotcmd)

				if err != nil {
					log.Println("Err is ", err)
					break
				}
				activePlotter.Write(data)
			}

		}
		// {

		// for {
		// 	msg := make([]byte, 1024)
		// 	n, _ := ws.Read(msg)
		// 	if n > 0 {
		// 		log.Printf("Read Message  %s", msg[:n])

		// 		var f FeedInfo
		// 		jerr := json.Unmarshal(msg[0:n], &f)
		// 		if jerr == nil {
		// 			fa = append(fa, f)
		// 			SubscriberLists[ws] = fa
		// 			log.Printf("Updated subscriptions for %v with %v ", ws.RemoteAddr(), fa)
		// 		} else {
		// 			fmt.Println("Error in Unmarshalling ", jerr, " See text ", string(msg[:n]))
		// 			f.ID = 0
		// 			f.FieldNames = []string{"Dummy"}
		// 			fa = append(fa, f)
		// 			SubscriberLists[ws] = fa
		// 			log.Printf("DUMMY Updated subscriptions (%v) for %v with %v ", ws, ws.RemoteAddr(), fa)

		// 		}
		// 	}

		// }
		// }

	}
}
