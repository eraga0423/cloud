package main

import (
	"flag"

	i "triple-s/internal"
)

func main() {
	flag.Parse()
	i.Server()
}
