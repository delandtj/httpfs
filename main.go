// webshare serves the current directory on port 3000.
package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.Formatter = new(logrus.JSONFormatter)
	log.Formatter = new(logrus.TextFormatter)
	log.Level = logrus.DebugLevel
}

func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Method, r.URL.Path)

		fn := path.Join(".", r.URL.Path)
		file, err := os.Open(fn)
		if err == nil {
			defer file.Close()
			fi, err := file.Stat()
			if err == nil {
				log.WithFields(logrus.Fields{
					filename: file.Name(),
					filesize: fi.Size(),
				}).Debug("file served")
			}
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	port := flag.Int("p", 3000, "port to listen on")
	flag.Parse()
	http.Handle("/", loggingHandler(http.FileServer(http.Dir("."))))

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				log.Printf("http://%s:%d", ipnet.IP.String(), *port)
			}
		}
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatal(err)
	}
}
