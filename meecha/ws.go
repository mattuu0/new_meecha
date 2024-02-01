package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"

	"meecha/location"
)

// コマンドメッセージ
type CommandMessage struct {
	//コマンド
	Command string

	//ペイロード
	Payload interface{}
}

// レスポンスメッセージ
type ResponseMessage struct {
	//コマンド
	Command string

	//ペイロード
	Payload interface{}
}

// Websocket 関数
func handle_ws(wsconn *websocket.Conn, userid string) {
	// TODO
	defer func () {
		if err := recover(); err != nil {
			log.Println(err)
		}

		//Websocket接続削除
		delete(wsconns,userid)

		//接続を閉じる
		wsconn.Close()	
	}()

	//ユーザIDを返す
	write_data := ResponseMessage{
		Command: "Auth_Complete",
		Payload: userid,
	}

	//書き込み
	err := wsconn.WriteJSON(write_data)

	//エラー処理
	if err != nil {
		log.Println(err)
		return
	}

	//別スレッドで開始
	go send_location_token(wsconn, userid)

	for {
		//メッセージ
		readmsg := &CommandMessage{}

		//メッセージ受信
		err := wsconn.ReadJSON(readmsg)

		//エラー処理
		if err != nil {
			log.Println(err)
			break
		}

		//コマンド処理
		switch (readmsg.Command) {
			case "location":
				payload := readmsg.Payload.(map[string]interface{})

				//位置情報
				token_userid,err := location.VerifyToken(payload["token"].(string))

				//エラー処理
				if err != nil {
					log.Println(err)
					break
				}

				//ユーザID比較
				if userid != token_userid {
					log.Println("userid error")
					break
				}

				log.Println(payload["lat"].(float64))
				log.Println(payload["lng"].(float64))
		}
	}
}

// 位置情報用トークンを送る
func send_location_token(wsconn *websocket.Conn, userid string) {
	for {
		//トークン生成
		token, err := location.GenToken(userid)

		//エラー処理
		if err != nil {
			log.Println(err)
			break
		}

		//トークンメッセージ
		rmsg := ResponseMessage{
			Command: "Location_Token",
			Payload: token,
		}

		//データを送る
		wresult := wsconn.WriteJSON(rmsg)

		//エラー処理
		if wresult != nil {
			log.Println(wresult)
			break
		}

		//5秒待機
		time.Sleep(time.Second * 3)
	}
}
