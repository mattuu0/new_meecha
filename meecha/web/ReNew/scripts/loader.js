const loading_div = document.getElementById("loading");
const main_content = document.getElementById("main_content");

function hide_load() {
    loading_div.style.visibility = "hidden";
    main_content.style.visibility = "";
}

function show_load() {
    loading_div.style.visibility = "";
    main_content.style.visibility = "hidden";
}

//TODO デバッグ用なので本番に消す
hide_load();