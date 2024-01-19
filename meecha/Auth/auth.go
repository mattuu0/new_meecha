package auth

import (
	"golang.org/x/crypto/bcrypt"

	"meecha/database"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"errors"
)

var (
	Cost int = 10
)


//ユーザを作成する
func CreateUser(username string,password string) (database.User,error) {
	//空のユーザを作成する
	user := database.User{}

	//初期化されていなかったらエラー
	if !isinit {
		return user,Get_Init_Error()
	}

	//IDを生成する
	uid,err := Genid()

	//エラー処理
	if err != nil {
		return user,err
	}

	//ぱすわーどをハッシュ化する
	spass,err := SecurePass(password)

	//エラー処理
	if err != nil {
		return user,err
	}

	//ユーザ情報を設定する
	user.UID = uid
	user.Name = username
	user.HashPass = spass

	//データベースに書き込む
	dbconn.Create(&user)

	//コミットする
	//dbconn.Commit()

	//ユーザデータを返す
	return user,nil
}

//ユーザ名でユーザを取得する
func GetUser_ByName(uname string) (FindResult,error) {
	//空のユーザを作成する
	fuser := database.User{}

	//結果
	result := FindResult{IsFind: false}

	//初期化されていなかったらエラー
	if !isinit {
		return result,Get_Init_Error()
	}

	//ユーザを取得する
	find_result := dbconn.Preload(clause.Associations).First(&fuser,"name = ?",uname)
	
	if err := find_result.Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return result,gorm.ErrRecordNotFound
	}
	
	//見つかった設定にする
	result.IsFind = true

	//情報をセットする
	result.UserData = fuser

	return result,nil
}

//ぱすわーどをハッシュ化する
func SecurePass(password string) (string,error) {
	//ぱすわーど文字列をバイナリにする
	binary_path := []byte(password)

	//ハッシュ化
	hashed, err := bcrypt.GenerateFromPassword(binary_path,Cost)

	//エラー処理
	if err != nil {
		return "",err
	}

	//ぱすわーどを返す
	return string(hashed),nil
}

/*
//ログイン
func Login(uname string,password string) (LoginResult,error) {
	//ログイン結果
	result := LoginResult{IsFind: false}

	_ = result
}

*/