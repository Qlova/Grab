package main

import (
	"net"
	"fmt"
	"os"
)

//Need to start up the server and listen for connections.
//We also need to open the framebuffer.
func main() {

	listener, err := net.Listen("tcp", ":222")
	if err != nil {
		fmt.Println("Cannot open socket, please run as root!", err)
		os.Exit(1)
	}
	
	go Events()
	go Renderer()
	
	fmt.Println("[GRAB] Running!")
	for {
        connection, err := listener.Accept()
        if err != nil {
        	fmt.Println("Error accepting client: "+err.Error())
        	continue
        }

        go Serve(connection)
    }
	
}
