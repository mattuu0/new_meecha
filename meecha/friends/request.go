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

//フレンド申請送信
func Send(Sender_id string, Receiver_id string) (string, error){
	//senduid取得
	suid,err := auth.Genid()

	//センドトークンの情報
	Stoken := database.Sent{
		Sender_id:   Sender_id,        	//送った側のID
		Receiver_id: Receiver_id,		//受け取る側のID
		UID: 		 suid,				//uuid
		SendTime:    time.Now().Unix(),	//送った時間
	}

	//ネームトークンフィルター
	named_filter := database.Sent{}

	//自分にフレンドを送った場合
	if(Sender_id == Receiver_id) {
		return "",errors.New("send_to_myself")
	}

	//既に存在していたら1、存在していなかったら0を代入(and)
	Sender   := dbconn.Where(database.Sent{Sender_id: Sender_id}).Where(database.Sent{Receiver_id: Receiver_id}).First(&named_filter).RowsAffected
	//既に存在していたら1、存在していなかったら0を代入(and)(逆)
	Receiver := dbconn.Where(database.Sent{Sender_id: Receiver_id}).Where(database.Sent{Receiver_id: Sender_id}).First(&named_filter).RowsAffected

	//既に存在している場合
	if (Sender + Receiver != 0){
		return "",errors.New("request_is_already_existing")
	}

	//データベースに書き込む
	dbconn.Create(&Stoken)
	
	return suid,err
}

//受信済み取得
func Get_Received(Receiver_id string) (map[string]map[string]string, error){

	//ネームトークンフィルター
	var named_filter []database.Sent

	//受信側IDが同じ申請の配列を取得
	acquisition := dbconn.Where(database.Sent{Receiver_id: Receiver_id}).Find(&named_filter)
	//受信側が受信した要素の数を代入
	length := acquisition.RowsAffected

	//map宣言
	maps := map[string]map[string]string{}
	
	//受信した配列をすべてmapに代入
	for i := 0; i < int(length); i++ {
		uinfo,err := auth.GetUser_ByID(named_filter[i].Receiver_id)

		//ユーザー情報取得に失敗
		if err != nil{
			continue
		}

		maps[named_filter[i].UID] = map[string]string{
			"name":uinfo.UserData.Name,					//送信した側の名前
			"time":strconv.Itoa(int(named_filter[i].SendTime)), //リクエストした時間
		}
	}
	
	//受信した配列の個数がが0の時
	if(length == 0){
		return map[string]map[string]string{},nil
	}
	
	return maps,nil
}


//送信済み取得
func Get_Sent(Sender_id string) (map[string]map[string]string, error){

	//ネームトークンフィルター
	var named_filter []database.Sent

	//送信側IDが同じ申請の配列を取得
	acquisition := dbconn.Where(database.Sent{Sender_id: Sender_id}).Find(&named_filter)
	//送信側が受信した要素の数を代入
	length := acquisition.RowsAffected

	//map宣言
	maps := map[string]map[string]string{}

	//送信した配列をすべてmapに代入
	for i := 0; i < int(length); i++ {
		//受信者の情報取得
		uinfo,err := auth.GetUser_ByID(named_filter[i].Receiver_id)

		//ユーザー情報取得に失敗
		if err != nil{
			continue
		}

		maps[named_filter[i].UID] = map[string]string{
			"name":uinfo.UserData.Name,//受信した側の名前
			"time":strconv.Itoa(int(named_filter[i].SendTime)), //リクエストした時間
		}
	}
	
	//送信した配列の個数がが0の時
	if(length == 0){
		return map[string]map[string]string{},nil
	}
	
	return maps,nil
}

//識別子をもとに、userを2人を返す。存在しなければエラー
func Get_Request(UID string)(string,string,error) {

	//ネームトークンフィルター
	var named_filter database.Sent

	//UIDが空ならばエラー返す
	if UID == ""{
		return "","",errors.New("UID_does_not_exist")
	}

	//識別子が存在しているか
	result := dbconn.Where(database.Sent{UID: UID}).First(&named_filter)

	log.Println("named_filter")
	log.Println(result)

	//エラーならば0とエラー型を返す
	if  result.Error != nil{
		log.Println("0000")
	
		return "","",result.Error
	}

	//Sender_id
	Sid := named_filter.Sender_id
	//Receiver_id
	Rid := named_filter.Receiver_id

	return Sid,Rid,nil
}

//承認の一連メソッド
func Accept(UID string,Receiver_id string) (string,error){

	//リクエストが存在しているか
	Sid,Rid,err := Get_Request(UID)

	//リクエストでエラーならば
	if err != nil {	
		log.Println(err)
		return "",err
	}
	
	//受信者側が一致していない時
	if Rid != Receiver_id {
		return "",errors.New("user_mismatch_existing")
	}

	//フレンドであることをDBに登録
	fuid,err := Record_Friends(Sid,Rid)

	//フレンドDB登録でエラーならば
	if err != nil {	
		log.Println("errors")
		return "",err
	}

	log.Println("承認の取り消し")
	log.Println(UID)

	//フレンドリクエストをDBから消去
	dbconn.Delete(database.Sent{},database.Sent{UID: UID})

	return fuid,err
}

//フレンドリクエストを取り消し
func Delete_Request(UID string,Sender_id string) (error){

	//リクエストが存在しているか
	Sid,Rid,err := Get_Request(UID)
	_=Rid

	//リクエストでエラーならば
	if err != nil {	
		log.Println("error")
		return err
	}
	
	//送信者側が一致していない時
	if Sid != Sender_id {
		return errors.New("user_mismatch_existing")
	}

	log.Println("リクエストの取り消し")
	log.Println(UID)

	//フレンドリクエストをDBから消去
	dbconn.Delete(database.Sent{},database.Sent{UID: UID})

	return err
}

//フレンドリクエスト拒否
func Rejection(UID string,Receiver_id string) (error){

	//リクエストが存在しているか
	Sid,Rid,err := Get_Request(UID)
	_ = Sid

	//リクエストでエラーならば
	if err != nil {	
		log.Println("error")
		return err
	}

	//受信者側が一致していない時
	if Rid != Receiver_id {
		return errors.New("user_mismatch_existing")
	}

	log.Println("拒否")
	log.Println(UID)

	//フレンドリクエストをDBから消去
	dbconn.Delete(database.Sent{},database.Sent{UID: UID})
	return err
}



