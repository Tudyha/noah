<script lang="ts" setup>
import { getAgentPage, getV2raySubscribe } from "@/api/agent";
import { usePageList } from "@/composables/usePageList";
import Bind from './components/bind.vue'
import Agent from './components/agent.vue'
import type { SearchItem } from '@/types'

const searchItems: SearchItem[] = [
  {
    type: "input",
    label: "主机名称",
    key: "hostname",
  },
  {
    type: "input",
    label: "IP地址",
    key: "ip",
    placeholder: "请输入IP地址"
  },
  {
    type: "select",
    label: "在线状态",
    key: "status",
    width: "w-32",
    placeholder: "请选择状态",
    default: "1",
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

const defaultFilters = searchItems.reduce<Record<string, any>>((acc, item) => {
  if (item.default !== undefined) {
    acc[item.key] = item.default
  }
  return acc
}, {})

const { data, loading, currentPage, pageSize, search, changePage } = usePageList(getAgentPage, {
  defaultFilters
})

// 多选逻辑
const selectedIds = ref<number[]>([])
const isAllSelected = computed(() => {
  if (!data.value?.list.length) return false
  return data.value.list.every(item => selectedIds.value.includes(item.id))
})

const toggleSelect = (id: number) => {
  const index = selectedIds.value.indexOf(id)
  if (index > -1) {
    selectedIds.value.splice(index, 1)
  } else {
    selectedIds.value.push(id)
  }
}

const toggleSelectAll = () => {
  if (isAllSelected.value) {
    selectedIds.value = []
  } else {
    selectedIds.value = data.value?.list.map(item => item.id) || []
  }
}

const clearSelection = () => {
  selectedIds.value = []
}

const handlePageChange = ({ page, pageSize }: { page: number; pageSize: number }) => {
  changePage(page, pageSize)
}

const handleSearch = (form: Record<string, any>) => {
  search(form)
}

const v2rayModalRef = ref<HTMLDialogElement | null>(null)
const v2rayLink = ref<string>("")
const handleV2raySubscribe = async () => {
  v2rayLink.value = await getV2raySubscribe(selectedIds.value)
  v2rayModalRef.value?.showModal()
}

</script>

<template>
  <div class="space-y-4">
    <!-- header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-4">
        <h2 class="text-2xl font-bold text-base-content">主机列表</h2>
        <div v-if="data?.list.length"
          class="flex items-center gap-2 px-3 py-1 bg-base-200/50 rounded-lg border border-base-300">
          <input type="checkbox" class="checkbox checkbox-primary checkbox-xs" :checked="isAllSelected"
            @change="toggleSelectAll" />
          <span class="text-xs font-medium text-base-content/70 select-none">全选</span>
        </div>
      </div>
      <button class="btn btn-primary btn-sm shrink-0 shadow-sm" onclick="bind_dialog.showModal()" aria-label="绑定主机">
        <Icon icon="mdi:link" class="w-4 h-4" /> 绑定主机
      </button>
    </div>

    <!-- 搜索栏 -->
    <Search :items="searchItems" @search="handleSearch" />

    <!-- 绑定模态框 -->
    <dialog id="bind_dialog" class="modal">
      <div class="modal-box">
        <Bind />
        <div class="modal-action">
          <form method="dialog">
            <button class="btn">关闭</button>
          </form>
        </div>
      </div>
    </dialog>

    <!-- 主机列表 -->
    <div>
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
            <Agent v-for="item in data?.list" :key="item.id" :item="item" :selected="selectedIds.includes(item.id)"
              @toggle-select="toggleSelect(item.id)" @refresh="handleSearch" />
          </TransitionGroup>

          <div v-else class="grid place-items-center h-96 bg-base-100 rounded-lg border border-base-200 border-dashed">
            <div class="text-center">
              <Icon icon="mdi:server-off" class="w-20 h-20 text-base-content/20 mx-auto mb-4" />
              <h3 class="text-lg font-semibold text-base-content/70 mb-2">暂未发现已绑定的主机</h3>
              <p class="text-sm text-base-content/50 mb-6">绑定主机后，您可以在此处查看和管理主机状态</p>
              <button class="btn btn-primary btn-md shadow-sm" onclick="bind_dialog.showModal()" aria-label="立即绑定主机">
                <Icon icon="mdi:link" /> 立即绑定
              </button>
            </div>
          </div>
        </template>
      </Transition>
    </div>
    <Pagination
      v-if="!loading"
      :total="data?.total || 0"
      :current-page="currentPage"
      :page-size="pageSize"
      @change="handlePageChange"
    />

    <!-- 批量操作工具栏 -->
    <Transition name="slide-up">
      <div v-if="selectedIds.length > 0"
        class="fixed bottom-8 left-1/2 -translate-x-1/2 z-40 bg-base-100 border border-base-300 shadow-2xl rounded-2xl px-6 py-3 flex items-center gap-6 animate-in fade-in zoom-in duration-300">
        <div class="flex items-center gap-3 border-r border-base-300 pr-2 lg:pr-6">
          <div
            class="w-8 h-8 rounded-full bg-primary text-primary-content flex items-center justify-center text-xs lg:text-sm font-bold">
            {{ selectedIds.length }}
          </div>
          <span class="text-xs lg:text-sm font-medium text-base-content/70">项已选择</span>
        </div>

        <div class="flex items-center gap-2">
          <button class="btn btn-primary btn-xs lg:btn-sm gap-2" @click="handleV2raySubscribe">
            <Icon icon="mdi:link" class="w-4 h-4" />
            v2ray订阅
          </button>
          <button class="btn btn-primary btn-xs lg:btn-sm gap-2" onclick="batch_command_modal.showModal()">
            <Icon icon="mdi:console-line" class="w-4 h-4" />
            执行命令
          </button>
        </div>

        <div class="border-l border-base-300 pl-4">
          <button class="btn btn-ghost btn-xs lg:btn-sm text-base-content/50 hover:text-base-content" @click="clearSelection">
            取消选择
          </button>
        </div>
      </div>
    </Transition>

    <dialog id="v2ray_modal" ref="v2rayModalRef" class="modal">
      <div class="modal-box">
        <form method="dialog">
          <button class="btn btn-sm btn-circle btn-ghost absolute right-2 top-2">✕</button>
        </form>
        <p class="py-4">{{ v2rayLink }}</p>
        <div class="modal-action">
          <CopyButton :text=v2rayLink />
        </div>
      </div>
    </dialog>

    <!-- 批量执行命令模态框 -->
    <dialog id="batch_command_modal" class="modal">
      <div class="modal-box max-w-2xl">
        <h3 class="font-bold text-lg mb-4 flex items-center gap-2">
          <Icon icon="mdi:console-line" class="text-primary w-6 h-6" />
          批量执行命令 ({{ selectedIds.length }} 台主机)
        </h3>
        <div class="space-y-4">
          <div class="form-control">
            <label class="label">
              <span class="label-text">请输入要在已选主机上执行的命令</span>
            </label>
            <textarea class="textarea textarea-bordered h-32 font-mono text-sm"
              placeholder="例如: uptime, df -h, ls -la..."></textarea>
          </div>

          <div class="alert alert-info shadow-sm py-2">
            <Icon icon="mdi:information-outline" class="w-5 h-5" />
            <span class="text-xs">命令将依次在所有选中的在线主机上执行，结果将实时显示。</span>
          </div>

          <div class="flex flex-wrap gap-2">
            <div v-for="id in selectedIds" :key="id" class="badge badge-outline badge-sm gap-1">
              {{data?.list.find(i => i.id === id)?.hostname || id}}
            </div>
          </div>
        </div>
        <div class="modal-action">
          <form method="dialog">
            <button class="btn btn-ghost">取消</button>
            <button class="btn btn-primary ml-2">开始执行</button>
          </form>
        </div>
      </div>
      <form method="dialog" class="modal-backdrop">
        <button>close</button>
      </form>
    </dialog>
  </div>
</template>

<style scoped>
.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.slide-up-enter-from {
  opacity: 0;
  transform: translate(-50%, 40px) scale(0.95);
}

.slide-up-leave-to {
  opacity: 0;
  transform: translate(-50%, 40px) scale(0.95);
}

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
