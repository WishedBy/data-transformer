package main

import "flag"

func main() {
	server := server.Server{
		Port: flag.Int("http-port", 80, "the port the api will run on"),
	}
	server.Run()
}
