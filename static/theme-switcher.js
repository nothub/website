const themes = new Map([
    ["light", new Map([
        ["--bg-figure", "#DDD"],
        ["--bg", "#FDFDFD"],
        ["--secondary", "#444"],
        ["--text", "#000"],
        ["--link-new", "#0000EE"],
        ["--link-vis", "#551A8B"],
    ])],
    ["dark", new Map([
        ["--bg-figure", "#151515"],
        ["--bg", "#050709"],
        ["--secondary", "#9EAC30"],
        ["--text", "#E6E2BD"],
        ["--link-new", "#9EAC30"],
        ["--link-vis", "#7E8926"],
    ])],
]);

const input = document.body.appendChild(
    document.createElement('input'));
input.setAttribute("class", "theme-switcher")
input.setAttribute("id", "theme-switcher")
input.setAttribute("type", "checkbox")

const label = document.body.appendChild(
    document.createElement('label'));
label.setAttribute("class", "theme-switcher-label")
label.setAttribute("for", "theme-switcher")

const sun = label.appendChild(document.createElement('i'))
sun.setAttribute("class", "theme-switcher-sun")
sun.style.fontSize = "15px"
sun.innerText = '🌞'

const moon = label.appendChild(document.createElement('i'))
moon.setAttribute("class", "theme-switcher-moon")
moon.style.fontSize = "15px"
moon.innerText = '🌛'

const knob = label.appendChild(document.createElement('span'))
knob.setAttribute("class", "theme-switcher-knob")

let theme = localStorage.getItem("theme");
if (theme === null) {
    if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
        theme = "dark"
        localStorage.setItem("theme", "dark");
    } else {
        theme = "light"
        localStorage.setItem("theme", "light");
    }
}

function setTheme(str) {
    theme = str
    localStorage.setItem("theme", str);
    for (const [name, color] of themes.get(str)) {
        document.documentElement.style.setProperty(name, color);
    }
}

const checkbox = document.getElementById("theme-switcher");

// ensure correct state
if (theme === "light") {
    checkbox.checked = !checkbox.checked
    setTheme("light")
} else {
    setTheme("dark")
}

checkbox.addEventListener("change", () => {
    switch (theme) {
        case "light":
            setTheme("dark")
            break;
        case "dark":
            setTheme("light")
            break;
    }
});
