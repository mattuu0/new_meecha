package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"


	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type CommandMessage struct {
	//コマンド
	Command string

	//ペイロード
	Payload interface{}
}

type ResponseMessage struct {
	//コマンド
	Command string

	//ペイロード
	Payload interface{}
}

func main() {
	//test20155

	//環境変数読み込み
	loadEnv()

	//1000回ループ
	for i := 0; i < 100000; i++ {
		log.Println("test" + strconv.Itoa(i))

		if true {
			url := "https://wao2server.tail6cf7b.ts.net/meecha/auth/signup"

			payload := strings.NewReader("{\n  \"name\" : \"" + "test" + strconv.Itoa(i) + "\",\n  \"password\" : \"password\"\n}")

			req, _ := http.NewRequest("POST", url, payload)

			req.Header.Add("Accept", "*/*")
			req.Header.Add("User-Agent", "Thunder Client (https://www.thunderclient.com)")
			req.Header.Add("Content-Type", "application/json")

			res, _ := http.DefaultClient.Do(req)

			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)

			fmt.Println(res)
			fmt.Println(string(body))

			/*
			user, err := auth.Login("test"+strconv.Itoa(i), "password")

			if err != nil {
				log.Println(err)
				continue
			}

			fid, err := friends.Record_Friends(user.Userid, "0f11d95c-a600-4ca5-9bd4-9248b5dba35b")

			if err != nil {
				log.Println(err)
				continue
			}

			log.Println(fid)
			*/
		}
	}

}
