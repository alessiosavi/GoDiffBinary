package api

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	fileutils "github.com/alessiosavi/GoGPUtils/files"
	stringutils "github.com/alessiosavi/GoGPUtils/string"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
)

// InitAPI is delegated to initialize the API and bind the connection the the related resources
func InitAPIGin(hostname, port string) {

	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.POST("/check", uploadGin(""))
	router.Run(hostname + ":" + port)
}

func uploadGin(parameterName string) gin.HandlerFunc {

	fn := func(c *gin.Context) {
		file, err := c.FormFile(parameterName)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, "/tmp/"+filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully", file.Filename))
	}
	return fn
}

// HandleRequests is the hook the real function/wrapper for expose the API. It's main scope it's to map the url to the function that have to do the work.
// It take in input the pointer to the list of file to server; The pointer to the datastructure.Configuration in order to change the parameter at runtime;the channel used for thread safety
func InitAPIFasthttp(hostname, port string) {
	log.Println("HandleRequests | START")
	m := func(ctx *fasthttp.RequestCtx) { // Hook to the API methods "magilogically"
		ctx.Response.Header.Set("GoLog-Viewer", "v0.0.1$/beta") // Set an header just for track the version of the software
		log.Println("REQUEST --> ", ctx, " | Headers: ", ctx.Request.Header.String())
		tmpChar := "============================================================"
		switch string(ctx.Path()) {
		case "/":
			FastHomePage(ctx, hostname, port) // Simply print some link
			log.Println(tmpChar)

		case "/upload":
			uploadFasthttp(ctx, "file1", "file2")

		default:
			ctx.WriteString("The url " + string(ctx.URI().RequestURI()) + " does not exist :(\n")
			FastHomePage(ctx, hostname, port) // Simply print some link
			log.Println(tmpChar)
		}
	}

	// The gzipHandler will serve a compress request only if the client request it with headers (Content-Type: gzip, deflate)
	gzipHandler := fasthttp.CompressHandlerLevel(m, fasthttp.CompressBestCompression) // Compress data before sending (if requested by the client)

	s := &fasthttp.Server{
		Handler: gzipHandler,
		// Every response will contain 'Server: Check diff binary server' header.
		Name: "Check diff binary server",
		// Other Server settings may be set here.
		MaxRequestBodySize: 4 * 1024 * 1024 * 1024,
	}

	err := s.ListenAndServe(hostname + ":" + port) // Try to start the server with input "host:port" received in input
	if err != nil {                                // No luck, connection not successfully. Probably port used ...
		log.Fatalln("Unable to spawn server")
	}
	log.Println("HandleRequests | STOP")
}

func uploadFasthttp(ctx *fasthttp.RequestCtx, file1, file2 string) {

	var f1, f2 string
	fh, err := ctx.FormFile(file1)
	if err != nil {
		panic(err)
	}

	f1 = "/tmp/" + fh.Filename + "_" + stringutils.RandomString(5)
	if err := fasthttp.SaveMultipartFile(fh, f1); err != nil {
		panic(err)
	}

	fh, err = ctx.FormFile(file2)
	if err != nil {
		panic(err)
	}
	f2 = "/tmp/" + fh.Filename + "_" + stringutils.RandomString(5)
	if err := fasthttp.SaveMultipartFile(fh, f2); err != nil {
		panic(err)
	}

	ret := compareBinaryFile(f1, f2)
	if ret == 0 {
		ctx.WriteString("Files are equal\n")
	} else {
		ctx.WriteString("Files are different\n")
	}
}

// FastHomePage is the methods for serve the home page. It print the list of file that you can query with the complete link in order to copy and paste easily
func FastHomePage(ctx *fasthttp.RequestCtx, hostname, port string) {
	log.Println("FastHomePage | START")
	ctx.Response.Header.SetContentType("text/plain; charset=utf-8")
	ctx.WriteString("Welcome to the GoDiffBinary!\n" + "API List!\n" +
		"http://" + hostname + ":" + port + "/check -> Check binary file\n")

	log.Println("FastHomePage | STOP")
}

func compareBinaryFile(file1, file2 string) int {
	var size1, size2 int64
	var err, err1, err2 error

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
	data1 := make([]byte, 1024)
	data2 := make([]byte, 1024)
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
