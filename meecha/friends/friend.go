package friends

import (
	"errors"
	"log"
	_ "log"
	"meecha/auth"
	"meecha/database"
	"strconv"
	"time"
)
var(
	Sid = "Sid"
	Rid = "Rid"
)

//フレンドであることをDBに登録
func Record_Friends(Sid string,Rid string)(string,error){

	//フレンドUIDを取得
	fuid,err := auth.Genid()

	//フレンドトークンの情報
	Ftoken := database.Friends{
		UID: 		 fuid,					//uuid
		Sender_id:   Sid,        			//送った側のID
		Receiver_id: Rid,		    		//受け取る側のID
		SendTime:    time.Now().Unix(),	    //送った時間
	}

	//ネームトークンフィルター
	var named_filter database.Friends

	//フレンド欄に既に存在していたら1、存在していなかったら0を代入(and)
	Sender   := dbconn.Where(database.Friends{Sender_id: Sid}).Where(database.Friends{Receiver_id: Rid}).First(&named_filter).RowsAffected
	//フレンド欄に既に存在していたら1、存在していなかったら0を代入(and)(逆)
	Receiver := dbconn.Where(database.Friends{Sender_id: Rid}).Where(database.Friends{Receiver_id: Sid}).First(&named_filter).RowsAffected

	log.Println(Sender)
	log.Println(Receiver)

	//既に存在している場合
	if(Sender + Receiver != 0){
		return "",errors.New("already_friend_registered")
	}

	//データベースに書き込む
	dbconn.Create(&Ftoken)

	return fuid,err
}

//フレンド一覧取得
func Get_Friends(Acquirer_id string) (map[string]map[string]string,error){

	//ネームトークンフィルター
	var named_filter []database.Friends

	//フレンドのなかに存在する配列を取得(or)
	acquisition := dbconn.Where(database.Friends{Sender_id: Acquirer_id}).Or(database.Friends{Receiver_id: Acquirer_id}).Find(&named_filter)
	//送信側が受信した要素の数を代入
	length := acquisition.RowsAffected

	//map宣言
	maps := map[string]map[string]string{}

	//送信した配列をすべてmapに代入
	for i := 0; i < int(length); i++ {
		maps[named_filter[i].UID] = map[string]string{
			Sid:named_filter[i].Sender_id,					//取得者の名前
			Rid:named_filter[i].Receiver_id,					//受信した側の名前
			"time":strconv.Itoa(int(named_filter[i].SendTime)), //フレンドになった時間
		}
	}
	
	//送信した配列の個数がが0の時
	if(length == 0){
		return map[string]map[string]string{},nil
	}
	
	return maps,nil
}

//フレンドを消す
func Delete_Friend(UID string,Deleter_id string) (error){

	//フレンドが存在しているか
	existence,err := Get_Friends(Deleter_id)

	if  err != nil{
		return errors.New("friends_not_existence")
	}

	user1 := (existence[UID][Rid])
	user2 := (existence[UID][Sid])


	//リクエストでエラーならば
	if err != nil {	
		log.Println("error")
		return err
	}
	
	log.Println(user1)
	log.Println(user2)

	//リクエスト者とUIDが一致していない時
	if !(user1 == Deleter_id || user2 == Deleter_id) {
		return errors.New("user_mismatch_existing")
	}

	log.Println("フレンドの消去")
	log.Println(UID)

	//フレンドをDBから消去
	dbconn.Delete(database.Friends{},database.Friends{UID: UID})

	return err
}


