package location

import (
	"github.com/vmihailenco/msgpack/v5"
	"errors"
	"log"
	"meecha/auth"
	"meecha/database"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm/clause"

	"context"

	"github.com/tidwall/geodesic"
)

func Update_Geo(uid string, Latitude float64, Longitude float64) error {
	if !isinit {
		//初期化されていなかったらエラーを返す
		return errors.New("Not Initialized")
	}

	var ctx = context.Background()

	if location_rdb.Exists(ctx, uid).Val() > 0 {
		//有効期限設定
		if err := location_rdb.Expire(ctx, uid, time.Duration(5)*time.Minute).Err(); err != nil {
			return err
		}
	}

	//位置情報更新
	if err := location_rdb.GeoAdd(ctx, uid, &redis.GeoLocation{Longitude: Longitude, Latitude: Latitude, Name: LocationKey, GeoHash: 1}).Err(); err != nil {
		log.Println(err)
		return err
	}

	//有効期限設定
	if err := location_rdb.Expire(ctx, uid, time.Duration(5)*time.Minute).Err(); err != nil {
		return err
	}

	return nil
}

type Ignore_point struct {
	//緯度
	Latitude float64

	//経度
	Longitude float64

	//距離
	Distance int64
}

func Add_Ignore_Point(uid string, point Ignore_point) (string, error) {
	if !isinit {
		//初期化されていなかったらエラーを返す
		return "", errors.New("Not Initialized")
	}

	//ID生成
	randid, err := auth.Genid()

	//エラー処理
	if err != nil {
		log.Println(err)
		return "", err
	}

	//除外ポイント
	Ignore_Data := database.Ignore_Point{
		Randid:      randid,
		UID:         uid,
		Distance:    point.Distance,
		Latiubetude: point.Latitude,
		Longitude:   point.Longitude,
	}

	//除外ポイント追加
	result := dbconn.Save(&Ignore_Data)

	//エラー処理
	if result.Error != nil {
		log.Println(result.Error)
		return "", result.Error
	}

	return randid, nil
}

//除外ポイント更新
func Refresh_Ignore_Point(uid string) error {
	if !isinit {
		//初期化されていなかったらエラーを返す
		return errors.New("Not Initialized")
	}

	//除外ポイント取得
	var ignore_datas []database.Ignore_Point
	result := dbconn.Preload(clause.Associations).Where(database.Ignore_Point{UID: uid}).Find(&ignore_datas)

	//エラー処理
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}

	//コンテキスト
	var ctx = context.Background()

	//設定するデータ
	setmap := map[string]Ignore_point{}

	//除外ポイント追加
	for i := 0; i < len(ignore_datas); i++ {
		Longitude := ignore_datas[i].Longitude  //位置情報
		Latitude := ignore_datas[i].Latiubetude //位置情報
		randid := ignore_datas[i].Randid        //ID
		distance := ignore_datas[i].Distance    //距離

		//ポイント情報設定
		setmap[randid] = Ignore_point{
			Latitude:  Latitude,
			Longitude: Longitude,
			Distance:  distance,
		}
	}

	//データをバイナリに
	marshal_data, err := msgpack.Marshal(setmap)

	//エラー処理
	if err != nil {
		log.Println(err)
		return err
	}

	//データリセット
	if err := Ignore_rdb.Set(ctx, uid, marshal_data, time.Duration(5)*time.Minute).Err(); err != nil {
		log.Println(err)
		return err
	}
	

	return nil
}

//除外ポイント取得
func Get_Ignore_Points(uid string) (map[string]Ignore_point,error) {
	//でコード
	setmap := map[string]Ignore_point{}

	if !isinit {
		//初期化されていなかったらエラーを返す
		return setmap,errors.New("Not Initialized")
	}

	//コンテキスト
	var ctx = context.Background()

	//存在しない場合
	if Ignore_rdb.Exists(ctx, uid).Val() == 0 {
		//有効期限設定
		err := Refresh_Ignore_Point(uid)

		//エラー処理
		if err != nil {
			log.Println(err)
			return setmap,err
		}
	}

	//情報取得
	marshal_data, err := Ignore_rdb.Get(ctx, uid).Result()

	//エラー処理
	if err != nil {
		log.Println(err)
		return setmap,err
	}

	//デコード
	err = msgpack.Unmarshal([]byte(marshal_data), &setmap)

	//エラー処理
	if err != nil {
		log.Println(err)
		return setmap,err
	}

	return setmap,nil
}

//除外ポイント更新
func Update_Ignore_Point(uid string,pointid string, point Ignore_point) (string, error) {
	if !isinit {
		//初期化されていなかったらエラーを返す
		return "", errors.New("Not Initialized")
	}

	//データ
	point_data := database.Ignore_Point{}
	//取得
	result := dbconn.Where(database.Ignore_Point{Randid: pointid, UID: uid}).First(&point_data)

	//エラー処理
	if result.Error != nil {
		log.Println(result.Error)
		return "", result.Error
	}

	//情報更新
	point_data.Longitude = point.Longitude
	point_data.Latiubetude = point.Latitude
	point_data.Distance = point.Distance

	//保存
	result = dbconn.Save(&point_data)

	//エラー処理
	if result.Error != nil {
		log.Println(result.Error)
		return "", result.Error
	}

	//データ更新
	err := Refresh_Ignore_Point(uid)

	//エラー処理
	if err != nil {
		log.Println(err)
		return "", err
	}

	return pointid, nil
}

//除外ポイント全削除
func Remove_All_Ignore_Point(uid string) (error) {
	if !isinit {
		//初期化されていなかったらエラーを返す
		return errors.New("Not Initialized")
	}

	//ユーザIDが同じものを一括削除
	result := dbconn.Where(database.Ignore_Point{UID: uid}).Delete(&database.Ignore_Point{})

	//エラー処理
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}

	return nil
}

//除外ポイント削除
func Remove_Ignore_Point(uid string, pointid string) (string, error) {
	if !isinit {
		//初期化されていなかったらエラーを返す
		return "", errors.New("Not Initialized")
	}

	var ctx = context.Background()

	//ID生成
	randid, err := auth.Genid()

	//エラー処理
	if err != nil {
		log.Println(err)
		return "", err
	}

	//除外ポイント追加
	result := dbconn.Delete(&database.Ignore_Point{}, database.Ignore_Point{Randid: pointid, UID: uid})

	//エラー処理
	if result.Error != nil {
		log.Println(result.Error)
		return "", result.Error
	}

	//ポイント追加
	if err := location_rdb.ZRem(ctx, uid, pointid).Err(); err != nil {
		log.Println(err)
		return "", err
	}

	return randid, nil
}

//位置情報を検証する
func Validate_Geo(uid string,point Point) (bool,error) {
	if !isinit {
		//初期化されていなかったらエラーを返す
		return false,errors.New("Not Initialized")
	}

	//除外ポイント取得
	setmap,err := Get_Ignore_Points(uid)

	//エラー処理
	if err != nil {
		log.Println(err)
		return false,err
	}

	valid_result := true	
	//除外ポイントを検証
	for _, ignore_point := range setmap {
		//距離
		distance := Get_Distance(point, Point{Lat: ignore_point.Latitude, Lon: ignore_point.Longitude})

		//距離判定
		if distance < ignore_point.Distance {
			//除外時範囲に入っている場合
			valid_result = false
			break
		}
	}

	return valid_result, nil
}

//位置情報取得
func GetLocation(uid string) (Point, error) {
	if !isinit {
		//初期化されていなかったらエラーを返す
		return Point{}, errors.New("Not Initialized")
	}

	//位置情報取得
	getopos, err := location_rdb.GeoPos(context.Background(), uid, LocationKey).Result()

	//エラー処理
	if err != nil {
		log.Println(err)
		return Point{}, err
	}

	//結果設定
	result_point := Point{}
	result_point.Lat = getopos[0].Latitude
	result_point.Lon = getopos[0].Longitude

	return result_point,nil
}

//位置情報消す
func RemoveLocation(uid string) error {
	if err := location_rdb.ZRem(context.Background(), uid, LocationKey).Err(); err != nil {
		log.Println(err)
		return err
	}
	
	return nil
}

// 地点を表す構造体
type Point struct {
	Lat float64 // 緯度（度数法）
	Lon float64 // 経度（度数法）
}

// 2地点間の距離を求める
func Get_Distance(point1 Point,point2 Point) int64 {
	var dist float64
	geodesic.WGS84.Inverse(point1.Lat, point1.Lon, point2.Lat, point2.Lon, &dist, nil, nil)

	return int64(dist)
}


//通知距離を更新する
func Update_Notification_Distance(uid string, distance int64) error {
	if !isinit {
		//初期化されていなかったらエラーを返す
		return errors.New("Not Initialized")
	}

	//通知距離更新
	if err := distance_rdb.Set(context.Background(), uid, distance,Distance_exp).Err(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
