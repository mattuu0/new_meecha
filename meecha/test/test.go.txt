package main

import (
	"log"
	auth "meecha/auth"
	"meecha/database"
	"meecha/friends"

	"sync"

	"time"
)

func main() {
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

	nya, err := friends.Send("q", "b")
	log.Println("1")
	log.Println(nya)
	log.Println(err)

	whaa, err := friends.Send("q", "qa")
	log.Println("2")
	log.Println(whaa)
	log.Println(err)

	three_token, err := friends.Send("|||||", "|||")
	log.Println("3")
	log.Println(three_token)
	log.Println(err)

	fuck, err := friends.Send("qq", "a")
	log.Println("4")
	log.Println(fuck)
	log.Println(err)

	tee, err := friends.Send("a", "qqq")
	log.Println("5")
	log.Println(tee)
	log.Println(err)

	three_token, err = friends.Send("qqqb", "qqq")
	log.Println("6")
	log.Println(three_token)
	log.Println(err)

	forth_token, err := friends.Get_Received("qqq")
	log.Println("7")
	log.Println(forth_token)
	log.Println(err)

	five_token, err := friends.Get_Sent("a")
	log.Println("8")
	log.Println(five_token)
	log.Println(err)
	log.Println("")

	log.Println("test承認")

	aaaa, err1 := friends.Accept(fuck, "a")
	_ = aaaa

	bbbb, err1 := friends.Accept(tee, "qqq")
	_ = err1
	_ = bbbb

	log.Println("test送信取り消し")
	err2 := friends.Delete_Request(nya, "q")
	_ = err2

	log.Println("test拒否")
	log.Println(whaa)
	err3 := friends.Rejection(whaa, "qa")
	log.Println(err3)

	log.Println("testフレンド欄")
	six_token, err := friends.Get_Friends("a")
	log.Println(six_token)
	log.Println(err)

	log.Println("testフレンド消去")
	err4 := friends.Delete_Friend(aaaa, "a")
	log.Println(err4)

	/*
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			adddb()
			wg.Done()
		}()
	}

	wg.Wait()

	//100000回ループ
	*/

	
}

func adddb() {
	dbconn := database.GetDB()

	for i := 0; i < 100000; i++ {
		randid1,_ := auth.Genid()
		randid2,_ := auth.Genid()

		fuid,_ := auth.Genid()

		//フレンドトークンの情報
		Ftoken := database.Friends{
			UID: 		 fuid,					//uuid
			Sender_id:   randid1,        			//送った側のID
			Receiver_id: randid2,		    		//受け取る側のID
			SendTime:    time.Now().Unix(),	    //送った時間
		}

		_ = dbconn.Create(&Ftoken)
		log.Println(i)
	}
}
