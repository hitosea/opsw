import {computed, watch, ref} from 'vue'
import {createPinia} from 'pinia';
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate';
import {ConfigProviderProps, createDiscreteApi, darkTheme, useOsTheme} from 'naive-ui'
import result from "../utils/result";
import local from "../utils/local";

const themeNameRef = ref('light')
const themeRef = computed(() => {
    const {value} = themeNameRef
    return value === 'dark' ? darkTheme : null
})
const configProviderPropsRef = computed<ConfigProviderProps>(() => ({
    theme: themeRef.value
}))


watch(themeNameRef, name => {
    local.save("themeName", name)
})

export function useThemeName() {
    return themeNameRef
}

export function dialogProvider() {
    const {dialog} = createDiscreteApi(['dialog'], {
        configProviderProps: configProviderPropsRef,
    })
    return dialog
}

export function siteSetup() {
    return {
        resultCode: result.code(),
        resultMsg: result.msg(),
        themeName: themeNameRef,
        theme: themeRef,
    }
}

export function init() {
    return new Promise<void>(async (resolve) => {
        themeNameRef.value = <string>await local.string("themeName")
        if (['light', 'dark'].indexOf(themeNameRef.value) === -1) {
            themeNameRef.value = useOsTheme().value
        }
        resolve()
    })
}

const pinia = createPinia();
pinia.use(piniaPluginPersistedstate);
export default pinia;
