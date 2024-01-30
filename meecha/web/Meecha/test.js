const get_user_info = document.getElementById("get_user_info");
const icon_img = document.getElementById("icon_img");

get_user_info.addEventListener("click",async function(evt){
    const req = await AccessPost(server_url + "/user_info",{});

    console.log(req.status);

    const res_json = await req.json();
    console.log(res_json);

    const ShowId = res_json["userid"];
    const iconurl = `${server_url}/geticon/${ShowId}?${new Date().getTime()}`;

    icon_img.src = iconurl;
})

const logout_btn = document.getElementById("logout_btn");

logout_btn.addEventListener("click",async function(evt){
    const req = await RefreshPost(logout_url,{});

    console.log(req.status);

    console.log(await req.json());
})

const icon_upload = document.getElementById("icon_upload");

const upload_icon_btn = document.getElementById("upload_icon");

upload_icon_btn.addEventListener("click",async function(evt){
    const updata = new FormData();
    updata.append("imgfile",icon_upload.files[0]);

    console.log(updata.getAll("file"));

    const icon_post = await AccessPost(server_url + "/upicon",updata,{},false);

    console.log(icon_post);
})