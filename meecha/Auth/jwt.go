package auth

import (
	"errors"
	"fmt"
	"log"
	"meecha/database"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 設定
var (
	SignMethod jwt.SigningMethod = jwt.SigningMethodHS512
	Secret     []byte
)

// アクセストークンデータ
type Atoken_Info struct {
	Token   string //アクセストークン
	TokenId string //トークンID
	exp     int64  //有効期限
}

//アクセストークン更新
func Refresh(uid string) (string,error) {
	//ユーザIDからアクセストークンID取得
	atokenid, err := Get_AToken_ByUID(uid)

	//エラー処理
	if err != nil {
		return "", err
	}

	//既存のトークン無効化
	if atokenid != "" {
		//トークンを無効化出来ない場合エラー
		if err := DisableAToken(atokenid); err != nil {
			log.Println(err)
			return "", err
		}
	}

	//アクセストークン生成
	Atoken,token_result,err := Gen_Access_Token(uid)

	//エラー処理
	if err != nil {
		log.Println(err)
		return "",err
	}

	//リフレッシュトークン取得
	rtokenid,err := Get_Token_ByUID(uid)

	//エラー処理
	if err != nil {
		log.Println(err)
		return "",err
	}

	//トークン登録
	err = register_token(uid,token_result.TokenId,rtokenid,token_result.exp)

	//エラー処理
	if err != nil {
		return "",err
	}

	return Atoken,nil
}

// アクセストークン生成
func Gen_Access_Token(uid string) (string, Atoken_Info, error) {
	result := Atoken_Info{}

	//アクセストークン生成
	//アクセストークンのID
	AtokenID, err := Genid()

	//有効期限 (72時間)
	exp := time.Now().Add(time.Hour * 72).Unix()

	//エラー処理
	if err != nil {
		return "", result, err
	}

	//リフレッシュトークン情報
	Aclaims := jwt.MapClaims{
		"userid":   uid,
		"tokenid":  AtokenID,
		"exp":      exp,
		"IsAccess": true,
	}

	//トークン生成
	Atoken := jwt.NewWithClaims(SignMethod, Aclaims)

	//トークン署名
	Access_token, err := Atoken.SignedString(Secret)

	//エラー処理
	if err != nil {
		return "", result, err
	}

	//情報設定
	result.exp = exp
	result.Token = Access_token
	result.TokenId = AtokenID

	return Access_token, result, nil
}

// ユーザIDを指定してリフレッシュトークン生成
func Gen_Refresh_Token(uid string) (Tokens, error) {
	//空の結果を生成
	tokens := Tokens{AccessToken: "", RefreshToken: ""}

	//初期化されていなかったらエラーを返す
	if !isinit {
		return tokens, init_error()
	}

	//リフレッシュトークンのID
	RtokenID, err := Genid()

	//エラー処理
	if err != nil {
		return tokens, err
	}

	//リフレッシュトークン情報
	Rclaims := jwt.MapClaims{
		"userid":   uid,
		"tokenid":  RtokenID,
		"IsAccess": false,
	}

	//トークン生成
	Rtoken := jwt.NewWithClaims(SignMethod, Rclaims)

	//トークン署名
	Refresh_token, err := Rtoken.SignedString(Secret)

	//エラー処理
	if err != nil {
		return tokens, err
	}

	//アクセストークン生成
	//アクセストークンのID
	Atoken, token_info, err := Gen_Access_Token(uid)

	//エラー処理
	if err != nil {
		return tokens, err
	}

	//トークンを登録
	err = register_token(uid, token_info.TokenId, RtokenID, token_info.exp)

	//エラー処理
	if err != nil {
		return tokens, err
	}

	//結果に設定
	tokens.AccessToken = Atoken
	tokens.RefreshToken = Refresh_token

	return tokens, nil
}

// トークンを登録する
func register_token(uid string, AccessId string, RefreshId string, exp int64) error {
	//アクセストークンの情報
	Atoken := database.AccessToken{
		UID:     uid,
		TokenID: AccessId,
		Exp:     exp,
	}

	dbconn.Save(&Atoken)
	
	//リフレッシュトークンの情報
	Rtoken := database.RefreshToken{
		UID:      uid,
		TokenID:  RefreshId,
		AccessID: AccessId,
	}

	//登録
	dbconn.Save(&Rtoken)
	
	//コミットする
	//dbconn.commit()

	return nil
}

// リフレッシュトークン無効化
func DisableRToken(tokenid string) error {
	//初期化されていなかったらエラーを返す
	if !isinit {
		return init_error()
	}

	//トークン
	filter_token := database.RefreshToken{}

	//トークン取得
	result := dbconn.Preload(clause.Associations).First(&filter_token, database.RefreshToken{TokenID: tokenid})

	//トークンが見つからない場合戻る
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println("Token Not found")
		return nil
	}

	//エラー処理
	if result.Error != nil {
		return result.Error
	}

	//アクセストークン削除
	dbconn.Unscoped().Delete(&database.AccessToken{}, &database.AccessToken{TokenID: filter_token.AccessID})

	//トークン削除
	dbconn.Unscoped().Delete(&filter_token, database.RefreshToken{TokenID: tokenid})

	//コミットする
	//dbconn.commit()

	return nil
}

// アクセストークン無効化
func DisableAToken(tokenid string) error {
	//初期化されていなかったらエラーを返す
	if !isinit {
		return init_error()
	}

	log.Println(tokenid)

	//トークン
	filter_token := database.AccessToken{}

	//トークン取得
	result := dbconn.First(&filter_token, database.AccessToken{TokenID: tokenid})

	//トークンが見つからない場合戻る
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println("Token Not found")
		return nil
	}

	//エラー処理
	if result.Error != nil {
		return result.Error
	}

	filter_token = database.AccessToken{}

	//アクセストークン削除
	dbconn.Unscoped().Delete(&filter_token, database.AccessToken{TokenID: tokenid})

	//コミットする
	//dbconn.commit()

	return nil
}

// ユーザIDからリフレッシュトークンID取得
func Get_Token_ByUID(uid string) (string, error) {
	//リフレッシュトークンフィルター
	rtoken_filter := database.RefreshToken{}

	//トークン取得
	Rresult := dbconn.First(&rtoken_filter, database.RefreshToken{UID: uid})

	//トークンが見つからない場合戻る
	if errors.Is(Rresult.Error, gorm.ErrRecordNotFound) {
		return "", nil
	}

	//エラー処理
	if Rresult.Error != nil {
		return "", Rresult.Error
	}

	//リフレッシュトークン
	return rtoken_filter.TokenID, nil
}

// ユーザIDからアクセストークンID取得
func Get_AToken_ByUID(uid string) (string, error) {
	//リフレッシュトークンフィルター
	atoken_filter := database.AccessToken{}

	//トークン取得
	Aresult := dbconn.First(&atoken_filter, database.AccessToken{UID: uid})

	//トークンが見つからない場合戻る
	if errors.Is(Aresult.Error, gorm.ErrRecordNotFound) {
		return "", nil
	}

	//エラー処理
	if Aresult.Error != nil {
		return "", Aresult.Error
	}

	//リフレッシュトークン
	return atoken_filter.TokenID, nil
}

// トークンをバリデーション (ユーザID,トークンID,エラー) を返す
func Valid_Token(token_str string) (TokenResult, error) {
	//結果
	tresult := TokenResult{
		Tokenid:   "",
		Userid:    "",
		IsRefresh: false,
	}

	//パース
	token, err := jwt.Parse(token_str, func(token *jwt.Token) (interface{}, error) {
		//署名方法確認
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		//鍵を返す
		return Secret, nil
	})

	//エラー処理
	if err != nil {
		return tresult, err
	}

	//Claimをデコード
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//トークンID
		tokenid := string(claims["tokenid"].(string))

		var result *gorm.DB

		//アクセストークンかどうか
		if claims["IsAccess"].(bool) {
			//リフレッシュトークンフィルター
			Atoken_filter := database.AccessToken{}

			//トークン取得
			result = dbconn.First(&Atoken_filter, database.AccessToken{TokenID: tokenid})
		} else {
			//リフレッシュトークンフィルター
			rtoken_filter := database.RefreshToken{}

			//トークン取得
			result = dbconn.First(&rtoken_filter, database.RefreshToken{TokenID: tokenid})
		}

		//見つからない場合も
		//エラー処理
		if result.Error != nil {
			log.Println(result.Error)
			return tresult, result.Error
		}

		tresult.Userid = string(claims["userid"].(string))
		tresult.IsRefresh = !claims["IsAccess"].(bool)
		tresult.Tokenid = string(claims["tokenid"].(string))

		return tresult, nil
	} else {
		//検証に失敗したらエラーを返す
		return tresult, err
	}
}
