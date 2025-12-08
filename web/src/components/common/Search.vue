<script setup lang="ts">
import type { SearchProps } from '@/types';

const props = defineProps<SearchProps>()

// 计算初始搜索表单状态，用于重置
const initialSearchState = computed(() => {
    return props.items.reduce((acc, item) => {
        acc[item.key] = '';
        return acc;
    }, {} as Record<string, any>);
});

const searchForm = ref<Record<string, any>>(initialSearchState.value);

const emit = defineEmits<{
    (e: 'search', form: Record<string, any>): void
    (e: 'reset'): void
}>()

const handleSearch = () => {
    emit('search', searchForm.value);
};

const handleReset = () => {
    searchForm.value = initialSearchState.value;
    emit('reset');
};

</script>

<template>
  <div class="card mb-4 rounded-box border border-base-content/5 bg-base-100">
    <div class="card-body">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-2">
        <div v-for="item in props.items" :key="item.key" class="form-control">
          <label class="label">
            <span class="label-text">{{ item.label }}</span>
          </label>

          <!-- 输入框 -->
          <input
            v-if="item.type === 'input'"
            type="text"
            :placeholder="item.placeholder || `请输入${item.label}`"
            class="input input-bordered input-sm"
            v-model="searchForm[item.key]"
          />

          <!-- 下拉选择 -->
          <select
            v-else-if="item.type === 'select'"
            class="select select-bordered select-sm"
            v-model="searchForm[item.key]"
          >
            <option v-for="option in item.options" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>

          <!-- 日期选择 -->
          <input
            v-else-if="item.type === 'date'"
            type="date"
            class="input input-bordered select-sm"
            v-model="searchForm[item.key]"
          />

          <!-- 日期范围 -->
          <!-- <div v-else-if="item.type === 'date-range'" class="flex gap-2">
            <input
              type="date"
              class="input input-bordered w-full"
              v-model="searchForm[item.key]?.[0]"
            />
            <span class="flex items-center">至</span>
            <input
              type="date"
              class="input input-bordered w-full"
              v-model="searchForm[item.key]?.[1]"
            />
          </div> -->
        </div>
      </div>

      <div class="card-actions justify-end mt-4">
        <button class="btn btn-ghost btn-sm" @click="handleReset">
          重置
        </button>
        <button class="btn btn-primary btn-sm" @click="handleSearch">
          搜索
        </button>
      </div>
    </div>
  </div>
</template>