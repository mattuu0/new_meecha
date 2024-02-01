//ユーザアイコン
const usericon = document.getElementById("usericon");

//ユーザ情報取得処理
async function get_userinfo() {
    try {
        //りくえすと　
        const req = await AccessPost(uinfo_url, {});

        //403の時
        if (req.status == 403) {
            //ログインに戻る
            window.location.href = "./login.html";
            return;
        }

        //ユーザデータ取得
        const userinfo = await req.json();

        //アイコンURL
        usericon.src = GetIconUrl(userinfo["userid"]);

        //ユーザid設定
        UserID = userinfo["userid"];
    } catch (error) {
        //エラー処理
        console.log(error);

        //ログインに戻る
        window.location.href = "./login.html";
        return;
    }
}

get_userinfo();

//WebSocket
let wsconn = null;

//接続しているか
let ws_connected = false;

function send_command(Command,payload,seriarize = true) {
    //送信するデータ
    let send_payload = payload;

    //シリアライズするか
    if (seriarize) {
        //Jsonをシリアライズ
        send_payload = JSON.stringify(payload);
    }

    //リクエストを飛ばす
    wsconn.send(JSON.stringify({
        "Command": Command,
        //アクセストークン
        "Payload": payload,
    }))
}

//WebSocket接続
function connect_ws() {
    wsconn = new WebSocket("ws://" + ServerIp + "/ws");

    wsconn.onopen = function () {
        //認証コマンド
        send_command("auth",get_access_token(),false);
    }

    //メッセージが来たとき
    wsconn.onmessage = function (evt) {
        //JSONに変換
        const load_json = JSON.parse(evt.data);

        //コマンドに応じて処理
        switch (load_json.Command) {
            case "Auth_Complete":
                //接続済みにする
                ws_connected = true;
                break;
            case "Location_Token":
                //トークンが来たときに位置情報を送る
                send_command("location",{
                    "token":load_json.Payload,
                    "lat" : myself_position[0],
                    "lng" : myself_position[1],
                });
                break;
            default:
                console.log(load_json);
                break;
        }
    }

    //切断時
    wsconn.onclose = function () {
        console.log("close");

        //接続済み解除
        ws_connected = false;
    }
}

window.addEventListener("load", function (evt) {
    connect_ws();
})

