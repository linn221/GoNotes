function $(selector) {
    return document.querySelector(selector)
}

function myscroll(selector) {
    $(selector).scrollIntoView({ block: 'start', behavior: 'smooth' })
}

function scrolltomain() {
    myscroll("#main")
}