<template>
  <div>
    <div v-if="!userStore.isLogined">
      <button class="btn btn-sm btn-ghost p-1" @click="login">
        <Icon icon="mdi:login" class="w-5 h-5" />
      </button>
    </div>
    <div v-else class="dropdown dropdown-end">
      <div tabindex="0" class="avatar btn btn-ghost btn-sm p-1">
        <div class="w-6 rounded-full">
          <img :src="userStore.user?.avatar" />
        </div>
      </div>
      <div tabindex="-1">
        <ul tabindex="0" class="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-28">
          <li>
            <RouterLink :to="{ name: 'Dashboard' }">进入工作台</RouterLink>
          </li>
          <li>
            <a href="#">Settings</a>
          </li>
          <li>
            <a href="#">Logout</a>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from "@/stores/auth"
const userStore = useUserStore()
const router = useRouter()

const login = () => {
  router.push({ name: "Login" })
}
</script>
