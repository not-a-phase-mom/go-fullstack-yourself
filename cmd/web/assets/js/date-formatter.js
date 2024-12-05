const dateElements = document.querySelectorAll("#date");
const locale = navigator.language;

dateElements.forEach((dateElement) => {
  const prefix = dateElement.getAttribute("data-prefix");
  const date = dateElement.getAttribute("data-date");

  const formattedDate = new Date(date).toLocaleDateString(locale, {
    year: "numeric",
    month: "long",
    day: "numeric",
  });

  // create a p element with the prefix
  const prefixElement = document.createElement("p");
  prefixElement.textContent = prefix;
  // create a p element with the formatted date
  const formatedDateElement = document.createElement("p");
  formatedDateElement.className = "text-black";
  formatedDateElement.textContent = formattedDate;

  // append the elements to the dateElement
  dateElement.appendChild(prefixElement);
  dateElement.appendChild(formatedDateElement);

  dateElement.classList.add("flex", "items-center", "space-x-1");
});
