<template>
  <div class="flex items-center justify-center h-full bg-base-200">
    <div class="card bg-base-100 w-full max-w-md shadow-xl">
      <div class="card-body">
        <div class="flex flex-col items-center mb-6">
          <Logo class="w-12 h-12 mb-2" />
          <h2 class="text-2xl font-bold">欢迎回来</h2>
          <p class="text-base-content/60 text-sm">请登录您的账户以继续</p>
        </div>

        <div class="tabs tabs-lift">
          <!-- Phone Login Tab -->
          <input type="radio" name="login_tabs" class="tab" :aria-label="t('login.loginByPhone')" 
                 :checked="form.login_type === 1" @change="form.login_type = 1" />
          <div class="tab-content bg-base-100 border-base-300 rounded-box p-6">
            <form class="space-y-4">
              <div class="form-control">
                <label class="label">
                  <span class="label-text">{{ t('common.phone') }}</span>
                </label>
                <input v-model="form.username" type="text" class="input input-bordered w-full" :placeholder="t('common.phone')" required />
              </div>

              <div class="form-control">
                <label class="label">
                  <span class="label-text">{{ t('login.code') }}</span>
                </label>
                <div class="join w-full">
                  <input v-model="form.code" type="text" class="input input-bordered join-item w-full" :placeholder="t('login.code')" required />
                  <button class="btn btn-primary join-item" type="button" @click="sendCode" :disabled="remaining != 0">
                    {{ remaining || t('login.sendCode') }}
                  </button>
                </div>
              </div>
              <button class="btn btn-primary w-full mt-4" type="button" @click="login" :disabled="loginLoading">
                <span v-if="loginLoading" class="loading loading-spinner"></span>
                {{ t('login.loginButton') }}
              </button>
            </form>
          </div>

          <!-- Password Login Tab -->
          <input type="radio" name="login_tabs" class="tab" :aria-label="t('login.loginByPassword')" 
                 :checked="form.login_type === 2" @change="form.login_type = 2" />
          <div class="tab-content bg-base-100 border-base-300 rounded-box p-6">
            <form class="space-y-4">
              <div class="form-control">
                <label class="label">
                  <span class="label-text">{{ t('common.username') }}</span>
                </label>
                <input v-model="form.username" type="text" class="input input-bordered w-full" :placeholder="t('common.username')" required />
              </div>

              <div class="form-control">
                <label class="label">
                  <span class="label-text">{{ t('common.password') }}</span>
                </label>
                <input v-model="form.password" type="password" class="input input-bordered w-full" :placeholder="t('common.password')" required />
              </div>
              <button class="btn btn-primary w-full mt-4" type="button" @click="login" :disabled="loginLoading">
                <span v-if="loginLoading" class="loading loading-spinner"></span>
                {{ t('login.loginButton') }}
              </button>
            </form>
          </div>
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
  login_type: 1,
  username: '',
  password: '',
  code: '',
})

async function login() {
  loginLoading.value = true
  const data = form.value
  if (data.login_type === 1) {
    delete data.password
  }
  if (data.login_type === 2) {
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
