package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/l-pavlova/image-master/app"
)

func main() {

	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) == 0 {
		err := fmt.Errorf("%s", "Pass correct command line arguments to the program\n")
		fmt.Println(err)
		return
	}
	command := argsWithoutProg[0]
	imageMaster := app.NewImageMaster()

	if strings.Contains(command, "-find=") {
		objectToFind := strings.Split(command, "=")[1]
		err := imageMaster.Find(objectToFind)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	switch command {
	case "-grayscale":
		//imageMaster.GrayScale()
	case "-smoothen":
		imageMaster.Smoothen(10)
	case "-sharpen":
		//imageMaster.GrayScale()
	case "-denoise":

	}
}

//docker run  --name cont -it --mount type=bind,source="$(pwd)",target=/app/images imasterps
