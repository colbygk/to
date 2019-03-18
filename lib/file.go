package lib

import (
	"os"
	"syscall"
)

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
