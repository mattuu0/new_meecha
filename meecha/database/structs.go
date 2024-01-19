package database

import "gorm.io/gorm"

//ユーザーアカウント
type User struct {
	UID  string `gorm:"primaryKey"` //ユーザID
	Name string //ユーザ名

	HashPass string //ぱすわーど
}

//アクセストークン
type AccessToken struct {
	gorm.Model
	TokenID string `gorm:"primaryKey"` //トークンID
	UID     string //トークンのユーザID
	Exp     int64  //トークンの有効期限
}

//アクセストークン
type RefreshToken struct {
	gorm.Model
	TokenID  string `gorm:"primaryKey"` //トークンID
	UID      string //トークンのユーザID
	AccessID string //アクセストークンID
}
