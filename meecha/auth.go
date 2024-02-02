package main

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"meecha/auth"

	"log"

	"net/http"

	"path/filepath"

	"fmt"
)

//アクセストークン更新
func token_refresh(ctx *gin.Context) {
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
		//403を返す
		ctx.AbortWithStatus(403)
		return
	}

	//アクセストークンかどうか
	if !Auth_Data.IsRefresh {
		//アクセストークンの場合エラー
		ctx.AbortWithStatus(400)
		return
	}

	//アクセストークン更新
	atoken,err := auth.Refresh(Auth_Data.UserId)

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(500)
		return
	}

	//トークンを返す
	ctx.JSON(http.StatusOK,gin.H{
		"token" : atoken,
	})
}

//ログアウト
func logout(ctx *gin.Context) {
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
		//403を返す
		ctx.AbortWithStatus(403)
		return
	}

	//リフレッシュトークンかどうか
	if !Auth_Data.IsRefresh {
		//アクセストークンの場合エラー
		ctx.AbortWithStatus(403)
		return
	}

	//ログアウト処理
	if err := auth.Logout(Auth_Data.Token); err != nil {
		log.Println(err)
		//ログアウトに失敗したときエラーを返す
		ctx.AbortWithStatus(500)
		return
	}

	//websocketを切断
	ws_disconnect(Auth_Data.UserId)
	//成功メッセージ
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

//ログイン
func login(ctx *gin.Context) {
	//データを受け取る
	var login_data LoginData

	//データを紐付ける
	err := ctx.BindJSON(&login_data)

	//エラー処理
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Bad Request",
		})
		return
	}

	//ログインを試行
	result, err := auth.Login(login_data.Name, login_data.Password)

	//見つからない場合
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(403, gin.H{
			"message": "Incorrect username or password",
		})

		return
	}

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.JSON(500, gin.H{
			"message": "Server Error",
		})
		return
	}

	//失敗した場合
	if !result.Success {
		//失敗レスポンス
		ctx.JSON(400, gin.H{
			"message":      "Login failed",
			"RefreshToken": "",
			"AccessToken":  "",
		})
		return
	}

	//成功レスポンス
	ctx.JSON(200, gin.H{
		"message":      "Login successful",
		"RefreshToken": result.RefreshToken,
		"AccessToken":  result.AccessToken,
	})
}

//ユーザ登録
func signup(ctx *gin.Context) {
	//データを受け取る
	var login_data LoginData

	//データを紐付ける
	err := ctx.BindJSON(&login_data)

	//エラー処理
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Bad Request",
		})
		return
	}

	//ユーザネームとパスワード検証
	if !auth.Validate_Name_Password(login_data.Name, login_data.Password) {
		ctx.JSON(400, gin.H{
			"message": "Bad Request",
		})
		return
	}

	//すでに存在するか確認
	fresult, _ := auth.GetUser_ByName(login_data.Name)

	//見つかった場合
	if fresult.IsFind {
		ctx.JSON(409, gin.H{
			"message": "user already exists",
		})
		return
	}

	//さいんあっぷを試行
	result, err := auth.CreateUser(login_data.Name, login_data.Password)

	//エラー処理
	if err != nil {
		//エラーを返す
		ctx.JSON(500, gin.H{
			"message": "Sign up failed",
		})
		return
	}

	//デフォルトアイコン生成
	savepath := filepath.Join(IconDir, fmt.Sprintf("%s.jpg", result.UID))
	copyfile(DefaultIcon, savepath)

	//成功レスポンス
	ctx.JSON(200, gin.H{
		"message": "Sign up successful",
		"userid":  result.UID,
	})
}