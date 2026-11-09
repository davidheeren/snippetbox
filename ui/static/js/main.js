var navLinks = document.querySelectorAll("nav a");
var hasNav = false
var title = document.title.trimStart().split(" ")[0]
for (var i = 0; i < navLinks.length; i++) {
    var link = navLinks[i]
    if (link.textContent == title) {
        link.classList.add("live");
        hasNav = true
        break;
    }
}

// theres probably a better way to handle dynamic nav tabs
if (!hasNav) {
    navlink = document.createElement("a")
    navlink.textContent = title
    navlink.classList.add("live");
    document.querySelector("nav div").appendChild(navlink)
}

function localizeTimes(root) {
    root.querySelectorAll("time").forEach(el => {
        el.textContent = new Date(el.textContent).toLocaleString(undefined, { dateStyle: "short", timeStyle: "short" })
    })
}

localizeTimes(document)
document.addEventListener("htmx:afterSwap", e => localizeTimes(e.detail.target))
