// TODO: tag cloud

const sections = Array.from(document.getElementById("reads")
    .getElementsByTagName("li"))
    .map(element => element.getElementsByTagName("section")[0]);

function select_tag(selected) {
    for (const section of sections) {
        if (Array.from(section.getElementsByClassName("tag"))
            .map(element => element.textContent)
            .some(element => element === selected)) {
            section.parentElement.style.removeProperty('opacity');
        } else {
            section.parentElement.style.opacity = '0.33';
        }
    }
}
