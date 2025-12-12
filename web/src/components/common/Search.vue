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
    <div class="">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-2">
        <div v-if="items && items.length > 0" v-for="item in props.items" :key="item.key" class="form-control">
          <label class="label">
            <span class="label-text">{{ item.label }}</span>
          </label>

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

          <!-- 日期选择 -->
          <input v-else-if="item.type === 'date'" type="date" class="input input-bordered select-sm"
            v-model="searchForm[item.key]" />

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
        <div class="flex gap-1">
          <button class="btn btn-ghost btn-sm btn-outline">
            重置
          </button>
          <button class="btn btn-primary btn-sm btn-outline" @click="handleSearch">
            搜索
          </button>
        </div>
      </div>
    </div>
  </div>
</template>