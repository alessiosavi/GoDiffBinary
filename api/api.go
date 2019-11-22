package api

import (
	"log"
	"os"

	"github.com/alessiosavi/GoDiffBinary/core"
	helper "github.com/alessiosavi/GoGPUtils/helper"
	stringutils "github.com/alessiosavi/GoGPUtils/string"
	"github.com/valyala/fasthttp"
)

// HandleRequests is the hook the real function/wrapper for expose the API. It's main scope it's to map the url to the function that have to do the work.
// It take in input the pointer to the list of file to server; The pointer to the datastructure.Configuration in order to change the parameter at runtime;the channel used for thread safety
func InitAPIFasthttp(hostname, port string, size int) {
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
			uploadFasthttp(ctx, "file1", "file2", size)

		default:
			_, err := ctx.WriteString("The url " + string(ctx.URI().RequestURI()) + " does not exist :(\n")
			core.Check(err)
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
		MaxRequestBodySize: 100 * 1024 * 1024,
	}

	log.Println("Max size allowed (per file) -> ", helper.ConvertSize(float64(s.MaxRequestBodySize), "MB"), " MB")
	err := s.ListenAndServe(hostname + ":" + port) // Try to start the server with input "host:port" received in input
	if err != nil {                                // No luck, connection not successfully. Probably port used ...
		log.Fatalln("Unable to spawn server")
	}
	log.Println("HandleRequests | STOP")
}

func uploadFasthttp(ctx *fasthttp.RequestCtx, file1, file2 string, size int) {

	var f1, f2 string
	var err error
	fh, err := ctx.FormFile(file1)
	if err != nil {
		panic(err)
	}

	f1 = "/tmp/" + fh.Filename + "_" + stringutils.RandomString(5)
	if err = fasthttp.SaveMultipartFile(fh, f1); err != nil {
		panic(err)
	}

	fh, err = ctx.FormFile(file2)
	if err != nil {
		panic(err)
	}
	f2 = "/tmp/" + fh.Filename + "_" + stringutils.RandomString(5)
	if err = fasthttp.SaveMultipartFile(fh, f2); err != nil {
		panic(err)
	}

	ret := core.CompareBinaryFile(f1, f2, size)
	if ret == 0 {
		_, err = ctx.WriteString("Files are equal\n")
		core.Check(err, "Writing to stdout result of comparing")
	} else {
		_, err = ctx.WriteString("Files are different\n")
		core.Check(err, "Writing to stdout result of comparing")
	}

	err = os.Remove(f1)
	core.Check(err, "Removing ["+f1+"]")

	err = os.Remove(f2)
	core.Check(err, "Removing ["+f2+"]")
}

// FastHomePage is the methods for serve the home page. It print the list of file that you can query with the complete link in order to copy and paste easily
func FastHomePage(ctx *fasthttp.RequestCtx, hostname, port string) {
	log.Println("FastHomePage | START")
	ctx.Response.Header.SetContentType("text/plain; charset=utf-8")
	_, err := ctx.WriteString("Welcome to the GoDiffBinary!\n" + "API List!\n" +
		"http://" + hostname + ":" + port + "/check -> Check binary file\n")

	core.Check(err, "Writing to stdout result FastHomePage")
	log.Println("FastHomePage | STOP")
}
