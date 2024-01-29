//ユーザアイコン
const usericon = document.getElementById("usericon");

//ユーザ情報取得処理
async function get_userinfo() {
    //りくえすと　
    const req = await AccessPost(uinfo_url,{});

    //ユーザデータ取得
    const userinfo = await req.json();

    //アイコンURL
    usericon.src = GetIconUrl(userinfo["userid"]);
}

get_userinfo();