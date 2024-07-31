//WebSocket
let wsconn = null;

//接続しているか
let ws_connected = false;

//Flutter Support
let is_flutter = false;

//イベント用div
const ws_event_div = document.createElement("div");

//イベント名
const ws_event_key = "ws_msg_event";

function send_command(Command,payload,seriarize = true) {
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
    //アプリの場合戻る
    if (Is_App()) {
        return
    }

    //Websocket 接続
    wsconn = new WebSocket("wss://" + ServerIp + "/ws");

    //接続時
    wsconn.onopen = function () {
        //認証コマンド
        send_command("auth",get_access_token(),false);
    }

    //メッセージが来たとき
    wsconn.onmessage = function (evt) {
        //シリアライズ
        const load_json = JSON.parse(evt.data);

        on_recved(load_json);
    }

    //切断時
    wsconn.onclose = function () {
        //Flutter 環境の場合戻る
        if (is_flutter) {
            return;
        }

        toastr.error('通信が切断されました<br>通知をタップすると再接続します',"通知",{
            "closeButton": true,
            "debug": false,
            "newestOnTop": true,
            "progressBar": true,
            "positionClass": "toast-top-center",
            "preventDuplicates": false,
            "onclick": function(evt){
                //再接続
                connect_ws();
            },
            "showDuration": "300",
            "timeOut": "0",
            "extendedTimeOut": "1000",
            "showEasing": "swing",
            "hideEasing": "linear",
            "showMethod": "fadeIn",
            "hideMethod": "fadeOut"
        });
        //接続済み解除  
        ws_connected = false;
    }
}

function on_recved(data) {
    //コマンドに応じて処理
    switch (data.Command) {
        case "Auth_Complete":
            //接続済みにする
            ws_connected = true;
            toastr["success"]("接続しました", "通知")
            break;
        case "Location_Token":
            //トークンが来たときに位置情報を送る
            send_command("location",{
                "token":data.Payload,
                "lat" : myself_position[0],
                "lng" : myself_position[1],
            });
            break;
        case "recv_request":
            toastr.info(`${data.Payload}さんからフレンドリクエストを受信しました`, "通知",{
                "closeButton": false,
                "debug": false,
                "newestOnTop": true,
                "progressBar": true,
                "positionClass": "toast-top-center",
                "preventDuplicates": false,
                "showDuration": "300",
                "timeOut": "0",
                "showEasing": "swing",
                "hideEasing": "linear",
                "showMethod": "fadeIn",
                "hideMethod": "fadeOut"
            })
            break;
        case "accept_request":
            toastr.info(`${data.Payload}さんとフレンドになりました`, "通知",{
                "closeButton": false,
                "debug": false,
                "newestOnTop": true,
                "progressBar": true,
                "positionClass": "toast-top-center",
                "preventDuplicates": false,
                "showDuration": "300",
                "timeOut": "0",
                "showEasing": "swing",
                "hideEasing": "linear",
                "showMethod": "fadeIn",
                "hideMethod": "fadeOut"
            })
            break;
        case "Notify_Disconnect":
            switch (data.Payload["code"]) {
                case "409":
                    //再読み込み
                    window.location.reload();
                    break;
            }
            break;
        default:
            //イベント作成
            const event = new CustomEvent(ws_event_key, { detail: data });
            ws_event_div.dispatchEvent(event);
            break;
    }
}

window.addEventListener("load", function (evt) {
    connect_ws();
    notify_init();
})
