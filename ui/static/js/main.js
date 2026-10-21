var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
    var link = navLinks[i]
    if (link.getAttribute('href') == window.location.pathname) {
        link.classList.add("live");
        break;
    }
}

document.querySelectorAll("time").forEach(el => {
    el.textContent = new Date(el.textContent).toLocaleString()
});
