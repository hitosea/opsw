<template>
    <n-config-provider :hljs="hljs">
        <div class="log">
            <n-scrollbar ref="nRef" class="log-scrollbar" x-scrollable>
                <n-code :code="content" language="bash" />
            </n-scrollbar>
            <div class="footer">
                <n-button :loading="loading" @click="getData(true)">刷新</n-button>
            </div>
        </div>
    </n-config-provider>
</template>

<style lang="less">
.log {
    .log-scrollbar {
        max-height: 300px;
        .__code__ {
            font-size: 14px;
            line-height: 1.35;
        }
    }
}
</style>
<style lang="less" scoped>
.log {
    .footer {
        display: flex;
        align-items: center;
        justify-content: center;
        margin-top: 26px;
    }
}
</style>
<script lang="ts">
import {defineComponent, onBeforeUnmount, nextTick, ref} from 'vue'
import hljs from 'highlight.js/lib/core'
import bash from 'highlight.js/lib/languages/bash'
import {ResultDialog} from "../api";
import {getServerLog} from "../api/modules/server";

hljs.registerLanguage('bash', bash)

export default defineComponent({
    props: {
        ip: {
            type: String,
            required: true
        },
        show: {
            type: Boolean,
        },
    },
    setup(props, {emit}) {
        const nRef = ref(null);
        const dLog = ref(null);
        const loading = ref(false);
        const content = ref("");

        const scrollToBottom = () => {
            const {scrollbarInstRef} = nRef.value
            const {containerRef, contentRef} = scrollbarInstRef
            if (containerRef && contentRef) {
                const containerHeight = containerRef.offsetHeight
                const containerScrollTop = containerRef.scrollTop
                const contentHeight = contentRef.offsetHeight
                const scrollBottom = contentHeight - containerScrollTop - containerHeight
                return scrollBottom < 20
            }
            return true
        }

        const getData = (manual=false) => {
            if (loading.value) {
                return
            }
            loading.value = true
            getServerLog({ip: props.ip})
                .then(({data}) => {
                    const autoBottom = manual === true || scrollToBottom()
                    content.value = data.log
                    autoBottom && nextTick(() => {
                        nRef.value?.scrollTo({position: 'bottom'})
                    })
                })
                .catch(res => {
                    if (dLog.value) {
                        dLog.value.destroy()
                        dLog.value = null
                    }
                    dLog.value = ResultDialog(res, {
                        onPositiveClick: () => {
                            emit("update:show", false)
                        }
                    })
                })
                .finally(() => {
                    loading.value = false
                })
        }
        getData()
        const getInter = setInterval(getData, 15 * 1000)

        onBeforeUnmount(() => {
            clearInterval(getInter)
        })

        return {
            hljs,
            content,
            loading,
            nRef,

            getData
        }
    }
})
</script>
