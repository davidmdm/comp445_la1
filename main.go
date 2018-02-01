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

func main() {

	method := flag.String("v", "GET", "The http verb")
	port := flag.String("p", "80", "The connection port")
	uri := flag.String("u", "", "The uri of the request")

	flag.Parse()

	*method = strings.ToUpper(*method)

	if *method != "GET" && *method != "POST" {
		log.Fatalf("Method not supported: %v\n", *method)
	}

	if *uri == "" {
		log.Fatalln("uri is required")
	}

	r, _ := regexp.Compile("/")

	pathIndex := r.FindStringIndex(*uri)

	var host string
	var path string
	if len(pathIndex) == 0 {
		host = *uri
		path = "/"
	} else {
		temp := *uri
		host = temp[:pathIndex[0]]
		path = temp[pathIndex[0]:]
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, *port))
	if err != nil {
		log.Fatalf("connection could no be established: %v", err)
	}

	fmt.Fprintf(conn, fmt.Sprintf("%s %s HTTP/1.0\r\nHost: %s\r\n\r\n", *method, path, host))

	bs, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Fatalf("Could not read response: %v", err)
	}

	fmt.Println(string(bs))

}
