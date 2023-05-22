import {computed} from 'vue'
import {createPinia, defineStore} from 'pinia';
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate';
import {ConfigProviderProps, createDiscreteApi, darkTheme, useOsTheme} from 'naive-ui'
import result from "../utils/result";
import {GlobalState} from "./interface";
import piniaPersistConfig from "./config/pinia-persist";

export const GlobalStore = defineStore({
    id: 'GlobalState',
    state: (): GlobalState => ({
        isLoading: 0,
        themeName: '',
    }),
    actions: {
        setLoading() {
            this.isLoading++;
        },
        cancelLoading() {
            this.isLoading--;
        },
        setThemeName(name: string) {
            this.themeName = name
        },
        async init() {
            this.isLoading = 0
            if (['light', 'dark'].indexOf(this.themeName) === -1) {
                this.themeName = useOsTheme().value
            }
        },
        appSetup() {
            return {
                resultCode: result.code(),
                resultMsg: result.msg(),
                themeName: this.themeName,
                theme: computed(() => {
                    return this.themeName === 'dark' ? darkTheme : null
                }),
            }
        },
        dialogProvider() {
            const {dialog} = createDiscreteApi(['dialog'], {
                configProviderProps: computed<ConfigProviderProps>(() => ({
                    theme: this.themeName === 'dark' ? darkTheme : null
                })),
            })
            return dialog
        }
    },
    persist: piniaPersistConfig('GlobalState'),
});

const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);
export default pinia;
