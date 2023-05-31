<template>
    <n-form
            ref="formRef"
            :model="formData"
            :rules="formRules"
            size="large"
            label-placement="left"
            label-width="auto"
            require-mark-placement="right-hanging">
        <n-form-item path="userids" label="用户">
            <n-select v-model:value="formData.userids" multiple :options="options" :loading="loadIng"/>
        </n-form-item>
        <n-row :gutter="[0, 24]">
            <n-col :span="24">
                <div class="buttons">
                    <n-button :loading="loadIng" round type="primary" @click="handleSubmit">
                        确定
                    </n-button>
                </div>
            </n-col>
        </n-row>
    </n-form>
</template>

<style lang="less" scoped>
.buttons {
    display: flex;
    justify-content: flex-end;

    > button + button {
        margin-left: 12px;
    }
}
</style>
<script lang="ts">
import {computed, defineComponent, ref} from 'vue'
import {FormInst, FormItemRule, FormRules, useMessage} from 'naive-ui'
import utils from "../utils/utils";
import {getUserShareOptions, userShareServer} from "../api/modules/user";

import {ResultDialog} from "../api";

export default defineComponent({
    props: {
        shareServerId: {
            type: Number,
            required: true
        },
    },
    setup(props, {emit}) {
        const message = useMessage()
        const loadIng = ref(false);
        const formRef = ref<FormInst>()
        const formData = ref({})
        const ipsRef = ref(null)
        const options = ref([])
        const formRules: FormRules = {
            userids: [
                {
                    required: true,
                    validator(rule: FormItemRule, value: string) {
                        if (utils.isArray(value)) {
                            if (value.length === 0) {
                                formData.value.userids = ""
                                return new Error('请选择用户')
                            }
                        }
                        return true
                    },
                    trigger: ['input', 'blur']
                }
            ],
        }
        const handleSubmit = (e: MouseEvent) => {
            e.preventDefault()
            formRef.value?.validate((errors) => {
                if (errors) {
                    return;
                }
                if (loadIng.value) {
                    return
                }
                console.log(formData.value)
                const data = utils.cloneJSON(formData.value)
                let userids: any = data.userids
                delete data.userids
                if (!utils.isArray(userids)) {
                    userids = [userids]
                }

                userShareServer({user_ids:userids, server_id:props.shareServerId})
                .then(({msg}) => {
                    message.success(msg)
                    emit('onShareDone')
                })
                .catch(res => {
                    loadIng.value = false
                    ResultDialog(res)
                })
            })
        }

        const onLoad = () => {
            if (loadIng.value) {
                // if (showLoad === true) {
                //     loadShow.value = tip
                // }
                return
            }
            loadIng.value = true
            // loadShow.value = showLoad
            //
            // const params = {page: page.value, page_size: 10, key: searchKey.value}
            getUserShareOptions()
                .then(({data}) => {
                    options.value = data
                })
                .catch(res => {
                    // if (tip) {
                    //     if (dLog.value) {
                    //         dLog.value.destroy()
                    //         dLog.value = null
                    //     }
                    //     dLog.value = ResultDialog(res)
                    // }
                }).finally(() => {
                    loadIng.value = false
                })
        }
        onLoad()

        return {
            loadIng,
            formRef,
            formData,
            formRules,
            options,
            handleSubmit
        }
    }
})
</script>
