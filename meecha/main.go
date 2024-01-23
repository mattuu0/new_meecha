package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"meecha/auth"
	"meecha/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/gin-contrib/cors"

	imgupload "github.com/olahol/go-imageupload"
)

var (
	//デフォルトアイコン
	DefaultIcon string  = "./assets/default_icon.jpg"
	//ユーザアイコンフォルダ
	IconDir string = "./UserIcons"
)

func main() {
	//データベース初期化
	database.DBpath = "./meecha.db"
	database.Init()

	//認証初期化
	auth.Init()

	router := gin.Default()

	//すべてのオリジンを承認
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	//ミドルウェア設定
	auth.Auth_Init(router)
	router.Use(auth.Auth_Middleware())

	//アイコン取得
	router.GET("/geticon/:uid", func(ctx *gin.Context) {
		//ユーザを検索
		result,err := auth.GetUser_ByID(ctx.Param("uid"))

		//エラー処理
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(500)
			return
		}

		//ユーザが見つからない場合
		if (!result.IsFind) {
			//404を返す
			ctx.AbortWithStatus(404)
			return
		}

		//画像のパス
		response_path := filepath.Join(IconDir, fmt.Sprintf("%s.jpg",result.UserData.UID))
		imgbin,err := ioutil.ReadFile(response_path)

		//エラー処理
		if err != nil {
			//サーバエラー
			ctx.AbortWithStatus(500)
			return
		}

		//データ返却
		ctx.Data(200,"image/jpeg",imgbin)
	})
	
	//アイコンを変更するエンドポイント
	router.POST("/upicon", func(ctx *gin.Context) {
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
		savepath := filepath.Join(IconDir, fmt.Sprintf("%s.jpg",Auth_Data.UserId))
		thumb.Save(savepath)
	})

	//ping
	router.POST("/user_info", func(ctx *gin.Context) {
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
	})

	//認証関連のグループ
	authg := router.Group("/auth")

	//ログアウト
	authg.POST("/logout", func(ctx *gin.Context) {
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

		//成功メッセージ
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Logout successful",
		})
	})

	//ログイン
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
			ctx.JSON(403, gin.H{
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
			"message":      "Login successful",
			"RefreshToken": result.RefreshToken,
			"AccessToken":  result.AccessToken,
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

		//ユーザネームとパスワード検証
		if !auth.Validate_Name_Password(login_data.Name,login_data.Password) {
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
		savepath := filepath.Join(IconDir, fmt.Sprintf("%s.jpg",result.UID))
		copyfile(DefaultIcon,savepath)

		//成功レスポンス
		ctx.JSON(200, gin.H{
			"message": "Sign up successful",
			"userid":  result.UID,
		})
	})

	router.Run("127.0.0.1:12222") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

//ファイルをコピーする関数
func copyfile(srcName string,dstName string) error {
	//元ファイルを開く
	src, err := os.Open(srcName)
    if err != nil {
        return err
    }

	//コピー元ファイルを閉じる
    defer src.Close()

	//コピー先ファイルを作成
    dst, err := os.Create(dstName)
    if err != nil {
        return err
    }

	//コピー先ファイルを閉じる
    defer dst.Close()

	//コピー
    _, err = io.Copy(dst, src)
    if  err != nil {
        return err
    }

	return nil
}