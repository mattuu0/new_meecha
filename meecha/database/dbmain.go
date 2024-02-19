package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"os"
)

var (
	DBpath	string
	dbconn	*gorm.DB
	IsInit 	bool = false
)

func Init() error {
	debug_mode := os.Getenv("DEBUG_MODE")

	var err error
	//デバッグモード
	if debug_mode == "true" {
		DBpath = "./meecha.db"
		dbconn,err = gorm.Open(sqlite.Open(DBpath))

		//エラー処理
		if err != nil {
			return err
		}

	} else {
		uname := os.Getenv("POSTGRES_USER")
		dbname := os.Getenv("POSTGRES_DB")
		password := os.Getenv("POSTGRES_PASSWORD")
		host := os.Getenv("POSTGRES_HOST")
		port := os.Getenv("POSTGRES_PORT")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",host,uname,password,dbname,port)
		dbconn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		//エラー処理
		if err != nil {
			return err
		}
	}

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


	//マイグレーション
	err = dbconn.AutoMigrate(&Sent{})

	//エラー処理
	if err != nil {
		return err
	}

	//マイグレーション
	err = dbconn.AutoMigrate(&Friends{})

	//エラー処理
	if err != nil {
		return err
	}

	//マイグレーション
	err = dbconn.AutoMigrate(&User_Data{})

	//エラー処理
	if err != nil {
		return err
	}

	//マイグレーション
	err = dbconn.AutoMigrate(&Ignore_Point{})

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