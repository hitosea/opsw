<template>
    <div class="main">
        <!-- 头部 -->
        <Header/>

        <!-- 搜索 -->
        <div class="search">
            <div class="wrapper" :class="{loading: loadIng && loadShow}">
                <div class="search-box">
                    <div class="input-box">
                        <n-input round v-model:value="searchKey" placeholder="">
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
                    <n-button type="success" :render-icon="addIcon" @click="createModal = true">
                        添加服务器
                    </n-button>
                </div>
                <n-divider/>
            </div>
        </div>

        <!-- 列表 -->
        <div class="list">
            <div class="wrapper">
                <n-empty v-if="searchList.length === 0" class="empty" size="huge" description="没有服务器"/>
                <template v-else>
                    <div class="item nav">
                        <div class="name">服务器</div>
                        <div class="release">系统版本</div>
                        <div class="state">状态</div>
                        <div class="menu">操作</div>
                    </div>
                    <n-list hoverable :show-divider="false">
                        <n-list-item v-for="item in searchList">
                            <div class="item">
                                <div class="name">
                                    <ul>
                                        <li>{{ item['ip'] }}</li>
                                        <li v-if="item['remark'] || item['hostname']" class="remark">{{ item['remark'] || item['hostname'] }}</li>
                                    </ul>
                                </div>
                                <div class="release">{{ item['platform'] }}-{{ item['platform_version'] }}</div>
                                <div class="state">
                                    <div v-if="stateJudge(item, 'Online')" class="run">
                                        <n-badge type="success" show-zero dot />
                                    </div>
                                    <div v-else-if="stateJudge(item, 'Offline')" class="run">
                                        <n-badge color="grey" show-zero dot />
                                    </div>
                                    <div v-else-if="stateLoading(item)" class="load">
                                        <Loading/>
                                    </div>
                                    <div class="text" :style="stateStyle(item)">{{ stateText(item) }}</div>
                                </div>
                                <n-dropdown
                                    trigger="click"
                                    :show-arrow="true"
                                    :options="operationMenu"
                                    :render-label="operationLabel"
                                    @updateShow="operationShow($event, item)"
                                    @select="operationSelect($event, item)">
                                    <n-button quaternary class="menu">
                                        <template #icon>
                                            <n-icon>
                                                <EllipsisVertical/>
                                            </n-icon>
                                        </template>
                                    </n-button>
                                </n-dropdown>
                            </div>
                        </n-list-item>
                    </n-list>
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
import {computed, defineComponent, h, ref, VNodeChild} from "vue";
import Header from "../components/Header.vue";
import {useDialog, useMessage} from "naive-ui";
import {AddOutline, EllipsisVertical, Reload, SearchOutline} from "@vicons/ionicons5";
import Loading from "../components/Loading.vue";
import Create from "../components/Create.vue";
import call from "../store/call";
import utils from "../store/utils";
import Log from "../components/Log.vue";

export default defineComponent({
    components: {
        Log,
        Create,
        Loading,
        Header,
        Reload,
        EllipsisVertical,
    },
    computed: {
        SearchOutline() {
            return SearchOutline
        }
    },
    setup() {
        const message = useMessage()
        const dialog = useDialog()
        const dLog = ref(null);
        const createModal = ref(false);
        const logModal = ref(false);
        const logIp = ref("");
        const loadIng = ref(false);
        const loadShow = ref(false);
        const items = ref([])
        const searchKey = ref("");
        const searchList = computed(() => {
            if (searchKey.value === "") {
                return items.value
            }
            return items.value.filter(item => `${item.ip} ${item.remark}`.indexOf(searchKey.value) !== -1)
        })

        const setItemState = (ip, state) => {
            items.value.forEach(item => {
                if (item.ip === ip) {
                    item.state = state
                }
            })
        }
        const operationItem = ref({})
        const operationMenu = ref([
            {
                label: '详情',
                key: 'info',
            }, {
                label: '日志',
                key: 'log',
            }, {
                type: 'divider',
                key: 'd1'
            }, {
                label: '升级',
                key: 'upgrade',
                disabled: false,
            }, {
                label: '删除',
                key: "delete",
            }
        ])

        const operationLabel = (option) => {
            if (option.disabled === true) {
                return option.label as VNodeChild
            }
            if (option.key === 'upgrade') {
                return `${option.label} (${operationItem.value.upgrade})` as VNodeChild
            } else if (option.key === 'delete') {
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

        const operationShow = (show: boolean, item) => {
            if (show) {
                operationItem.value = item
                operationSetDisabled('upgrade', item.upgrade === "")
            }
        }

        const operationSetDisabled = (key: string, disabled: boolean) => {
            operationMenu.value.forEach(item => {
                if (item.key === key) {
                    item['disabled'] = disabled
                }
            })
        }

        const operationSelect = (key: string | number, item) => {
            if (key === 'info') {
                message.warning("查看详情")
            } else if (key === 'log') {
                logIp.value = item.ip
                logModal.value = true
            } else if (key === 'upgrade') {
                const dd = dialog.warning({
                    title: '升级服务器',
                    content: '确定要升级服务器吗？',
                    positiveText: '确定',
                    negativeText: '取消',
                    onPositiveClick: () => {
                        dd.loading = true
                        return operationInstance('upgrade', item.ip)
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
            }
        }

        const operationInstance = (operation, ip) => {
            return new Promise((resolve) => {
                call.get("server/operation", {operation, ip})
                    .then(({data, msg}) => {
                        message.success(msg)
                        setItemState(ip, data.state)
                        onLoad(false, true)
                    })
                    .catch(call.dialog)
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
            call.get("server/list").then(({data}) => {
                if (!utils.isArray(data.list)) {
                    if (tip === true) {
                        message.warning("暂无数据")
                    }
                    items.value = []
                    return
                }
                items.value = data.list
            }).catch(res => {
                if (tip) {
                    if (dLog.value) {
                        dLog.value.destroy()
                        dLog.value = null
                    }
                    dLog.value = call.dialog(res)
                }
            }).finally(() => {
                loadIng.value = false
            })
        }
        onLoad(false, true)

        const addIcon = () => {
            return h(AddOutline);
        }
        const createDone = () => {
            createModal.value = false
            onLoad(true, true)
        }

        const systemsFormat = (systems, key, def = "-") => {
            systems = utils.jsonParse(systems)
            return systems[key] || def
        }

        const stateJudge = (item, state) => {
            return stateText(item) === state
        }
        const stateLoading = (item) => {
            const state = stateText(item)
            return state.slice(-3) === 'ing'
        }
        const stateStyle = (item) => {
            const state = stateText(item)
            switch (state) {
                case 'Failed':
                case 'Unknown':
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

        return {
            createModal,
            createDone,

            logModal,
            logIp,

            loadIng,
            loadShow,

            searchKey,
            searchList,

            operationMenu,
            operationLabel,
            operationShow,
            operationSelect,

            onLoad,
            addIcon,
            systemsFormat,

            stateJudge,
            stateLoading,
            stateStyle,
            stateText,
        };
    }
})
</script>

<style lang="less" scoped>
.main {
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
                    width: 40%;
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
                                &:hover {
                                    opacity: 1;
                                }
                            }
                        }
                    }
                }

                .release {
                    width: 30%;
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

                .menu {
                    min-width: 32px;
                    padding: 0;
                }
            }
        }
    }
}
</style>
