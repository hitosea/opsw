import {computed, watch, ref} from 'vue'
import {ConfigProviderProps, createDiscreteApi, darkTheme, useOsTheme} from 'naive-ui'
import utils from "./utils";
import call from "./call";

interface UserInfoModel {
    id?: number
    email?: string
    name?: string
    token?: string
    avatar?: string
    created_at?: string
    updated_at?: string
}

const userInfoRef = ref<UserInfoModel>({})
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
        call.get<UserInfoModel>('user/info')
            .then(({data}) => {
                if (utils.isEmpty(data.name)) {
                    data.name = data.email
                }
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
