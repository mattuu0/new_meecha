function notify_init() {
    try {
        try {
            //WebSocket切断
            wsconn.close();
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