import {createApp} from 'vue'
import App from './App.vue'
import {init} from './store'
import {router} from "./routes/router";
import createDemoRouter from './routes'
import './styles/common.less'

const app = createApp(App)
const route = createDemoRouter(app, router)
app.use(route)

init().then(() => {
    route.isReady().then(() => {
        app.mount('#app')
    })
})

