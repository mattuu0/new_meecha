//設定ファイル

//サーバIP
const ServerIp = "meecha.tail6cf7b.ts.net/app";

//サーバURL
const server_url = "https://" + ServerIp;

//ログインURL
const login_url = server_url + "/auth/login";

//サインアップURL
const signup_url = server_url + "/auth/signup";

//ログアウトURL
const logout_url = server_url + "/auth/logout";

//ユーザ情報URL
const uinfo_url = server_url + "/user_info";

//メインURL
const main_url = server_url;

//パソコンURL
const desktop_url = "./desktop_show.html";

//アクセストークンキー
const access_token_key = "access_token";

//リフレッシュトークンキー
const refresh_token_key = "refresh_token";

//アイコンURL
const IconBaseUrl = server_url + "/geticon/"

//リフレッシュURL
const refresh_url = server_url + "/auth/refresh";

//フレンド一覧URL
const friends_url = server_url + "/friend/getall";

//フレンド検索URL
const friend_search_url = server_url + "/friend/search";

//リクエスト送信URL
const send_request_url = server_url + "/friend/request";

//送信済み取得URL
const get_sent_url = server_url + "/friend/get_sent";

//受信済み取得URL
const get_recved_url = server_url + "/friend/get_recved";

//リクエスト承認URL
const accept_request_url = server_url + "/friend/accept_request";

//リクエスト拒否URL
const reject_request_url = server_url + "/friend/reject_request";

//フレンド削除URL
const remove_friend_url = server_url + "/friend/remove_friend";

//フレンド削除URL
const cancel_request_url = server_url + "/friend/cancel_request";

//無視リスト保存URL
const save_ignore_point_url = server_url + "/location/save_ignore_point";

//無視リスト保存URL
const load_ignore_point_url = server_url + "/location/load_ignore_point";

//通知距離設定URL
const set_notify_distance_url = server_url + "/location/set_notify_distance";

//通知距離取得URL
const get_notify_distance_url = server_url + "/location/get_notify_distance";

//アイコンURL取得
function GetIconUrl(userid) {
    // 現在時刻
    let ntime = new Date(); 
    //アイコンのURLを返す
    return IconBaseUrl + userid + "?time=" + ntime.getTime();
}

//TODO 本番にはfalseで設定する
//デバッグ設定
const debug = true;

//ユーザID
let UserID = "";