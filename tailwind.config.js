/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
		"./views/**/*.{html,js,templ,go}",
		"./components/**/*.{html,js,templ,go}",
	],
  theme: {
    screens: {
      xsm: '480px',
      sm: '640px',
      md: '768px',
      lg: '1024px',
      xl: '1280px',
      xxl: '1536px',
    },
    fontFamily: {
      sans: ['Atkinson Hyperlegible', 'sans-serif']
    },
    extend: {},
  },
  plugins: [],
}
