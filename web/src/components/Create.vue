<template>
    <n-form
            ref="formRef"
            :model="formData"
            :rules="formRules"
            size="large"
            label-placement="left"
            label-width="auto"
            require-mark-placement="right-hanging">
        <n-form-item path="ips" label="服务器IP">
            <n-dynamic-tags v-if="typeof formData.ips === 'object'" v-model:value="formData.ips" />
            <n-input v-else v-model:value="formData.ips" type="text" placeholder="请输入服务器IP地址" @blur="ipsBlur" @keyup.enter="ipsBlur"/>
        </n-form-item>
        <n-form-item v-show="advancedShow" path="username" label="用户名">
            <n-input v-model:value="formData.username" placeholder="SSH登录用户名，默认：root"/>
        </n-form-item>
        <n-form-item path="password" label="密码">
            <n-input v-model:value="formData.password" placeholder="SSH登录密码，默认：(空)"/>
        </n-form-item>
        <n-form-item v-show="advancedShow" path="port" label="端口">
            <n-input v-model:value="formData.port" placeholder="SSH登录端口，默认：22"/>
        </n-form-item>
        <n-form-item v-show="advancedShow" path="remark" label="备注">
            <n-input v-model:value="formData.remark" placeholder=""/>
        </n-form-item>
        <n-row :gutter="[0, 24]">
            <n-col :span="24">
                <div class="buttons">
                    <n-button round quaternary type="default" @click="advancedShow=!advancedShow">
                        {{advancedShow ? '收起' : '高级选项'}}
                    </n-button>
                    <n-button :loading="loadIng > 0" round type="primary" @click="handleSubmit">
                        添加
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
import {
    FormInst,
    FormItemRule,
    FormRules, useDialog, useMessage
} from 'naive-ui'
import call from "../store/call";
import utils from "../store/utils";

interface ModelType {
    ips: string | object
    username?: string | null
    password?: string | null
    port?: string | null
    remark?: string | null
}

export default defineComponent({
    setup(props, {emit}) {
        const message = useMessage()
        const loadIng = ref<number>(0)
        const formRef = ref<FormInst | null>(null)
        const formData = ref<ModelType>({
            ips: "",
        })
        const ipsRef = ref(null)
        const ipsComputed = computed(() => {
            if (ipsRef.value == null) {
                return []
            }
            return ipsRef.value.map(item => {
                return {
                    label: item['ip'],
                    value: item['ip']
                }
            })
        })
        const advancedShow = ref<boolean>(false)
        const formRules: FormRules = {
            ips: [
                {
                    required: true,
                    validator(rule: FormItemRule, value: string) {
                        if (utils.isArray(value)) {
                            if (value.length === 0) {
                                formData.value.ips = ""
                                return new Error('请输入服务器IP地址')
                            }
                            for (let i = 0; i < value.length; i++) {
                                if (!utils.isIpv4(value[i])) {
                                    return new Error(`第${i + 1}个IP地址不合法`)
                                }
                            }
                        } else if (!utils.isIpv4(value)) {
                            return new Error('请输入有效的IP地址')
                        }
                        return true
                    },
                    trigger: ['input', 'blur']
                }
            ],
        }
        const ipsBlur = () => {
            const value = `${formData.value.ips}`.trim()
            if (utils.isIpv4(value)) {
                formData.value.ips = [value]
            }
        }
        const handleSubmit = (e: MouseEvent) => {
            e.preventDefault()
            formRef.value?.validate((errors) => {
                if (errors) {
                    return;
                }
                //
                if (loadIng.value > 0) {
                    return
                }
                const data = utils.cloneJSON(formData.value)
                let ips: any = data.ips
                delete data.ips
                if (!utils.isArray(ips)) {
                    ips = [ips]
                }
                ips.forEach(ip => {
                    loadIng.value++
                    call.post('server/create', Object.assign(data, {ip}))
                        .then(({msg}) => {
                            loadIng.value--
                            if (loadIng.value === 0) {
                                message.success(msg)
                                emit('onDone')
                            }
                        })
                        .catch(res => {
                            loadIng.value--
                            if (loadIng.value === 0) {
                                call.dialog(res)
                            }
                        })
                })
            })
        }
        return {
            advancedShow,
            loadIng,
            formRef,
            formData,
            formRules,
            ipsComputed,
            ipsBlur,
            handleSubmit
        }
    }
})
</script>
