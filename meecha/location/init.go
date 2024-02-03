package location

import (
	"errors"
	"github.com/redis/go-redis/v9"

	"meecha/database"
)

var (
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: 1000,
	})
	isinit bool = false
	secret []byte
	tokens map[string]string = map[string]string{}
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

	return nil
}

//初期化していないエラーを返す
func init_error() error {
	return errors.New("package not initialized")
}