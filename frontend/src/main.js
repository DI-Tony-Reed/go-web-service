import './assets/main.css'

import {createApp} from 'vue'
import App from './App.vue'
import router from './router'
import EnvironmentVariables from './helpers/EnvironmentVariables'

const app = createApp(App)

app.use(router)

app.mount('#app')

let applicationUrl = localStorage.getItem('applicationUrl')
if (!applicationUrl) {
  const vars = new EnvironmentVariables()
  const applicationUrl = `${vars.applicationProtocol()}${vars.applicationUrl()}:${vars.applicationPort()}`
  localStorage.setItem('applicationUrl', applicationUrl)
}
