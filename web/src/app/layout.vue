<template>
    <n-layout class="root-layout">
        <Header/>
        <router-view/>
    </n-layout>
</template>

<script lang="ts">
import {onMounted, defineComponent} from 'vue'
import {useLoadingBar} from 'naive-ui'
import {loadingBarApiRef} from "../routes";
import Header from "../components/Header.vue";

export default defineComponent({
    components: {Header},
    setup() {
        const loadingBar = useLoadingBar()
        onMounted(() => {
            loadingBarApiRef.value = loadingBar
            loadingBar.finish()
        })
        return {}
    }
})
</script>
