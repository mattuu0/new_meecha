package main

import (
	"log"
	"meecha/auth"
	"meecha/database"
)

func main() {
	database.DBpath = "./aaa.db"
	database.Init()

	auth.Init()

	fresult, _ := auth.GetUser_ByName("mattuu")

	if !fresult.IsFind {
		_, err := auth.CreateUser("mattuu", "password")
		log.Println(err)
	}	

	result, _ := auth.GetUser_ByName("mattuu")

	if !result.IsFind {
		return
	}

	log.Println(result.UserData.UID)

}
