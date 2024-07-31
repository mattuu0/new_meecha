function notify_init() {
    //Flutter環境を確認
    if (!Is_App()) {
        return;
    }

    //Flutter環境
    is_flutter = true;

    try {
        try {
            //WebSocket切断
            wsconn.close();
        } catch (error) {
            console.log(error);
        }

        try {
            //位置情報の監視をやめる
            clear_watch();
        } catch (error) {
            console.log(error);
        }

        //アクセストークンを通知
        window.flutter_inappwebview.callHandler('web_inited', JSON.stringify({
            "token": get_access_token(),                        //アクセストークン
            "wsurl": "wss://" + ServerIp + "/ws",              //WebSocketのURL
        }));
    } catch (error) {
        console.log(error);
    }
}

//WebSocket切断
function stop_ws() {
    //Flutter環境を確認
    if (!Is_App()) {
        return;
    }
    
    try {
        //アクセストークンを通知
        window.flutter_inappwebview.callHandler('stop_ws', "");
    } catch (error) {
        console.log(error);
    }
}

//Flutterなら true を返す
function Is_App() {  
    return window.flutter_inappwebview != undefined
}