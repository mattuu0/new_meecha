//設定ファイル

//サーバURL
const server_url = "http://127.0.0.1:12222";

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

//アイコンURL取得
function GetIconUrl(userid) {
    //アイコンのURLを返す
    return IconBaseUrl + userid;
}

//TODO 本番にはfalseで設定する
//デバッグ設定
const debug = true;

//ユーザID
let UserID = "";