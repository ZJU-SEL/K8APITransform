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

	err := fti.Wartoimage("mytest-tomcat-log")
	if err != nil {
		panic(err)
	}

}
