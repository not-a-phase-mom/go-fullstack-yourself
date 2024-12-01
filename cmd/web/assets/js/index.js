const date = document.getElementById("current-date");
const weather = document.getElementById("weather-desc");
const navigationList = document.getElementById("navigation-list");

// get user locale
const locale = navigator.language;

if (!date || !weather) {
  console.error("Couldn't find the date or weather element");
} else {
  const forecast = [
    "Cloudy with a chance of merge conflicts",
    "Sunny with occasional cache invalidation",
    "Heavy npm install storms expected",
    "Clear skies, perfect for cloud computing",
    "Foggy, visibility reduced to 404 errors",
    "High pressure system of looming deadlines",
    "Scattered showers of semicolons",
    "Windy, with gusts of rapid API changes",
    "Overcast, 90% chance of coffee spills",
    "Heat wave of overclocking CPUs",
    "Partly cloudy with intermittent Wi-Fi outages",
    "Thunderstorms of stack overflow questions",
    "Mild, with a refreshing breeze of code refactoring",
    "Expect heavy downpours of data packets",
    "Chilly, don't forget your Docker layers",
  ];

  // get current date
  const options = {
    weekday: "long",
    year: "numeric",
    month: "long",
    day: "numeric",
  };

  const today = new Date();
  const dateParts = today.toLocaleDateString(locale, options).split(",");
  date.innerHTML = `${dateParts[0]}, ${dateParts[1]}`;

  // get weather forecast
  const weatherIndex = Math.floor(Math.random() * forecast.length);
  //content should be: Weather: <span id="weather-desc" class="text-stone-600"></span>
  weather.innerHTML = `Weather: <span class="text-stone-600">${forecast[weatherIndex]}</span>`;
}

// based on the current url like: /articles or /about. We can highlight the current page in the navigation. use the li elements inside the ul with the id navigation-list
const currentPath = window.location.pathname;
const navigationLinks = navigationList.querySelectorAll("li");

navigationLinks.forEach((link) => {
  const anchor = link.querySelector("a");
  if (anchor && anchor.href === currentPath) {
    link.classList.add("bg-indigo-500");
  }
});

if (currentPath === "/articles") {
  const dateElements = document.querySelectorAll("#date-published");
  console.log(dateElements);

  dateElements.forEach((element) => {
    const dataDate = (element.attributes["data-date"].value =
      element.textContent);
    const date = new Date(dataDate);
    const options = {
      weekday: "long",
      year: "numeric",
      month: "long",
      day: "numeric",
    };
    element.textContent = date.toLocaleDateString(locale, options);
  });
}
