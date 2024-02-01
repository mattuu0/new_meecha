//ログインフォーム取得
const login_form = document.getElementById("login_form");

//ユーザー名入力エリア
const id_login = document.getElementById("id_login");

//パスワード入力エリア
const id_password = document.getElementById("id_password");

//ログイン処理
function submit_login(evt){
    //イベント中止
    evt.preventDefault();

    id_login.value;

    id_password.value;

    

}

//ロード完了イベント
window.addEventListener("DOMContentLoaded",function(evt){
    //送信イベント
    login_form.addEventListener("submit",submit_login);
});