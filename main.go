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
	fmt.Println("we running")

	command := argsWithoutProg[0]
	imageMaster := app.NewImageMaster()

	if strings.Contains(command, "-find=") {
		//run tensorflow
	}
	switch command {
	case "-grayscale":
		//imageMaster.GrayScale()
	case "-smoothen":
		//imageMaster.GrayScale()
	case "-sharpen":
		//imageMaster.GrayScale()
	case "-denoise":

	}

	err := imageMaster.Find("./images/sandwich.jpg", "./images/otuput", 2)

	if err != nil {
		fmt.Print(err)
	}
}

//docker run  --name cont -it --mount type=bind,source="$(pwd)",target=/app/images imasterps
