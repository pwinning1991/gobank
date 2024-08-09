package main

import "fmt"

func main() {
	fmt.Println("Yeah Buddy!")
	server := NewApiServer("3000")
	server.Run()
}
