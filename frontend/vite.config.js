import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'

// https://vitejs.dev/config/
const config = defineConfig({
  plugins: [react()],
  build: {
		rollupOptions: {
			input: {
        main: 'index.html',
        login: 'login.html',
      },
		},
		outDir: '../dist',
		emptyOutDir: true,
		minify: false, // Disable minification
	},
})

export default Object.freeze(config)