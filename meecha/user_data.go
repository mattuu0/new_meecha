package main

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"

	"meecha/auth"

	"log"

	"fmt"

	"path/filepath"

	imgupload "github.com/olahol/go-imageupload"

	"net/http"
)

func geticon(ctx *gin.Context) {
	//ユーザを検索
	result, err := auth.GetUser_ByID(ctx.Param("uid"))

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(500)
		return
	}

	//ユーザが見つからない場合
	if !result.IsFind {
		//404を返す
		ctx.AbortWithStatus(404)
		return
	}

	//画像のパス
	response_path := filepath.Join(IconDir, fmt.Sprintf("%s.jpg", result.UserData.UID))
	imgbin, err := ioutil.ReadFile(response_path)

	//エラー処理
	if err != nil {
		//サーバエラー
		ctx.AbortWithStatus(500)
		return
	}

	//データ返却
	ctx.Data(200, "image/jpeg", imgbin)
}

//画像アップロード
func uploadimg(ctx *gin.Context) {
	//認証情報を取得
	result, exits := ctx.Get(auth.KeyName)

	//設定されていないとき戻る
	if !exits {
		//403を返す
		ctx.AbortWithStatus(403)
		return
	}

	//型を変換
	Auth_Data := result.(auth.Auth_Result)

	//認証に失敗してるとき戻る
	if !Auth_Data.Success {
		//エラーのHTMLを返す
		ctx.AbortWithStatus(403)
		return
	}

	//画像を受け取る
	img, err := imgupload.Process(ctx.Request, "file")
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(500)
		return
	}

	//アイコンをリサイズ
	thumb, err := imgupload.ThumbnailJPEG(img, 300, 300, 50)
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(500)
		return
	}

	//保存するパス
	savepath := filepath.Join(IconDir, fmt.Sprintf("%s.jpg", Auth_Data.UserId))
	thumb.Save(savepath)
}

//ユーザ情報取得
func get_user_info(ctx *gin.Context) {
	//認証情報を取得
	result, exits := ctx.Get(auth.KeyName)

	//設定されていないとき戻る
	if !exits {
		//403を返す
		ctx.AbortWithStatus(403)
		return
	}

	//型を変換
	Auth_Data := result.(auth.Auth_Result)

	//認証に失敗してるとき戻る
	if !Auth_Data.Success {
		//エラーのHTMLを返す
		ctx.AbortWithStatus(403)
		return
	}

	//ユーザデータ取得
	uresult, err := auth.GetUser_ByID(Auth_Data.UserId)

	//エラー処理
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	//ユーザが見つからないとき
	if !uresult.IsFind {
		ctx.AbortWithStatus(404)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"userid": Auth_Data.UserId,
		"name":   uresult.UserData.Name,
	})
}