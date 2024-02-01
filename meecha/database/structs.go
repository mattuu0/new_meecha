package database

import (
)

// ユーザーアカウント
type User struct {
	UID  string `gorm:"primaryKey"` //ユーザID
	Name string //ユーザ名

	HashPass string //パスワード
}

// アクセストークン
type AccessToken struct {
	TokenID string `gorm:"primaryKey"` //トークンID
	UID     string //トークンのユーザID
	Exp     int64  //トークンの有効期限
}

// リフレッシュトークン
type RefreshToken struct {
	TokenID  string `gorm:"primaryKey"` //トークンID
	UID      string //トークンのユーザID
	AccessID string //アクセストークンID
}

// フレンド一覧
type Friends struct {
	UID        string  `gorm:"primaryKey"` //トークンID
	Sender_id    string //一人目のユーザー
	Receiver_id string //二人目のユーザー
	SendTime int64
}

// フレンド申請一覧
type Sent struct {
	UID        string `gorm:"primaryKey"` //トークンID
	Sender_id  string
	Receiver_id string

	SendTime int64
}
