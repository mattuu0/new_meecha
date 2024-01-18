package database

//ユーザーアカウント
type User struct {
	UID		string	`gorm:"primaryKey"`	//ユーザID
	Name	string  					//ユーザ名

	HashPass	string 					//ぱすわーど
}