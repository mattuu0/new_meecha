//ユーザー名エリア
const username_input = document.getElementById("login_username_input");

//パスワードエリア
const password_input = document.getElementById("login_password_input");

//ログインフォーム
const login_form = document.getElementById("login_form");

//イベント設定
login_form.addEventListener("submit", async function (evt) {
    //formのイベントキャンセル
    evt.preventDefault();

    //ユーザ名
    const username = username_input.value;

    //パスワード
    const password = password_input.value;

    //ログインするa
    await login(username,password);
})

//ログイン関数
async function login(username,password) {
    //送信するデータ
    const login_data = {
        "name": username,
        "password": password
    };

    //送信
    const req = await fetch(login_url,{
        method: "post",
        body : JSON.stringify(login_data)
    });

    //ステータスコード確認
    if (req.status != 200) {
        //200以外の時
        alert("ログインに失敗しました");
        return;
    }

    //データ取得
    const result = await req.json();

    //結果からトークンを取り出す
    const AccessToken = result["AccessToken"];
    const RefreshToken = result["RefreshToken"];

    //トークン保存
    save_token(AccessToken,RefreshToken);
}
