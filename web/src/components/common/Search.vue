<template>
    <div class="grid grid-cols-2 gap-4 md:grid-cols-4 lg:grid-cols-5 mb-4">
        <div v-for="item in props.items" :key="item.key" class="flex flex-row gap-2">
            <label class="label">
                <span class="label-text">{{ item.label }}</span>
            </label>

            <!-- 渲染搜索项 -->
            <template v-if="item.type === 'input'">
                <input type="text" :placeholder="item.placeholder || '请输入'" class="input input-sm"
                    v-model="searchForm[item.key]" />
            </template>

            <template v-else-if="item.type === 'select'">
                <select class="select select-sm" v-model="searchForm[item.key]">
                    <option value="">{{ item.placeholder || "请选择" }}</option>
                    <option v-for="option in (item.options || [])" :key="option.value" :value="option.value">
                        {{ option.label }}
                    </option>
                </select>
            </template>

            <template v-else-if="item.type === 'date-range'">
            </template>
        </div>

        <!-- 搜索/重置按钮 -->
        <div class="flex items-end space-x-2 mt-2">
            <button class="btn btn-primary btn-sm shadow-md" @click="handleSearchClick">
                搜索
            </button>
            <button class="btn btn-ghost btn-sm shadow-md" @click="handleResetClick">
                重置
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import type { SearchProps } from '@/types';

const props = defineProps<SearchProps>();

const emit = defineEmits<{
    (e: 'search', form: Record<string, any>): void
    (e: 'reset'): void
}>()

// 计算初始搜索表单状态，用于重置
const initialSearchState = computed(() => {
    return props.items.reduce((acc, item) => {
        acc[item.key] = '';
        return acc;
    }, {} as Record<string, any>);
});

const searchForm = ref<Record<string, any>>(initialSearchState.value);

// 搜索配置变化时，重置表单
watch(initialSearchState, (newInitial) => {
    searchForm.value = newInitial;
}, { immediate: true });

const handleSearchClick = () => {
    emit('search', searchForm.value);
};

const handleResetClick = () => {
    searchForm.value = initialSearchState.value;
    emit('reset');
};
</script>