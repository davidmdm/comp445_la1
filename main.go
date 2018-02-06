package main

import (
	"bufio"
	"comp445/la1/httpc/http"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	var headers http.HeaderMap = map[string]string{}
	flag.Var(headers, "h", "header key:value pair")
	p := flag.String("p", "80", "The connection port")
	v := flag.Bool("v", false, "Verbose mode")
	d := flag.Bool("d", false, "Data to transmit")
	f := flag.String("f", "", "File to transmit")
	o := flag.String("o", "", "Destination file to write response to. If not specified defaults to standard output")
	r := flag.Bool("r", false, "Follow redirects")

	flag.Parse()

	method := strings.ToUpper(flag.Arg(0))
	if method != "GET" && method != "POST" {
		log.Fatalf("Expecting method to be GET or POST, got: %s", method)
	}

	uri := flag.Arg(1)
	if uri == "" {
		log.Fatal("Uri is required")
	}

	if method != "GET" && method != "POST" {
		log.Fatalf("Method not supported: %v\n", method)
	}

	if method == "GET" && *d == true {
		log.Fatalln("Can only use data option with POST request")
	}

	if method == "GET" && *f != "" {
		log.Fatalln("Can only use file option with POST request")
	}

	if *d == true && *f != "" {
		log.Fatalln("data and file options cannot be used together pick one")
	}

	var data string
	if *d {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Enter the data to transmit")
		fmt.Print("-> ")
		scanner.Scan()
		data = scanner.Text()
		headers["Content-Length"] = strconv.Itoa(len([]byte(data)))
	}

	if *f != "" {
		file, err := ioutil.ReadFile(*f)
		if err != nil {
			log.Fatalf("Could not read file: %v", err)
		}
		headers["Content-Length"] = strconv.Itoa(len(file))
		data = string(file)
	}

	if *o != "" && *o == *f {
		log.Fatalln("Input and output file can not be the same.")
	}

	var w io.Writer
	if *o != "" {
		file, err := os.Create(*o)
		defer file.Close()
		if err != nil {
			log.Fatalf("Could not create file: %v", err)
		}
		w = file
	} else {
		w = os.Stdout
	}

	options := http.RequestOptions{
		Uri:            uri,
		Port:           *p,
		Headers:        headers,
		Verbose:        *v,
		Data:           data,
		W:              w,
		FollowRedirect: *r,
	}

	switch method {
	case "GET":
		if err := http.Get(options); err != nil {
			log.Fatal(err)
		}
	case "POST":
		if err := http.Post(options); err != nil {
			log.Fatal(err)
		}
	}
}
