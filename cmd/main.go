package main

import "github.com/ashtonx86/mocker/internal/server"

func main() {
	server := server.NewWebServer()
	server.App.Listen(":3000")
}