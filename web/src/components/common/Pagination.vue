<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps({
  currentPage: {
    type: Number,
    default: 1
  },
  pageSize: {
    type: Number,
    default: 10
  },
  total: {
    type: Number,
    required: true
  },
  pageSizes: {
    type: Array,
    default: () => [10, 20, 50, 100]
  },
  maxPages: {
    type: Number,
    default: 7 // 最多显示的页码按钮数量
  }
})

const emit = defineEmits(['update:currentPage', 'update:pageSize', 'change'])

const totalPages = computed(() => Math.ceil(props.total / props.pageSize))

const pageRange = computed(() => {
  const range = []
  const half = Math.floor(props.maxPages / 2)
  let start = Math.max(1, props.currentPage - half)
  let end = Math.min(totalPages.value, start + props.maxPages - 1)
  
  if (end - start < props.maxPages - 1) {
    start = Math.max(1, end - props.maxPages + 1)
  }
  
  for (let i = start; i <= end; i++) {
    range.push(i)
  }
  
  return range
})

const handlePageChange = (page) => {
  if (page < 1 || page > totalPages.value || page === props.currentPage) return
  emit('update:currentPage', page)
  emit('change', { page, pageSize: props.pageSize })
}

const handlePageSizeChange = (size) => {
  emit('update:pageSize', size)
  emit('update:currentPage', 1)
  emit('change', { page: 1, pageSize: size })
}
</script>

<template>
  <div class="flex flex-col sm:flex-row items-center justify-between gap-4 mt-4">
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
        @change="handlePageSizeChange(Number($event.target.value))"
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