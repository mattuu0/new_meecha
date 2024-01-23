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

                switch (msg_data.msgtype) {
                    case "success_update_distance":
                        toastr["info"]("距離設定を更新しました");

                        show_distance(msg_data.distance);
                        break;
                    case "now_distance":
                        show_distance(msg_data.distance);
                        break;
                    case "request_accpet_error":
                        toastr["error"]("承認に失敗しました");
                        break;
                    case "search_response":
                        show_search_result(msg_data.match_users);
                        break;
                    case "recv_friend_request":
                        show_recving_request(msg_data.request_data);
                        break;
                    case "recved_friend_request":
                        show_recved_requests(msg_data.recved_requests);
                        break;
                    case "sended_friend_request":
                        show_sended_friend_requests(msg_data.sended_requests);
                        break;
                    case "friends_response":
                        show_friend(msg_data.friends);
                        break;
                    case "request_already_sended":
                        request_already_sended(msg_data.send_id);
                        break;
                    case "already_friend":
                        toastr["error"]("既にフレンドです");
                        break;
                    case "success_update_memo":
                        toastr["info"]("メモを更新しました");
                        break;
                    case "success_remove_friend":
                        toastr["info"]("フレンドを削除しました");
                        break;
                    case "accepted_friend_request":
                        toastr["info"]("フレンドリクエストが承認されました");
                        break;
                    case "success_reject_request":
                        toastr["info"]("フレンドリクエストを拒否しました");
                        break;
                    case "accept_friend_request":
                        toastr["info"]("フレンドリクエストを承認しました");
                        success_accept_request(msg_data);
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
        toastr["info"]("接続しました","通知");
        
        send_command("get_now_distance",{});
    }
}

connect_ws();

function change_distance(evt) {
    const num = evt.target.selectedIndex;
	const select_val = evt.target.options[num].value;

    send_command("update_distance",{"distance" : select_val});
}

function show_distance(distance) {
    distance_show_area.textContent = "現在は" + distance + "mで通知します"
}

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

//受信したフレンドリクエストを表示する
function show_recving_request(result) {
    toastr["info"]("フレンドリクエストを受信しました","通知",{disableTimeOut: true, closeButton:true,timeOut : "0",extendedTimeOut : "0"});
    window.navigator.vibrate(200); 
    send_command("get_sended_request",{});
    //show_friend_request(result.sender_id,result.username,result.requestid,recved_request_show_area);
}

//受信済みフレンドリクエストを表示する
function show_recved_requests(result) {
    clear_child_elems(recved_request_show_area);

    for (let requestid in result) {
        try {
            delete_sended_request(requestid);
        } catch (ex) {
            console.log(ex.message)
        }
        //送信者ID
        let senderid = result[requestid].sender_userid;  
        
        //送信者名
        let sender_name = result[requestid].user_name;

        show_friend_request(senderid,sender_name,requestid,recved_request_show_area);
    }
}
 
//既に送信済みの場合エラーを出す
function request_already_sended(result) {
    toastr["error"]("フレンドリクエストを既に送信しています");
}

//フレンドリクエストを表示する
function show_friend_request(senderid,username,requestid,showdiv) {
    //結果のdiv
    let add_div = document.createElement("div");
    add_div.classList.add("recved_request_area");

    /*ID表示
    let sender_userid_area = document.createElement("p");
    sender_userid_area.textContent = "ID : " + senderid;

    //追加
    add_div.appendChild(sender_userid_area);
    */

    //ID表示
    let username_area = document.createElement("p");
    username_area.textContent = "ユーザー名 : " + username;

    //追加
    add_div.appendChild(username_area);

    //承認ボタン
    let accept_btn = document.createElement("input");
    accept_btn.type = "button";
    accept_btn.value = "承認";
    accept_btn.requestid = requestid;
    accept_btn.base_div = add_div;
    accept_btn.classList.add("meecha_button");

    //イベント登録
    accept_btn.addEventListener("click",accept_friend_request);
    add_div.append(accept_btn);

    //フレンドリクエストボタン
    let reject_btn = document.createElement("input");
    reject_btn.type = "button";
    reject_btn.value = "拒否";
    reject_btn.requestid = requestid;
    reject_btn.base_div = add_div;
    reject_btn.classList.add("meecha_button");

    //イベント登録
    reject_btn.addEventListener("click",reject_friend_request);
    add_div.append(reject_btn);

    //承認コード入力
    let verify_code_area = document.createElement("input");
    verify_code_area.type = "text";
    verify_code_area.base_div = add_div;
    verify_code_area.placeholder = "承認コード"

    add_div.append(verify_code_area);

    //ボタンなどを追加する
    showdiv.appendChild(add_div);

    try {
        recved_requests[String(requestid)] = {div : add_div,verify_area : verify_code_area};
    } catch (ex) {
        console.log(ex.message);
    }
}


//取得した送信済みフレンドリクエストを表示する
function show_sended_friend_requests(result) {
    clear_child_elems(sended_request_show_area);

    for (let requestid in result) {
        let verify_code = result[requestid].verify_code;
        let username = result[requestid].user_name;

        //結果のdiv
        let add_div = document.createElement("div");
        add_div.classList.add("sended_request_area");
        
        //ID表示
        let username_area = document.createElement("p");
        username_area.textContent = "ユーザー名 : " + username;

        //追加
        add_div.appendChild(username_area);


        //承認コード表示
        let verify_code_userid_area = document.createElement("p");
        verify_code_userid_area.textContent = "承認コード : " + verify_code;
        add_div.appendChild(verify_code_userid_area);

        //承認ボタン
        let cancel_btn = document.createElement("input");
        cancel_btn.type = "button";
        cancel_btn.value = "取り消し";
        cancel_btn.requestid = requestid;
        cancel_btn.base_div = add_div;
        cancel_btn.classList.add("meecha_btn");

        //イベント登録
        cancel_btn.addEventListener("click",cancel_friend_request);
        add_div.append(cancel_btn);

        sended_request_show_area.appendChild(add_div);
    }
}

//フレンドを表示する
function show_friend(result) {
    clear_child_elems(friend_show_area);

    result.forEach((val,index) => {
        //結果のdiv
        let add_div = document.createElement("div");
        add_div.classList.add("friend_info");

        //ID表示
        let username_area = document.createElement("p");
        username_area.textContent = "ユーザー名 : " + val["friend_username"];

        //追加
        add_div.appendChild(username_area);

        //削除ボタン
        let remove_btn = document.createElement("input");
        remove_btn.type = "button";
        remove_btn.value = "フレンド削除";
        remove_btn.friendid = val["friendid"];
        remove_btn.base_div = add_div;
        remove_btn.classList.add("meecha_btn");

        //イベント登録
        remove_btn.addEventListener("click",remove_friend);
        add_div.append(remove_btn);

        //メモ編集エリア
        let friend_memo_area = document.createElement("div");
        friend_memo_area.classList.add("friend_memo");

        //フレンドリクエストボタン
        let memo_area = document.createElement("textarea");
        memo_area.friendid = val["friendid"];
        memo_area.friend_userid = val["friend_userid"];
        memo_area.value = val["friend_memo"];
        memo_area.classList.add("memo_textarea");
        friend_memo_area.append(memo_area);

        //フレンドリクエストボタン
        let memo_btn = document.createElement("input");
        memo_btn.type = "button";
        memo_btn.value = "メモ更新";
        memo_btn.friendid = val["friendid"];
        memo_btn.friend_userid = val["friend_userid"];
        memo_btn.memo_text_area = memo_area;
        memo_btn.classList.add("meecha_btn");
        //イベント登録
        memo_btn.addEventListener("click",update_memo);
        //friend_memo_area.append(memo_btn);
        
        add_div.append(memo_btn);

        add_div.appendChild(friend_memo_area);


        friend_show_area.appendChild(add_div);
    })
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

//フレンドリクエストを承認
function accept_friend_request(evt) {
    //ユーザーID
    let send_id = evt.target.requestid;
    let send_code = String(recved_requests[send_id]["verify_area"].value);

    send_command("accept_request",{requestid:send_id,verify_code:send_code});
}

//フレンドを削除
function remove_friend(evt) {
    //ユーザーID
    let send_id = evt.target.friendid;
    evt.target.base_div.remove();

    send_command("remove_friend",{friendid:send_id});
}


//メモ更新
function update_memo(evt) {
    //ユーザーID
    let send_id = evt.target.friend_userid;
    let memo_value = evt.target.memo_text_area;

    send_command("update_memo",{friendid:send_id,uodate_memo : memo_value.value});
}


//フレンドリクエストを拒否
function reject_friend_request(evt) {
    //ユーザーID
    let send_id = evt.target.requestid;
    evt.target.base_div.remove();

    send_command("reject_request",{requestid:send_id});
}

//フレンドリクエストをキャンセル
function cancel_friend_request(evt) {
    //ユーザーID
    let send_id = evt.target.requestid;
    evt.target.base_div.remove();

    send_command("cancel_request",{requestid:send_id});
}


//フレンド検索結果を削除する
function clear_friend_search(evt) {
    clear_child_elems(search_result_area);
}

//フレンドリクエストを承認したとき
function success_accept_request(accept_data) {
    delete_sended_request(accept_data["request_id"]);
}

function delete_sended_request(request_id) {
    var delete_div = recved_requests[request_id]["div"];
    delete_div.remove()

    delete recved_requests[request_id]
}

//検索結果表示
function show_search_result(result) {
    clear_child_elems(search_result_area);

    for (let userid in result) {
        //結果のdiv
        let add_div = document.createElement("div");

        let user_div = document.createElement("div");
        user_div.classList.add("accpet_area");
        //ID表示
        let username_area = document.createElement("p");
        username_area.textContent = "ユーザー名 : " + result[userid].user_name;

        //追加
        user_div.appendChild(username_area);
        
        let request_btn = document.createElement("input");
        request_btn.type = "button";
        
        request_btn.userid = userid;
        request_btn.value = "フレンドリクエスト";
        request_btn.classList.add("meecha_btn");

        //イベント登録
        if (result[userid].is_friend == "0") {
            //フレンドリクエストボタン
            request_btn.addEventListener("click",send_firend_req);
        } else {
            //フレンドならボタンを無効にする
            request_btn.disabled = true;
        }

        user_div.append(request_btn);

        add_div.append(user_div);
        search_result_area.appendChild(add_div);
    }
}

//フレンドリクエストを送る
function send_firend_req(evt) {
    //ユーザーID
    let send_id = evt.target.userid;

    send_command("friend_request",{userid:send_id});
}

//受信済みフレンドリクエストを取得する
function get_friend_request(evt) {
    recved_request_area.classList.toggle("is-show")
    send_command("get_recved_request",{});
}

//送信済みフレンドリクエストを取得する
function get_sended_friend_request(evt) {
    sended_request_area.classList.toggle("is-show");
    send_command("get_sended_request",{});
}

//フレンド一覧を取得する
function get_friends(evt) {
    pupup_friends_show_area.classList.toggle("is-show");

    send_command("get_friends",{});

    friend_show_view.style.display = "absolute";
}