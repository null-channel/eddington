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
      backgroundImage: {
        "hero-bg": "url('./src/assets/img/hero-bg.png')",
        "sign-up-bg": "url('./src/assets/img/sign-up.png')",
      },
      fontFamily: {
        inter: ["Inter", "sans-serif"],
        roboto: ["Roboto", "sans-serif"],
      },
    },
  },
  plugins: [formKitTailwind],
};
