package http

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"regexp"
	"strings"
)

var colon, _ = regexp.Compile(":")
var slash, _ = regexp.Compile("/")

// RequestOptions struct containing fields needed for request options
type RequestOptions struct {
	Uri     string
	Port    string
	Headers HeaderMap
	Verbose bool
	Data    string
	W       io.Writer
}

// HeaderMap is a key value map for http headers that implements the flag.Value interface.
type HeaderMap map[string]string

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

func uriToHostAndPath(uri string) (host, path string) {
	indexes := slash.FindStringIndex(uri)
	if len(indexes) == 0 {
		host = uri
		path = "/"
	} else {
		host = uri[:indexes[0]]
		path = uri[indexes[0]:]
	}
	return
}

func send(host, port, protocol string, verbose bool, w io.Writer) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return fmt.Errorf("Error establishing connection: %v", err)
	}

	if verbose {
		fmt.Fprintf(w, "\n%s\n\n", protocol)
	}

	_, err = fmt.Fprint(conn, protocol)
	if err != nil {
		return fmt.Errorf("Error writing to connection: %v", err)
	}

	response, err := ioutil.ReadAll(conn)
	if err != nil {
		return fmt.Errorf("Error reading response: %v", err)
	}

	sections := strings.Split(string(response), "\r\n\r\n")

	if verbose {
		fmt.Fprintln(w, sections[0])
	}

	_, err = fmt.Fprintf(w, "\n%s\n", strings.Join(sections[1:], "\r\n\r\n"))

	if err != nil {
		return fmt.Errorf("Could not write response: %v", err)
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
	return send(host, options.Port, protocol, options.Verbose, options.W)
}

// Post creates an http post request given a uri, headers, port and data.
// A verbose boolean can be set to print the http protocol string.
func Post(options RequestOptions) error {
	host, path := uriToHostAndPath(options.Uri)
	options.Headers["Host"] = host
	request := fmt.Sprintf("POST %s HTTP/1.0", path)
	protocol := fmt.Sprintf("%s\r\n%s\r\n%s", request, options.Headers, options.Data)
	return send(host, options.Port, protocol, options.Verbose, options.W)
}
