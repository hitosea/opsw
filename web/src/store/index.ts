import {computed, watch, ref} from 'vue'
import {darkTheme, useOsTheme} from 'naive-ui'
import utils from "./utils";
import call from "./call";

const userInfoRef = ref({name: '', avatar_url: ''})
const themeNameRef = ref('light')
const themeRef = computed(() => {
    const {value} = themeNameRef
    return value === 'dark' ? darkTheme : null
})
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
        call({
            method: "get",
            url: 'user/info',
        }).then(({data}) => {
            userInfoRef.value = data
            resolve(data)
        }).catch(err => {
            reject(err)
        })
    })
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
