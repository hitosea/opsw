import {defineStore} from 'pinia';
import {UserState} from './interface';
import piniaPersistConfig from "./config/pinia-persist";
import {getUserInfo} from "../api/modules/user";
import utils from "../utils/utils";
import {ref, watch} from "vue";
import {WsStore} from "./ws";

const wsWatch = ref(false)

export const UserStore = defineStore({
    id: 'UserState',
    state: (): UserState => ({
        info: {},
    }),
    actions: {
        refresh() {
            if (!wsWatch.value) {
                wsWatch.value = true
                const wsStore = WsStore()
                watch(
                    _ => this.info,
                    info => {
                        if (wsStore.uid !== info.id) {
                            wsStore.uid = info.id
                            wsStore.connection()
                        }
                    }
                )
            }
            return new Promise((resolve, reject) => {
                getUserInfo()
                    .then(({data}) => {
                        if (utils.isEmpty(data.name)) {
                            data.name = data.email
                        }
                        this.info = data
                        resolve(data)
                    })
                    .catch(err => {
                        reject(err)
                    })
            })
        }
    },
    persist: piniaPersistConfig('UserState'),
});
