package http

import (
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
)

// HeaderMap is a key value map for http headers that implements the flag.Value interface.
type HeaderMap map[string]string

var colon, _ = regexp.Compile(":")

func (h HeaderMap) String() string {
	s := ""
	for key, value := range h {
		s += fmt.Sprintf("%s: %v\r\n", key, value)
	}
	return s
}

// Set deconstructs strings into key value pairs separated by a ":"
func (h HeaderMap) Set(s string) error {
	indexes := colon.FindStringIndex(s)
	if len(indexes) < 2 {
		return fmt.Errorf("Header value must contain a key:value pair; recieved: %s", s)
	}
	h[s[:indexes[0]]] = s[indexes[1]:]
	return nil
}

func uriToHostAndPath(uri string) (string, string) {
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
	return host, path
}

// RequestOptions struct containing fields needed for request options
type RequestOptions struct {
	Uri     string
	Port    string
	Headers HeaderMap
	Verbose bool
	Data    string
}

func send(host, port, protocol string, verbose bool) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return fmt.Errorf("Error establishing connection: %v", err)
	}

	if verbose {
		fmt.Println()
		fmt.Printf("Request:\n\n%s\n\n", protocol)
		fmt.Print("Reponse:\n\n")
	}

	_, err = fmt.Fprint(conn, protocol)
	if err != nil {
		return fmt.Errorf("Error writing to connection: %v", err)
	}

	if _, err = io.Copy(os.Stdout, conn); err != nil {
		return fmt.Errorf("Error reading response: %v", err)
	}

	return nil
}

// Get creates an http get request given a uri, headers and port.
// A verbose boolean can be set to print the http protocol string.
func Get(options RequestOptions) error {
	host, path := uriToHostAndPath(options.Uri)
	options.Headers["Host"] = host
	request := fmt.Sprintf("GET %s HTTP/1.0", path)
	protocol := fmt.Sprintf("%s\r\n%s\r\n", request, options.Headers)

	return send(host, options.Port, protocol, options.Verbose)
}

// Post creates an http post request given a uri, headers, port and data.
// A verbose boolean can be set to print the http protocol string.
func Post(options RequestOptions) error {
	host, path := uriToHostAndPath(options.Uri)
	options.Headers["Host"] = host

	request := fmt.Sprintf("POST %s HTTP/1.0", path)
	protocol := fmt.Sprintf("%s\r\n%s\r\n%s", request, options.Headers, options.Data)

	return send(host, options.Port, protocol, options.Verbose)

}
