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
    <div class="navbar p-2 border-b border-base-content/5">
      <div class="flex-1">
        <span>主机列表</span>
      </div>
      <div class="flex-none">
        <button class="btn btn-outline btn-primary btn-sm" onclick="bing_dialog.showModal()">
          <Icon icon="mdi:link" />绑定
        </button>
        <dialog id="bing_dialog" class="modal">
          <div class="modal-box">
            <Bind />
            <div class="modal-action">
              <form method="dialog">
                <!-- if there is a button in form, it will close the modal -->
                <button class="btn">关闭</button>
              </form>
            </div>
          </div>
        </dialog>
      </div>
    </div>

    <Search :items="searchItems" @search="handleSearch" />


    <div class="p-2 border-b border-base-content/5">
      <template v-if="loading">
        <div class="grid place-items-center">
          <span>数据加载中</span>
          <span class="loading loading-xl loading-spinner text-primary" />
        </div>
      </template>
      <template v-else>
        <div v-if="data && data.list.length > 0"
          class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-6">
          <template v-for="item in data?.list">
            <Client :item="item" @refresh="handleSearch" />
          </template>
        </div>

        <div v-else class="grid place-items-center">
          <span>暂无数据，立即<button class="btn btn-link btn-lg" onclick="bing_dialog.showModal()">绑定</button></span>
        </div>

      </template>
    </div>

    <Pagination v-if="!loading" :total="data?.total || 0" @change="handlePageChange" />
  </div>
</template>