package main

import (
	"fmt"
	"os"

	"github.com/l-pavlova/image-master/imageParse"
)

func main() {

	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) == 0 {
		err := fmt.Errorf("%s", "Pass correct command line arguments to the program\n")
		fmt.Println(err)
		return
	}

	imageParse.ReadFrom()
}