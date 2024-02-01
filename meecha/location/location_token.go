package location

import (
	"github.com/golang-jwt/jwt/v5"

	"meecha/database"
)

func GenToken() {
	//初期化されているか
	if (!database.IsInit) {
		//初期化されていなかったらエラーを返す
		return
	}

	//トークン作成
	//リフレッシュトークン情報
	Rclaims := jwt.MapClaims{
		"userid":   uid,
		"tokenid":  RtokenID,
	}

	//トークン生成
	Rtoken := jwt.NewWithClaims(SignMethod, Rclaims)

	//トークン署名
	Refresh_token, err := Rtoken.SignedString(Secret)
}