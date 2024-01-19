package main

import (
	"log"
	auth "meecha/Auth"
	"meecha/database"
)

func main() {
	database.DBpath = "./test.db"
	database.Init()

	auth.Init()

	fresult, _ := auth.GetUser_ByName("mattuua")

	if !fresult.IsFind {
		_, err := auth.CreateUser("mattuua", "password")
		log.Println(err)
	}

	one_token, err := auth.Login("mattuua", "password")
	//log.Println(one_token)
	_ = one_token
	//log.Println(err)

	second_token, err := auth.Login("mattuu", "password")
	//log.Println(second_token)
	_ = second_token
	//log.Println(err)

	_ = err
}
