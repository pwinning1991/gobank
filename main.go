package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Yeah Buddy!")
	store, er := NewPostgresStore()
	if er != nil {
		log.Fatal(er)
	}

	fmt.Printf("%v\n", store)

	//server := NewApiServer(":3000", store)
	//server.Run()
}
