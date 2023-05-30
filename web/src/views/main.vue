<template>
    <div class="main">
        <!-- 搜索 -->
        <div class="search">
            <div class="wrapper" :class="{loading: loadIng && loadShow}">
                <div class="search-box">
                    <div class="input-box">
                        <n-input round v-model:value="searchKey" clearable placeholder="">
                            <template #prefix>
                                <n-icon :component="SearchOutline"/>
                            </template>
                        </n-input>
                        <div class="reload" @click="onLoad(true, true)">
                            <Loading v-if="loadIng && loadShow"/>
                            <n-icon v-else>
                                <Reload/>
                            </n-icon>
                        </div>
                    </div>
                    <div class="interval"></div>
                    <n-button type="success" @click="createModal = true">
                        <template #icon>
                            <n-icon>
                                <AddOutline />
                            </n-icon>
                        </template>
                        添加服务器
                    </n-button>
                </div>
                <n-divider/>
            </div>
        </div>

        <!-- 列表 -->
        <div class="list">
            <div class="wrapper">
                <n-empty v-if="servers.length === 0" class="empty" size="huge" :description="emptyPrompt"/>
                <template v-else>
                    <div class="item nav">
                        <div class="name">服务器</div>
                        <div class="release">系统版本</div>
                        <div class="state">状态</div>
                        <div class="uptime">运行时间</div>
                        <div class="menu">操作</div>
                    </div>
                    <n-list hoverable :show-divider="false">
                        <n-list-item v-for="item in servers">
                            <div class="item">
                                <div class="name">
                                    <ul>
                                        <li>{{ item.ip }}</li>
                                        <li v-if="item.remark || item.hostname" class="remark">
                                            <EditText :value="item.remark" size="small" :params="{ip:item.ip}" :on-update="remarkUpdate">
                                                <div class="remark-edit">
                                                    <span>{{ item.remark || item.hostname }}</span>
                                                    <n-icon class="remark-icon" :component="Pencil"/>
                                                </div>
                                            </EditText>
                                        </li>
                                    </ul>
                                </div>
                                <div class="release">{{ item.platform }}-{{ item.platform_version }}</div>
                                <div class="state">
                                    <div v-if="stateJudge(item, 'Online')" class="run">
                                        <n-badge type="success" show-zero dot/>
                                    </div>
                                    <div v-else-if="stateJudge(item, 'Offline')" class="run">
                                        <n-badge color="grey" show-zero dot/>
                                    </div>
                                    <div v-else-if="stateLoading(item)" class="load">
                                        <Loading/>
                                    </div>
                                    <div class="text" :style="stateStyle(item)">{{ stateText(item) }}</div>
                                </div>
                                <div class="uptime">{{stateJudge(item, 'Online') ? uptime(item.current_info.time_since_uptime) : '-'}}</div>
                                <n-dropdown
                                    trigger="click"
                                    :show-arrow="true"
                                    :options="operationMenu"
                                    :render-label="operationLabel"
                                    @updateShow="operationShow($event, item)"
                                    @select="operationSelect($event, item)">
                                    <n-badge :show="!!item.upgrade" dot type="warning">
                                        <n-button quaternary :focusable="false" class="menu">
                                            <template #icon>
                                                <n-icon>
                                                    <EllipsisVertical/>
                                                </n-icon>
                                            </template>
                                        </n-button>
                                    </n-badge>
                                </n-dropdown>
                            </div>
                        </n-list-item>
                    </n-list>
                    <div v-if="pageCount > 1" class="page">
                        <n-pagination v-model:page="page" :page-count="pageCount" :disabled="loadIng" @update-page="onLoad(true, true)">
                            <template #prev></template>
                            <template #next></template>
                        </n-pagination>
                    </div>
                </template>
            </div>
        </div>

        <!-- 添加服务器 -->
        <n-modal v-model:show="createModal" :auto-focus="false">
            <n-card
                style="width:600px;max-width:90%"
                title="添加服务器"
                :bordered="false"
                size="huge"
                closable
                @close="createModal=false">
                <Create @onDone="createDone"/>
            </n-card>
        </n-modal>

        <!-- 日志 -->
        <n-modal v-model:show="logModal" :auto-focus="false">
            <n-card
                style="width:600px;max-width:90%"
                title="日志"
                :bordered="false"
                size="huge"
                closable
                @close="logModal=false">
                <Log :ip="logIp" v-model:show="logModal"/>
            </n-card>
        </n-modal>
    </div>
</template>

<script lang="ts">
import {defineComponent, h, ref, VNodeChild, watch} from "vue";
import Header from "../components/Header.vue";
import {useDialog, useMessage} from "naive-ui";
import {AddOutline, Pencil, EllipsisVertical, Reload, SearchOutline} from "@vicons/ionicons5";
import Loading from "../components/Loading.vue";
import Create from "../components/Create.vue";
import {ResultDialog} from "../api";
import utils from "../utils/utils";
import Log from "../components/Log.vue";
import {CONST} from "../store/constant";
import {getServerList, getServerOne, operationServer} from "../api/modules/server";
import {WsStore} from "../store/ws";
import EditText from "../components/EditText.vue";
import {Server} from "../api/interface/server";
import {GlobalStore, uptime} from "../store";

export default defineComponent({
    components: {
        EditText,
        Log,
        Create,
        Loading,
        Header,
        Reload,
        AddOutline,
        EllipsisVertical,
    },
    setup() {
        const message = useMessage()
        const dialog = useDialog()
        const globalStore = GlobalStore()
        const wsStore = WsStore()
        const dLog = ref(null);
        const createModal = ref(false);
        const logModal = ref(false);
        const logIp = ref("");
        const loadIng = ref(false);
        const loadShow = ref(false);
        const servers = ref<Server.Item[]>([])
        const page = ref(1)
        const pageCount = ref(0)
        const searchKey = ref("");
        const searchLast = ref("");
        const emptyPrompt = ref("暂无数据");

        watch(searchKey, (val) => {
            globalStore.timeout(600, "search").then(() => {
                if (val != searchLast.value) {
                    searchLast.value = val
                    onLoad(true, true)
                }
            })
        })

        const setServerItem = (ip, key, value) => {
            servers.value.forEach(item => {
                if (item.ip === ip) {
                    item[key] = value
                }
            })
        }
        const operationItem = ref({})
        const operationMenu = ref([
            {
                label: '应用商店',
                key: 'manage/panel/apps/all',
            }, {
                label: '网站',
                key: 'manage/panel/websites',
            }, {
                label: '数据库',
                key: 'manage/panel/databases/mysql',
            }, {
                label: '容器',
                key: 'manage/panel/containers/container',
            }, {
                label: '计划任务',
                key: 'manage/panel/cronjobs',
            }, {
                type: 'divider',
                key: 'd1'
            }, {
                label: '文件',
                key: 'manage/panel/hosts/files',
            }, {
                label: '监控',
                key: 'manage/panel/hosts/monitor/monitor',
            }, {
                label: '终端',
                key: 'manage/panel/hosts/terminal',
            }, {
                label: '防火墙',
                key: 'manage/panel/hosts/firewall/port',
            }, {
                type: 'divider',
                key: 'd1'
            }, {
                label: '日志',
                key: 'log',
            }, {
                label: '升级',
                key: 'upgrade',
                show: true,
            }, {
                label: '重置',
                key: 'reset',
                show: false,
            }, {
                label: '删除',
                key: "delete",
                color: 'rgb(248,113,113)',
            }
        ])

        const operationLabel = (option) => {
            if (option.disabled === true) {
                return option.label as VNodeChild
            }
            if (option.key === 'upgrade') {
                return `${option.label} (${operationItem.value.upgrade})` as VNodeChild
            } else if (typeof option.color === 'string') {
                return h(
                    'span',
                    {
                        style: `color:${option.color}`,
                    },
                    {
                        default: () => option.label as VNodeChild
                    }
                )
            }
            return option.label as VNodeChild
        }

        const operationShow = (show: boolean, item) => {
            if (show) {
                operationItem.value = item
                const disabled = !stateJudge(item, 'Online')
                operationMenu.value.forEach(v => {
                    if (utils.leftExists(v.key, 'manage/')) {
                        v['disabled'] = disabled
                    }
                })
                operationSetShow('upgrade', item.upgrade !== "")
                operationSetShow('reset', item.upgrade === "")
            }
        }

        const operationSetShow = (key: string, show: boolean) => {
            operationMenu.value.forEach(item => {
                if (item.key === key) {
                    item['show'] = show
                }
            })
        }

        const operationSelect = (key: string | number, item) => {
            if (key === 'log') {
                logIp.value = item.ip
                logModal.value = true
            } else if (key === 'upgrade') {
                const dd = dialog.warning({
                    title: '升级服务器',
                    content: () => h('div', [
                        h('div', '确定要升级服务器吗？'),
                        h('div', `当前版本：${item.version}，最新版本：${item.upgrade}`),
                    ]),
                    positiveText: '确定',
                    negativeText: '取消',
                    onPositiveClick: () => {
                        dd.loading = true
                        return operationInstance('upgrade', item.ip)
                    }
                })
            } else if (key === 'reset') {
                const dd = dialog.warning({
                    title: '重置服务器',
                    content: '确定要重置服务器吗？',
                    positiveText: '确定',
                    negativeText: '取消',
                    onPositiveClick: () => {
                        dd.loading = true
                        return operationInstance('reset', item.ip)
                    }
                })
            } else if (key === 'delete') {
                const dd = dialog.warning({
                    title: '删除服务器',
                    content: '确定要删除此服务器吗？',
                    positiveText: '确定',
                    negativeText: '取消',
                    onPositiveClick: () => {
                        dd.loading = true
                        return operationInstance('delete', item.ip)
                    }
                })
            } else if (utils.leftExists(key, 'manage/')) {
                const url = `${key}?ip=${item.ip}&theme=${globalStore.themeName}`
                window.open(url)
            } else {
                message.warning(`未知操作：${key}`)
            }
        }

        const operationInstance = (operation, ip) => {
            return new Promise((resolve) => {
                operationServer({operation, ip})
                    .then(({data, msg}) => {
                        message.success(msg)
                        setServerItem(ip, 'state', data.state)
                        onLoad(false, true)
                    })
                    .catch(ResultDialog)
                    .finally(() => {
                        resolve()
                    })
            })
        }

        const onLoad = (tip, showLoad) => {
            if (loadIng.value) {
                if (showLoad === true) {
                    loadShow.value = tip
                }
                return
            }
            loadIng.value = true
            loadShow.value = showLoad
            //
            const params = {page: page.value, page_size: 10, key: searchKey.value}
            getServerList(params)
                .then(({data}) => {
                    if (params.key != searchKey.value) {
                        return
                    }
                    if (utils.parseInt(data.total) === 0) {
                        if (tip === true) {
                            emptyPrompt.value = params.key ? "暂无搜索结果" : "暂无数据"
                            message.warning(emptyPrompt.value)
                        }
                        page.value = 1
                        pageCount.value = 0
                        servers.value = []
                        return
                    }
                    page.value = data.page
                    pageCount.value = data.page_count
                    servers.value = data.data
                })
                .catch(res => {
                    if (tip) {
                        if (dLog.value) {
                            dLog.value.destroy()
                            dLog.value = null
                        }
                        dLog.value = ResultDialog(res)
                    }
                }).finally(() => {
                loadIng.value = false
            })
        }
        onLoad(false, true)

        const createDone = () => {
            createModal.value = false
            onLoad(true, true)
        }

        const systemsFormat = (systems, key, def = "-") => {
            systems = utils.jsonParse(systems)
            return systems[key] || def
        }

        const stateJudge = (item, state) => {
            const ss = stateText(item)
            if (utils.isArray(state)) {
                return state.includes(ss)
            }
            return ss === state
        }
        const stateLoading = (item) => {
            const state = stateText(item)
            return state.slice(-3) === 'ing'
        }
        const stateStyle = (item) => {
            const state = stateText(item)
            switch (state) {
                case 'Timeout':
                case 'Error':
                    return {
                        color: 'rgb(248,113,113)'
                    }
                default:
                    return {}
            }
        }
        const stateText = (item) => {
            return item.state || 'Unknown'
        }

        const remarkUpdate = (value, params) => {
            return new Promise((resolve, reject) => {
                operationServer(Object.assign(params, {
                    operation: 'remark',
                    remark: value
                })).then(_ => {
                    message.success("修改成功")
                    setServerItem(params.ip, 'remark', value)
                    resolve()
                }).catch(res => {
                    ResultDialog(res)
                    reject()
                })
            })
        }

        wsStore.listener("main", data => {
            if (data.type === CONST.WsIsServer) {
                globalStore.timeout(3000, "main", data.cid).then(_ => {
                    if (servers.value.find(item => item.id === data.cid)) {
                        getServerOne({
                            id: data.cid
                        }).then(({data}) => {
                            const index = servers.value.findIndex(item => item.id === data.id)
                            if (index !== -1) {
                                servers.value.splice(index, 1, data)
                            }
                        })
                    }
                })
            }
        })

        return {
            Pencil,
            SearchOutline,

            createModal,
            createDone,

            logModal,
            logIp,

            loadIng,
            loadShow,

            page,
            pageCount,
            servers,
            searchKey,
            emptyPrompt,

            operationMenu,
            operationLabel,
            operationShow,
            operationSelect,

            onLoad,
            systemsFormat,

            stateJudge,
            stateLoading,
            stateStyle,
            stateText,

            remarkUpdate,

            uptime,
        };
    }
})
</script>

<style lang="less" scoped>
.main {
    padding-bottom: 36px;
    .search {
        display: flex;
        flex-direction: row;
        justify-content: center;
        align-items: center;

        .wrapper {
            flex: 1;

            .search-box {
                display: flex;
                align-items: center;
                flex-direction: row;
                justify-content: space-between;
            }

            &.loading,
            &:hover {
                .input-box {
                    .reload {
                        > i,
                        .loading {
                            opacity: 1;
                        }
                    }
                }
            }

            .input-box {
                display: flex;
                align-items: center;

                .reload {
                    margin: 0 32px 0 16px;
                    width: 30px;
                    height: 30px;
                    display: flex;
                    align-items: center;
                    justify-items: center;

                    > i,
                    .loading {
                        transition: all 0.3s;
                        opacity: 0.5;
                        font-size: 20px;
                        width: 20px;
                        height: 20px;
                    }
                }
            }

            .interval {
                flex: 1;
            }
        }
    }

    .list {
        display: flex;
        flex-direction: row;
        justify-content: center;
        align-items: center;

        .wrapper {
            flex: 1;

            .empty {
                margin: 120px 0;
            }

            > ul {
                background-color: transparent;

                > li {
                    border-radius: 18px;
                }
            }

            .item {
                display: flex;
                align-items: center;
                list-style: none;
                white-space: nowrap;
                justify-content: space-between;
                padding: 12px 0;

                &.nav {
                    font-size: 16px;
                    font-weight: 600;
                    margin-bottom: 8px;
                    padding-left: 20px;
                    padding-right: 20px;
                }

                .name {
                    width: 30%;

                    &:hover {
                        .remark-icon {
                            display: block !important;
                        }
                    }

                    ul {
                        display: flex;
                        flex-direction: column;
                        justify-content: center;
                        padding: 0;
                        margin: 0;

                        > li {
                            list-style: none;
                            padding: 0 6px 0 0;
                            margin: 0;
                            overflow: hidden;
                            text-overflow: ellipsis;
                            white-space: nowrap;

                            &.remark {
                                font-weight: normal;
                                opacity: 0.5;
                                user-select: auto;

                                .remark-edit {
                                    display: flex;
                                    align-items: center;
                                    .remark-icon {
                                        display: none;
                                        margin-left: 6px;
                                    }
                                }
                            }
                        }
                    }
                }

                .release {
                    width: 20%;
                }

                .state {
                    width: 20%;
                    display: flex;
                    align-items: center;

                    .load,
                    .run {
                        flex-shrink: 0;
                        width: 16px;
                        height: 16px;
                        margin-right: 6px;
                        display: flex;
                        align-items: center;
                        justify-content: center;
                    }

                    .run {
                        .n-badge {
                            transform: scale(1.2);
                        }
                    }

                    .text {
                        flex: 1;
                    }
                }

                .uptime {
                    width: 20%;
                }

                .menu {
                    min-width: 32px;
                    padding: 0;
                }
            }

            .page {
                margin-top: 30px;
                display: flex;
                flex-direction: column;
                align-items: center;
            }
        }
    }
}
</style>
