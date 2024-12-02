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

  dateElement.textContent = `${prefix} ${formattedDate}`;
});
