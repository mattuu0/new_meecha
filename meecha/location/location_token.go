package location

import (
	"github.com/golang-jwt/jwt/v5"

	"meecha/database"

	"meecha/auth"

	"fmt"
)


func GenToken(uid string) (string,error) {
	//初期化されているか
	if (!database.IsInit) {
		//初期化されていなかったらエラーを返す
		return "", auth.Get_Init_Error()
	}

	//トークンID
	tokenID, err := auth.Genid()

	//エラー処理
	if err != nil {
		return "", err
	}

	//トークン作成
	claims := jwt.MapClaims{
		"userid":   uid,
		"tokenid":  tokenID,
	}

	//トークン生成
	token := jwt.NewWithClaims(auth.SignMethod, claims)

	//トークン署名
	signed_token, err := token.SignedString(secret)

	//エラー処理
	if err != nil {
		return "", err
	}

	//トークンを登録
	err = registerToken(uid,tokenID)

	//エラー処理
	if err != nil {
		return "", err
	}

	//トークンを返す
	return signed_token, nil
}

//トークン検証
func VerifyToken(token string) (string, error) {
	//初期化されているか
	if (!isinit) {
		//初期化されていなかったらエラーを返す
		return "", auth.Get_Init_Error()
	}

	//トークン
	parse_token, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		//署名方法確認
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		//鍵を返す
		return secret, nil
	})

	//エラー処理
	if err != nil {
		return "", err
	}

	//Claimをデコード
	if claims, ok := parse_token.Claims.(jwt.MapClaims); ok && parse_token.Valid {
		//トークンID
		tokenid := string(claims["tokenid"].(string))
		userid := string(claims["userid"].(string))

		//トークンを取得
		findid,err := Get_Token_ByUID(userid)

		//エラー処理
		if err != nil {
			return "", err
		}

		//トークンIDを比較
		if tokenid != findid {
			return "", fmt.Errorf("invalid token")
		}

		//成功
		return userid, nil
	}

	//失敗
	return "", fmt.Errorf("invalid token")
}

//トークンを取得
func Get_Token_ByUID(uid string) (string, error) {
	//初期化されているか
	if (!isinit) {
		//初期化されていなかったらエラーを返す
		return "", auth.Get_Init_Error()
	}

	//トークン
	result,found := tokens[uid]

	//トークンが見つからなかったらエラーを返す
	if !found {
		return "", fmt.Errorf("token not found")
	}

	return result,nil
}

//トークン無効か
func Disable_Geo_Token(uid string) error {
	//トークン削除
	delete(tokens,uid)

	return nil
}

//トークンを登録
func registerToken(uid string,tokenid string) error {
	//初期化されているか
	if (!database.IsInit) {
		//初期化されていなかったらエラーを返す
		return auth.Get_Init_Error()
	}

	//トークン保存
	tokens[uid] = tokenid

	return nil
}