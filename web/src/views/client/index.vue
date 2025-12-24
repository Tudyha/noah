<script lang="ts" setup>
import { useRequest } from 'vue-hooks-plus'
import { getClientPage } from "@/api/client";
import Bind from './components/bind.vue'
import Client from './components/client.vue'
import type { SearchItem } from '@/types'

const searchItems: SearchItem[] = [
  {
    type: "input",
    label: "主机名称",
    key: "username",
  },
  {
    type: "input",
    label: "IP地址",
    key: "username",
    placeholder: "请输入IP地址"
  },
  {
    type: "select",
    label: "在线状态",
    key: "status",
    width: "w-20",
    options: [
      {
        value: "1",
        label: "在线"
      },
      {
        value: "2",
        label: "离线"
      }
    ]
  }
]

const { data, loading, run } = useRequest(getClientPage)

const handlePageChange = ({ page, pageSize }: { page: number, pageSize: number }) => {
  console.log(page, pageSize)
}

const handleSearch = () => {
  run()
}

</script>

<template>
  <div>
    <div class="navbar border-b border-base-content/5">
      <div class="flex-1">
        <h1>主机列表</h1>
      </div>
      <button
        class="btn btn-primary btn-sm btn-soft flex-shrink-0 transition-all duration-200 focus-visible:ring focus-visible:ring-primary/30"
        onclick="bing_dialog.showModal()" aria-label="绑定主机">
        <Icon icon="mdi:link" />绑定
      </button>
    </div>

    <Search :items="searchItems" @search="handleSearch" />

    <!-- 绑定模态框 -->
    <dialog id="bing_dialog" class="modal">
      <div class="modal-box">
        <Bind />
        <div class="modal-action">
          <form method="dialog">
            <button class="btn">关闭</button>
          </form>
        </div>
      </div>
    </dialog>

    <div class="p-2">
      <Transition name="fade" mode="out-in">
        <template v-if="loading">
          <div class="grid place-items-center h-48">
            <span class="text-lg font-semibold text-base-content/70 mb-4">正在获取主机数据...</span>
            <span class="loading loading-spinner loading-lg text-primary animate-pulse" />
          </div>
        </template>
        <template v-else>
          <TransitionGroup v-if="data && data.list.length > 0" name="list-fade" tag="div"
            class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-6">
            <Client v-for="item in data?.list" :key="item.id" :item="item" @refresh="handleSearch" />
          </TransitionGroup>

          <div v-else class="grid place-items-center h-48">
            <Icon icon="mdi:server-off" class="w-16 h-16 text-base-content/30 mb-4" />
            <span class="text-lg font-semibold text-base-content/70 mb-4">暂未发现已绑定的主机</span>
            <button
              class="btn btn-primary btn-md btn-soft transition-all duration-200 focus-visible:ring focus-visible:ring-primary/30"
              onclick="bing_dialog.showModal()" aria-label="立即绑定主机">
              <Icon icon="mdi:link" />立即绑定
            </button>
          </div>
        </template>
      </Transition>
    </div>

    <Pagination v-if="!loading" :total="data?.total || 0" @change="handlePageChange" />
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity .2s ease, transform .2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(4px);
}

.list-fade-enter-active,
.list-fade-leave-active {
  transition: opacity .25s ease, transform .25s ease;
}

.list-fade-enter-from,
.list-fade-leave-to {
  opacity: 0;
  transform: scale(0.98) translateY(6px);
}
</style>
