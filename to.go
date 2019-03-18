package main

import (
	"os"
	"flag"
	"fmt"
    "strings"
    "time"
	toLib "cgk.sh/to/lib"
)


// start flags
var CurrentTime = time.Now()
var Year = flag.Int("y", CurrentTime.Year(), "year")
var Mon = flag.String("m", strings.ToLower(CurrentTime.Month().String()[:3]), "month")
var Day = flag.Int("d", CurrentTime.Day(), "day")
var printHelp = flag.Bool("h", false, "Show this helpful information.")
var HowFarBack = flag.Int("b", 30, "How far back to look for previous tracking files in days")
var ProjectName = flag.String("p", "tracking", "Name of file used for tracking")
var TrackingPath *toLib.TemporalLocation
var defPathFormat string
var GitHash string
var GitFetchURL string
var TimeBuilt string

func handleFlags() {
	version := flag.Bool("version", false, "Show version")
	flag.Parse()
	if *version {
		fmt.Printf("Version: %s\n", GitHash)
        fmt.Printf("Built: %s\n", TimeBuilt)
        fmt.Printf("Origin: %s\n", GitFetchURL)
        os.Exit(0)
	}
}


func getArgs() []string {
	return os.Args[1:]
}

func parseCommands() {
    args := getArgs()
    if len(args) < 1 || *printHelp {
        flag.PrintDefaults()
        os.Exit(1)
    }
    switch args[0] {
    case "do":
        fmt.Println("do-ray")
    case "day":
        fmt.Println("hoo-day")
    case "morrow":
        fmt.Println("hoo-ray")
    }
}


func init() {    
	handleFlags()
    parseCommands()

    TrackingPath = toLib.GetTrackingPath(*Year, *Mon, *Day, *ProjectName)
}


func main() {
	var trackingFile *toLib.ToFile
	var err error

    toLib.CopyMostRecentTrackingFile(TrackingPath, *HowFarBack)
	trackingFile, err =
        toLib.OpenTracking(toLib.GetFQTrackingPath(TrackingPath))
	if err != nil {
        fmt.Printf("Unable to open %s\n*%s!\n",
            toLib.GetFQTrackingPath(TrackingPath),
            err)
        os.Exit(-1)
	}

	trackingTree := toLib.ParseTracking(trackingFile)
	toLib.PrintTree(trackingTree)
}
