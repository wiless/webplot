package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	// _ "github.com/wiless/webplot"

	"go/build"

	"github.com/nightlyone/lockfile"
)

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
	// Simple static webserver:
	fmt.Println("Starting server at 8888")

	portid = ":8888"
	if len(os.Args) == 1 {
		servedir = wwwroot + "/dashboard"
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
