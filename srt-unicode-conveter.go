package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	var srtPath = flag.String("p", "nil", "Input srt path to convert")
	flag.Parse()

	if srtPath == nil || len(*srtPath) == 0 || *srtPath == "nil" {
		panic("set path to file with -p")
	}

	fileBytes, err := ioutil.ReadFile(*srtPath)

	if err != nil {
		panic("what the fuck is this file?")
	}

	text := string(fileBytes)

	if len(text) == 0 {
		panic("file is empty")
	}

	result, err := detectEncoding(text)
	if err != nil {
		panic("cannot find unicode of file")
	}

	fmt.Printf(
		"Detected charset is %s, language is %s\n",
		result.Charset,
		result.Language)

	convertResult := convertToUTF8(text, result.Charset)

	result, err = detectEncoding(convertResult)

	if err != nil {
		panic("converted encoding not recognized!")
	}

	fmt.Printf(
		"New Detected charset is %s, language is %s\n",
		result.Charset,
		result.Language)

	srtPathStrong := *srtPath
	srtFileName := srtPathStrong[0 : len(srtPathStrong)-len(filepath.Ext(*srtPath))]

	newName := fmt.Sprintf("%s-utf-8%s", srtFileName, filepath.Ext(*srtPath))

	// open output file
	outputFile, err := os.Create(newName)
	if err != nil {
		panic("cannot create output file!")
	}

	defer func() {
		if err := outputFile.Close(); err != nil {
			panic(err)
		}
	}()

	err = ioutil.WriteFile(newName, []byte(convertResult), 0644)

	if err != nil {
		panic("cannot write to file")
	}

	fmt.Printf(
		"Srt converted successfully : %s \n",
		newName)
}

func detectEncoding(str string) (r *chardet.Result, err error) {
	detector := chardet.NewTextDetector()
	return detector.DetectBest([]byte(str))
}

func convertToUTF8(str string, origEncoding string) string {
	strBytes := []byte(str)
	byteReader := bytes.NewReader(strBytes)
	reader, _ := charset.NewReaderLabel(origEncoding, byteReader)
	strBytes, _ = ioutil.ReadAll(reader)
	return string(strBytes)
}
