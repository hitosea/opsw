<template>
    <div class="main">
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
        </div>

        <!-- 列表 -->
        <div class="list">
            <div class="wrapper">
                <n-empty v-if="searchList.length === 0" class="empty" size="huge" description="没有服务器"/>
                <template v-else>
                    <div class="item nav">
                        <div class="name">工作区名称</div>
                        <div class="release">系统版本</div>
                        <div class="state">状态</div>
                        <div class="menu">操作</div>
                    </div>
                    <n-list hoverable :show-divider="false">
                        <n-list-item v-for="item in searchList">
                            <div class="item">

                            </div>
                        </n-list-item>
                    </n-list>
                </template>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import {computed, defineComponent, h, ref} from "vue";
import Header from "../components/Header.vue";
import {useDialog, useMessage} from "naive-ui";
import {AddOutline, Reload, SearchOutline} from "@vicons/ionicons5";
import Loading from "../components/Loading.vue";
import Create from "../components/Create.vue";

export default defineComponent({
    components: {Create, Loading, Header, Reload},
    computed: {
        SearchOutline() {
            return SearchOutline
        }
    },
    setup() {
        const message = useMessage()
        const dialog = useDialog()
        const createModal = ref(false);
        const loadIng = ref(false);
        const loadShow = ref(false);
        const items = ref([])
        const searchKey = ref("");
        const searchList = computed(() => {
            if (searchKey.value === "") {
                return items.value
            }
            return items.value.filter(item => item.name.indexOf(searchKey.value) !== -1)
        })

        const onLoad = (show = false, ing = false) => {
            loadShow.value = show
            loadIng.value = ing
            setTimeout(() => {
                loadShow.value = !show
                loadIng.value = !ing
            }, 2000)
        }
        const addIcon = () => {
            return h(AddOutline);
        }
        const createDone = () => {
            createModal.value = false
            onLoad(true, true)
        }

        return {
            createModal,
            loadIng,
            loadShow,
            searchKey,
            searchList,

            onLoad,
            addIcon,
            createDone
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

            }

        }
    }
}
</style>
