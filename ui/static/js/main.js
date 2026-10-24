var navLinks = document.querySelectorAll("nav a");
var hasNav = false
for (var i = 0; i < navLinks.length; i++) {
    var link = navLinks[i]
    if (link.getAttribute('href') == window.location.pathname) {
        link.classList.add("live");
        hasNav = true
        break;
    }
}

// theres probably a better way to handle dynamic nav tabs
if (!hasNav) {
    navlink = document.createElement("a")
    navlink.textContent = document.title.trimStart().split(" ")[0];
    navlink.classList.add("live");
    document.querySelector("nav div").appendChild(navlink)
}

document.querySelectorAll("time").forEach(el => {
    el.textContent = new Date(el.textContent).toLocaleString()
});
