# GoDiffBinary

A simple cli/http tool for check if the given file are different

[![Go Report Card](https://goreportcard.com/badge/github.com/alessiosavi/GoDiffBinary)](https://goreportcard.com/report/github.com/alessiosavi/GoDiffBinary) [![GoDoc](https://godoc.org/github.com/alessiosavi/GoDiffBinary?status.svg)](https://godoc.org/github.com/alessiosavi/GoDiffBinary) [![License](https://img.shields.io/github/license/alessiosavi/GoDiffBinary)](https://img.shields.io/github/license/alessiosavi/GoDiffBinary) [![Version](https://img.shields.io/github/v/tag/alessiosavi/GoDiffBinary)](https://img.shields.io/github/v/tag/alessiosavi/GoDiffBinary) [![Code size](https://img.shields.io/github/languages/code-size/alessiosavi/GoDiffBinary)](https://img.shields.io/github/languages/code-size/alessiosavi/GoDiffBinary) [![Repo size](https://img.shields.io/github/repo-size/alessiosavi/GoDiffBinary)](https://img.shields.io/github/repo-size/alessiosavi/GoDiffBinary) [![Issue open](https://img.shields.io/github/issues/alessiosavi/GoDiffBinary)](https://img.shields.io/github/issues/alessiosavi/GoDiffBinary)
[![Issue closed](https://img.shields.io/github/issues-closed/alessiosavi/GoDiffBinary)](https://img.shields.io/github/issues-closed/alessiosavi/GoDiffBinary)

## Getting Started

This tool is developed in order to compare two different file and verify that they are equals.

It can be run as a CLI (just pass the file as parameter) or in HTTP mode (upload file trought the `/upload` API)

## Prerequisites

Golang have to be installed in the system in order to compile yourself. No other dependencies/packages are need

## Install Golang

In order to install golang in your machine, you have to run the following commands:

- NOTE:  
  - It's preferable __*to don't run these command as root*__. Simply *`chown`* the *`root_foolder`* of golang to be compliant with your user and run the script;  
  - Run this "installer" script only once;  

```bash
golang_version="1.13.4"
golang_link="https://dl.google.com/go/go$golang_version.linux-amd64.tar.gz"
root_foolder="/opt/GOLANG" # Set the tree variable needed for build the enviroinment
go_source="$root_foolder/go"
go_projects="$root_foolder/go_projects"

# Check if this script was alredy run
if [ -d "$root_foolder" ] || [ -d "$go_source" ] || [ -d "$go_projects" ]; then
  echo "Golang is alredy installed!"
  exit 1
fi

# Be sure that golang is not alredy installed
command -v go >/dev/null 2>&1 && { echo >&2 "Seems that go is alredy installed in $(which go)"; exit 2 }

mkdir -p $root_foolder # creating dir for golang source code
cd $root_foolder # entering dir
wget $golang_link #downloading golang
tar xf $(ls | grep "tar") # extract only the tar file
mkdir $go_projects # creating folder for dependencies/project

# Add Go to the current user path
echo '
export GOPATH="$go_projects"
export GOBIN="$GOPATH/bin"
export GOROOT="$go_source"
export PATH="$PATH:$GOROOT/bin:$GOBIN"
' >> /home/$(whoami)/.bashrc

# Load the fresh changed .bashrc env file
source /home/$(whoami)/.bashrc

# Print the golang version
go version
```

After running these command, you have to be able to see the Go version installed and run/compile Go code.

### Post Prerequisites

__*NOTE*__:

- *It's preferable to logout and login to the system for a fresh reload of the configuration after have installed all the packaged listed above.*  

## Installing

GoDiffBinary use the new golang modules manager. You can retrieve the source code by the following command:

```bash
  go get -v -u github.com/alessiosavi/GoDiffBinary
```

In case of problem, you have to download it manually

```bash
  cd $GOPATH/src
  git clone --depth=1 https://github.com/alessiosavi/GoDiffBinary.git
  cd GoDiffBinary
  go clean
  go build
```

## Documentation

### HELP

For print the simple documentation

- Without compile: `go run main.go -help`  
- Compiling: `go build; ./GoDiffBinary -help`  

```text
  -file1 string
        File1 to compare
  -file2 string
        File2 to compare against
  -help
        Print usage
  -http
        Mode, if true will spawn http api instead of command line tool
  -size int
        Dimension that have be compared at each iteration (1k as default)
```

#### Example

Check two equal files in CLI mode:  
`go build; ./GoDiffBinary -file1 GoDiffBinary -file2 GoDiffBinary`

Check two different files in CLI mode:
`go build; ./GoDiffBinary -file1 GoDiffBinary -file2 main.go`

Spawn the HTTP API for upload and check files:  
`go build; ./GoDiffBinary -http`  
Compare same files:  
`curl -F "file1=@GoDiffBinary" -F "file2=@GoDiffBinary"  localhost:8080/upload`  
Compare different files:  
`curl -F "file1=@GoDiffBinary" -F "file2=@main.go"  localhost:8080/upload`

## Built With

- [FastHTTP](https://github.com/valyala/fasthttp) - HTTP Framework | Tuned for high performance. Zero memory allocations in hot paths. Up to 10x faster than net/http  
- [GoGPUtils](https://github.com/alessiosavi/GoGPUtils) - A set of Go methods for enhance productivity  

## Contributing

- Feel free to open issue in order to __*require new functionality*__;  
- Feel free to open issue __*if you discover a bug*__;  
- New idea/request/concept are very appreciated!;  

## Versioning

We use [SemVer](http://semver.org/) for versioning.

## Authors

- **Alessio Savi** - *Initial work & Concept* - [IBM Client Innovation Center [CIC]](https://github.ibm.com/Alessio-Savi)  

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
