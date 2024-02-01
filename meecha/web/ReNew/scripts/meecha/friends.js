//フレンド一覧を取得する
async function get_all_friends() {
    //フレンド一覧取得
    const res = await AccessPost(friends_url, {});

    //JSONに変換
    const friends = await res.json();

    console.log(friends);
}

//送信済みフレンドリクエストを取得する
async function get_sent_requests() {
    //フレンド一覧取得
    const res = await AccessPost(get_send_url, {});

    //JSONに変換
    const sent_requests = await res.json();

    console.log(sent_requests);
}

