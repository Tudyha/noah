<script setup lang="ts">
import type { PaginationProps } from '@/types'

const props = defineProps<PaginationProps>()

const emit = defineEmits(['change'])

const pageSizes = [10, 20, 30, 40, 50]
const currentPage = ref(1)
const pageSize = ref(10)
const totalPages = computed(() => Math.ceil(props.total / pageSize.value))

const pageRange = computed(() => {
  const range = []
  const half = Math.floor(totalPages.value / 2)
  let start = Math.max(1, currentPage.value - half)
  let end = Math.min(totalPages.value, start + totalPages.value - 1)
  
  if (end - start < totalPages.value - 1) {
    start = Math.max(1, end - totalPages.value + 1)
  }
  
  for (let i = start; i <= end; i++) {
    range.push(i)
  }
  
  return range
})

const handlePageChange = (page: number) => {
  if (page < 1 || page > totalPages.value || page === currentPage.value) return
  currentPage.value = page
  emit('change', { page, pageSize: pageSize.value })
}
</script>

<template>
  <div class="flex flex-col sm:flex-row items-center justify-between gap-4 p-2">
    <!-- 分页信息 -->
    <div class="text-sm text-base-content/60">
      共 <span class="font-semibold text-base-content">{{ total }}</span> 条数据
      <span class="ml-2">
        第 <span class="font-semibold text-base-content">{{ currentPage }}</span> / 
        <span class="font-semibold text-base-content">{{ totalPages }}</span> 页
      </span>
    </div>

    <!-- 分页器 -->
    <div class="flex items-center gap-2">
      <!-- 每页条数选择 -->
      <select 
        class="select select-bordered select-sm"
        :value="pageSize"
      >
        <option v-for="size in pageSizes" :key="size" :value="size">
          {{ size }} 条/页
        </option>
      </select>

      <!-- 分页按钮 -->
      <div class="join">
        <!-- 首页 -->
        <button 
          class="join-item btn btn-sm"
          :disabled="currentPage === 1"
          @click="handlePageChange(1)"
        >
          «
        </button>

        <!-- 上一页 -->
        <button 
          class="join-item btn btn-sm"
          :disabled="currentPage === 1"
          @click="handlePageChange(currentPage - 1)"
        >
          ‹
        </button>

        <!-- 页码 -->
        <button
          v-for="page in pageRange"
          :key="page"
          class="join-item btn btn-sm"
          :class="{ 'btn-active': page === currentPage }"
          @click="handlePageChange(page)"
        >
          {{ page }}
        </button>

        <!-- 下一页 -->
        <button 
          class="join-item btn btn-sm"
          :disabled="currentPage === totalPages"
          @click="handlePageChange(currentPage + 1)"
        >
          ›
        </button>

        <!-- 末页 -->
        <button 
          class="join-item btn btn-sm"
          :disabled="currentPage === totalPages"
          @click="handlePageChange(totalPages)"
        >
          »
        </button>
      </div>
    </div>
  </div>
</template>