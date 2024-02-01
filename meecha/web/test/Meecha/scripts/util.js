//アクセストークンでポストする
async function AccessPost(posturl,body,headers = {},seriarize = true) {
    //リクエストを送信
    const first_request = await TokenPost(posturl,body,get_access_token(),headers,seriarize);

    //リクエストが403の時
    if (first_request.status == 403) {
        //トークンを更新
        const req = await TokenPost(refresh_url,body,get_refresh_token(),{},seriarize);

        //Jsonに変換
        const Atoken = await req.json();

        //トークン抽出
        const AccessToken = Atoken["token"];

        //トークン保存
        save_token(AccessToken,get_refresh_token());

        //もう一度リクエストを送信
        return await TokenPost(posturl,body,AccessToken,headers,seriarize);
    }
    

    //リクエスト飛ばす
    return first_request;
}

//リフレッシュトークンでポストする
async function RefreshPost(posturl,body,headers = {},seriarize = true) {
    //リクエスト飛ばす
    return await TokenPost(posturl,body,get_refresh_token(),headers,seriarize);
}

//トークン付きでポストする
async function TokenPost(posturl,body,token,headers,seriarize = true) {
    //送信するデータ
    let body_data = body;

    //シリアライズするか
    if (seriarize) {
        //Jsonをシリアライズ
        body_data = JSON.stringify(body);
    }

    //トークンを設定
    headers["Authorization"] = token;
    //リクエストを飛ばす
    const req = await fetch(posturl,{
        method: "POST",
        headers,
        body: body_data
    })

    return req;
}


//トークンを保存する関数
function save_token(AccessToken,RefreshToken) {
    //ローカルストレージ
    const storage = window.localStorage;

    //アクセストークン保存
    storage.setItem(access_token_key,AccessToken);

    //リフレッシュトークン保存
    storage.setItem(refresh_token_key,RefreshToken);
}

//アクセストークンを取得する
function get_access_token() {
    //ローカルストレージ
    const storage = window.localStorage;

    //トークンを取得
    let token = storage.getItem(access_token_key);

    //トークンが存在しない場合空文字にする
    if (!token) {
        token = "";
    }

    return token; 
}

//リフレッシュトークンを取得する
function get_refresh_token() {
    //ローカルストレージ
    const storage = window.localStorage;

    //トークンを取得
    let token = storage.getItem(refresh_token_key);

    //トークンが存在しない場合空文字にする
    if (!token) {
        token = "";
    }

    return token; 
}