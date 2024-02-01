package main

import (
	"errors"
	"io"

	"log"

	"os"
	"path/filepath"
	"time"

	"meecha/auth"
	"meecha/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/gin-contrib/cors"

	"github.com/gorilla/websocket"
)


var (
	//デフォルトアイコン
	DefaultIcon string = "./assets/default_icon.jpg"
	//ユーザアイコンフォルダ
	IconDir string = "./UserIcons"
	
	//ウェブソケット
	wsconns = make(map[string]*websocket.Conn)
)

func getFileNameWithoutExt(path string) string {
    // Fixed with a nice method given by mattn-san
    return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

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

	//フォルダ開く
	Icons,err := os.Open(IconDir)

	//エラー処理
	if err != nil {
		log.Println("エラーです")
		return
	}

	//アイコンファイル
	IconFiles,err := Icons.ReadDir(0)

	//エラー処理
	if err != nil {
		log.Println("エラーです")
		return
	}

	//残っているユーザアイコンを消す
	for _,val := range IconFiles {
		//ユーザを取得する
		_,err := auth.GetUser_ByID(getFileNameWithoutExt(val.Name()))

		//エラー処理
		//見つからないとき
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//削除する
			//エラー処理
			if err := os.Remove(filepath.Join(IconDir,val.Name())); err != nil {
				log.Println(err)
				continue
			}
			continue
		}

		//エラー処理
		if err != nil {
			log.Println(err)
			continue
		}
	}

	router.Use(auth.Auth_Middleware())

	//アイコン取得
	router.GET("/geticon/:uid",geticon)

	//アイコンを変更するエンドポイント
	router.POST("/upicon",uploadimg)

	//ping
	router.POST("/user_info", get_user_info)

	//認証関連のグループ
	authg := router.Group("/auth")

	//アクセストークンリフレッシュ
	authg.POST("/refresh",token_refresh)

	//ログアウト
	authg.POST("/logout", logout)

	//ログイン
	authg.POST("/login", login)

	//サインアップ
	authg.POST("/signup", signup)

	router.GET("/ws",func(ctx *gin.Context) {
		//Websocket接続
		wsconn,err := wsupgrader.Upgrade(ctx.Writer,ctx.Request,nil)

		//エラー処理
		if err != nil {
			log.Println(err)
			return
		}

		//コンティニュー
		ws_contiune := false

		//ユーザID
		userid := ""

		//無限ループ
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

			//認証
			if readmsg.Command == "auth" {
				//認証
				auth_result,err := auth.Valid_Token(readmsg.Payload.(string))

				//エラー処理
				if err != nil {
					log.Println(err)
					//切断
					break
				}

				//リフレッシュトークンの場合閉じる
				if (auth_result.IsRefresh) {
					break
				}
				
				//ユーザIDを設定
				userid = auth_result.Userid

				//接続を追加
				wsconns[auth_result.Userid] = wsconn
				ws_contiune = true

				//抜ける
				break
			}
		}

		//コンティニュー
		if (!ws_contiune) {
			//切断
			wsconn.Close()

			log.Println("切断")
			return
		}

		//スレッド作成
		go handle_ws(wsconn,userid)
	})

	router.Run("127.0.0.1:12222") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// ファイルをコピーする関数
func copyfile(srcName string, dstName string) error {
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
	if err != nil {
		return err
	}

	return nil
}
