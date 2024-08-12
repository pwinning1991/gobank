package main

import (
	"flag"
	"fmt"
	"log"
)

func seedAccount(store Storage, fname, lname, pw string) *Account {
	acc, err := NewAccount(fname, lname, pw)
	if err != nil {
		log.Fatal(err)
	}

	if err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	return acc
}

func seedAccounts(store Storage) {
	seedAccount(store, "phil", "w", "hunter8888")

}

func main() {
	seed := flag.Bool("seed", false, "seed the db")
	flag.Parse()
	fmt.Println("Yeah Buddy!")
	store, er := NewPostgresStore()
	if er != nil {
		log.Fatal(er)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("Seeding db")
		seedAccounts(store)
	}

	server := NewApiServer(":3000", store)
	server.Run()
}
