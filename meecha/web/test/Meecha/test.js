const get_user_info = document.getElementById("get_user_info");

get_user_info.addEventListener("click", async function (evt) {
    const req = await AccessPost(server_url + "/user_info", {});

    console.log(await req.json())
})

const logout_btn = document.getElementById("logout_btn");

logout_btn.addEventListener("click", async function (evt) {
    const req = await RefreshPost(logout_url, {});

    console.log(req.status);

    console.log(await req.json());
})

const icon_upload = document.getElementById("icon_upload");

const upload_icon_btn = document.getElementById("upload_icon");

upload_icon_btn.addEventListener("click", async function (evt) {
    const updata = new FormData();
    updata.append("file", icon_upload.files[0]);

    console.log(updata.getAll("file"));

    const icon_post = await AccessPost(server_url + "/upicon", updata, {}, false);

    console.log(icon_post);
})

const refresh_btn = document.getElementById("refresh_btn");

refresh_btn.addEventListener("click", async function (evt) {
    const req = await RefreshPost(refresh_url, {});

    console.log(req.status);

    console.log(await req.json());
})

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
                send_command("location",load_json.Payload);
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

