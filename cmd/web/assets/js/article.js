try {
    const articleContainer = document.getElementById("article");
    const articleContent = articleContainer.getAttribute("data-content");

    articleContainer.innerHTML = articleContent;
} catch (e) {}

document.addEventListener("DOMContentLoaded", (event) => {
    console.log("DOM fully loaded and parsed");
    document.body.addEventListener("htmx:beforeSwap", (e) => {
        if (
            e.detail.xhr.status === 422 ||
            e.detail.xhr.status === 500 ||
            e.detail.xhr.status === 400
        ) {
            e.detail.shouldSwap = true;
            e.detail.isError = false;
        }
    });
});
