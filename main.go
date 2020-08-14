package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var NotFound = ErrorPage(http.StatusNotFound, "Not Found")
var BadGateway = ErrorPage(http.StatusBadGateway, "Bad Gateway")

func main() {

	var err error
	port := 8080
	jsonString := "{\"http://www.example.com\":\"http://www.example.com/elsewhere\"}"
	file := ""
	mappingsString := make(map[string]string)

	flag.IntVar(&port, "port", port, "Port to listen on.")
	flag.StringVar(&jsonString, "json", jsonString, "Mapping JSON. Is overridden by -file.")
	flag.StringVar(&file, "file", file, "Path to mapping JSON file with key and value as URL and replacement URL. Overrides -json.")
	flag.Parse()

	if file != "" {
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		jsonString = string(bytes)
	}

	err = json.Unmarshal([]byte(jsonString), &mappingsString)
	if err != nil {
		panic(err)
	}

	mappings := make([]Mapping, 0)
	for key, value := range mappingsString {
		from, err := url.Parse(key)
		if err != nil {
			panic(fmt.Sprintf("Could not parse URL: %s", key))
		}
		to, err := url.Parse(value)
		if err != nil {
			panic(fmt.Sprintf("Could not parse URL: %s", value))
		}
		mappings = append(mappings, Mapping{
			From: *from,
			To:   *to,
		})
	}

	handler := CompileHandler(mappings, NotFound)

	log.Printf("Listening on http://0.0.0.0:%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), handler)

}

type Mapping struct {
	From url.URL
	To   url.URL
}

func CompileHandler(mappings []Mapping, next http.HandlerFunc) http.HandlerFunc {
	if len(mappings) == 0 {
		return next
	}
	next = CompileHandler(mappings[1:], next)
	mapping := mappings[0]
	log.Printf("Mapping %s -> %s", mapping.From.String(), mapping.To.String())
	baseUrl := url.URL{
		Scheme: mapping.To.Scheme,
		Host:   mapping.To.Host,
	}
	proxy := httputil.NewSingleHostReverseProxy(&baseUrl)
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Host == mapping.From.Host {
			// Host matches
			if len(req.URL.Path) >= len(mapping.From.Path) && req.URL.Path[0:len(mapping.From.Path)] == mapping.From.Path {
				// Path prefix matches

				newUrl := url.URL{
					// Scheme: mapping.To.Scheme,
					// Host:   mapping.To.Host,
					Path: mapping.To.Path + req.URL.Path[len(mapping.From.Path):],
				}

				log.Printf("%s %s http://%s%s -> %s://%s%s\n", req.RemoteAddr, req.Method, req.Host, req.URL.String(), mapping.To.Scheme, mapping.To.Host, newUrl.Path)
				req.Host = mapping.To.Host
				req.URL = &newUrl
				proxy.ServeHTTP(w, req)
				return
			}
		}
		next(w, req)
	})
}

func ErrorPage(code int, text string) http.HandlerFunc {
	bytes := []byte(fmt.Sprintf(`
		<html>
		<head><title>%d %s</title></head>
		<body>
		<center><h1>%d %s</h1></center>
		<hr><center>http-route/1.0.0</center>
		</body>
		</html>
		<!-- a padding to disable MSIE and Chrome friendly error page -->
		<!-- a padding to disable MSIE and Chrome friendly error page -->
		<!-- a padding to disable MSIE and Chrome friendly error page -->
		<!-- a padding to disable MSIE and Chrome friendly error page -->
		<!-- a padding to disable MSIE and Chrome friendly error page -->
		<!-- a padding to disable MSIE and Chrome friendly error page -->
	`, code, text, code, text))
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s http://%s%s -> %d %s\n", req.RemoteAddr, req.Method, req.Host, req.URL.String(), code, text)
		w.WriteHeader(http.StatusNotFound)
		w.Write(bytes)
	})
}
