package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"meecha/friends"
	"meecha/location"

	"meecha/auth"
)

func getid(ctx *gin.Context) (string, error) {
	//認証情報を取得
	result, exits := ctx.Get(auth.KeyName)

	//設定されていないとき戻る
	if !exits {
		//403を返す
		ctx.AbortWithStatus(403)
		return "", fmt.Errorf("error")
	}

	//型を変換
	Auth_Data := result.(auth.Auth_Result)

	//認証に失敗してるとき戻る
	if !Auth_Data.Success {
		//エラーのHTMLを返す
		ctx.AbortWithStatus(403)
		return "", fmt.Errorf("error")
	}

	//IDを取得
	return Auth_Data.UserId, nil
}

// フレンド一覧取得
func get_friends(ctx *gin.Context) {
	//認証情報を取得
	uid, err := getid(ctx)

	//エラー処理
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	//フレンド一覧取得
	result, err := friends.Get_Friends(uid)

	//エラー処理
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	ctx.JSON(200, result)
}

// ユーザ検索
// func search_user(ctx *gin.Context) {
// 	//ユーザ名
// 	user_name := ctx.GetHeader("username")

// 	log.Println(user_name)
// 	//ユーザを検索
// 	result, err := auth.GetUser_ByName(user_name)

// 	//エラー処理
// 	if err != nil {
// 		log.Println(err)
// 		ctx.AbortWithStatus(500)
// 		return
// 	}

// 	//ユーザが見つからない場合
// 	if !result.IsFind {
// 		ctx.AbortWithStatus(404)
// 		return
// 	}

// 	result_data := map[string]string{}
// 	result_data["uid"] = result.UserData.UID
// 	result_data["name"] = result.UserData.Name

// 	//データ返却
// 	ctx.JSON(200, result_data)
// }

type SearchData struct {
	UserName string
}

func search_user(ctx *gin.Context) {
	var data SearchData

	// 値をバインドする
	if err := ctx.ShouldBind(&data); err != nil {
		ctx.AbortWithStatus(400)
	}

	//ユーザ名
	// user_name := ctx.GetHeader("username")
	user_name := data.UserName

	log.Println(user_name)
	//ユーザを検索
	result, err := auth.GetUser_ByName(user_name)

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(500)
		return
	}

	//ユーザが見つからない場合
	if !result.IsFind {
		ctx.AbortWithStatus(404)
		return
	}

	result_data := map[string]string{}
	result_data["uid"] = result.UserData.UID
	result_data["name"] = result.UserData.Name

	//データ返却
	ctx.JSON(200, result_data)
}

type SendRequestData struct {
	Targetid string //送信先
}

// リクエスト送信
func send_request(ctx *gin.Context) {
	//送信情報を取得
	var request_data SendRequestData

	//データを紐付ける
	err := ctx.BindJSON(&request_data)

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(500)
		return
	}

	//ユーザID取得
	uid, err := getid(ctx)

	//エラー処理
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	//IDが同じかを比較する
	if uid == request_data.Targetid {
		ctx.AbortWithStatus(400)
		return
	}

	//リクエスト送信
	requestid, err := friends.Send(uid, request_data.Targetid)

	//エラー処理
	if err != nil {
		log.Println(err)
		//重複エラー
		if err.Error() == "request_is_already_existing" {
			ctx.AbortWithStatus(409)
			return
		}

		ctx.AbortWithStatus(500)
		return
	}

	//送信者取得
	sender,err := auth.GetUser_ByID(uid)

	//エラー処理
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	//送信者名
	sender_name := sender.UserData.Name

	//通知を飛ばす
	err = Send_ws(request_data.Targetid, "recv_request",sender_name)

	//エラー処理
	if err != nil {
		log.Println(err)
	}

	//成功
	ctx.JSON(200, map[string]string{"requestid": requestid})
}

type RequestData struct {
	Requestid string //リクエストID
}

// リクエストキャンセル
func cancel_request(ctx *gin.Context) {
	//送信情報を取得
	var request_data RequestData

	//データを紐付ける
	err := ctx.BindJSON(&request_data)

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(500)
		return
	}

	//ユーザID取得
	uid, err := getid(ctx)

	//エラー処理
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	//リクエストキャンセル
	err = friends.Delete_Request(request_data.Requestid,uid)

	//エラー処理
	if err != nil {
		//エラー
		ctx.AbortWithStatus(500)
		return
	}

	//成功
	ctx.JSON(200, nil)
}

// リクエスト拒否
func reject_request(ctx *gin.Context) {
	//送信情報を取得
	var request_data RequestData

	//データを紐付ける
	err := ctx.BindJSON(&request_data)

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(500)
		return
	}

	//ユーザID取得
	uid, err := getid(ctx)

	//エラー処理
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	//リクエスト拒否
	err = friends.Rejection(request_data.Requestid,uid)

	//エラー処理
	if err != nil {
		//エラー
		ctx.AbortWithStatus(500)
		return
	}

	//成功
	ctx.JSON(200, nil)
}

// リクエスト承認
func accept_request(ctx *gin.Context) {
	//送信情報を取得
	var request_data RequestData

	//データを紐付ける
	err := ctx.BindJSON(&request_data)

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(500)
		return
	}

	//ユーザID取得
	uid, err := getid(ctx)

	//エラー処理
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	log.Println(request_data)
	//リクエスト承認
	fid,sid,err := friends.Accept(request_data.Requestid,uid)

	//エラー処理
	if err != nil {
		if (err.Error() == "already_friend_registered") {
			ctx.AbortWithStatus(409)
			return
		}
		//エラー
		ctx.AbortWithStatus(500)
		return
	}

	//承認者取得
	accepter,err := auth.GetUser_ByID(uid)

	//エラー処理
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	//送信者名
	accepter_name := accepter.UserData.Name

	//通知を飛ばす
	Send_ws(sid, "accept_request",accepter_name)

	//自信のキャッシュ更新
	location.Update_Cache(uid)

	//相手のキャッシュ更新
	location.Update_Cache(sid)

	//成功
	ctx.JSON(200, gin.H{
		"friendid": fid,
	})
}

// 送信済み取得
func get_sent_request(ctx *gin.Context) {
	//認証情報を取得
	uid, err := getid(ctx)

	//エラー処理
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	//送信済み取得
	result, err := friends.Get_Sent(uid)

	//エラー処理
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	ctx.JSON(200, result)
}

// 受信済み取得
func get_recved_request(ctx *gin.Context) {
	//認証情報を取得
	uid, err := getid(ctx)

	//エラー処理
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	//受信済み取得
	result, err := friends.Get_Received(uid)

	//エラー処理
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	ctx.JSON(200, result)
}

// フレンドデータ
type FriendData struct {
	Friendid string //フレンドID
}

//フレンド削除
func remove_friend(ctx *gin.Context) {
	//認証情報を取得
	uid, err := getid(ctx)

	//エラー処理
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	//送信情報を取得
	var friend_data FriendData

	//データを紐付ける
	err = ctx.BindJSON(&friend_data)

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(500)
		return
	}

	//フレンド削除 (相手のID取得)
	aiteid,err := friends.Delete_Friend(friend_data.Friendid,uid)

	//エラー処理
	if err != nil {
		//エラー
		ctx.AbortWithStatus(500)
		return
	}

	//キャッシュ更新
	location.Update_Cache(uid)

	//相手のキャッシュ更新
	location.Update_Cache(aiteid)

	//成功
	ctx.JSON(200, nil)
}