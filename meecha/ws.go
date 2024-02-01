package main

import (
	"log"

	"github.com/gorilla/websocket"
)

//コマンドメッセージ
type CommandMessage struct {
	//コマンド
	Command string

	//ペイロード
	Payload interface{}
}

//レスポンスメッセージ
type ResponseMessage struct {
	//コマンド
	Command string

	//ペイロード
	Payload interface{}
}

// Websocket 関数
func handle_ws(wsconn *websocket.Conn,userid string) {
	// TODO
	defer wsconn.Close()

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

		log.Println(readmsg)
	}
}

//位置情報用トークンを送る
func send_location_token(wsconn *websocket.Conn, userid string) {
	
}