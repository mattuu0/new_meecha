package auth

import (
	"log"

	"github.com/gin-gonic/gin"
)

// 認証結果
type Auth_Result struct {
	Success   bool   //認証成功か
	UserId    string //認証ユーザー
	IsRefresh bool   //リフレッシュトークンか

	Token string //トークン
}

// 設定
var (
	HeaderName string = "Authorization"
	KeyName    string = "Auth" //認証結果のキー
)

// 認証ミドルウェア
func Auth_Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//空の認証結果
		Set_Result := Auth_Result{
			Success:   false,
			UserId:    "",
			IsRefresh: false,
			Token: "",
		}

		//ヘッダー取得
		AuthToken := ctx.Request.Header.Get(HeaderName)

		//トークンが空文字の場合戻る
		if len(AuthToken) == 0 {
			//結果を設定
			ctx.Set(KeyName, Set_Result)
			return
		}

		//トークン検証
		result, err := Valid_Token(AuthToken)

		//エラー処理
		if err != nil {
			log.Print(err)
			//結果を設定
			ctx.Set(KeyName, Set_Result)

			//処理を終了
			ctx.Next()
			return
		}

		//情報を設定
		Set_Result.Success = true
		Set_Result.IsRefresh = result.IsRefresh
		Set_Result.UserId = result.Userid
		Set_Result.Token = AuthToken
		ctx.Set(KeyName, Set_Result)

		//処理続行
		ctx.Next()
	}
}

// 認証ミドルウェア
func Auth_Require_Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//ヘッダー取得
		AuthToken := ctx.Request.Header.Get(HeaderName)

		//トークンが空文字の場合戻る
		if len(AuthToken) == 0 {
			//結果を設定
			ctx.AbortWithStatus(401)
			return
		}

		//トークン検証
		result, err := Valid_Token(AuthToken)

		//エラー処理
		if err != nil {
			log.Print(err)
			//結果を設定
			ctx.AbortWithStatus(401)

			//処理を終了
			ctx.Next()
			return
		}

		//リフレッシュトークンか
		if result.IsRefresh {
			//結果を設定
			ctx.AbortWithStatus(400)
			return
		}

		//ユーザが存在するか
		find_result,err := GetUser_ByID(result.Userid)

		//エラー処理
		if err != nil {
			log.Print(err)
			//結果を設定
			ctx.AbortWithStatus(403)
			return;
		}

		//ユーザが見つからない場合
		if !find_result.IsFind {
			//結果を設定
			ctx.AbortWithStatus(403)
			return
		}

		//処理続行
		ctx.Next()
	}
}