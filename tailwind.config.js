/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./views/**/*.templ",
    "./views/**/*.go",
    "./cmd/**/*.go",
    "./internal/**/*.go",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
