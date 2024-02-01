//自動追尾設定
let auto_tracking = true;

//初期化されているか
let myself_init = false;

//自分のピン
let myself_marker = null;

//自分の位置情報データ
let myself_position = [0, 0];

//位置情報監視ID
let watchid = null;

//マップ設定
const main_map = L.map('show_map', { zoomControl: false })

//レイヤー設定
/*
var tileLayer = L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© <a href="http://osm.org/copyright">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>',
});
*/

/*
var tileLayer = L.tileLayer('https://{s}.google.com/vt/lyrs=s,h&x={x}&y={y}&z={z}', {
    subdomains: ['mt0', 'mt1', 'mt2', 'mt3'],
    attribution: "<a href='https://developers.google.com/maps/documentation' target='_blank'>Google Map</a>",
});

//マップに追加
tileLayer.addTo(main_map);

//デバッグ
if (debug) {
    //位置設定
    setView(34.70653432424858, 135.50369152261842, 20);
}

*/

//マップを移動する
function setView(latitude, longitude, zoom = main_map.getZoom()) {
    //マップ設定
    main_map.setView([latitude, longitude], zoom);
}

//位置情報取得イベント
function success(pos) {
    //位置情報
    var crd = pos.coords;

    try {
        //マーカーが設定されていなかったら現在地に打つ
        if (myself_marker == null) {
            setView(crd.latitude, crd.longitude, 20);

            if (UserID == "") {
                return;
            }

            //マーカーを設定
            var myicon = L.icon({
                iconUrl: GetIconUrl(UserID),
                iconSize: [50, 50],
                iconAnchor: [37, 75],
                popupAnchor: [0, -70],
                className: "MapIcon",
            });

            myself_marker = L.marker([crd.latitude, crd.longitude], { icon: myicon })
                .bindPopup('<p>あなたの現在地</p>');
            main_map.addLayer(myself_marker);
        } else {
            //マーカーがあったら移動する
            myself_marker.setLatLng([crd.latitude, crd.longitude]);
        }

        //トラッキング
        if (auto_tracking) {
            setView(crd.latitude, crd.longitude);
        }
    } catch (err) {
        console.log(err);
    }

    //自身の位置を更新する
    myself_position[0] = crd.latitude;
    myself_position[1] = crd.longitude;

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
options = {
    enableHighAccuracy: true,
    timeout: 5000,
    maximumAge: 5000
};

//位置情報の監視開始
function start_gps(evt) {
    watchid = navigator.geolocation.watchPosition(success, error, options);
}

start_gps(null);