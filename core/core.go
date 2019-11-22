package core

import (
	"bytes"
	"io"
	"log"
	"os"

	fileutils "github.com/alessiosavi/GoGPUtils/files"
)

// CompareBinaryFile is delegated to compare two files using chunks of byte
func CompareBinaryFile(file1, file2 string, nByte int) int {
	var size1, size2 int64
	var err, err1, err2 error
	if nByte < 1 {
		log.Println("Chunks of bytes size not provided, using 1k byte")
		nByte = 1024
	}

	if !fileutils.FileExists(file1) {
		log.Fatal("File [", file1, "] does not exist!")
	}

	if !fileutils.FileExists(file2) {
		log.Fatal("File [", file2, "] does not exist!")
	}

	// Get file size of file1
	size1, err = fileutils.GetFileSize(file1)
	if err != nil {
		log.Fatal("Unable to read file [" + file1 + "]")
	}

	// Get file size of file2
	size2, err = fileutils.GetFileSize(file2)
	if err != nil {
		log.Fatal("Unable to read file [" + file2 + "]")
	}

	// Compare file size (disabled)

	if size1 != size2 {
		log.Println("Size of ["+file1+"]-> ", size1)
		log.Println("Size of ["+file2+"]-> ", size2)
		log.Println("Files are not equals! Dimension mismatch!")
		return 1
	}

	// Open first file
	fdFile1, err := os.Open(file1)
	if err != nil {
		log.Fatal("Error while opening file", err)
	}
	// Close file at return
	defer fdFile1.Close()

	log.Printf("%s opened\n", file1)

	// Open second file
	fdFile2, err := os.Open(file1)
	if err != nil {
		log.Fatal("Error while opening file", err)
	}

	defer fdFile2.Close()

	log.Printf("%s opened\n", file2)

	// Read 1k bytes at iteration
	data1 := make([]byte, nByte)
	data2 := make([]byte, nByte)
	for err1 != io.EOF || err2 != io.EOF {
		_, err1 = fdFile1.Read(data1)
		if err1 != nil && err1 != io.EOF {
			log.Fatal("Error on file [" + file1 + "] -> [" + err1.Error() + "]")
		}

		_, err2 = fdFile2.Read(data2)
		if err2 != nil && err2 != io.EOF {
			log.Fatal("Error on file [" + file2 + "] -> [" + err2.Error() + "]")
		}

		if !bytes.Equal(data1, data2) {
			var pos1, pos2 int64
			pos1, _ = fdFile1.Seek(0, 1)
			pos2, _ = fdFile1.Seek(0, 1)
			log.Println("Files are not equals! At position [Pos1:", pos1, "Pos2:", pos2, "]")
			return 1
		}
	}

	log.Println("Files [", file1, "-", file2, "] are equal!")
	return 0
}

func Check(err error, description ...string) {
	if err != nil {
		log.Println("Error occurred in [", description, "] -> ["+err.Error()+"]")
	}
}
