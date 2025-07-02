package main

import (
	"fmt"

	"github.com/yankeguo/rg"
)

func main() {
	var err error
	defer rg.Guard(&err)
	fmt.Println("Hello, World!")
}
