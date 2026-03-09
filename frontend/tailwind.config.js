/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{vue,js,ts}"],
  theme: {
    extend: {
      colors: {
        primary: "#00bfa3",
        dark: "#080817",
        lightBg: "#F2F4F8",
      },
      fontFamily: {
        sans: ["Inter", "sans-serif"],
      },
      boxShadow: {
        card: "0 10px 30px rgba(0,0,0,0.05)",
      },
      borderRadius: {
        xl2: "24px",
      },
    },
  },
}