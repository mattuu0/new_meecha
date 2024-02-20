package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"math/rand"

	"github.com/gorilla/websocket"

	"time"
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

func random(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()*(max-min) + min
}

var (
	now_count int64 = 0
)

func main() {
	//環境変数読み込み
	loadEnv()

	scanner := bufio.NewScanner(os.Stdin)
	//1000000人のユーザを作成
	for i := 0; i < 20000; i++ {
		if i%50 == 0 {
			log.Println(now_count)
			log.Println(i)
			log.Println("wait key")
			scanner.Scan()
			scanner.Text()
		}

		time.Sleep(time.Duration(200) * time.Millisecond)
		log.Println(now_count)
		now_count++
		go send_location(i)
	}

	log.Println(now_count)
	log.Println("wait key")
	scanner.Scan()
	scanner.Text()

	log.Println(login("wao", "password"))
}

type Response struct {
	AccessToken  string
	RefreshToken string
	message      string
}

// https://wao2server.tail6cf7b.ts.net:13333/static/meecha/index.html
func login(uname, password string) string {
	url := "https://wao2server.tail6cf7b.ts.net/meecha/auth/login"

	payload := strings.NewReader("{\n  \"name\" : \"" + uname + "\",\n  \"password\" : \"" + password + "\"\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", "Thunder Client (https://www.thunderclient.com)")
	req.Header.Add("Content-Type", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	sslcli := &http.Client{Transport: tr}

	res, err := sslcli.Do(req)

	if err != nil {
		log.Println(err)
		return ""
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var resdata Response
	json.Unmarshal(body, &resdata)

	return resdata.AccessToken
}

func send_location(loginid int) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}

		now_count--

		log.Println(now_count)
	}()

	atoken := login("test"+strconv.Itoa(loginid), "password")

	wsurl := url.URL{Scheme: "wss", Host: "wao2server.tail6cf7b.ts.net", Path: "/meecha/ws"}
	//log.Printf("connecting to %s", wsurl.String())

	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	wsconn, _, err := dialer.Dial(wsurl.String(), nil)

	if err != nil {
		return
	}

	defer wsconn.Close()

	wsconn.WriteJSON(CommandMessage{Command: "auth", Payload: atoken})

	for {
		var recvmsg = ResponseMessage{}
		err = wsconn.ReadJSON(&recvmsg)

		if err != nil {
			log.Println("error : " + err.Error())
			wsconn.Close()
			return
		}

		switch recvmsg.Command {
		case "Location_Token":
			//log.Println(loginid)
			send_data := CommandMessage{
				Command: "location",
				Payload: map[string]interface{}{
					"token": recvmsg.Payload.(string),
					"lat":   34.7088768 + random(-0.005, 0.005),
					"lng":   135.4969214 + random(-0.005, 0.005),
				},
			}

			err = wsconn.WriteJSON(send_data)

			if err != nil {
				log.Println("error : " + err.Error())
				wsconn.Close()
				return
			}

			break

		case "near_friend":
			//log.Println(recvmsg.Payload)
			break
		}
	}
}
