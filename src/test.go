package main

import (
	"log"
	"transmission"
)

func main() {
	/*
		res, err := http.Get("http://www.google.de/robots.txt")
		if err != nil {
			log.Fatal(err)
		}
		robots, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", robots)
	*/
	c, err := transmission.NewClient(false, "localhost", 9091)
	if err != nil {
		log.Fatalln(err)
	}
	err = c.Add("http://cdimage.debian.org/debian-cd/6.0.6/amd64/bt-cd/debian-6.0.6-amd64-CD-1.iso.torrent")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Done.")
}
