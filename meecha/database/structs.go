package database

import "gorm.io/gorm"

//ユーザーアカウント
type User struct {
	UID		string	`gorm:"primaryKey"`	//ユーザID
	Name	string  					//ユーザ名

	HashPass	string 					//ぱすわーど
}

//トークン
type Tokens struct {
	gorm.Model
	TokenID	string `gorm:"primaryKey"`	//トークンID
	UID		string						//トークンのユーザID
	IsAccessT bool						//アクセストークンか
}