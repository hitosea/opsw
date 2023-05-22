<template>
    <div class="header">
        <n-layout-header class="header-menu" bordered>
            <div class="wrapper">
                <n-menu mode="horizontal" v-model:value="menuAction" :options="menuOptions" />
                <div class="user">
                    <n-button
                        size="small"
                        quaternary
                        class="name"
                        @click="handleThemeUpdate">
                        {{ themeLabelMap[themeName] }}
                    </n-button>
                    <n-dropdown
                        v-if="userInfo.name"
                        trigger="click"
                        :show-arrow="true"
                        :options="userMenuOptions"
                        :render-label="renderDropdownLabel"
                        @select="handleMenuSelect">
                        <n-button
                            size="small"
                            quaternary
                            class="name">
                            {{userInfo.name}}
                        </n-button>
                    </n-dropdown>
                    <n-dropdown
                        v-if="userInfo.avatar"
                        trigger="click"
                        :show-arrow="true"
                        :options="userMenuOptions"
                        :render-label="renderDropdownLabel"
                        @select="handleMenuSelect">
                        <n-avatar class="avatar" round :size="28" :src="userInfo.avatar"></n-avatar>
                    </n-dropdown>
                </div>
            </div>
        </n-layout-header>
        <template v-if="menuMap[menuAction].title">
            <div class="header-banner">
                <div class="wrapper">
                    <div class="title">{{menuMap[menuAction].title}}</div>
                    <div class="sub-title">{{menuMap[menuAction].subtitle}}</div>
                </div>
            </div>
            <n-divider/>
        </template>
    </div>
</template>

<style lang="less">
.header {
    .header-menu {
        .wrapper {
            .n-menu-item-content {
                padding: 0;
                margin: 2px 20px 0 0;
            }
        }
    }
}
</style>
<style lang="less" scoped>
.header {
    .header-menu {
        display: flex;
        flex-direction: row;
        justify-content: center;
        align-items: center;
        height: 64px;

        .wrapper {
            flex: 1;
            display: flex;
            flex-direction: row;
            justify-content: space-between;
            align-items: center;
        }

        .user {
            display: flex;
            align-items: center;

            .name {
                margin-left: 14px;
                cursor: pointer;
            }

            .avatar {
                margin-left: 14px;
                cursor: pointer;
            }
        }
    }
    .header-banner {
        display: flex;
        flex-direction: row;
        justify-content: center;
        align-items: center;

        .wrapper {
            flex: 1;
            padding: 24px 0 8px;

            .title {
                font-weight: 700;
                font-size: 38px;
                line-height: 64px;
            }

            .sub-title {
                font-size: 16px;
                line-height: 28px;
            }
        }
    }
}
</style>

<script lang="ts">
import {defineComponent, computed, h, ref, VNodeChild} from "vue";
import {useMessage, NButton} from 'naive-ui'
import { RouterLink } from 'vue-router'
import {EllipsisVertical} from "@vicons/ionicons5";
import {useThemeName} from '../store'
import cookie from "../utils/cookie";
import {UserStore} from "../store/user";

const userStore = UserStore()
const message = useMessage()

export default defineComponent({
    components: {EllipsisVertical},
    setup() {
        userStore.refresh()
        const userInfo = userStore.info
        const menuMap = computed(() => ({
            main: {
                label: '服务器管理',
                title: '服务器',
                subtitle: '在线管理和添加服务器',
            },
        }))
        const menuAction = ref('main')
        const menuChildren = (key, name = undefined) => {
            if (name === undefined) {
                name = key
            }
            let children: any = {
                default: () => menuMap.value[key]['label']
            }
            if (menuAction.value === key) {
                children = [
                    h(NButton, {
                        tertiary: true,
                        round: true,
                    }, children),
                ]
            }
            return {
                key,
                label: () => h(
                    RouterLink,
                    {
                        to: {
                            name,
                        }
                    },
                    children
                ),
            }
        }
        const menuOptions = [
            menuChildren('main'),
        ]
        const themeLabelMap = computed(() => ({
            dark: "浅色",
            light: "深色"
        }))
        const themeName = useThemeName()
        const handleThemeUpdate = () => {
            if (themeName.value === 'dark') {
                themeName.value = 'light'
            } else {
                themeName.value = 'dark'
            }
        }
        const userMenuOptions = ref([{
            label: '退出登录',
            key: "logout",
        }])
        const renderDropdownLabel = (option) => {
            if (option.key === 'logout') {
                return h(
                    'span',
                    {
                        style: 'color:rgb(248,113,113);height:34px;display:flex;align-items:center',
                    },
                    {
                        default: () => option.label as VNodeChild
                    }
                )
            }
            return option.label as VNodeChild
        }
        const handleMenuSelect = (key: string) => {
            if (key === 'logout') {
                cookie.remove('result_token')
                window.location.href = "/api/user/logout"
            } else {
                message.warning('未知操作')
            }
        }
        return {
            // user
            userInfo,
            // menu
            menuMap,
            menuAction,
            menuOptions,
            // theme
            themeName,
            themeLabelMap,
            handleThemeUpdate,
            //
            userMenuOptions,
            handleMenuSelect,
            renderDropdownLabel,
        }
    }
})
</script>
