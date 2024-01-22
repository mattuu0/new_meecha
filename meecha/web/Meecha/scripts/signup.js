//ユーザー名エリア
const signup_username_input = document.getElementById("signup_username_input");

//パスワードエリア
const signup_password_input = document.getElementById("signup_password_input");

//サインアップフォーム
const signup_form = document.getElementById("signup_form");

//イベント設定
signup_form.addEventListener("submit", async function (evt) {
    //formのイベントキャンセル
    evt.preventDefault();

    //ユーザ名
    const username = signup_username_input.value;

    //パスワード
    const password = signup_password_input.value;

    //送信するデータ
    const signup_data = {
        "name": username,
        "password": password
    };

    //送信
    const req = await fetch(signup_url,{
        method: "post",
        body : JSON.stringify(signup_data)
    });

    switch (req.status) {
        case 200:
            //成功したので何もしない
            break;
        case 409:
            //ユーザ名が重複したとき
            alert("既に登録されています");
            return;
        default:
            alert("サインアップに失敗しました");
            return;
    }

    //データ取得
    await login(username,password);
})