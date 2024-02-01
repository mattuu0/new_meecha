package main

import (
	"log"
	"meecha/database"

	"meecha/location"
)

func main() {
	database.DBpath = "./test.db"
	database.Init()

	location.Init("pMTpmD3N7qGdY4JSjc1fhBaOZyZXGh1e")

	username := "mattuua"

	token, err := location.GenToken(username)

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(token)
}
