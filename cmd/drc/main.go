package main

import (
	"fmt"
	"github.com/dantin/mysql-tools/drc"
)

// main is the bootstrap.
func main() {
	cfg := drc.NewConfig()
	fmt.Println(cfg)
	fmt.Println("hello drc.")
}
