package auth

import (
	"errors"
	"gorm.io/gorm"

	"meecha/database"
)

var (
	dbconn *gorm.DB
	isinit bool = false
)

// 初期化
func Init() error {
	//データベースが初期化されているか
	if !database.IsInit {
		//初期化されていなかったらエラーを返す
		return errors.New("database not initialized")
	}

	//データベース接続を取得
	dbconn = database.GetDB()

	//初期化済みにする
	isinit = true

	return nil
}

//初期化していないエラーを返す
func init_error() error {
	return errors.New("package not initialized")
}