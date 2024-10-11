async function change_distance(evt) {
    const num = evt.target.selectedIndex;
	const select_val = evt.target.options[num].value;
    
    //通知距離設定
    const req = await AccessPost(set_notify_distance_url,{"distance":select_val});

    if (req.status != 200) {
        toastr["error"]("通知距離設定に失敗しました");
        return;
    }

    toastr["success"]("通知距離を設定しました");
    show_distance(select_val);
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

//除外設定ボタン
const ignore_setting_button = document.getElementById("ignore_setting_button");

//フレンド取得中表示
const friend_log_area = document.getElementById("friend_log_area");

//送信済み取得表示
const sended_log_area = document.getElementById("sended_log_area");

//受信済み取得表示
const recved_log_area = document.getElementById("recved_log_area");

function init(evt) {
    //イベント関連
    async function search_user(evt){
        //検索
        // const req = await AccessPost(friend_search_url,{},{"username":search_value.value});
        const req = await AccessPost(friend_search_url,{"UserName":search_value.value},{
            "content-type" : "application/json",
        });

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

//フレンド検索結果を削除する
function clear_friend_search(evt) {
    clear_child_elems(search_result_area);
}

//受信済みフレンドリクエストを取得する
function refresh_recved_request() {
    //受信済みリクエストをリセット
    clear_child_elems(recved_request_show_area);
    recved_log_area.innerText = "データを取得しています";

    get_recved_requests().then((result) => {
        //受信済みリクエストを表示
        for (let requestid in result) {
            add_recved_request(requestid,result[requestid].aite_name);
        }

        recved_log_area.innerText = "";
    });
}

//受信済みフレンドリクエストを取得する
function get_friend_request(evt) {
    //受信済みリクエストを取得
    refresh_recved_request();

    recved_request_area.classList.toggle("is-show");
}

//送信済みフレンドリクエストを取得する
function refresh_sent() {
    clear_child_elems(sended_request_show_area);
    sended_log_area.innerText = "データを取得しています";

    get_sent_requests().then((result) => {
        //送信済みリクエストを表示
        for (let requestid in result) {
            add_sent_request(requestid,result[requestid].name,result[requestid].uid);
        }

        sended_log_area.innerText = "";
    });
}
//送信済みフレンドリクエストを取得する
async function get_sended_friend_request(evt) {
    refresh_sent();

    sended_request_area.classList.toggle("is-show");
}

function refresh_friend() {
    clear_child_elems(friend_show_area);
    //取得中表示にする
    friend_log_area.innerText = "データを取得しています";

    get_all_friends().then((result) => {
        for (let key in result) {
            console.log(result[key]);
            //追加するdiv
            const adddiv = document.createElement("div");
            adddiv.id = key;

            const dirty = `
                <div class="friend_data_area">
                    <img class="user_icon" src="${GetIconUrl(result[key].aiteid)}">
                    <p class="username_area">${result[key].aite}</p>
                    <button id="remove_button" class="friend_btn">削除</button>
                </div>
            `
            const clean = DOMPurify.sanitize(dirty, { USE_PROFILES: { html: true } });
            //divに書き込む
            adddiv.insertAdjacentHTML("beforeend",clean);

            adddiv.querySelector("#remove_button").addEventListener("click",async function(evt) {
                //イベントキャンセル
                evt.preventDefault();

                //削除
                const req = await AccessPost(remove_friend_url,{"Friendid":key});

                //200以外
                if (req.status != 200) {
                    toastr["error"]("削除に失敗しました");
                    return;
                }

                //通知表示
                toastr["success"]("フレンドを削除しました");

                refresh_friend();
            })

            //追加
            friends_show_area.appendChild(adddiv);
        }

        //取得中表示を消す
        friend_log_area.innerText = ""
    });
}

//フレンド一覧を取得する
async function get_friends(evt) {
    refresh_friend();

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
    const result_div = document.createElement("div");
    result_div.classList.add("search_result");
    
    const dirty = `
        <img class="search_result_icon result" src="${GetIconUrl(uid)}">
        <p class="search_result_name result">${name}</p>
        <button class="send_request_button result" id="${uid}">送信</button>
    `

    const clean = DOMPurify.sanitize(dirty, { USE_PROFILES: { html: true } });
    //検索結果を表示
    result_div.insertAdjacentHTML("beforeend",clean);

    //送信ボタン取得
    const request_btn = result_div.querySelector(".send_request_button");
    
    request_btn.addEventListener("click",async function(evt) {
        //ID取得
        const sendid = evt.target.id;

        //送信
        const req = await AccessPost(send_request_url,{"Targetid":sendid});

        //200以外
        if (req.status != 200) {
            toastr["error"]("送信に失敗しました");
            return;
        }

        //通知表示
        toastr["success"]("フレンドリクエストを送信しました");
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
        user_icon.style.display = "";

        //ユーザ名設定
        user_name.textContent = userinfo["name"];

        //ユーザid設定
        UserID = userinfo["userid"];

        //除外ポイント取得
        const ignore_req = await AccessPost(load_ignore_point_url,{});

        //200以外
        if (ignore_req.status != 200) {
            toastr["warning"]("除外ポイント取得に失敗しました");
            return;
        }

        //除外ポイント
        const ignore_point = await ignore_req.json();

        for (let pointid in ignore_point) {
            //除外マーカー作成
            const point_data = ignore_point[pointid];

            add_pin(pointid,point_data["Latitude"], point_data["Longitude"],point_data["Distance"]);
        }

        //通知距離取得
        const distance_req = await AccessPost(get_notify_distance_url,{});

        //200以外
        if (distance_req.status != 200) {
            toastr["warning"]("通知距離取得に失敗しました");
            return;
        }

        //通知距離
        const distance = await distance_req.json();

        //通知距離設定
        show_distance(distance["distance"]);

        return;

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

//リクエストを追加
function add_recved_request(requestid,name) {
    //追加するdiv
    const adddiv = document.createElement("div");
    adddiv.id = requestid;

    const dirty = `
        <div class="data_area">
            <p class="username_area">${name}</p>
            <button id="accept_btn" class="accept_request_btn">承認</button>
            <button id="reject_btn" class="reject_request_btn">拒否</button>
        </div>
    `
    const clean = DOMPurify.sanitize(dirty, { USE_PROFILES: { html: true } });
    //divに書き込む
    adddiv.insertAdjacentHTML("beforeend",clean);

    //ボタン取得
    const acceptbtn = adddiv.querySelector("#accept_btn");
    const rejectbtn = adddiv.querySelector("#reject_btn");

    //イベント登録
    acceptbtn.addEventListener("click",async function(evt) {
        //イベントキャンセル
        evt.preventDefault();

        //承認
        const req = await AccessPost(accept_request_url,{"requestid":requestid});

        //200以外
        if (req.status != 200) {
            toastr["error"]("承認に失敗しました");
            return;
        }

        toastr["success"]("フレンドリクエストを承認しました");

        refresh_recved_request();
    })

    //イベント登録
    rejectbtn.addEventListener("click",async function(evt) {
        //イベントキャンセル
        evt.preventDefault();

        //拒否
        const req = await AccessPost(reject_request_url,{"requestid":requestid});

        //200以外
        if (req.status != 200) {
            toastr["error"]("拒否に失敗しました");
            return;
        }

        //通知表示
        toastr["success"]("フレンドリクエストを拒否しました");

        //リクエスト欄更新
        refresh_recved_request();
    })
    //追加
    recved_request_show_area.appendChild(adddiv);
}


//リクエストを追加
function add_sent_request(requestid,name,uid) {
    //追加するdiv
    const adddiv = document.createElement("div");
    adddiv.id = requestid;

    const dirty = `
        <div class="request_data_area">
            <img class="user_icon" src="${GetIconUrl(uid)}">
            <p class="username_area">${name}</p>
            <button id="cancel_btn" class="friend_btn">取り消し</button>
        </div>
    `
    const clean = DOMPurify.sanitize(dirty, { USE_PROFILES: { html: true } });
    //divに書き込む
    adddiv.insertAdjacentHTML("beforeend",clean);

    //ボタン取得
    const acceptbtn = adddiv.querySelector("#cancel_btn");

    //イベント登録
    acceptbtn.addEventListener("click",async function(evt) {
        //イベントキャンセル
        evt.preventDefault();

        //承認
        const req = await AccessPost(cancel_request_url,{"requestid":requestid});

        //200以外
        if (req.status != 200) {
            toastr["error"]("キャンセルに失敗しました");
            return;
        }

        toastr["success"]("リクエストをキャンセルしました");

        refresh_sent();
    })

    //追加
    sended_request_show_area.appendChild(adddiv);
}

ignore_setting_button.addEventListener("click",function(evt){
    //イベントキャンセル
    evt.preventDefault();

    ignore_map_area.style.visibility = "visible";
})

ws_event_div.addEventListener(ws_event_key, function (evt) {
    const payload = evt.detail["Payload"];

    switch (evt.detail["Command"]) {
        case "near_friend": {        
            if (payload["is_first"] && !payload["is_self"]) {
                toastr.info(`${payload["unane"]}さんが近くにいます`, "通知",{
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

            break;
        }
        case "stop_notify": {
            break;
        }
    };
})

//アイコン選択
const icon_select = document.getElementById("icon_select");

icon_select.addEventListener("change",async function(evt){
    const updata = new FormData();
    updata.append("file",icon_select.files[0]);

    const icon_post = await AccessPost(server_url + "/upicon",updata,{},false);

    if (icon_post.status != 200) {
        toastr["error"]("アイコンの更新に失敗しました");
        return;
    }

    //アイコンを更新
    user_icon.src = GetIconUrl(UserID);
})

user_icon.addEventListener("click",function(evt){
    //イベントキャンセル
    evt.preventDefault();

    icon_select.click();
})