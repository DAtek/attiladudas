/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
	theme: {
		extend: {
			colors: {
				primary: "#3b82f6",
				danger: "#dc2626",
				active: "#1d4ed8"
			}
		},
	},
	plugins: [],
}
