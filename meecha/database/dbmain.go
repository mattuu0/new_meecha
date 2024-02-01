package database

import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

var (
	DBpath	string
	dbconn	*gorm.DB
	IsInit 	bool = false
)

func Init() error {
	//データベースを開く
	db,err := gorm.Open(sqlite.Open(DBpath))

	//エラー処理
	if err != nil {
		return err
	}

	//グローバル変数に保存
	dbconn = db

	//マイグレーション
	err = dbconn.AutoMigrate(&User{})

	//エラー処理
	if err != nil {
		return err
	}

	//マイグレーション
	err = dbconn.AutoMigrate(&AccessToken{})

	//エラー処理
	if err != nil {
		return err
	}

	//マイグレーション
	err = dbconn.AutoMigrate(&RefreshToken{})

	//エラー処理
	if err != nil {
		return err
	}

	//マイグレーション
	err = dbconn.AutoMigrate(&User_Location{})

	//エラー処理
	if err != nil {
		return err
	}


	IsInit = true
	return nil
}

//DB接続を取得
func GetDB() *gorm.DB {
	return dbconn
}