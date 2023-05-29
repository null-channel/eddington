/** @type {import('tailwindcss').Config} */
const formKitTailwind = require("@formkit/themes/tailwindcss");

export default {
  content: [
    "./public/**/*.html",
    "./src/**/*.{js,jsx,ts,tsx,vue,html}",
    "./node_modules/@formkit/themes/dist/tailwindcss/genesis/index.cjs",
  ],
  theme: {
    extend: {
    },
  },
  plugins: [formKitTailwind],
};
