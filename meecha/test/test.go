package main

import (
	"log"
	"meecha/Auth"
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

	lresult,err := auth.Login("mattuu","password")

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(lresult)
}
