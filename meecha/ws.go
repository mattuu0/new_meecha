package main

import (
	"errors"
	"log"
	"time"

	"github.com/gorilla/websocket"

	"meecha/auth"
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

func Send_ws(uid string, command string, payload interface{}) error {
	//接続されていなかったらエラーを返す
	if wsconns[uid] == nil {
		return errors.New("not connected")
	}

	//送信データ
	write_data := ResponseMessage{
		Command: command,
		Payload: payload,
	}

	//書き込み
	err := wsconns[uid].WriteJSON(write_data)

	//エラー処理
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Websocket切断
func ws_disconnect(uid string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	location.Disable_Geo_Token(uid)
	//Websocket接続を閉じる
	wsconns[uid].Close()
	//Websocket接続削除
	delete(wsconns, uid)
}

// Websocket 関数
func handle_ws(wsconn *websocket.Conn, userid string) {
	// TODO
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}

		//トークン無効化
		location.Disable_Geo_Token(userid)

		//位置情報削除
		location.RemoveLocation(userid)

		//Websocket接続削除
		delete(wsconns, userid)

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

	//ステータス更新
	//長さ取得
	distance, err := location.Get_Notify_distance(userid)

	//エラー処理
	if err != nil {
		log.Println(err)
	}

	//ステータス更新
	location.Update_data(userid, distance, "online")

	//別スレッドで開始
	go send_location_token(wsconn, userid)

	//通知者一覧
	notify_dict := map[string]string{}

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
		switch readmsg.Command {
		case "location":
			payload := readmsg.Payload.(map[string]interface{})

			//位置情報
			token_userid, err := location.VerifyToken(payload["token"].(string))

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

			//位置情報更新
			err = location.Update_Geo(userid, payload["lat"].(float64), payload["lng"].(float64))

			//エラー処理
			if err != nil {
				log.Println(err)
				continue
			}

			//既存のトークン無効か
			err = location.Disable_Geo_Token(userid)

			//エラー処理
			if err != nil {
				log.Println(err)
				continue
			}

			//自分の位置情報
			point_data := location.Point{Lat: payload["lat"].(float64), Lon: payload["lng"].(float64)}
			//位置情報検証
			result, err := location.Validate_Geo(userid, point_data)

			//エラー処理
			if err != nil {
				for key := range notify_dict {
					if _, ok := notify_dict[key]; ok {
						stop_notify(notify_dict, key, userid)
					}
				}

				log.Println(err)
				continue
			}

			//除外エリアに入っている場合
			if !result {
				//全フレンドに通知
				for key := range notify_dict {
					if _, ok := notify_dict[key]; ok {
						stop_notify(notify_dict, key, userid)
					}
				}

				log.Println("除外エリアです")
				//戻る
				continue
			}

			//フレンド取得
			friends, err := location.Get_Friends(userid)

			//エラー処理
			if err != nil {
				log.Println(err)
				continue
			}

			//自信のデータ取得
			mydata, err := auth.GetUser_ByID(userid)

			//エラー処理
			if err != nil {
				log.Println(err)
				continue
			}

			for _, val := range friends {

				//フレンドの位置を取得
				friend_point_data, err := location.GetLocation(val)

				//エラー処理
				if err != nil {
					//離れたことを通知
					if _, ok := notify_dict[val]; ok {
						stop_notify(notify_dict, val, userid)
					}
					log.Println(err)
					continue
				}

				//フレンドが除外範囲にいるか
				result, err := location.Validate_Geo(val, friend_point_data)

				//エラー処理
				if err != nil {
					//離れたことを通知
					if _, ok := notify_dict[val]; ok {
						stop_notify(notify_dict, val, userid)
					}
					log.Println(err)
					continue
				}

				//除外範囲にいる場合
				if !result {
					//離れたことを通知
					if _, ok := notify_dict[val]; ok {
						stop_notify(notify_dict, val, userid)
					}
					continue
				}

				//自分の設定距離取得
				jibunn_distance, err := location.Get_Notify_distance(userid)

				//エラー処理
				if err != nil {
					//離れたことを通知
					if _, ok := notify_dict[val]; ok {
						stop_notify(notify_dict, val, userid)
					}
					log.Println(err)
					continue
				}

				//相手の設定距離取得
				aite_distance, err := location.Get_Notify_distance(val)

				//エラー処理
				if err != nil {
					//離れたことを通知
					if _, ok := notify_dict[val]; ok {
						stop_notify(notify_dict, val, userid)
					}
					log.Println(err)
					continue
				}

				//通知距離を自分の距離に設定
				check_distance := jibunn_distance

				//自分よりあいてのほうが短い場合
				if aite_distance < jibunn_distance {
					//相手の距離にする
					check_distance = aite_distance
				}

				//距離取得
				distance := location.Get_Distance(point_data, friend_point_data)

				//距離検証
				if distance > check_distance {
					//離れたことを通知
					if _, ok := notify_dict[val]; ok {
						stop_notify(notify_dict, val, userid)
					}
					//距離より大きい場合
					continue
				}
				
				log.Println(check_distance)
				log.Println(distance)

				log.Println(check_distance)
				log.Println(distance)

				//初回かどうか
				is_first := false

				//初回か判定
				if _, ok := notify_dict[val]; !ok {
					is_first = true
				}

				//相手に通知
				Send_ws(val, "near_friend", map[string]interface{}{
					"userid":   userid,
					"unane":    mydata.UserData.Name,
					"is_first": is_first,
					"is_self":  false,
					"point":    point_data,
				})

				//相手のデータ取得
				aite_data, err := auth.GetUser_ByID(val)

				//エラー処理
				if err != nil {
					//離れたことを通知
					if _, ok := notify_dict[val]; ok {
						stop_notify(notify_dict, val, userid)
					}
					log.Println(err)
					continue
				}

				//自分に通知
				Send_ws(userid, "near_friend", map[string]interface{}{
					"userid":   val,
					"unane":    aite_data.UserData.Name,
					"is_first": is_first,
					"is_self":  true,
					"point":    friend_point_data,
				})

				//通知者に設定
				notify_dict[val] = ""
			}
		}
	}

	//ステータス更新
	//長さ取得
	distance, err = location.Get_Notify_distance(userid)

	//エラー処理
	if err != nil {
		log.Println(err)
	}

	//ステータス更新
	location.Update_data(userid, distance, "offline")

	//トークン無効化
	location.Disable_Geo_Token(userid)

	//位置情報削除
	location.RemoveLocation(userid)

	//切断通知
	for key := range notify_dict {
		if _, ok := notify_dict[key]; ok {
			stop_notify(notify_dict, key, userid)
		}
	}
}

func stop_notify(notify_dict map[string]string, uid string, myid string) {
	//通知削除
	delete(notify_dict, uid)

	Send_ws(uid, "stop_notify", map[string]string{
		"userid": myid,
	})
}

// 位置情報用トークンを送る
func send_location_token(wsconn *websocket.Conn, userid string) {
	// TODO
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}

	}()

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
		time.Sleep(location.TokenExp)
	}
}
