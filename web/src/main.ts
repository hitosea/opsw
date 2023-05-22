import {createApp} from 'vue'
import App from './App.vue'
import pinia, {GlobalStore} from './store'
import {routes} from "./routes/routes";
import createDemoRouter from './routes'
import './styles/common.less'

const app = createApp(App)
const route = createDemoRouter(app, routes)
app.use(route)
app.use(pinia)

const globalStore = GlobalStore()

globalStore.init().then(() => {
    route.isReady().then(() => {
        app.mount('#app')
    })
})

