package main

import (
	"bufio"
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

	handleCommand(command, imageMaster)
	for {
		consoleReader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter new operation:")

		input, _ := consoleReader.ReadString('\n')

		input = strings.ToLower(input)
		input = strings.TrimSpace(input)

		handleCommand(input, imageMaster)
		if strings.HasPrefix(input, "bye") {
			fmt.Println("exiting now.")
			os.Exit(0)
		}
	}

}

// function that executes the corresponding operation to the command
func handleCommand(command string, imageMaster *app.ImageMaster) {
	if strings.Contains(command, "-find=") {
		fmt.Println("performing search operation")
		objectToFind := strings.Split(command, "=")[1]
		err := imageMaster.Find(objectToFind)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	switch command {
	case "-help":
		imageMaster.ShowHelp()
		break
	case "-grayscale":
		fmt.Println("performing grayscale operation")
		imageMaster.GrayScale()
		break
	case "-smoothen":
		fmt.Println("performing smoothen operation")
		imageMaster.Smoothen(10)
	case "-sharpen":
		fmt.Println("performing sharpen operation")
		imageMaster.Sharpen()
		break
	case "-list-result-images":
		fmt.Println("here's the images we have data about: ")
		imageMaster.ListAll()
	}
}

//docker run  --name cont -it --mount type=bind,source="$(pwd)",target=/app/images imasterps
//mockgen -source="D:/go/src/GoCourse/image-master/image-master/app/app.go" -destination="D:/go/src/GoCourse/image-master/image-master/mocks/app.go"
