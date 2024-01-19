package main

import (
	"log"
	auth "meecha/Auth"
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

	lresult, err := auth.Login("mattuu", "password")

	if err != nil {
		log.Println(err)
		return
	}

	token_data, err := auth.Valid_Token(lresult.RefreshToken)

	if err != nil {
		log.Println("トークンの検証に失敗しました")
		log.Println(err)
		return
	}

	log.Println("トークン検証成功")

	log.Println(auth.Logout(lresult.RefreshToken))

	token_data, err = auth.Valid_Token(lresult.AccessToken)

	if err != nil {
		log.Println("トークンの検証に失敗しました")
		return
	}

	log.Println(token_data.Userid)
}
