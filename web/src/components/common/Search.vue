<script setup lang="ts">
import type { SearchProps } from '@/types';

const props = defineProps<SearchProps>()
const searchForm = ref<Record<string, any>>({});

const emit = defineEmits<{
  (e: 'search', form: Record<string, any>): void
}>()

const handleSearch = () => {
  emit('search', searchForm.value);
};

</script>

<template>
  <div class="border-b border-base-content/5 p-2">
    <div class="grid grid-cols-4 auto-rows-min gap-x-4 gap-y-2" style="grid-auto-flow: dense;">
      <div v-if="items && items.length > 0" v-for="item in props.items" :key="item.key" class="flex items-center gap-1">
        <div class="w-20">
          <span>{{ item.label }}</span>
        </div>

        <div :class="item.width ? item.width : 'w-full'">
          <!-- 输入框 -->
          <input v-if="item.type === 'input'" type="text" :placeholder="item.placeholder || `请输入${item.label}`"
            class="input input-bordered input-sm" v-model="searchForm[item.key]" />

          <!-- 下拉选择 -->
          <select v-else-if="item.type === 'select'" class="select select-bordered select-sm"
            v-model="searchForm[item.key]">
            <option v-for="option in item.options" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>
        </div>
      </div>
      <div class="flex items-center gap-1">
        <button class="btn btn-sm btn-soft btn-primary btn-square" @click="handleSearch">
          <Icon icon="mdi:search" class="w-6 h-6" />
        </button>
        <button class="btn btn-sm btn-soft btn-warning btn-square">
          <Icon icon="mdi:refresh" class="w-6 h-6" />
        </button>
      </div>
    </div>
  </div>
</template>
