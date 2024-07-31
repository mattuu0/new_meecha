const ignore_map = L.map('ignore_show_map', { zoomControl: false })

//レイヤー設定
/*
var tileLayer = L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© <a href="http://osm.org/copyright">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>',
});
*/

const ignore_distances_select = document.getElementById("ignore_distances_select");

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
tileLayer.addTo(ignore_map);

setView(34.70653432424858, 135.50369152261842, 12);

//マップを移動する
function setView(latitude, longitude, zoom = ignore_map.getZoom()) {
    //マップ設定
    ignore_map.setView([latitude, longitude], zoom);
}

//ピン一覧
let map_pins = {}
let selectid = "";

const addpin_btn = document.getElementById("addpin_btn");

addpin_btn.addEventListener("click", function (evt) {
    const point = ignore_map.getCenter();

    //ピンのIDを生成
    const pinid = crypto.randomUUID();

    //ピンを追加
    console.log(point);
    add_pin(pinid, point.lat, point.lng);
})

//ピンを追加する関数
function add_pin(add_pinid, lat, lng,distance = 3000) {
    //デフォルト距離

    //除外マーカー作成
    const ignore_point = L.marker([lat, lng], {
        draggable: true,
    });

    //範囲円を作成
    const ignore_circle = L.circle([lat, lng], {
        radius: distance,
    })

    //マーカーがドラッグされたとき
    ignore_point.on("drag", function (evt) {
        //マーカーを取得
        var marker = evt.target;

        //円を移動
        ignore_circle.setLatLng(marker._latlng);
    });

    ignore_point.on("dragend", function (evt) {
        //マーカーを取得
        var marker = evt.target;

        const point = marker._latlng;
    });

    //マーカー登録
    map_pins[add_pinid] = {
        "circle": ignore_circle,
        "point": ignore_point,
        "distance": distance,
    }

    function restore_old() {
        //もともと選択されているか
        if (selectid != "") {
            //青に戻す
            try {
                map_pins[selectid]["point"].setIcon(L.spriteIcon('blue'));
            } catch (error) {
                console.log(error);
            }
        }
    }

    //マーカーがクリックされたとき
    ignore_point.on("click", function (evt) {
        restore_old();

        //選択状態にする
        evt.target.setIcon(L.spriteIcon('red'));

        //選択されたマーカーのIDを設定
        selectid = evt.target.id;

        //選択されたマーカーの距離を設定
        ignore_distances_select.value = String(map_pins[selectid]["distance"]);
    })

    //マーカーのIDを設定
    ignore_point.id = add_pinid;

    //地図に追加
    ignore_map.addLayer(ignore_circle);
    ignore_map.addLayer(ignore_point);
}

//削除ボタン
const removepin_btn = document.getElementById("removepin_btn");

//削除ボタンを押したとき
removepin_btn.addEventListener("click", function (evt) {
    //選択されていないとき
    if (selectid == "") {
        //戻る
        return;
    }

    //選択されているマーカーを削除
    ignore_map.removeLayer(map_pins[selectid]["point"]);

    //マーカーに関連する円を削除
    ignore_map.removeLayer(map_pins[selectid]["circle"]);

    //登録を解除
    delete map_pins[selectid];
})

//距離設定が更新されたとき
ignore_distances_select.addEventListener("change", function (evt) {
    //選択されていなかったら戻る
    if (selectid == "") {
        return;
    }

    //マーカーが登録されていなかったら戻る
    if (map_pins[selectid] == undefined) {
        return;
    }

    //変更後の値を取得
    const setval = Number(ignore_distances_select.value);

    //登録情報を更新
    map_pins[selectid]["distance"] = setval;

    //円の半径を更新
    map_pins[selectid]["circle"].setRadius(setval);
})

const ignore_map_area = document.getElementById("ignore_map_area");
//保存ボタン
const savepin_btn = document.getElementById("savepin_btn");

savepin_btn.addEventListener("click",function(evt){
    ignore_map_area.style.visibility = "hidden";

    let result_list = [];

    //登録情報を取得
    for (const key in map_pins) {
        const point = map_pins[key]["point"].getLatLng();
        result_list.push({
            "lat": point.lat,
            "lng": point.lng,
            "distance": map_pins[key]["distance"]
        });
    }

    //保存
    AccessPost(save_ignore_point_url,{"points" : JSON.stringify(result_list)}).then(async (result) => {
        if (result.status == 200) {
            toastr["success"]("保存しました");
        } else {
            toastr["error"]("保存に失敗しました");
        }
    }).catch((err) => {
        console.log(err);
    });;
})