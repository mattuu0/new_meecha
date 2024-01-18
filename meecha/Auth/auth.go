package auth

import (
	"crypto/sha512"
	"encoding/hex"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"meecha/database"

	"gorm.io/gorm"

	"errors"
)

var (
	Cost int = 10
	dbconn *gorm.DB
	isinit bool = false
)

type FindResult struct {
	IsFind bool
	User database.User
}

//初期化
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

// IDを生成する
func genid() (string, error) {

	//UUIDを生成する
	uuid_obj, err := uuid.NewRandom()

	//エラー処理
	if err != nil {
		return "", err
	}

	//UUID文字列をSHA512文字列にする
	hash_byte := sha512.Sum512([]byte(uuid_obj.String()))
	hex_string := hex.EncodeToString(hash_byte[:])

	//文字列を返す
	return hex_string, nil
}

//ユーザを作成する
func CreateUser(username string,password string) (database.User,error) {
	//空のユーザを作成する
	user := database.User{}

	//初期化されていなかったらエラー
	if !isinit {
		return user,Get_Init_Error()
	}

	//IDを生成する
	uid,err := genid()

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
	dbconn.Commit()

	//ユーザデータを返す
	return user,nil
}

//ユーザ名でユーザを取得する
func GetUser_ByName(uname string) (FindResult,error) {
	//空のユーザを作成する
	user := database.User{}

	//結果
	result := FindResult{IsFind: false,User: user}

	//初期化されていなかったらエラー
	if !isinit {
		return result,Get_Init_Error()
	}

	//名前を設定する
	user.Name = uname

	//ユーザを取得する
	if err := dbconn.First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return result,gorm.ErrRecordNotFound
	}

	//見つかった設定にする
	result.IsFind = true

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

//初期化されていないエラーを返す
func Get_Init_Error() error {
	return errors.New("not initialized")
}