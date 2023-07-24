const themes = new Map([
    ["light", new Map([
        ["--bg-figure", "#DDD"],
        ["--bg", "#FFF"],
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

let theme = localStorage.getItem("theme");
if (theme === null) {
    if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
        localStorage.setItem("theme", "dark");
        theme = "dark"
    } else {
        localStorage.setItem("theme", "light");
        theme = "light"
    }
}

const checkbox = document.getElementById("theme-switcher");

// init checkbox to correct setting
if (theme === "light") {
    checkbox.checked = !checkbox.checked
}

checkbox.addEventListener("change", () => {
    switch (theme) {
        case "light":
            theme = "dark"
            for (const [name, color] of themes.get("dark")) {
                document.documentElement.style.setProperty(name, color);
            }
            break;
        case "dark":
            theme = "light"
            for (const [name, color] of themes.get("light")) {
                document.documentElement.style.setProperty(name, color);
            }
            break;
    }
});
