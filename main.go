package main

import (
	"GraphGen/internal"
	"net/http"
)

const port = ":8087"

func main() {
	var g GraphGen.Graph
	g.NewGraph()
	http.HandleFunc("/", g.Handler)
	http.ListenAndServe(port, nil)
}
