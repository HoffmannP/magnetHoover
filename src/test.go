package main

import (
	"log"
	"transmission"
)

func main() {
	c, err := transmission.NewClient(false, "127.0.0.1", 9091)
	if err != nil {
		log.Fatalln(err)
	}
	err = c.Add("http://cdimage.debian.org/debian-cd/6.0.6/amd64/bt-cd/debian-6.0.6-amd64-CD-1.iso.torrent")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Done.")
}
