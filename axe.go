package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	errorNoFile = errors.New("no filename given")
)

func main() {
	log.SetLevel(log.InfoLevel)

	splitFlag := flag.Bool("split", false, "Split a file")
	joinFlag := flag.Bool("join", false, "Join a file")
	fileNameFlag := flag.String("file", "", "File to split/join")
	b64encodeflag := flag.Bool("b64", false, "Encode to base64")
	splitIntoFlag := flag.Int("number", 2, "Split into N files")
	flag.Parse()

	if !*splitFlag && !*joinFlag {
		log.Fatal("I need an operation mode split/join")
	}

	if *fileNameFlag == "" {
		log.Fatalf(errorNoFile.Error())
	}

	if *splitFlag {
		bufA, err := ioutil.ReadFile(*fileNameFlag)
		checkError(err)

		log.Debugf("File length: %d", len(bufA))

		var cnt int
		for i := 0; i < len(bufA); i += (len(bufA) / *splitIntoFlag) {
			top := i + (len(bufA) / *splitIntoFlag)
			if top > len(bufA) {
				top = len(bufA)
			}

			log.Debugf("%d:%d", i, top)

			var bufB []byte
			if *b64encodeflag {
				bufB = []byte(base64.StdEncoding.EncodeToString(bufA[i:top]))
			} else {
				bufB = bufA[i:top]
			}

			ioutil.WriteFile(fmt.Sprintf("%s.split%d", *fileNameFlag, cnt), bufB, 0600)

			cnt++
		}

		os.Exit(0)
	}

	if *joinFlag {
		_, err := os.Stat(*fileNameFlag)
		if err == nil {
			log.Fatal("File already exists, plese give me the output file name")
		}

		var i int
		var b []byte
		for {
			buf, err := ioutil.ReadFile(fmt.Sprintf("%s.split%d", *fileNameFlag, i))
			if err != nil {
				break
			}

			var data []byte
			if *b64encodeflag {
				data, err = base64.StdEncoding.DecodeString(string(buf))
				checkError(err)
			} else {
				data = buf
			}

			b = append(b, data[:]...)

			i++
		}

		ioutil.WriteFile(*fileNameFlag, b, 0600)

		os.Exit(0)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}
