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