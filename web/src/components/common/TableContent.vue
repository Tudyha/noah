<script setup lang="ts">
import type { TableProps, TableColumn } from '@/types';

const props = defineProps<TableProps>();
const renderCell = (column: TableColumn, row: any, index: number) => {
  if (column.render) {
    return column.render(column, row, index)
  }
};
</script>

<template>
  <div class="overflow-x-auto rounded-box border border-base-content/5 bg-base-100">
    <table class="table">
      <thead>
        <tr>
          <th v-for="column in columns" :key="column.key" :style="{ width: column.width }">
            <div class="flex items-center gap-2">
              <span>{{ column.label }}</span>
              <span v-if="column.sortable" class="flex flex-col">
                <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                  <path
                    d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z" />
                </svg>
                <svg class="w-3 h-3 -mt-1" fill="currentColor" viewBox="0 0 20 20">
                  <path
                    d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z" />
                </svg>
              </span>
            </div>
          </th>
        </tr>
      </thead>
      <tbody>
        <!-- 加载状态 -->
        <tr v-if="props.isLoading">
          <td :colspan="columns.length" class="text-center py-8">
            <span class="loading loading-spinner loading-lg text-primary"></span>
            <p class="mt-2 text-sm text-base-content/60">加载中...</p>
          </td>
        </tr>

        <!-- 空数据 -->
        <tr v-else-if="!data || data.length === 0">
          <td :colspan="columns.length" class="text-center py-8">
            <svg class="w-16 h-16 mx-auto text-base-content/20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
            </svg>
            <p class="mt-2 text-sm text-base-content/60">暂无数据</p>
          </td>
        </tr>

        <!-- 数据行 -->
        <tr v-else v-for="(row, index) in data" :key="row.id || index">
          <td v-for="column in columns" :key="column.key">
            <component v-if="column.render && typeof column.render === 'function'"
              :is="renderCell(column, row, index)" />
              <template v-else>{{ row[column.key] }}</template>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>