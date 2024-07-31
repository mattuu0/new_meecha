//自動追尾設定
let auto_tracking = false;

//マップ設定
const main_map = L.map('show_map', { zoomControl: false })

//レイヤー設定

var tileLayer = L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© <a href="http://osm.org/copyright">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>',
});

/*
var tileLayer = L.tileLayer('https://{s}.google.com/vt/lyrs=s,h&x={x}&y={y}&z={z}', {
    subdomains: ['mt0', 'mt1', 'mt2', 'mt3'],
    attribution: "<a href='https://developers.google.com/maps/documentation' target='_blank'>Google Map</a>",
});
*/

//マップに追加
tileLayer.addTo(main_map);

//マップを移動する
function setView(latitude, longitude, zoom = main_map.getZoom()) {
    //マップ設定
    main_map.setView([latitude, longitude], zoom);
}

//位置情報取得イベント
function change_location(latitude, longitude) {
    try {
        //マーカーが設定されていなかったら現在地に打つ
        if (myself_marker == null) {
            setView(latitude, longitude, 20);

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

            myself_marker = L.marker([latitude, longitude], {
                icon: myicon,
            }).bindPopup('<p>あなたの現在地</p>');

            main_map.addLayer(myself_marker);
        } else {
            //マーカーがあったら移動する
            myself_marker.setLatLng([latitude, longitude]);
        }

        
    } catch (err) {
        console.log(err);
    }
}

event_elem.addEventListener(event_name, function (evt) {
    change_location(evt.detail.latitude, evt.detail.longitude);
})