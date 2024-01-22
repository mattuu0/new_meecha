//アクセストークンでポストする
async function AccessPost(posturl,body) {
    //リクエスト飛ばす
    return await TokenPost(posturl,body,get_access_token());
}

//リフレッシュトークンでポストする
async function RefreshPost(posturl,body) {
    //リクエスト飛ばす
    return await TokenPost(posturl,body,get_refresh_token());
}

//トークン付きでポストする
async function TokenPost(posturl,body,token) {
    //リクエストを飛ばす
    const req = await fetch(posturl,{
        method: "POST",
        headers: {
            "Authorization" : token, 
        },
        body: JSON.stringify(body)
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