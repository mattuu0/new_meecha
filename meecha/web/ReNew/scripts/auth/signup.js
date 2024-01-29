async function signup(username,password) {
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
            return 409;
        default:
            return 400;
    }

    //ログイン試行
    if (!await login(username,password)) {
        //ログイン失敗
        return 400;
    }
    //データ取得
    return 200;
}