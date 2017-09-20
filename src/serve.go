package main

import "net"
import "fmt"
import "strings"

func Serve(connection net.Conn) {
	var buffer = make([]byte, 256)
	for {
		count, err := connection.Read(buffer)
		if err != nil || count == 0 {
			fmt.Println("Could not read from client! :O")
			return
		}
		
		request := string(buffer[:count])
		split := strings.SplitN(request, " ", 2)
		
		command := split[0]
		
		split = strings.SplitN(split[1], ";", 2)
		option := split[0]
		data := split[1]
		
		switch command {
			case "GRAB":
				switch option {
					case "button":
						fmt.Println("Button request!", data)
				}
		}	
	}
}
