package auth

import (
	"errors"
	"meecha/database"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

//設定
var (
	SignMethod	jwt.SigningMethod = jwt.SigningMethodHS512
	Secret []byte	
)


//ユーザIDを指定してリフレッシュトークン生成
func Gen_Refresh_Token(uid string) (Tokens,error) {
	//空の結果を生成
	tokens := Tokens{AccessToken: "",RefreshToken: ""}

	//初期化されていなかったらエラーを返す
	if !isinit {
		return tokens,init_error()
	}


	//リフレッシュトークンのID
	RtokenID,err := Genid()

	//エラー処理
	if err != nil {
		return tokens,err
	}

	//リフレッシュトークン情報
	Rclaims := jwt.MapClaims{
		"userid" : uid,
		"tokenid" : RtokenID,
	}

	//トークン生成
	Rtoken := jwt.NewWithClaims(SignMethod,Rclaims)

	//トークン署名
	Refresh_token,err := Rtoken.SignedString(Secret)

	//エラー処理
	if err != nil {
		return tokens,err
	}
	

	//アクセストークン生成
	//アクセストークンのID
	AtokenID,err := Genid()

	//有効期限 (72時間)
	exp := time.Now().Add(time.Hour * 72).Unix()

	//エラー処理
	if err != nil {
		return tokens,err
	}

	//リフレッシュトークン情報
	Aclaims := jwt.MapClaims{
		"userid" : uid,
		"tokenid" : AtokenID,
		"exp" : exp,
	}

	//トークン生成
	Atoken := jwt.NewWithClaims(SignMethod,Aclaims)

	//トークン署名
	Access_token,err := Atoken.SignedString(Secret)

	//エラー処理
	if err != nil {
		return tokens,err
	}

	//トークンを登録
	err = register_token(uid,AtokenID,RtokenID,exp)

	//エラー処理
	if err != nil {
		return tokens,err
	}

	//結果に設定
	tokens.AccessToken = Access_token
	tokens.RefreshToken = Refresh_token

	return tokens,nil
}

//トークンを登録する
func register_token(uid	string,AccessId string,RefreshId string,exp int64) error {
	//初期化されていなかったらエラーを返す
	if !isinit {
		return init_error()
	}

	//アクセストークンの情報
	Atoken := database.AccessToken{
		UID: uid,
		TokenID: AccessId,
		Exp: exp,
	}

	//リフレッシュトークンの情報
	Rtoken := database.RefreshToken{
		UID: uid,
		TokenID: RefreshId,
		AccessID: AccessId,
	}

	//登録
	dbconn.Create(&Atoken)
	dbconn.Create(&Rtoken)

	return nil
}

//リフレッシュトークン無効化
func DisableRToken(tokenid string) error {
	//初期化されていなかったらエラーを返す
	if !isinit {
		return init_error()
	}

	//トークン
	filter_token := database.RefreshToken{TokenID: tokenid}

	//トークン取得
	result := dbconn.First(&filter_token)

	//トークンが見つからない場合戻る
	if errors.Is(result.Error,gorm.ErrRecordNotFound) {
		return nil
	}

	//エラー処理
	if result.Error != nil {
		return result.Error
	}

	//トークン削除
	dbconn.Delete(&filter_token)

	return nil
}