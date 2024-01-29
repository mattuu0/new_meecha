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

function change_distance(evt) {
    const num = evt.target.selectedIndex;
	const select_val = evt.target.options[num].value;

    send_command("update_distance",{"distance" : select_val});
}

function show_distance(distance) {
    distance_show_area.textContent = "現在は" + distance + "mで通知します"
}

//設定
var recved_requests = {}

let select_box = document.getElementById("notify_distances_select");
let distance_show_area = document.getElementById("distance_show");

select_box.addEventListener("change",change_distance);

//popup.classList.toggle('is-show');
//追加機能

//検索ボタン
const search_button = document.getElementById("user_search_button");

//検索するユーザー名
const search_value = document.getElementById("search_username_value");

//ユーザー検索結果
const search_result_area = document.getElementById("user_searh_result");

//受信済みリクエスト表示場所
const recved_request_show_area = document.getElementById("recved_request_show_area");

//受信済みリクエスト取得ボタン
const get_recved_request_button = document.getElementById("get_request_button");

//受信済みリクエスト表示エリア
const recved_request_area = document.getElementById("show_recved_request_area");

//送信済みリクエスト取得ボタン
const get_sended_request_button = document.getElementById("sended_request_button");

//送信済みリクエスト表示
const sended_request_area = document.getElementById("show_sended_request_area");

//送信済みリクエスト表示場所
const sended_request_show_area = document.getElementById("sended_request_show_area");

//フレンド取得ボタン
const get_friends_btn = document.getElementById("get_friends_button");

//フレンド一覧表示場所
const pupup_friends_show_area = document.getElementById("pupup_friends_show_area");

//フレンド表示場所
const friend_show_area = document.getElementById("friends_show_area");

//フレンドビュー
const friend_show_view = document.getElementById("friends_area");

//マップを追跡させるか
const clear_search_btn = document.getElementById("search_clear_btn");

//通知距離表示ボタン
const setting_notify_distance = document.getElementById("setting_notify_distance");

//通知距離表示エリア
const setting_distance_area = document.getElementById("setting_distance_area");

//ユーザー検索ボタン
const search_user_button = document.getElementById("search_user_button");

//検索エリア
const pupup_search_area = document.getElementById("pupup_search_area");

function init(evt) {
    //オブジェクト取得
    
    //イベント関連
    function search_user(evt){
        send_command("search_user",{username : search_value.value});
    }

    //イベント登録
    search_button.addEventListener("click",search_user);
    get_recved_request_button.addEventListener("click",get_friend_request);
    get_sended_request_button.addEventListener("click",get_sended_friend_request);
    get_friends_btn.addEventListener("click",get_friends);
    clear_search_btn.addEventListener("click",clear_friend_search);
    setting_notify_distance.addEventListener("click",show_setting_distance);
    search_user_button.addEventListener("click",show_pupup_search_area)
}

window.onload = init;

function show_setting_distance(evt) {
    setting_distance_area.classList.toggle("is-show");
}

function show_pupup_search_area(evt) {
    pupup_search_area.classList.toggle("is-show");
}


//ユーザー検索関連
function clear_child_elems(elem) {
    //結果を削除する
    while (elem.lastChild) {
        elem.removeChild(elem.lastChild);
    }
}

//受信済みフレンドリクエストを表示する
function show_recved_requests(result) {
    clear_child_elems(recved_request_show_area);
}

//取得した送信済みフレンドリクエストを表示する
function show_sended_friend_requests(result) {
    clear_child_elems(sended_request_show_area);
}

//フレンドを表示する
function show_friend(result) {
    clear_child_elems(friend_show_area);
}

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

//フレンド検索結果を削除する
function clear_friend_search(evt) {
    clear_child_elems(search_result_area);
}


//受信済みフレンドリクエストを取得する
function get_friend_request(evt) {
    recved_request_area.classList.toggle("is-show");
}

//送信済みフレンドリクエストを取得する
function get_sended_friend_request(evt) {
    sended_request_area.classList.toggle("is-show");
}

//フレンド一覧を取得する
function get_friends(evt) {
    pupup_friends_show_area.classList.toggle("is-show");

    friend_show_view.style.display = "absolute";
}