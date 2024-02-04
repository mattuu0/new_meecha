package main

import (
	"encoding/json"
	"log"
	"meecha/location"

	"github.com/gin-gonic/gin"
)

type Ignore_Point_Datas struct {
	Points string //除外設定
}

type Ignore_Point_Data struct {
	Distance	int64
	Lat 		float64
	Lng 		float64
}


func save_ignore_point(ctx *gin.Context) {

	//認証情報を取得
	uid, err := getid(ctx)

	//エラー処理
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}

	//送信情報を取得
	var ignore_point_data Ignore_Point_Datas

	//データを紐付ける
	err = ctx.BindJSON(&ignore_point_data)

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(500)
		return
	}

	//除外設定
	var ignore_data []Ignore_Point_Data

	//JSON文字列を配列に変換
	err = json.Unmarshal([]byte(ignore_point_data.Points), &ignore_data)

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(500)
		return
	}

	//古いポイントを削除
	err = location.Remove_All_Ignore_Point(uid)

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(500)
		return
	}

	for _,val := range ignore_data {
		//データ生成
		point_data := location.Ignore_point{
			Latitude: val.Lat,
			Longitude: val.Lng,
			Distance: val.Distance,
		}

		//除外ポイント追加
		_,err := location.Add_Ignore_Point(uid,point_data)

		//エラー処理
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(500)
			return
		}
	}

	//データ更新
	err = location.Refresh_Ignore_Point(uid)

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(500)
		return
	}

	//データ返却
	ctx.JSON(200, nil)
}

//除外ポイント取得
func get_ignore_point(ctx *gin.Context) {

	//認証情報を取得
	uid, err := getid(ctx)

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(500)
		return
	}

	//除外ポイント取得
	result, err := location.Get_Ignore_Points(uid)

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(500)
		return
	}

	//データ返却
	ctx.JSON(200, result)
}