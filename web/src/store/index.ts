import {computed, watch, ref} from 'vue'
import {ConfigProviderProps, createDiscreteApi, darkTheme, useOsTheme} from 'naive-ui'
import utils from "./utils";
import call from "./call";

const userInfoRef = ref()
const themeNameRef = ref('light')
const themeRef = computed(() => {
    const {value} = themeNameRef
    return value === 'dark' ? darkTheme : null
})
const configProviderPropsRef = computed<ConfigProviderProps>(() => ({
    theme: themeRef.value
}))
watch(themeNameRef, name => {
    utils.IDBSave("themeName", name)
})

export function useThemeName() {
    return themeNameRef
}

export function useUserInfo() {
    return userInfoRef
}

export function loadUserInfo() {
    return new Promise((resolve, reject) => {
        call.get('user/info')
            .then(({data}) => {
                userInfoRef.value = data
                resolve(data)
            })
            .catch(err => {
                reject(err)
            })
    })
}

export function dialogProvider() {
    const {dialog} = createDiscreteApi(['dialog'], {
        configProviderProps: configProviderPropsRef,
    })
    return dialog
}

export function siteSetup() {
    return {
        resultCode: utils.resultCode(),
        resultMsg: utils.resultMsg(),
        themeName: themeNameRef,
        theme: themeRef,
    }
}

export function init() {
    return new Promise<void>(async (resolve) => {
        themeNameRef.value = <string>await utils.IDBString("themeName")
        if (['light', 'dark'].indexOf(themeNameRef.value) === -1) {
            themeNameRef.value = useOsTheme().value
        }
        resolve()
    })
}
