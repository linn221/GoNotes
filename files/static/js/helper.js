function $(selector) {
    return document.querySelector(selector)
}

function myscroll(selector) {
    $(selector).scrollIntoView({ block: 'start', behavior: 'smooth' })
}

function scrolltomain() {
    myscroll("#main")
}
function autoResize(textarea) {
    textarea.style.height = 'auto'; // Reset the height
    textarea.style.height = textarea.scrollHeight + 'px'; // Set height based on content
}

// resize a to height of b
function resizeAb(a, b) {
    a.style.height = 'auto'; // Reset the height
    a.style.height = b.scrollHeight + 'px'; // Set height based on content
}
document.body.addEventListener('htmx:responseError', function (event) {
    // Only handle 500 errors
    if (event.detail.xhr.status === 500) {
        const message = event.detail.xhr.responseText;

        const alert = document.createElement('div');
        alert.setAttribute('role', 'alert');
        alert.className = 'alert';
        alert.style.padding = '1rem';
        alert.style.border = '1px solid var(--pico-color-red)';
        alert.style.backgroundColor = 'var(--pico-color-red-100)';
        alert.style.color = 'var(--pico-color-red)';
        alert.style.marginBottom = '1rem';

        alert.textContent = message;

        const container = document.getElementById('alert-container');
        container.appendChild(alert);

        // Remove alert after 5 seconds
        setTimeout(() => {
            alert.remove();
        }, 4000);
    }
});