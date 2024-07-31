//サインインフォーム取得
const signup_form_ = document.getElementById("signup_form");

//ユーザー名入力エリア
const signin_id_login = document.getElementById("id_username");

//パスワード入力エリア
const signin_id_password1 = document.getElementById("id_password1");

//パスワード確認入力エリア
const signin_id_password2 = document.getElementById("id_password2");

//サインアップ処理
async function submit_signup(evt){
    //イベント中止
    evt.preventDefault();

    //ユーザー名取得 
    const username = signin_id_login.value;

    //パスワード取得
    const password1 = signin_id_password1.value;

    //確認パスワード取得
    const password2 = signin_id_password2.value;

    //パスワード確認
    if (password1 != password2) {
        //パスワードが一致しない場合
        alert("パスワードが一致しません");
        return;
    }

    //結果
    const login_result = await signup(username,password1);

    //失敗したとき
    if (login_result != 200) {
        alert("サインアップに失敗しました");
        return;
    }

    //ログイン成功したら飛ばす
    if (login_result) {
        window.location.href = "./index.html";
    }
}

//ユーザ情報取得処理
async function get_userinfo() {
    try {
        //りくえすと　
        const req = await AccessPost(uinfo_url, {});

        //403の時
        if (req.status == 200) {
            //ログインに戻る
            window.location.href = "./index.html";
            return;
        }
    } catch (error) {
        //エラー処理
        console.log(error);
        return;
    }
}

//ロード完了イベント
window.addEventListener("DOMContentLoaded",function(evt){
    get_userinfo();
    //送信イベント
    signup_form_.addEventListener("submit",submit_signup);
});