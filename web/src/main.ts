import {createApp} from 'vue'
import App from './App.vue'
import {init} from './store'
import {routes} from './routes/routes'
import createDemoRouter from './routes/router'
import './style.less'

const app = createApp(App)
const router = createDemoRouter(app, routes)
app.use(router)

init().then(() => {
    router.isReady().then(() => {
        app.mount('#app')
    })
})

