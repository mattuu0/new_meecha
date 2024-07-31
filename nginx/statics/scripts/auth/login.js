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
        return false;
    }

    //データ取得
    const result = await req.json();

    //結果からトークンを取り出す
    const AccessToken = result["AccessToken"];
    const RefreshToken = result["RefreshToken"];

    //トークン保存
    save_token(AccessToken,RefreshToken);

	return true;
}
