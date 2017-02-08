package main

import "github.com/dpolansky/ci/server"

func main() {
	serv := server.New()
	serv.Serve()
}
