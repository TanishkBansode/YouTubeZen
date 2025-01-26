import App from './App.svelte'
import './styles/styles.css';
import './styles/theme-light.css';
import './styles/theme-dark.css';

const app = new App({
  target: document.getElementById('app')
})

export default app
