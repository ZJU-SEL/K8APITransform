package main

import (
	"Fti/Ftitool"
	//	"fmt"
	//	"github.com/fsouza/go-dockerclient"
	//  "io"
	//	"os"
	//  "reflect"
)

func main() {

	err := Ftitool.Wartoimage("mytest-tomcat-log")
	if err != nil {
		panic(err)
	}
	//Ftitool.Cleandir("testtemp")
	//	fmt.Println(str)

	//	image := "testimage-deploy"
	//	Ftitool.Createdir(image)
	//	fmt.Println("using the image")
	//	Ftitool.Cleandir(image)
}
