package main

import (
	"log"
	"meecha/database"

	"meecha/location"
)

func main() {
	database.DBpath = "./test.db"
	database.Init()

	username := "mattuua"

	token, err := location.GenToken(username)

	if err != nil {
		log.Println(err)
		return
	}

	_ = token

	token2, err := location.GenToken(username)

	if err != nil {
		log.Println(err)
		return
	}

	//トークン検証
	userid, err := location.VerifyToken(token2)

	if err != nil {
		log.Println(err)
		return
	}

	//トークン取得
	log.Println(userid)
}
