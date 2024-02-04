package location

import (
	"context"
	"errors"
	"log"
	"meecha/database"
)

func Init_user_data(uid string) error {
	//ユーザデータ検索
	search_data := dbconn.Where(database.User_Data{UID: uid}).Find(&database.User_Data{})

	//ユーザデータが見つかった場合
	if search_data.RowsAffected > 0 {
		return nil
	}
	//ユーザデータ初期化
	user_data := database.User_Data{
		UID: uid,
		Distance: 50,
		Status: "Offline",
	}

	//ユーザデータ保存
	result := dbconn.Save(&user_data)

	//エラー処理
	if result.Error != nil {
		return result.Error
	}

	return nil
}

//ステータス更新
func Update_data(uid string,distance int64,status ...string) error {
	result_data := &database.User_Data{}

	//距離を検証
	if distance < 0 {
		return errors.New("Invalid Distance")
	}

	//5kmを超えている時
	if distance > 5000 {
		return errors.New("Invalid Distance")
	}

	//ユーザデータ検索
	result := dbconn.Where(database.User_Data{UID: uid}).First(result_data)

	//エラー処理
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}

	//ステータス更新
	if (status != nil) {
		result_data.Status = status[0]
	}
	result_data.Distance = distance

	//ユーザデータ保存
	result = dbconn.Save(result_data)

	//エラー処理
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}

	//キャッシュ保存
	err := Update_Notification_Distance(uid, distance)

	//エラー処理
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//ステータス更新
func Get_Notify_distance(uid string) (int64,error) {
	//初期化されていなかったら
	if !isinit {
		//初期化されていなかったらエラーを返す
		return 0,errors.New("Not Initialized")
	}

	//キャッシュから取得
	distance,err := distance_rdb.Get(context.Background(), uid).Int64()

	//エラー処理
	if err != nil {
		log.Println(err)
		log.Println("dbから取得")
	} else {
		log.Println("キャッシュ取得")
		return distance,nil
	}

	result_data := &database.User_Data{}

	//ユーザデータ検索
	result := dbconn.Where(database.User_Data{UID: uid}).First(result_data)

	//エラー処理
	if result.Error != nil {
		log.Println(result.Error)
		return 0,result.Error
	}

	//キャッシュに保存
	err = Update_Notification_Distance(uid,result_data.Distance)

	//エラー処理
	if err != nil {
		log.Println(err)
		return 0,err
	}

	return result_data.Distance,nil
}
