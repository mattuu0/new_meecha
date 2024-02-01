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
    async function search_user(evt){
        //検索
        const req = await AccessPost(friend_search_url,{},{"username":search_value.value});

        //200以外
        if (req.status != 200) {
            console.log("ユーザー検索に失敗しました");
            return;
        }

        //検索結果
        const result = await req.json();

        //検索結果を表示
        console.log(result);
        clear_child_elems(search_result_area);
        add_search_result(result.uid,result.name);
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

//長さを選ぶポップアップ表示
function show_setting_distance(evt) {
    setting_distance_area.classList.toggle("is-show");
}

//検索ポップアップ表示
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
async function get_sended_friend_request(evt) {
    await get_sent_requests();

    sended_request_area.classList.toggle("is-show");
}

//フレンド一覧を取得する
async function get_friends(evt) {
    await get_all_friends();

    pupup_friends_show_area.classList.toggle("is-show");

    friend_show_view.style.display = "absolute";
}

//ログアウトリンク取得
const logout_link = document.getElementById("logout_link");

//ログアウト関数
async function logout(evt) {
    //イベントキャンセル
    evt.preventDefault();

    //ログアウト送信
    const req = await RefreshPost(logout_url,{});

    //200以外
    if (req.status != 200) {
        alert("ログアウトに失敗しました");
        return;
    }

    //トークン削除
    delete_token();

    //ログイン画面へ
    window.location.href = "./login.html";
}

//イベント登録
logout_link.addEventListener("click",logout);

function add_search_result(uid,name) {
    //検索結果を表示
    const result_div = document.createElement("div");
    result_div.classList.add("search_result");
    
    const dirty = `
        <img class="search_result_icon" src="${GetIconUrl(uid)}">
        <p class="search_result_name">${name}</p>
        <button class="send_request_button" id="${uid}">送信</button>
    `

    const clean = DOMPurify.sanitize(dirty, { USE_PROFILES: { html: true } });
    //検索結果を表示
    result_div.insertAdjacentHTML("beforeend",clean);

    //送信ボタン取得
    const request_btn = result_div.querySelector(".send_request_button");
    
    request_btn.addEventListener("click",async function(evt) {
        //ID取得
        const sendid = evt.target.id;

        await AccessPost(send_request_url,{"Targetid":sendid});
    })

    search_result_area.appendChild(result_div);
}

//ユーザー名
const user_name = document.getElementById("user_name");

//ユーザアイコンエリア
const user_icon = document.getElementById("user_icon");

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

        const userid = userinfo["userid"];
        //アイコンURL
        user_icon.src = GetIconUrl(userid);

        //ユーザ名設定
        user_name.textContent = userinfo["name"];

        //ユーザid設定
        UserID = userinfo["userid"];
    } catch (error) {
        //エラー処理
        console.log(error);

        //ログインに戻る
        window.location.href = "./login.html";
        return;
    }
}

//ロード完了イベント
window.addEventListener("load",function(evt) {
    //ゆーざー情報取得
    get_userinfo();
})