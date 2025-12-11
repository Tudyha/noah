<script setup lang="ts">
import Search from './Search.vue';
import TableContent from './TableContent.vue';
import Pagination from './Pagination.vue';
import type { TableProps } from '@/types';

const props = defineProps<TableProps>();

const emit = defineEmits<{
    (e: 'search', form: Record<string, any>): void
    (e: 'reset'): void
}>()

const handleSearch = (params: Record<string, any>) => {
  emit('search', params)
}

const handleReset = () => {
  emit('reset')
}

</script>

<template>
  <div>
    <div class="navbar mb-4 rounded-box border border-base-content/5 bg-base-100 p-6">
      <div class="flex-1">
        <span>{{ props.title }}</span>
      </div>
      <div class="flex-none">
        <button class="btn btn-square btn-ghost">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
            class="inline-block h-5 w-5 stroke-current">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M5 12h.01M12 12h.01M19 12h.01M6 12a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0z">
            </path>
          </svg>
        </button>
      </div>
    </div>

    <!-- 搜索栏 -->
    <Search v-if="searchItems" :items="searchItems" @search="handleSearch" @reset="handleReset" />

    <!-- 表格内容 -->
    <TableContent :columns="columns" :data="data" :is-loading="isLoading" />

    <!-- 分页器 -->
    <Pagination v-if="total && total > 0" :current-page="currentPage" :page-size="pageSize" :total="total" />
  </div>
</template>