package main

import (
	"errors"
	"net/http"

	"meecha/auth"
	"meecha/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	//データベース初期化
	database.DBpath = "./test.db"
	database.Init()

	//認証初期化
	auth.Init()

	router := gin.Default()
	
	//ミドルウェア設定
	auth.Auth_Init(router)
	router.Use(auth.Auth_Middleware())

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	//認証関連のグループ
	authg := router.Group("/auth")
	authg.POST("/login", func(ctx *gin.Context) {
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
			ctx.JSON(404, gin.H{
				"message": "Incorrect username or password",
			})

			return
		}

		//エラー処理
		if err != nil {
			return
		}

		//成功レスポンス
		ctx.JSON(200, gin.H{
			"message": "Login successful",
			"RefreshToken": result.RefreshToken,
			"AccessToken" : result.AccessToken,
		})
	})

	//サインアップ
	authg.POST("/signup", func(ctx *gin.Context) {
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

		//すでに存在するか確認
		fresult, _ := auth.GetUser_ByName(login_data.Name)

		//見つかった場合
		if fresult.IsFind {
			ctx.JSON(400, gin.H{
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

		//成功レスポンス
		ctx.JSON(200, gin.H{
			"message": "Sign up successful",
			"userid":  result.UID,
		})
	})

	router.Run("127.0.0.1:12222") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
