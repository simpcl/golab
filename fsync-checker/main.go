package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	CHECKER_FILE_SIZE      = 1 * 1024 * 1023
	CHECKER_FSYNC_SIZE     = 4 * 1024
	CHECKER_FSYNC_INTERVAL = 1000 // milliseconds
)

func main() {
	fsizePtr := flag.Int("fsize", CHECKER_FILE_SIZE, "file size")
	fsyncsizePtr := flag.Int("fsync-size", CHECKER_FSYNC_SIZE, "file sync size")
	intervalPtr := flag.Int("interval", CHECKER_FSYNC_INTERVAL, "file sync interval")
	var fpath string
	flag.StringVar(&fpath, "fpath", "./fsync.dat", "file path")
	flag.Parse()

	fmt.Printf("checker file size: %d\n", *fsizePtr)
	fmt.Printf("checker fsync size: %d\n", *fsyncsizePtr)
	fmt.Printf("checker fsync interval: %d s\n", *intervalPtr)

	f, err := os.Create(fpath)
	if err != nil {
		log.Printf("create file %s error: %s\n", fpath, err.Error())
		return
	}

	buf := make([]byte, *fsyncsizePtr)
	pos := 0
	for {
		begin := time.Now()
		wsize, err := f.Write(buf)
		if err != nil {
			log.Printf("write file error: %s\n", err.Error())
			return
		}
		err = f.Sync()
		if err != nil {
			log.Printf("sync file error: %s\n", err.Error())
			return
		}
		end := time.Now()
		log.Printf("fsync duration: %d us\n", end.Sub(begin).Microseconds())
		pos += wsize
		if pos >= *fsizePtr {
			_, err = f.Seek(0, 0)
			if err != nil {
				log.Printf("seek file error: %s\n", err.Error())
				return
			}
			pos = 0
		}
		time.Sleep(time.Duration(*intervalPtr) * time.Millisecond)
	}
}
