<script setup lang="ts">
import type { SearchProps } from '@/types';

const props = defineProps<SearchProps>()
const searchForm = ref<Record<string, any>>({});

// 初始化默认值
onMounted(() => {
  props.items?.forEach(item => {
    if (item.default !== undefined) {
      searchForm.value[item.key] = item.default;
    }
  });
});

const emit = defineEmits<{
  (e: 'search', form: Record<string, any>): void
}>()

const handleSearch = () => {
  emit('search', searchForm.value);
};

const handleReset = () => {
  // 重置为默认值或清空
  props.items?.forEach(item => {
    searchForm.value[item.key] = item.default !== undefined ? item.default : undefined;
  });
  emit('search', searchForm.value);
};
</script>

<template>
  <div class="bg-base-100 p-4 rounded-lg mb-4 shadow-sm border border-base-200">
    <div class="flex flex-wrap items-end gap-4">
      <div v-if="items && items.length > 0" v-for="item in props.items" :key="item.key" 
           class="form-control" :class="item.width ? item.width : 'w-full sm:w-64'">
        <label class="label py-1 px-0">
          <span class="label-text font-medium">{{ item.label }}</span>
        </label>
        
        <div class="w-full">
          <!-- 输入框 -->
          <input v-if="item.type === 'input'" type="text" :placeholder="item.placeholder || `请输入${item.label}`"
            class="input input-bordered input-sm w-full" v-model="searchForm[item.key]" />

          <!-- 下拉选择 -->
          <select v-else-if="item.type === 'select'" class="select select-bordered select-sm w-full"
            v-model="searchForm[item.key]">
            <option v-for="option in item.options" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>
        </div>
      </div>
      
      <div class="flex items-center gap-2 pb-0.5">
        <button class="btn btn-sm btn-primary shadow-sm" @click="handleSearch">
          <Icon icon="mdi:search" class="w-4 h-4" /> 搜索
        </button>
        <button class="btn btn-sm btn-ghost border border-base-300 shadow-sm" @click="handleReset">
          <Icon icon="mdi:refresh" class="w-4 h-4" /> 重置
        </button>
      </div>
    </div>
  </div>
</template>
