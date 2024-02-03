package database

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

//位置情報
type User_Location struct {
	UID string  `gorm:"primaryKey"` //トークンID
	Lat float64 //緯度
	Lng float64 //経度
}

// フレンド申請一覧
type User_Data struct {
	UID      string `gorm:"primaryKey"` //ユーザID
	Distance int64  //距離
	Status   string //ステータス
}

// フレンド申請一覧
type Ignore_Point struct {
	Randid      string  `gorm:"primaryKey"` //ポイントID
	UID         string  //ユーザID
	Distance    int64   //距離
	Latiubetude float64 //緯度
	Longitude   float64 //経度
}

// フレンド一覧
type Friends struct {
	UID         string `gorm:"primaryKey"` //トークンID
	Sender_id   string //一人目のユーザー
	Receiver_id string //二人目のユーザー
	SendTime    int64
}

// フレンド申請一覧
type Sent struct {
	UID         string `gorm:"primaryKey"` //トークンID
	Sender_id   string
	Receiver_id string

	SendTime int64
}
