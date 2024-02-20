package main

import (
	"log"
	auth "meecha/auth"
	"meecha/database"
	"meecha/friends"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	loadEnv()
	database.DBpath = "./test.db"
	database.Init()

	auth.Init()
	friends.Init()

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


	fuck, err := friends.Send("qq", "a")
	log.Println("1")
	log.Println(fuck)
	log.Println(err)

	aaaa,send1,err1 := friends.Accept(fuck,"a")
	_=err1
	_=aaaa
	_=send1

	
	abc, err := friends.Send("a", "qq")
	log.Println("3")
	log.Println(abc)
	log.Println(err)


	abcd, err := friends.Send("qq", "a")
	log.Println("4")
	log.Println(abcd)
	log.Println(err)
}

// log.Println("test送信取り消し")
// err2 := friends.Delete_Request(nya,"q")
// _=err2

// log.Println("test拒否")
// log.Println(whaa)
// err3 := friends.Rejection(whaa,"qa")
// log.Println(err3)

// log.Println("testフレンド欄")
// six_token,err := friends.Get_Friends("a")
// log.Println(six_token)
// log.Println(err)

// log.Println("testフレンド消去")
// id,err4 := friends.Delete_Friend(aaaa,"a")
// log.Println(err4)
// _=id
