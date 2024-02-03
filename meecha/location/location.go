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

func Add_Ignore_Point(uid string, Latitude float64, Longitude float64, distance int64) (string, error) {
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
		Distance:    distance,
		Latiubetude: Latitude,
		Longitude:   Longitude,
	}

	//除外ポイント追加
	result := dbconn.Save(&Ignore_Data)

	//エラー処理
	if result.Error != nil {
		log.Println(result.Error)
		return "", result.Error
	}

	//情報更新
	Refresh_Ignore_Point(uid)

	return randid, nil
}

type Ignore_point struct {
	//緯度
	Latitude float64

	//経度
	Longitude float64

	//距離
	Distance int64
}

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

	//有効期限設定
	if err := Ignore_rdb.Expire(ctx, uid, time.Duration(5)*time.Minute).Err(); err != nil {
		return err
	}

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
	if err := Ignore_rdb.Set(ctx, uid, marshal_data, 0).Err(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

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

func Get_Distance(point1 Point,point2 Point) int {
	var dist float64
	geodesic.WGS84.Inverse(point1.Lat, point1.Lon, point2.Lat, point2.Lon, &dist, nil, nil)

	return int(dist)
}