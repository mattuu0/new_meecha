package main

import (
	"errors"
	"io"
	"strconv"
	"sync"

	"log"

	"os"
	"path/filepath"
	"time"

	"meecha/auth"
	"meecha/database"
	"meecha/friends"
	"meecha/location"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/gin-contrib/cors"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}

var (
	//デフォルトアイコン
	DefaultIcon string = "./assets/default_icon.jpg"
	//ユーザアイコンフォルダ
	IconDir string = "./UserIcons"

	//ウェブソケット
	wsconns = sync.Map{} // make(map[string]*websocket.Conn)

)

func getFileNameWithoutExt(path string) string {
	// Fixed with a nice method given by mattn-san
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	//環境変数読み込み
	loadEnv()

	//データベース初期化
	database.Init()

	//認証初期化
	auth.Init()

	//位置情報初期化
	location.TokenExp = time.Duration(3) * time.Second
	location.Init(os.Getenv("LOCATION_TOKEN"))

	//フレンド初期化
	friends.Init()

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

	rate_limit, err := strconv.Atoi(os.Getenv("RateSec"))

	//エラー処理
	if err != nil {
		log.Println("無効なレート制限です")
		return
	}

	log.Println(os.Getenv("RedisUrl"))
	//レート制限
	store := ratelimit.RedisStore(&ratelimit.RedisOptions{
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     os.Getenv("RedisUrl"),
			Password: os.Getenv("RedisPass"), // no password set
			DB:       7,                      // use default DB
			PoolSize: 1000,
		}),
		Rate:  time.Second,
		Limit: uint(rate_limit),
	})

	mw := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})

	_ = mw
	//router.Use(mw)

	//フォルダ開く
	Icons, err := os.Open(IconDir)

	//エラー処理
	if err != nil {
		log.Println("エラーです")
		return
	}

	//アイコンファイル
	IconFiles, err := Icons.ReadDir(0)

	//エラー処理
	if err != nil {
		log.Println("エラーです")
		return
	}

	//残っているユーザアイコンを消す
	for _, val := range IconFiles {
		//ユーザを取得する
		_, err := auth.GetUser_ByID(getFileNameWithoutExt(val.Name()))

		//エラー処理
		//見つからないとき
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//削除する
			//エラー処理
			if err := os.Remove(filepath.Join(IconDir, val.Name())); err != nil {
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

	//ミドルウェア設定
	router.Use(auth.Auth_Middleware())

	//アイコン取得
	router.GET("/geticon/:uid", geticon)

	//アイコンを変更するエンドポイント
	router.POST("/upicon", uploadimg)

	//ping
	router.POST("/user_info", get_user_info)

	//位置情報グループ
	location_group := router.Group("/location")
	location_group.Use(auth.Auth_Middleware())
	location_group.Use(auth.Auth_Require_Middleware())
	//除外設定更新
	location_group.POST("/save_ignore_point", save_ignore_point)

	//除外設定取得
	location_group.POST("/load_ignore_point", get_ignore_point)

	//通知距離取得
	location_group.POST("/get_notify_distance", Get_notify_Distance)

	//通知距離設定
	location_group.POST("/set_notify_distance", Set_notify_Distance)

	//フレンド
	friendg := router.Group("/friend")

	//ミドルウェア設定
	friendg.Use(auth.Auth_Middleware())
	friendg.Use(auth.Auth_Require_Middleware())

	//フレンド一覧取得
	friendg.POST("/getall", get_friends)

	//ユーザ検索
	friendg.POST("/search", search_user)

	//リクエスト送信
	friendg.POST("/request", send_request)

	//リクエスト拒否
	friendg.POST("/reject_request", reject_request)

	//リクエスト承認
	friendg.POST("/accept_request", accept_request)

	//送信済み取得
	friendg.POST("/get_sent", get_sent_request)

	//受信済み取得
	friendg.POST("/get_recved", get_recved_request)

	//フレンド削除
	friendg.POST("/remove_friend", remove_friend)

	//リクエストきゃせる
	friendg.POST("/cancel_request", cancel_request)

	//認証関連のグループ
	authg := router.Group("/auth")

	//アクセストークンリフレッシュ
	authg.POST("/refresh", token_refresh)

	//ログアウト
	authg.POST("/logout", logout)

	//ログイン
	authg.POST("/login", login)

	//サインアップ
	authg.POST("/signup", signup)

	router.GET("/ws", func(ctx *gin.Context) {
		defer func() {
			if rcover := recover(); rcover != nil {
				log.Println("Panic : " + rcover.(string))
			}
		}()

		//Websocket接続
		wsconn, err := wsupgrader.Upgrade(ctx.Writer, ctx.Request, nil)

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
				auth_result, err := auth.Valid_Token(readmsg.Payload.(string))

				//エラー処理
				if err != nil {
					log.Println(err)
					//切断
					break
				}

				//リフレッシュトークンの場合閉じる
				if auth_result.IsRefresh {
					break
				}

				//ユーザIDを設定
				userid = auth_result.Userid

				//接続を追加
				wsconns.Store(auth_result.Userid, wsconn)
				ws_contiune = true

				//抜ける
				break
			}
		}

		//コンティニュー
		if !ws_contiune {
			//切断
			wsconn.Close()

			log.Println("切断")
			return
		}

		//スレッド作成
		go handle_ws(wsconn, userid)
	})

	debug_mode := os.Getenv("DEBUG_MODE")

	//デバッグモード
	if debug_mode == "true" {
		router.RunTLS("0.0.0.0:12222", "./keys/server.crt", "./keys/server.key") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	} else {
		router.Run("0.0.0.0:12222")
	}

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
