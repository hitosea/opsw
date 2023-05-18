<template>
    <div class="main">
        <Header/>
    </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import Header from "../components/Header.vue";

export default defineComponent({
    components: {Header},
    setup() {
        return {};
    }
})
</script>

<style lang="less" scoped>
.main {

}
</style>
