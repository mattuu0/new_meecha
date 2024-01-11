const loading_div = document.getElementById("loading");
const main_content = document.getElementById("main_content");

window.addEventListener("load",function(evt){
    loading_div.style.visibility = "hidden";
    main_content.style.visibility = "";
})
