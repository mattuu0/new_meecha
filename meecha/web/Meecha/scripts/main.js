//トースター初期化
toastr.options = {
    "closeButton": true,
    "debug": false,
    "newestOnTop": true,
    "progressBar": true,
    "positionClass": "toast-top-center",
    "preventDuplicates": false,
    "onclick": null,
    "showDuration": "300",
    "hideDuration": "1000",
    "timeOut": "3000",
    "extendedTimeOut": "1000",
    "showEasing": "swing",
    "hideEasing": "linear",
    "showMethod": "fadeIn",
    "hideMethod": "fadeOut"
}

//マップ設定
const main_map = L.map('show_map')

var tileLayer = L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© <a href="http://osm.org/copyright">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>',
});


/*
var tileLayer = L.tileLayer('https://{s}.google.com/vt/lyrs=s,h&x={x}&y={y}&z={z}', {
	subdomains: ['mt0','mt1','mt2','mt3'],
	attribution: "<a href='https://developers.google.com/maps/documentation' target='_blank'>Google Map</a>",
});
*/

tileLayer.addTo(main_map);
//各ユーザーのピン
var friends_pins = {};
var myself_marker = null;
var myself_position = [0,0];
var myself_init = false;

//情報を送るインターバル
var interval_id = null;

//送信済みリクエスト
var recved_requests = {}

//Websockeet
let ws_conn;
var ws_connected = false;

function connect_ws() {
    //Websocketに接続
    //ws_conn = new WebSocket('ws://' + window.location.host + '/ws/connect');
    ws_conn = new WebSocket('wss://' + window.location.host + '/ws');

    ws_conn.onmessage = function(evt) {
        const parse_data = JSON.parse(evt.data);
        
        switch (parse_data.msgtype) {
            case "server_msg":
                let msg_data = parse_data.data;

                console.log(msg_data.msgtype);
                switch (msg_data.msgtype) {
                    case "accepted_friend_request":
                        toastr["info"]("フレンドリクエストが承認されました");
                        break;
                    case "near_friend_notify":
                        show_near_friend(msg_data);
                        break;
                    case "remove_map_pin":
                        remove_mappin(msg_data.userid);
                        break;
                    default:
                        console.log(msg_data);
                        break;
                }
                break;
        }
    };
        
    ws_conn.onclose = function(evt) {
        console.error('socket closed unexpectedly');
        ws_connected = false;
        
        toastr["warning"]("サーバーとの接続が切断しました、再接続するにはリロードしてください","通知",{disableTimeOut: true, closeButton:true,timeOut : "0",extendedTimeOut : "0"});
    };

    ws_conn.onopen = function(evt) {
        ws_connected = true;
        console.log("接続しました");
        toastr["info"]("接続しました","通知");
    }
}

window.addEventListener("load",function(evt){
    connect_ws();
})

//サーバーにコマンドを送信する
function send_command(command,data) {
    if (ws_connected) {
        var packet = {
            "command":command,
            "data":data
        }
        
        var send_data = JSON.stringify(packet);

        ws_conn.send(send_data);
    }
}

//マップを追跡させるか
const auto_change_map_check = document.getElementById("auto_change_map");


function init(evt) {
    //オブジェクト取得
}

window.onload = init;

//ユーザー検索関連
function clear_child_elems(elem) {
    //結果を削除する
    while (elem.lastChild) {
        elem.removeChild(elem.lastChild);
    }
}

//ピンを削除する
function remove_mappin(userid) {
    if (userid in friends_pins) {
        var friend_data = friends_pins[userid];

        main_map.removeLayer(friend_data["pin"]);

        delete friends_pins[userid];
    }
}

//近くのフレンドを表示する
function show_near_friend(data) {
    if (data["userid"] in friends_pins) {
        var friend_data = friends_pins[data["userid"]];

        console.log(data["geodata"][0],data["geodata"][1]);
        friend_data["pin"].setLatLng([data["geodata"][0],data["geodata"][1]]);
    } else {
        var friend_data = {};

        var friend_pin = L.marker([data["geodata"][0],data["geodata"][1]],{icon: L.spriteIcon('red')}).addTo(main_map).bindPopup(data["username"]);

        friend_data["pin"] = friend_pin;
        
        friends_pins[data["userid"]] = friend_data;
    }

    if (data["show_notify"]) {
        toastr["info"](data["username"] + "さんが近くにいます！","通知",{disableTimeOut: true, closeButton:true,timeOut : "0",extendedTimeOut : "0"});
        if(window.navigator.vibrate){
            window.navigator.vibrate([200]);
        }else if(window.navigator.mozVibrate){
            window.navigator.mozVibrate([200]);
        }else if(window.navigator.webkitVibrate){
            window.navigator.webkitVibrate([200]);
        }else{
            
        }
    }
}

//受信したフレンドリクエストを表示する
function show_recving_request(result) {
    toastr["info"]("フレンドリクエストを受信しました","通知",{disableTimeOut: true, closeButton:true,timeOut : "0",extendedTimeOut : "0"});
    window.navigator.vibrate(200); 
}


//位置情報取得
var id, target, options;

function success(pos) {
    var crd = pos.coords;

    
    //マーカーが設定されていなかったら現在地に打つ
    if (myself_marker == null) {
        main_map.setView([crd.latitude,crd.longitude],20);

        var popup = L.popup();
        myself_marker = L.marker([crd.latitude,crd.longitude]).addTo(main_map).on('click', function (e) {
            popup
            .setLatLng(e.latlng)
            .setContent("あなたの現在地")
            .openOn(main_map);
        });
    } else {
        //マーカーがあったら移動する
        myself_marker.setLatLng([crd.latitude,crd.longitude])
    }

    if (auto_change_map_check.checked) {
        main_map.setView([crd.latitude,crd.longitude],main_map.getZoom());
    }

    myself_position[0] = crd.latitude;
    myself_position[1] = crd.longitude;

    if (myself_init) {
        return;
    }

    myself_init = true;
    
    interval_id = setInterval(() => {
        post_pos();
    }, 3000);
}

function post_pos() {
    send_command("post_location",{
        latitude : myself_position[0],
        longitude : myself_position[1]
    });
}

function error(err) {
    console.warn('ERROR(' + err.code + '): ' + err.message);
}

function clear_watch() {
    navigator.geolocation.clearWatch(id);
}

options = { 
    enableHighAccuracy: true,
    timeout: 5000,
    maximumAge: 5000
};

function start_gps() {
    id = navigator.geolocation.watchPosition(success, error, options);
}

start_gps();

screen.orientation.lock();
