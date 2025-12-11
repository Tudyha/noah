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