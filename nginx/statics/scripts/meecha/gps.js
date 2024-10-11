//初期化されているか
let myself_init = false;

//自分のピン
let myself_marker = null;

//自分の位置情報データ
let myself_position = [0, 0];

//位置情報監視ID
let watchid = null;

//イベント名
const event_name = "geo_change";
//イベント用の要素作成
const event_elem = document.createElement("div");

function call_event(latitude,longitude) {
    myself_position[0] = latitude;
    myself_position[1] = longitude;
    //イベント作成
    const event = new CustomEvent(event_name, {
        detail: {
            latitude: latitude,
            longitude: longitude
        }
    });

    //イベント発火
    event_elem.dispatchEvent(event);
}

function success(crd) {
    const latitude = crd.coords.latitude;
    const longitude = crd.coords.longitude;

    console.log(latitude,longitude);
    //自身の位置を更新する
    myself_position[0] = latitude;
    myself_position[1] = longitude;

    call_event(latitude, longitude);
    //自身が初期化されていたら
    if (myself_init) {
        //戻る
        return;
    }

    myself_init = true;
}

//位置情報の監視をやめる
function clear_watch() {
    navigator.geolocation.clearWatch(watchid);
}

//エラー処理
function error(err) {
    console.log(err.message);
}

//位置情報監視設定
options = {};

//位置情報の監視開始
function start_gps(evt) {
    watchid = navigator.geolocation.watchPosition(success, error, options);
}
