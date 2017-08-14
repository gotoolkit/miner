package db

import (
	"log"

	"github.com/boltdb/bolt"
)

var store *bolt.DB

func Close() {
	err := store.Close()
	if err != nil {
		log.Println(err)
	}
}
func Init() {
	if store != nil {
		return
	}
	var err error
	store, err = bolt.Open("system.db", 0666, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
