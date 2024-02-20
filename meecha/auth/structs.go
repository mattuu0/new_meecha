package auth

import "meecha/database"

// トークン
type Tokens struct {
	RefreshToken string
	AccessToken  string
}

// 検索結果
type FindResult struct {
	IsFind   bool
	UserData database.User
}

// 検索結果
type TokenResult struct {
	IsRefresh bool   //リフレッシュトークンか
	Userid    string //ユーザID
	Tokenid   string //トークンID
}

// ログイン結果
type LoginResult struct {
	Success bool //成功したか
	IsFind  bool //ユーザが見つかったか

	RefreshToken string //リフレッシュトークン
	AccessToken  string //あくせすトークン
	Userid       string //ユーザID
}
