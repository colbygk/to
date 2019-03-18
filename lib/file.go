package lib

import (
    "errors"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
    "syscall"
    "time"
)

var defPathFormat string

func init() {
    // LeadingPath/TemporalLocation
    defPathFormat = "%s/%s"
}

type ToFile struct {
	fullPath string
	fd       int
	rawData  []byte
	err      error
	stat     os.FileInfo
	fileP    *os.File
	lines    [][]byte
	current  int
}

type TemporalLocation struct {
    loc time.Time
    leadingPath string
    projectName string
    timeFormat string
}

func createTemporalLocation(t time.Time, leading, project, format string) (*TemporalLocation) {
    newloc := new(TemporalLocation)
    newloc.loc = t
    newloc.leadingPath = leading
    newloc.projectName = project
    newloc.timeFormat = format
    return newloc
}

func duplicateTemporalLocation(from *TemporalLocation) (*TemporalLocation) {
    newPath := new(TemporalLocation)
    newPath.loc = from.loc
    newPath.leadingPath = from.leadingPath
    newPath.projectName = from.projectName
    newPath.timeFormat = from.timeFormat
    return newPath
}

func updateTemporalLocation(on *TemporalLocation, newTime time.Time) {
    on.loc = newTime
    return
}

func DefaultTrackingPath() string {
	return Getenv("HOME", ".") + "/notes"
}

func GetTrackingPath(year int, mon string, day int, project string) *TemporalLocation {
    format := "2006/Jan/02"
    t,_ := time.Parse("2006/Jan/02", fmt.Sprintf("%d/%s/%02d", year, mon, day))
    newPath := createTemporalLocation(
        t,
        Getenv("TO_NOTES_PATH", DefaultTrackingPath()),
        project,
        format)
    
	return newPath
}

func getFileName(t *TemporalLocation) string {
    return t.projectName + ".to"
}

func getDirTrackingPath(t *TemporalLocation) string {
    return filepath.FromSlash(getLocation(t))
}

func GetFQTrackingPath(t *TemporalLocation) string {
    return filepath.FromSlash(getLocation(t) + "/" + getFileName(t))
}

func getLocation(on *TemporalLocation) string {
    return fmt.Sprintf(defPathFormat,
        Getenv("TO_NOTES_PATH", DefaultTrackingPath()),
        strings.ToLower(on.loc.Format(on.timeFormat)))
}

func CreateCurrentDayDir(TrackingPath *TemporalLocation) {
    if !DirExists(getDirTrackingPath(TrackingPath)) {
        err := os.MkdirAll(getDirTrackingPath(TrackingPath), 0755)
        if err != nil {
            panic(err)
        }
    }
}

func CopyMostRecentTrackingFile(TrackingPath *TemporalLocation, howFar int) error {
    CreateCurrentDayDir(TrackingPath)
    if !FileExists(GetFQTrackingPath(TrackingPath)) {
        fmt.Printf("Looking for tracking file\n")
        olde := MostRecentTracking(TrackingPath, howFar)
        if olde.loc != TrackingPath.loc {
            fmt.Printf("  found and copying from: %s\n", getDirTrackingPath(olde))
            return CopyFile(GetFQTrackingPath(olde),GetFQTrackingPath(TrackingPath))
        } else {
            return errors.New(fmt.Sprintf("Never found %s! Create it %s\n",
                getFileName(TrackingPath), getDirTrackingPath(TrackingPath)))
        }
    }
    return nil
}

func backADay(t *TemporalLocation) {
    t.loc = t.loc.AddDate(0,0,-1)
}

func MostRecentTracking(t *TemporalLocation, howFar int) *TemporalLocation {
    olde := duplicateTemporalLocation(t)
    back := time.Now()
    back = back.AddDate(0,0,-1 * howFar)
    for !olde.loc.Before(back) {
        if FileExists(GetFQTrackingPath(olde)) {
            return olde
        }
        backADay(olde)
    }
    return t
}

func Getenv(key, fallback string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		v = fallback
	}
	return v
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	return !os.IsNotExist(err) && info.Mode().IsDir()
}

func FileExists(path string) bool {
    info, err := os.Stat(path)
    return !os.IsNotExist(err) && !info.Mode().IsDir()
}

func OpenTracking(fullPath string) (*ToFile, error) {
	var tF *ToFile

	tF = new(ToFile)
	tF.fullPath = fullPath

	tF.stat, tF.err = os.Stat(fullPath)
	if tF.err != nil {
		return nil, tF.err
	}

	tF.fileP, tF.err = os.OpenFile(fullPath, os.O_CREATE|os.O_RDONLY, 0)
	if tF.err != nil {
		return nil, tF.err
	}
	tF.fd = int(tF.fileP.Fd())

	tF.rawData, tF.err = syscall.Mmap(tF.fd, 0, int(tF.stat.Size()), syscall.PROT_READ, syscall.MAP_PRIVATE)
	if tF.err != nil {
		return nil, tF.err
	}

	return tF, nil
}

/*
 * Copied from: https://github.com/mactsouk/opensource.com/blob/master/cp3.go
 */
func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file.", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = os.Stat(dst)
	if err == nil {
		return fmt.Errorf("File %s already exists.", dst)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	if err != nil {
		panic(err)
	}

    // Default buffer size 1000000
	buf := make([]byte, 1000000)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}
