package main

import (
  "flag"
  "fmt"
  "os"
  "strings"
  "time"

  toLib "cgk.sh/to/lib"
)

// CurrentTime as in right now
var CurrentTime = time.Now()

// Year year
var Year = flag.Int("y", CurrentTime.Year(), "year")

// Mon month
var Mon = flag.String("m", strings.ToLower(CurrentTime.Month().String()[:3]), "month")

// Day day
var Day = flag.Int("d", CurrentTime.Day(), "day")
var printHelp = flag.Bool("h", false, "Show this helpful information.")

// HowFarBack how far to look back
var HowFarBack = flag.Int("b", 30, "How far back to look for previous tracking files in days")

// ProjectName name of the project
var ProjectName = flag.String("p", "tracking", "Name of file used for tracking")

// TrackingPath where tracking files are kept
var TrackingPath *toLib.TemporalLocation
var defPathFormat string

// GitHash of the project for versioning
var GitHash string

// GitFetchURL of the project for versioning
var GitFetchURL string

// TimeBuilt of the project for versioning
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
  case "day":
  case "do":
  case "json":
  case "morrow":
  default:
    fmt.Println("Unknown command")
    flag.PrintDefaults()
    os.Exit(1)
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

  trackingTree, err := toLib.ParseTracking(trackingFile)
  if err != nil {
    fmt.Printf("Error: %s", err)
    os.Exit(-1)
  }
  toLib.PrintJSONTree(trackingTree)
  // toLib.PrintJSONTree(trackingTree)
}
