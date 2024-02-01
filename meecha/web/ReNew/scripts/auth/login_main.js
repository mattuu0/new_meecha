//ログインフォーム取得
const login_form_ = document.getElementById("login_form");

//ユーザー名入力エリア
const id_login = document.getElementById("id_login");

//パスワード入力エリア
const id_password = document.getElementById("id_password");

//ログイン処理
async function submit_login(evt){
    //イベント中止
    evt.preventDefault();

    //ユーザー名取得 
    const username = id_login.value;

    //パスワード取得
    const password = id_password.value;

    //結果
    const login_result = await login(username,password);

    //失敗したとき
    if (!login_result) {
        alert("サインインに失敗しました");
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
    //送信イベント
    login_form.addEventListener("submit",submit_login);

    get_userinfo();
});