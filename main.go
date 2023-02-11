package main

import (
	"fmt"
	"os"

	"github.com/l-pavlova/image-master/app"
)

func main() {

	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) == 0 {
		err := fmt.Errorf("%s", "Pass correct command line arguments to the program\n")
		fmt.Println(err)
		return
	}

	imMaster := app.NewImageMaster()
	//imMaster.GrayScale("C:\\Users\\Lyudmila\\Desktop\\images\\test (2).jpg", "C:\\Users\\Lyudmila\\Desktop\\images\\otuput")
	//imMaster.Smoothen("C:\\Users\\Lyudmila\\Desktop\\images\\test (2).jpg", "C:\\Users\\Lyudmila\\Desktop\\images\\otuput", 4)
	//imMaster.Sharpen("C:\\Users\\Lyudmila\\Desktop\\images\\test (2).jpg", "C:\\Users\\Lyudmila\\Desktop\\images\\otuput", 3)
	imMaster.Find("C:\\Users\\Lyudmila\\Desktop\\images\\laptop.jpg", "C:\\Users\\Lyudmila\\Desktop\\images\\otuput")
}
