package main

import (
	"log"
	"meecha/auth"
	"meecha/database"
)

func main() {
	log.Print()

	database.DBpath = "./aaa.db"
	database.Init()

	auth.Init()

	_, err := auth.CreateUser("mattuu", "password")
	log.Println(err)

	auth.GetUser_ByName("mattuu")
}
