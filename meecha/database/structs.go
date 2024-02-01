package database

//ユーザーアカウント
type User struct {
	UID  string `gorm:"primaryKey"` //ユーザID
	Name string //ユーザ名

	HashPass string //ぱすわーど
}

//アクセストークン
type AccessToken struct {
	TokenID string `gorm:"primaryKey"` //トークンID
	UID     string //トークンのユーザID
	Exp     int64  //トークンの有効期限
}

//アクセストークン
type RefreshToken struct {
	TokenID  string `gorm:"primaryKey"` //トークンID
	UID      string //トークンのユーザID
	AccessID string //アクセストークンID
}

//位置情報
type User_Location struct {
	UID      string  `gorm:"primaryKey"` //トークンID
	TokenID  string  //アクセストークンID
	Lat      float64 //緯度
	Lng      float64 //経度
	WaitTime int64   //待機時間
}
