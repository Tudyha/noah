<template>
  <div class="flex items-center justify-center h-full">
    <div class="card bg-base-100 w-xl shadow-sm">
      <div class="card-body">
        <div class="tabs tabs-lift">

          <!-- tab -->
          <label class="tab gap-2">
            <input type="radio" name="loginByPhone" value="1" v-model="form.loginType"
              :checked="form.loginType === 1" />
            <Icon icon="material-symbols:perm-phone-msg" />
            {{ t('login.loginByPhone') }}
          </label>

          <!-- form -->
          <div class="tab-content bg-base-100 border-base-300 p-2">
            <form class="fieldset bg-base-200 border-base-300 rounded-box border p-4">
              <fieldset class="fieldset">
                <label class="label">{{ t('common.phone') }}</label>
                <input v-model="form.username" type="text" class="input validator" :placeholder="t('common.phone')"
                  required />
                <p class="validator-hint hidden">{{ t('login.phoneValidator') }}</p>
              </fieldset>

              <fieldset class="fieldset">
                <label class="label">{{ t('login.code') }}</label>
                <div class="join">
                  <input v-model="form.code" type="text" class="input validator w-40" :placeholder="t('login.code')"
                    required />
                  <button class="btn btn-primary join-item" @click="sendCode" :disabled="remaining != 0">{{
                    remaining ||
                    t('login.sendCode') }}</button>
                </div>
              </fieldset>
              <button class="btn btn-primary mt-6" type="button" @click="login">{{ t('login.loginButton') }}</button>
            </form>
          </div>

          <!-- tab -->
          <label class="tab gap-2">
            <input type="radio" name="loginByPassword" value="2" v-model="form.loginType"
              :checked="form.loginType === 2" />
            <Icon icon="carbon:password" />
            {{ t('login.loginByPassword') }}
          </label>
          <!-- form -->
          <div class="tab-content bg-base-100 border-base-300 p-2">
            <form class="fieldset bg-base-200 border-base-300 rounded-box border p-4">
              <fieldset class="fieldset">
                <label class="label">{{ t('common.username') }}</label>
                <input v-model="form.username" type="text" class="input validator" :placeholder="t('common.username')"
                  required />
                <p class="validator-hint hidden">{{ t('login.usernameValidator') }}</p>
              </fieldset>

              <fieldset class="fieldset">
                <label class="label">{{ t('common.password') }}</label>
                <input v-model="form.password" type="password" class="input validator"
                  :placeholder="t('common.password')" required />
                <p class="validator-hint hidden">{{ t('login.passwordValidator') }}</p>
              </fieldset>
              <button class="btn btn-primary mt-6" type="button" @click="login">{{ t('login.loginButton') }}</button>
            </form>
          </div>

        </div>
        <div class="card-actions justify-end">
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { type LoginRequest } from '@/types'
import { useUserStore } from '@/stores/auth'

const userStore = useUserStore()
const router = useRouter()
const loginLoading = ref(false)
const countdown = shallowRef(0)
const { start, remaining } = useCountdown(countdown, {})

const { t } = useI18n()
const form = ref<LoginRequest>({
  loginType: 1,
  username: '',
  password: '',
  code: '',
})

async function login() {
  loginLoading.value = true
  const data = form.value
  if (data.loginType === 1) {
    delete data.password
  }
  if (data.loginType === 2) {
    delete data.code
  }
  userStore.login(data).then(() => {
    // 登录成功
    router.push({ name: 'Dashboard' })
  }).catch((e) => {
    // 登录失败
    console.error(e)
  })
  loginLoading.value = false
}

const sendCode = async () => {
  start(60)
}
</script>
