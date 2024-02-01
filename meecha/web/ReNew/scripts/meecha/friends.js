async function get_all_friends() {
    //フレンド一覧取得
    const res = await AccessPost(friends_url, {});

    //JSONに変換
    const friends = await res.json();

    console.log(friends);
}

