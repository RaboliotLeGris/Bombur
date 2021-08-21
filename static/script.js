// It's a terrible piece of Javascript. But for now, it will do.
document.addEventListener("DOMContentLoaded", function() {
    let form = document.getElementById("link-form")
    form.addEventListener('submit', (e) => {
        e.preventDefault();
        e.stopPropagation();

        const link = e.target[0].value
        const expireIn = +e.target[1].value

        const request = new XMLHttpRequest();
        request.open("POST", "/link");
        request.setRequestHeader("Content-Type", "application/json");
        request.onreadystatechange = function () {
            if (request.readyState === 4 && request.status === 200) {
                displayLink(JSON.parse(request.responseText).link)
            }
        };
        const data = {link}
        if (expireIn) {
            data.expire = `${expireIn*60}s`
        }
        request.send(JSON.stringify(data))
    });
});

function displayLink(link) {
    console.log(link)

    const result = document.getElementById("result");
    result.classList.remove("hidden")

    document.getElementById("output-link").value = link
}