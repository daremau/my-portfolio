// Command preview serves the generated site for local checks.
package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":4321", "listen address")
	dir := flag.String("dir", "dist", "directory to serve")
	flag.Parse()

	log.Printf("serving %s at http://localhost%s", *dir, *addr)
	log.Fatal(http.ListenAndServe(*addr, http.FileServer(http.Dir(*dir))))
}
