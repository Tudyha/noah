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
  <div class="overflow-x-auto">
    <table class="table table-zebra">
      <thead class="bg-base-200/50">
        <tr>
          <th v-for="column in columns" :key="column.key" :style="{ width: column.width }">
            <div class="flex items-center gap-2 font-semibold">
              <span>{{ column.label }}</span>
            </div>
          </th>
        </tr>
      </thead>
      <tbody>
        <!-- 加载状态 -->
        <tr v-if="props.isLoading">
          <td :colspan="columns.length" class="text-center py-12">
            <span class="loading loading-spinner loading-lg text-primary"></span>
            <p class="mt-2 text-sm text-base-content/60">加载中...</p>
          </td>
        </tr>

        <!-- 空数据 -->
        <tr v-else-if="!data || data.length === 0">
          <td :colspan="columns.length" class="text-center py-12">
            <div class="flex flex-col items-center justify-center opacity-60">
              <Icon icon="mdi:package-variant-closed" class="w-12 h-12 mb-2" />
              <p class="text-sm">暂无数据</p>
            </div>
          </td>
        </tr>

        <!-- 数据行 -->
        <tr v-else v-for="(row, index) in data" :key="row.id || index" class="hover">
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