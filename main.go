// EXAMPLE USAGE:
// The tool use 2 parameter [file1,file2]. You can use the relative/absolute path
// go run main.go -file1 /path/of/file1 -file2 /path/of/file2
package main

import (
	"flag"
	"log"

	"github.com/alessiosavi/GoDiffBinary/api"
	"github.com/alessiosavi/GoDiffBinary/core"
	stringutils "github.com/alessiosavi/GoGPUtils/string"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	var file1, file2 string
	var size int
	var mode bool
	file1, file2, size, mode = verifyCommandLineInput()

	if !mode {
		core.CompareBinaryFile(file1, file2, size)
		return
	}
	log.Println("Http enabled!")
	//api.InitAPIGin("localhost", "8080")
	api.InitAPIFasthttp("localhost", "8080", size)
}

// verifyCommandLineInput verify about the INPUT parameter passed as arg[]
func verifyCommandLineInput() (string, string, int, bool) {
	log.Println("VerifyCommandLineInput | START")
	file1 := flag.String("file1", "", "File1 to compare")
	file2 := flag.String("file2", "", "File2 to compare against")
	dimension := flag.Int("size", 0, "Dimension that have be compared at each iteration")
	mode := flag.Bool("http", false, "Mode, if true will spawn http api instead of command line tool")

	flag.Parse()
	if *mode {
		return "", "", *dimension, *mode
	}
	log.Println("Cli mode enabled - Going to work in cli mode (use -http true for enable HTTP API)")
	if stringutils.IsBlank(*file1) {
		flag.PrintDefaults()
		log.Fatal("file1 parameter not passed in input!")
	}
	if stringutils.IsBlank(*file2) {
		flag.PrintDefaults()
		log.Fatal("file2 parameter not passed in input!")
	}
	return *file1, *file2, *dimension, false
}
