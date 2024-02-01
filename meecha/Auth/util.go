package auth

import (
	_ "crypto/sha512"
	_ "encoding/hex"
	"errors"

	"github.com/google/uuid"
)

// IDを生成する
func Genid() (string, error) {

	//UUIDを生成する
	uuid_obj, err := uuid.NewRandom()

	//エラー処理
	if err != nil {
		return "", err
	}

	/*
	//UUID文字列をSHA512文字列にする
	hash_byte := sha512.Sum512([]byte(uuid_obj.String())) //128桁の文字列のid
	hex_string := hex.EncodeToString(hash_byte[:])        //16進数表記

	//文字列を返す
	return hex_string, nil
	*/

	return uuid_obj.String(),nil
}

// 初期化されていないエラーを返す
func Get_Init_Error() error {
	return errors.New("not initialized")
}
