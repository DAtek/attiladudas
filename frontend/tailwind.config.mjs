/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
	theme: {
		extend: {
			colors: {
				primary: "#3b82f6",
				danger: "#dc2626",
				active: "#1d4ed8",
				black: "#0f172a",
				"gray-light": "#6b7280",
				"gray-dark": "#1f2937",
			},
			fontFamily: {
				sans: ['Atkinson', 'sans-serif']
			}
		},
	},
	plugins: [],
}
