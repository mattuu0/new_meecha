const get_user_info = document.getElementById("get_user_info");

get_user_info.addEventListener("click",async function(evt){
    const req = await AccessPost(server_url + "/user_info",{});

    console.log(req.status);

    console.log(await req.json());
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
    updata.append("file",icon_upload.files[0]);

    console.log(updata.getAll("file"));

    const icon_post = await AccessPost(server_url + "/upicon",updata,{},false);

    console.log(icon_post);
})

const refresh_btn = document.getElementById("refresh_btn");

refresh_btn.addEventListener("click",async function(evt){
    const req = await RefreshPost(refresh_url,{});

    console.log(req.status);

    console.log(await req.json());
})