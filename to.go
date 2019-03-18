package main

import (
	"os"
	"strings"
	"flag"
	"fmt"
	"path/filepath"
	"time"
	toLib "cgk.sh/to/lib"
)

// start flags
var CurrentTime = time.Now()
var Year = flag.Int("y", CurrentTime.Year(), "year")
var Mon = flag.String("m", strings.ToLower(CurrentTime.Month().String()[:3]), "month")
var Day = flag.Int("d", CurrentTime.Day(), "day")
var NotesData = flag.String("f", "tracking.to", "Notes data file name")

func init() {
	handleFlags()
}

func GetNotesPath() string {
	return filepath.FromSlash(toLib.Getenv("TO_NOTES_PATH", DefaultNotesPath()))
}

func DefaultNotesPath() string {
	return fmt.Sprintf("%s/%d/%s/%02d/",
		toLib.Getenv("HOME", ".")+"/notes",
		*Year, *Mon, *Day)
}

func handleFlags() {
	version := flag.Bool("version", false, "Show version")
	flag.Parse()
	if *version {
		fmt.Println("Version: Alpha")
	}
}

func handleArgs() []string {
	return os.Args[1:]
}

func main() {
	var trackingFile *toLib.ToFile
	var err error
	var notesPath string
	var fileName string

	notesPath = GetNotesPath()

	if !toLib.DirExists(notesPath) {
		err = os.MkdirAll(notesPath, 0755)
		if err != nil {
			panic(err)
		}
	}

	fileName = *NotesData

	trackingFile, err = toLib.OpenTracking(GetNotesPath() + fileName)
	if err != nil {
		panic(err)
	}

	trackingTree := toLib.ParseTracking(trackingFile)
	toLib.PrintTree(trackingTree)
}
