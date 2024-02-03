package location

import (
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"meecha/database"
)

var (
	token_rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: 1000,
	})

	location_rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       1,  // use default DB
		PoolSize: 1000,
	})

	Ignore_rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       2,  // use default DB
		PoolSize: 1000,
	})

	LocationKey = "location"
	isinit bool = false
	secret []byte
	dbconn *gorm.DB = nil
	TokenExp time.Duration = time.Duration(5) * time.Second
)

// 初期化
func Init(token string) error {
	//データベースが初期化されているか
	if !database.IsInit {
		//初期化されていなかったらエラーを返す
		return errors.New("database not initialized")
	}

	//初期化済みにする
	isinit = true

	//シークレットを設定
	secret = []byte(token)

	//DB接続を取得
	dbconn = database.GetDB()

	return nil
}

//初期化していないエラーを返す
func init_error() error {
	return errors.New("package not initialized")
}