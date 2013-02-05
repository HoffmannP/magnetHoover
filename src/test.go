package main

import (
	"fmt"
	"transmission"
)

func main() {
	c, err := transmission.NewClient(false, "localhost", 9091)
	if err != nil {
		fmt.Println(err)
	}
	err = c.Add("http://cdimage.debian.org/debian-cd/6.0.6/amd64/bt-cd/debian-6.0.6-amd64-CD-1.iso.torrent")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Done.")
}
