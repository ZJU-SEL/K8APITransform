package main

import (
	"github.com/fsouza/go-dockerclient"
)

func main() {

	pushopts := docker.PushImageOptions{
		Name:         imageprefix + "/" + newimage,
		Tag:          "latest",
		Registry:     imageprefix,
		OutputStream: os.Stdout,
	}

	fmt.Println("Name:", newimage)

	fmt.Println("Registry:", imageprefix)

	err = client.PushImage(pushopts, auth)
	if err != nil {
		return "", err
	}
}
