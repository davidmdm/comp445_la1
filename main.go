package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"regexp"
	"strings"
)

type header map[string]string

var colon, _ = regexp.Compile(":")

func (h header) String() string {
	s := "{\n"
	for key, value := range h {
		s += fmt.Sprintf("  %s: %v\n", key, value)
	}
	s += "}"
	return s
}

func (h header) Set(s string) error {
	if !colon.MatchString(s) {
		return fmt.Errorf("Header value must contain a key:value pair; recieved: %s", s)
	}
	header := strings.Split(s, ":")
	h[header[0]] = header[1]
	return nil
}

func main() {

	var headerMap header

	flag.Var(headerMap, "h", "header key:value pair")

	port := flag.String("p", "80", "The connection port")
	verbose := flag.Bool("v", false, "Verbose mode")

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

	r, _ := regexp.Compile("/")

	pathIndex := r.FindStringIndex(uri)

	var host string
	var path string
	if len(pathIndex) == 0 {
		host = uri
		path = "/"
	} else {
		temp := uri
		host = temp[:pathIndex[0]]
		path = temp[pathIndex[0]:]
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, *port))
	if err != nil {
		log.Fatalf("connection could no be established: %v", err)
	}

	protocol := fmt.Sprintf("%s %s HTTP/1.0\r\nHost: %s\r\nContent-Type: application/x-www-form-urlencoded\r\nContent-Length: 13\r\n\r\nsay=Hi&to=Mom", method, path, host)

	if *verbose == true {
		fmt.Println(protocol, "\n\n")
	}

	fmt.Fprintf(conn, protocol)

	bs, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Fatalf("Could not read response: %v", err)
	}

	fmt.Println(string(bs))

}
