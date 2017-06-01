package main

import (
	"fmt"

	"strconv"

	"./stringutil"	
)

func main() {

	defer fmt.Printf("\n\n")

	for i := 0; i < 10; i++ {
		defer fmt.Printf(strconv.Itoa(i))
		fmt.Printf(stringutil.Reverse("\n!oG ,olleH"))
	}
}