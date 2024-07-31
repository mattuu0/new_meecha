
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

        //ユーザid設定
        UserID = userinfo["userid"];

        navigator.geolocation.getCurrentPosition(
            (position) => {
                //成功時のコールバック関数
                //位置情報
                const latitude = position.coords.latitude;
                const longitude = position.coords.longitude;

                console.log(latitude, longitude);

                change_location(latitude, longitude);

                start_gps(null);
            },
            (error) => {
                //失敗時のコールバック関数
                console.log(error);
                start_gps(null);
            },
            //位置情報監視設定
            {}
        );

        
    } catch (error) {
        //エラー処理
        console.log(error);

        //ログインに戻る
        window.location.href = "./login.html";
        return;
    }
}

get_userinfo();

//フレンドのピン
let friend_pins = {};

function set_marker(uid, name, latitude, longitude) {
    //既に存在する場合移動させる
    if (friend_pins[uid]) {
        friend_pins[uid].setLatLng([latitude, longitude]);
        return;
    }
    //マーカーを設定
    var friend_icon = L.icon({
        iconUrl: GetIconUrl(uid),
        iconSize: [50, 50],
        iconAnchor: [37, 75],
        popupAnchor: [0, -70],
        className: "MapIcon",
    });

    const friend_marker = L.marker([latitude, longitude], {
        icon: friend_icon,
    }).bindPopup(`<p>${name}さん</p>`);

    main_map.addLayer(friend_marker);

    friend_pins[uid] = friend_marker;
}

//イベント登録
ws_event_div.addEventListener(ws_event_key, function (evt) {
    const payload = evt.detail["Payload"];

    switch (evt.detail["Command"]) {
        case "near_friend": {
            if (payload["is_first"] && !payload["is_self"]) {
                toastr.info(`${payload["unane"]}さんが近くにいます`, "通知", {
                    "closeButton": false,
                    "debug": false,
                    "newestOnTop": true,
                    "progressBar": true,
                    "positionClass": "toast-top-center",
                    "preventDuplicates": true,
                    "showDuration": "300",
                    "timeOut": "0",
                    "showEasing": "swing",
                    "hideEasing": "linear",
                    "showMethod": "fadeIn",
                    "hideMethod": "fadeOut"
                });
            }

            const point_data = payload["point"];
            set_marker(payload["userid"], payload["unane"], point_data["Lat"], point_data["Lon"]);
            break;
        }
        case "stop_notify": {
            const userid = payload["userid"];

            try {
                main_map.removeLayer(friend_pins[userid]);
            } catch (error) {
                console.log(error);
            }

            //登録を解除
            delete friend_pins[userid];
            break;
        }
    };
})

//戻るピン
const pin = document.getElementById("pin");

//ピン押したとき
pin.addEventListener("click", () => {
    //マップ設定
    setView(myself_position[0], myself_position[1]);
})