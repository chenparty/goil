package zlog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// zWriter ...
type zWriter struct {
	baseDir  string
	dateDir  string
	dnFormat string
	fileName string
	fnFormat string
	ext      string
	maxBytes int64

	fd *os.File
	mu sync.Mutex
}

var logger *log.Logger

// NewDefault a logger
func NewDefault() {
	New("", log.Ldate|log.Ltime|log.Lshortfile)
}

// New a logger
func New(prefix string, flag int) {
	w := &zWriter{
		baseDir:  "./zlog",
		dateDir:  "",
		dnFormat: "200601",
		fileName: "",
		fnFormat: "02",
		ext:      ".log",
		maxBytes: 1024,
	}
	logger = log.New(w, prefix, flag)
}

// Write ...
func (w *zWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.dateDir = time.Now().Format(w.dnFormat)
	dirPath := filepath.Join(w.baseDir, w.dateDir)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, 0666)
		w.openFile()
	} else {
		if fi, err := w.fd.Stat(); err == nil {
			fName := time.Now().Format(w.fnFormat)
			if fName != w.fileName {
				w.fd.Close()
				w.openFile()
			} else if fi.Size() > w.maxBytes {
				w.fd.Close()
				w.rotate()
				w.openFile()
			}
		} else {
			w.openFile()
		}
	}
	return w.fd.Write(p)
}

func (w *zWriter) openFile() {
	w.fileName = time.Now().Format(w.fnFormat)
	fPath := filepath.Join(w.baseDir, w.dateDir, w.fileName+w.ext)
	w.fd, _ = os.OpenFile(fPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

func (w *zWriter) rotate() {
	oldPath := filepath.Join(w.baseDir, w.dateDir, w.fileName+w.ext)
	newFileName := fmt.Sprintf("%s_%d%s", w.fileName, time.Now().Unix(), w.ext)
	newPath := filepath.Join(w.baseDir, w.dateDir, newFileName)
	os.Rename(oldPath, newPath)
}

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	logger.Output(2, fmt.Sprint(v...))
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	logger.Output(2, fmt.Sprintf(format, v...))
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	logger.Output(2, fmt.Sprintln(v...))
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	logger.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	logger.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) {
	logger.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

// Panic is equivalent to Print() followed by a call to panic().
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	logger.Output(2, s)
	panic(s)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	logger.Output(2, s)
	panic(s)
}

// Panicln is equivalent to Println() followed by a call to panic().
func Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	logger.Output(2, s)
	panic(s)
}
