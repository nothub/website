// https://web.archive.org/web/20210312101244/http://www.kopimi.com/kopimi/

onmousemove = e => {
    const w = e.x / window.innerWidth;
    const h = e.y / window.innerHeight;
    let kopimi = document.getElementById("kopimi")
    let bounds = kopimi.getBoundingClientRect()
    let x = (bounds.left + bounds.right) / 2
    let y = (bounds.top + bounds.bottom) / 2

    if (w > 0.65 && h <= 0.25) kopimi.style.visibility = "visible"
    else kopimi.style.visibility = "hidden"

    let distance = Math.floor(Math.sqrt(Math.pow(e.clientX - (kopimi.offset().left + (kopimi.width() / 2)), 2) + Math.pow(e.clientY - (kopimi.offset().top + (kopimi.height() / 2)), 2)));
    console.log(distance)
}
