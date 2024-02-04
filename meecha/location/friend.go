package location

import (
	"context"
	"errors"
	"log"
	"meecha/friends"

	"github.com/vmihailenco/msgpack/v5"
)

//フレンド一覧取得
func Get_Friends(uid string) ([]string,error) {
	if (!isinit) {
		//初期化されていなかったらエラーを返す
		return nil, errors.New("Not Initialized")
	}

	//キャッシュ取得
	result,err := friend_rdb.Get(context.Background(), uid).Result()

	//エラー処理
	if err != nil {
		//キャッシュ更新
		err = Update_Cache(uid)

		//エラー処理
		if err != nil {
			log.Println(err)
			return nil,err
		}

		//キャッシュ取得
		result,err = friend_rdb.Get(context.Background(), uid).Result()

		//エラー処理
		if err != nil {
			log.Println(err)
			return nil,err
		}
	}

	//デシリアライズ
	var aiteid_list []string
	err = msgpack.Unmarshal([]byte(result), &aiteid_list)

	//エラー処理
	if err != nil {
		log.Println(err)
		return nil,err
	}

	//フレンド一覧
	return aiteid_list,nil
}

//キャッシュ更新
func Update_Cache(uid string) (error) {
	if (!isinit) {
		//初期化されていなかったらエラーを返す
		return errors.New("Not Initialized")
	}

	log.Println("キャッシュ更新")

	//フレンド一覧取得 (db)
	result,err := friends.Get_Friends(uid)

	//エラー処理
	if err != nil {
		log.Println(err)
		return err
	}

	//相手のリスト
	var aiteid_list []string


	for key := range result {
		//キャッシュ更新
		aiteid_list = append(aiteid_list, result[key]["aiteid"])
	}

	//シリアライズ
	data,err := msgpack.Marshal(aiteid_list)

	//エラー処理
	if err != nil {
		log.Println(err)
		return err
	}

	//キャッシュ更新
	err = friend_rdb.Set(context.Background(), uid, data, Friend_exp).Err()
	
	//エラー処理
	if err != nil {
		log.Println(err)
		return err
	}

	return nil 
}