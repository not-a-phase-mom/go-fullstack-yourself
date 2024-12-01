/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./cmd/web/templates/**/*.{html,templ}"],
  theme: {
    extend: {
      fontFamily: {
        //--font-sans: 'Inter', sans-serif;
        //--font-serif: 'IBM Plex Serif', serif;
        //--font-mono: 'IBM Plex Mono', monospace;
        //--font-display: 'Poppins', sans-serif;

        "font-sans": ["Inter", "sans-serif"],
        "font-serif": ["IBM Plex Serif", "serif"],
        "font-mono": ["IBM Plex Mono", "monospace"],
        "font-display": ["Poppins", "sans-serif"],
      },
    },
  },
  plugins: [],
};
