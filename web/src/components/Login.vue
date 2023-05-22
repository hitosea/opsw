<template>
    <div class="login">
        <div class="wrapper">
            <div class="title">
                {{formData.type === 'reg' ? '注册' : '登录'}}
            </div>
            <n-form
                ref="formRef"
                :model="formData"
                :rules="formRules"
                :show-label="false"
                size="large">
                <n-form-item path="email" label="邮箱">
                    <n-input v-model:value="formData.email" placeholder="请输入邮箱地址"></n-input>
                </n-form-item>
                <n-form-item path="password" label="密码">
                    <n-input v-model:value="formData.password" type="password" placeholder="请输入密码"></n-input>
                </n-form-item>
                <n-form-item v-if="formData.type === 'reg'" path="password2" label="确认密码">
                    <n-input v-model:value="formData.password2" type="password" placeholder="请输入确认密码"></n-input>
                </n-form-item>
                <n-grid :cols="1">
                    <n-grid-item class="buttons">
                        <n-button :loading="loadIng" round type="primary" @click="handleSubmit">
                            {{formData.type === 'reg' ? '注册' : '登录'}}
                        </n-button>
                        <n-button :loading="loadIng" round type="default" @click="formData.type=formData.type === 'reg' ? 'login' : 'reg'">
                            {{formData.type === 'reg' ? '登录' : '注册'}}
                        </n-button>
                    </n-grid-item>
                </n-grid>
            </n-form>
        </div>
        <div class="policy">登录即表示您同意我们的服务条款和隐私政策。</div>
    </div>
</template>

<style lang="less" scoped>
.login {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;

    .wrapper {
        flex: 1;
        display: flex;
        justify-content: center;
        flex-direction: column;
        width: 90%;
        max-width: 220px;

        .title {
            text-align: center;
            font-size: 24px;
            margin-bottom: 24px;
        }

        .buttons {
            margin-top: 2px;
            display: flex;
            justify-content: center;

            > button + button {
                margin-left: 12px;
            }
        }
    }

    .policy {
        padding: 12px 32px 32px;
    }
}
</style>
<script lang="ts">
import {defineComponent, ref} from "vue";
import {LogoGithub, AddCircleOutline} from "@vicons/ionicons5";
import {FormInst, FormItemRule, FormRules, useMessage} from 'naive-ui'
import call from "../api";
import utils from "../utils/utils";


export default defineComponent({
    components: {
        LogoGithub,
        AddCircleOutline
    },
    setup() {
        const message = useMessage()
        const loadIng = ref<boolean>(false)
        const formRef = ref<FormInst | null>(null)
        const formData = ref({
            type: "login",
            email: "",
            password: "",
            password2: "",
        })
        const formRules: FormRules = {
            email: [
                {
                    validator(rule: FormItemRule, value: string) {
                        if (value) {
                            if (!utils.isEmail(value)) {
                                return new Error('邮箱格式错误')
                            }
                        } else {
                            return new Error('邮箱不能为空')
                        }
                        return true
                    },
                    required: true,
                    trigger: ['input', 'blur']
                }
            ],
            password: [
                {
                    validator(rule: FormItemRule, value: string) {
                        if (value) {
                            if (value.length < 6 || value.length > 20) {
                                return new Error('密码长度必须大于6位小于20位')
                            }
                        } else {
                            return new Error('密码不能为空')
                        }
                        return true
                    },
                    required: true,
                    trigger: ['input', 'blur']
                }
            ],
            password2: [
                {
                    validator(rule: FormItemRule, value: string) {
                        if (formData.value.type === "reg") {
                            if (value) {
                                if (value !== formData.value.password) {
                                    return new Error('两次密码不一致')
                                }
                            } else {
                                return new Error('确认密码不能为空')
                            }
                        }
                        return true
                    },
                    required: formData.value.type === "reg",
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
                //
                if (loadIng.value) {
                    return
                }
                loadIng.value = true
                call.post(formData.value.type === 'reg' ? 'user/reg' : 'user/login', formData.value)
                    .then(({msg}) => {
                        message.success(msg);
                        setTimeout(() => {
                            window.location.href = utils.removeURLParameter(window.location.href, ['result_code', 'result_msg'])
                        }, 300)
                    })
                    .catch(call.dialog)
                    .finally(() => {
                        loadIng.value = false
                    })
            })
        }
        return {
            loadIng,
            formRef,
            formData,
            formRules,

            handleSubmit,
        }
    }
})
</script>
