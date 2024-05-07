/** @type {import('tailwindcss').Config} */
export default {
    content: ['./index.html', './src/**/*.{ts,tsx,js,jsx,svg}'],
    theme: {
        extend: {},
    },
    plugins: [require('@tailwindcss/forms')],
}
