import {defineConfig} from 'vite'
import {svelte} from '@sveltejs/vite-plugin-svelte'
import basicSsl from '@vitejs/plugin-basic-ssl'

// https://vitejs.dev/config/
export default defineConfig(({mode}) => {
  // Relative base so the built app works at any subpath
  // (e.g. GitHub Pages project sites at /<repo>/).
  // `npm run dev:phone` serves HTTPS so phones on the LAN can install the
  // PWA (service workers need a secure context).
  return {
    base: './',
    plugins: mode === 'phone' ? [svelte(), basicSsl()] : [svelte()],
  }
})
