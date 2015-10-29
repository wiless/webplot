package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/nightlyone/lockfile"
	"github.com/wiless/vlib"

	"go/build"

	"golang.org/x/net/websocket"
	// "websocket"

	// "github.com/gorilla/websocket"
)

var ch chan bool

var activePlotter *websocket.Conn

func init() {

	http.Handle("/", websocket.Handler(socketListener))
	http.Handle("/matsock", websocket.Handler(handleMatlabCommands))
	// http.HandleFunc("/series", FetchSeries)

	// 	ch <- true
	// }()

}

func startSocketServer() {
	// go func() {
	addr := "ws://localhost:9999/"
	log.Println("Plot Server listening at ", addr)
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		fmt.Println("Error ", err)
	}
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

// func main() {
// 	startSocketServer()
// 	/// Subscription test
// 	// go EventWatcher()
// 	// <-ch
// }

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

	for {

		log.Printf("Connection Opened from %v", ws.RemoteAddr())
		activePlotter = ws
		cnt := 0
		for {

			in := bufio.NewReader(os.Stdin)
			fmt.Printf("WebPlot > ")
			str, _ := in.ReadString('\n')
			var plotcmd PlotInfo
			// log.Println("Input Command ", str)
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

	}
}

// Handles commands comming from Matlab Session objects and writes to the corresponding
func handleMatlabCommands(ms *websocket.Conn) {
	log.Printf("New Mat Client %s", ms.RemoteAddr())
	io.Copy(activePlotter, ms)

	// for {
	// 	var msg []byte
	// 	msg = make([]byte, 1024)
	// 	// log.Printf("Trying to read..")
	// 	n, err := ms.Read(msg)
	// 	_ = err
	// 	if n > 0 {
	// 		log.Printf("WSServer Rx: %s", msg[0:n])
	// 		activePlotter.Write(msg)
	// 	}
	// }
	log.Println("Leaving session")
}

var servedir string
var portid string
var lock lockfile.Lockfile
var wwwroot string

func init() {
	basePkg := "github.com/wiless/webplot"
	p, err := build.Default.Import(basePkg, "", build.FindOnly)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't find webplot files: %v\n", err)
		fmt.Fprintf(os.Stderr, basePathMessage, basePkg)
		os.Exit(1)
	}
	wwwroot = p.Dir
	log.Printf("My Path %#v", p.Dir)
}

func checkLock() bool {
	flock, err := lockfile.New("/tmp/wmserve.now.lck")
	lock = flock
	if err != nil {
		log.Printf("Cannot init lock. reason: %v", err)
		panic(err)
	}
	err = lock.TryLock()
	// Error handling is essential, as we only try to get the lock.
	if err != nil {
		log.Printf("Cannot lock %v reason: %v", lock, err)
		panic(err)
	}

	//
	return true
}

func main() {
	checkLock()
	defer lock.Unlock()
	go startSocketServer()

	// Simple static webserver:f
	fmt.Println("Starting server at 8888")
	portid = ":8888"
	if len(os.Args) == 1 {
		servedir = wwwroot + "/wsserver"
	} else {
		servedir = os.Args[1]
	}

	if len(os.Args) > 2 {
		portid = ":" + os.Args[2]
	}

	adrs, err := net.InterfaceAddrs()
	// fmt.Printf("\n Folder Listed : %s", servedir)
	for _, adr := range adrs {
		fmt.Printf("\n Open http://%v%s", strings.Split(adr.String(), "/")[0], portid)
	}

	err = http.ListenAndServe(portid, http.FileServer(http.Dir(servedir)))

	if err == nil {
		fmt.Printf("\n Folder Listed : %s", servedir)
		for _, adr := range adrs {
			fmt.Printf("\n Open http://%v%s", strings.Split(adr.String(), "/")[0], portid)
		}
	} else {
		fmt.Println("Error Starting Listen ", err)
	}

}

const basePathMessage = `
By default, webplot locates the slide css,js files static content by looking for a %q package
in your Go workspaces (GOPATH).
You may use the -base flag to specify an alternate location.
`
