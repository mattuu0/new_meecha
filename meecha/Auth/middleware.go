package auth

import (
	"log"

	"github.com/gin-gonic/gin"
)

//設定
var (
	HeaderName	string	= "Authorization"
)

//設定
func Auth_Init(router *gin.Engine) {
	router.LoadHTMLFiles(
		"./templates/Auth/Auth_Error.html",
	)
}

//認証ミドルウェア
func Auth_Middleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        //ヘッダー取得
		AuthToken := ctx.Request.Header.Get(HeaderName)

		log.Println(len(AuthToken))
		//トークンが空文字の場合戻る
		if len(AuthToken) == 0 {
			return
		}

		//トークン検証
		result,err := Valid_Token(AuthToken)

		//エラー処理
		if err != nil {
			log.Print(err)
			//処理停止
			ctx.AbortWithStatus(403)
			return
		}

		log.Println(result)

		//処理続行
        ctx.Next()
    }
}
