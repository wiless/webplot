package wm

import (
	"encoding/json"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/wiless/vlib"

	"golang.org/x/net/websocket"
	// "websocket"

	// "github.com/gorilla/websocket"
)

var ch chan bool
var activePlotter *websocket.Conn
var SessionCommand chan PlotInfo

var (
	MATCOLORS      []string = []string{"c", "r", "b", "k", "g", "y"}
	MATCOLORS_full []string = []string{"cyan", "red", "blue", "black", "green", "yellow"}
	MATMARKERS     []string = []string{"+", "x", "*", "o", "s", "."}
	MATLINETYPES   []string = []string{"-", "--", ":"}
)

func init() {

	// http.Handle("/", websocket.Handler(socketListener))
	// http.Handle("/matsock", websocket.Handler(matsocketListener))
	// http.HandleFunc("/series", FetchSeries)
	// activePlotter = nil
	SessionCommand = make(chan PlotInfo)
	// go func() {
	// 	err := http.ListenAndServe(":9999", nil)
	// 	if err != nil {
	// 		fmt.Println("Error ", err)
	// 	}
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

type PlotOption struct {
	Marker    string
	LineWidth int
	Color     string
	LineType  string
	Title     string
	holdOn    bool
	handle    int
}

type PlotInfo struct {
	Type    string
	X, Y    vlib.VectorF
	Handle  int
	HoldOn  bool
	Options PlotOption
}

func RemoveString(a []string, i int) []string {
	a = append(a[:i], a[i+1:]...)
	return a
}

func FindAny(input string, option []string) (int, bool) {
	var result int
	found := false
	result = -1
	// substr = strings.ToUpper(substr)
	for i, v := range option {

		if strings.Contains(input, v) {
			result = i
			found = true
			break
		}
	}
	return result, found
}
func FindStringA(strA []string, substr string) (int, bool) {
	var result int
	found := false
	result = -1
	substr = strings.ToUpper(substr)
	for i, v := range strA {
		v = strings.ToUpper(v)
		if strings.Contains(v, substr) {
			result = i
			found = true

		}
	}
	return result, found
}

// example command : plot(y, "holdoff", "r*-","LineWidth:2")
func (p *PlotOption) Parse(param ...string) {
	L := len(param)
	// optionstring := strings.Join(param, ",")
	// for i := 0; i < L; i++ {
	// 	// proces holdon / holdoff
	// 	if strings.Contains(optionstring, "holdon") {
	// 		p.holdOn = true
	// 	}
	// 	if strings.Contains(optionstring, "holdoff") {
	// 		p.holdOn = false
	// 	}

	// 	//Process LineWidth

	// }
	findx, found := FindStringA(param, "Handle")
	if found {
		values := strings.Split(param[findx], "=")
		h, _ := (strconv.ParseInt(values[1], 10, 64))
		p.handle = int(h)
		L--
		param = RemoveString(param, findx)
	}
	/// Check for HOLD
	findx, found = FindStringA(param, "holdon")
	if found {
		p.holdOn = true
		L--
		param = RemoveString(param, findx)
	}
	findx, found = FindStringA(param, "holdoff")
	if found {
		p.holdOn = false
		param = RemoveString(param, findx)
		L--
	}

	// LineWidth
	findx, found = FindStringA(param, "LineWidth")
	if found {

		values := strings.Split(param[findx], "=")
		lw, _ := (strconv.ParseInt(values[1], 10, 64))
		p.LineWidth = int(lw)
		L--
		param = RemoveString(param, findx)
	}
	// LineType
	findx, found = FindStringA(param, "LineType")
	if found {
		log.Println("Found Line type !!")
		values := strings.Split(param[findx], "=")
		// lt, _ :=
		p.LineType = values[1]
		L--
		param = RemoveString(param, findx)
	}
	// Color
	findx, found = FindStringA(param, "Color")
	if found {

		values := strings.Split(param[findx], "=")
		// lt, _ :=
		p.Color = values[1]
		L--
		param = RemoveString(param, findx)
	}
	// Title
	findx, found = FindStringA(param, "Title")
	if found {

		values := strings.Split(param[findx], "=")
		// lt, _ :=
		p.Title = values[1]
		L--
		param = RemoveString(param, findx)
	}
	/// Remaining look for combination of Markers and predefined color
	// Color
	findx, found = FindStringA(param, "style")
	if found {
		// log.Printf("Found %v", param[findx])
		values := strings.Split(param[findx], "=")
		// lt, _ :=
		style := values[1]

		//
		// Look for color
		idx, found := FindAny(style, MATCOLORS)
		if found {
			p.Color = MATCOLORS_full[idx]
		}
		// Look for Marker
		idx, found = FindAny(style, MATMARKERS)
		if found {
			p.Marker = MATMARKERS[idx]
		}
		// Look for LineStyle
		idx, found = FindAny(style, MATLINETYPES)
		if found {
			p.LineType = MATLINETYPES[idx]
		}
		L--
		param = RemoveString(param, findx)
	}
	if len(param) > 0 {
		log.Printf("Unknown Parameters %#v", param)
	}

}

type MatlabSession struct {
	prefix      string
	CMDWindow   chan PlotInfo
	browsercn   *websocket.Conn
	matlabshell *websocket.Conn
}

// Starts a session and triggers the browser pointing to dashboard
func NewDashboard(name string) *MatlabSession {

	result := &MatlabSession{}
	result.CMDWindow = make(chan PlotInfo)
	result.Connect() // Connects to the MatSocket
	go result.Listen()
	// socket server not ready has to trigger or give wait time..
	cmd := exec.Command("xdg-open", "http://localhost:8888")
	cmd.Output()
	log.Println("Client Ready")
	result.prefix = name
	return result
}

func (m *MatlabSession) Connect() {

	origin := "http://localhost/"
	url := "ws://localhost:9999/matsock"
	// url = "ws://localhost:12345/echo"
	var err error
	m.matlabshell, err = websocket.Dial(url, "", origin)

	if err != nil {
		log.Fatal(err)
	}
	// Say hello to the matlab shell

	n, err := m.matlabshell.Write([]byte("Matlab Connecting... !"))

	log.Printf("Wrote %d writing Error %v ", n, err)
	// var msg = make([]byte, 512)
	// if n, err = m.matlabshell.Read(msg); err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Received: %s.\n", msg[:n])
	// var msg = make([]byte, 512)
	// var n int
	// if n, err = ws.Read(msg); err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Received: %s.\n", msg[:n])

	// m.matlabshell=
}
func (m *MatlabSession) Listen() {
	for {
		select {
		case plotcmd := <-m.CMDWindow:
			if m.matlabshell != nil {

				log.Printf("Listen :: Read from Channel %v", plotcmd.Options.Title)
				data, err := json.Marshal(plotcmd)
				if err != nil {
					log.Println("Json Marshal Err is ", err)
					break
				} else {

					// log.Println("Greeting again ..")
					n, err := m.matlabshell.Write(data)

					// m.matlabshell.Write([]byte("Hello again .. "))
					if n > 0 {

						log.Printf("Attempting %d bytes [Error %v]", n, err)
						if err != nil {
							m.matlabshell = nil
						}

					}
				}
			}
			// default:
			// log.Println("Waited for you")
		}
	}
}

func NewSession(name string) *MatlabSession {
	result := &MatlabSession{}
	result.CMDWindow = make(chan PlotInfo)
	result.Connect() // Connects to the MatSocket
	go result.Listen()
	// socket server not ready has to trigger or give wait time..

	if !CheckActiveSession() {
		// cmd := exec.Command("xdg-open", "http://localhost:8888")
		// cmd.Output()
		// log.Println("Client Ready")
	}
	result.prefix = name
	return result
}

func NewSessionOld(name string) *MatlabSession {
	result := &MatlabSession{}
	if !CheckActiveSession() {
		// cmd := exec.Command("xdg-open", "http://localhost:8888")
		// cmd.Output()
		// log.Println("Client Ready")
	}
	result.prefix = name
	return result
}

func CheckActiveSession() bool {
	/// Ideally check if there is already an active session
	// Use Ping-pong to test
	// log.Println("activePlotter : ", activePlotter)
	return activePlotter != nil /// Not fool-proof
}

func (m *MatlabSession) Plot(y vlib.VectorF, params ...string) int {
	var p PlotInfo
	p.Y = y
	p.Type = "plot"
	p.Options.Parse(params...)
	p.Handle = p.Options.handle
	p.HoldOn = p.Options.holdOn ///
	p.Options.Title = "[ " + m.prefix + " ] " + p.Options.Title
	log.Printf("Write to Channel %#v", p.Options.Title)
	m.CMDWindow <- p
	// SessionCommand <- p
	return p.Handle
}
func (m *MatlabSession) Scatter(y vlib.VectorF, params ...string) int {
	var p PlotInfo
	p.Y = y
	p.Type = "scatter"
	p.Options.Parse(params...)
	p.Handle = p.Options.handle
	p.HoldOn = p.Options.holdOn ///
	p.Options.Title = "[ " + m.prefix + " ] " + p.Options.Title
	log.Printf("Sending plot %#v", p.Options.Title)
	// m.CMDWindow <- p
	SessionCommand <- p
	return p.Handle
}
func (m *MatlabSession) PlotXY(x, y vlib.VectorF, params ...string) int {
	var p PlotInfo
	p.X = x
	p.Y = y
	p.Type = "plot"
	p.Options.Parse(params...)
	p.Handle = p.Options.handle
	p.HoldOn = p.Options.holdOn ///
	p.Options.Title = "[ " + m.prefix + " ] " + p.Options.Title
	log.Printf("Sending plot %#v", p.Options.Title)
	m.CMDWindow <- p
	//SessionCommand <- p
	return p.Handle
}
func (m *MatlabSession) ScatterC(c vlib.VectorC, params ...string) int {

	return m.PlotXY(c.Real(), c.Imag(), params...)

}

// func main() {

// 	log.Println("Reading after init")

// 	go func() {
// 		s := NewSession("HETNET")

// 		for i := 0; i < 10; i++ {

// 			s.plot(vlib.RandUFVec(10), "holdon", "title=CDF Plot of received signal", "style=b+", "LineWidth=2")
// 			time.Sleep(5 * time.Second)
// 		}
// 	}()

// 	s := NewSession("Single Cell")

// 	for i := 0; i < 10; i++ {

// 		s.plot(vlib.RandUFVec(10), "holdon", "title=CDF Plot of received signal", "style=b+", "LineWidth=2")
// 		time.Sleep(1 * time.Second)
// 	}

// 	// wait if someone closes
// 	<-ch

// }

// func socketListener(ws *websocket.Conn) {

// 	/// Allowing only one plotting session
// 	// if activePlotter != nil {
// 	// 	log.Printf("Denying  %v", ws.RemoteAddr())
// 	// 	return
// 	// }

// 	log.Printf("Connection Opened from %v", ws.RemoteAddr())

// 	// var fa FeedArray
// 	// SubscriberLists[ws] = fa

// 	activePlotter = ws
// 	// var msg []byte
// 	// go func() {
// 	// 	n, err := activePlotter.Read(msg)
// 	// 	if err != nil {
// 	// 		log.Println(err)
// 	// 	}
// 	// 	fmt.Printf("Found something %s", msg[0:n])
// 	// }()

// 	for {

// 		select {
// 		case plotcmd := <-SessionCommand:

// 			log.Printf("Received Matlab Command %#v", plotcmd.Options.Title)
// 			data, err := json.Marshal(plotcmd)
// 			log.Println("JSON Marshal Err is ", err)

// 			if err != nil {
// 				log.Println("Err is ", err)
// 				break
// 			} else {
// 				n, err := activePlotter.Write(data)
// 				if err != nil {
// 					activePlotter = nil

// 				} else {
// 					log.Printf("Wrote %d bytes [Error %v]", n, err)

// 				}
// 			}

// 		}
// 		if ws == nil {
// 			break
// 		}
// 	}

// }
