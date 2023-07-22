const button = document.body.appendChild(document.createElement('button'));
button.title = "top";
button.textContent = "☝";
button.style.position = "fixed";
button.style.right = "15px";
button.style.top = "10px";
button.style.display = "none";

button.onclick = () => {
    document.documentElement.scrollTop = 0;
    document.body.scrollTop = 0;
};

window.onscroll = () => {
    if (document.documentElement.scrollTop >= 10 || document.body.scrollTop >= 10) {
        button.style.display = "block";
    } else {
        button.style.display = "none";
    }
};
