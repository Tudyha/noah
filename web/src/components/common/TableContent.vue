<template>
    <div class="rounded-xl shadow-lg">
        <!-- 响应式表格容器，在小屏幕上允许横向滚动 -->
        <div class="overflow-x-auto">
            <table class="table">
                <!-- 表格头部 -->
                <thead>
                    <tr class="uppercase text-sm leading-normal">
                        <th v-for="col in columns" :key="col.key"
                            :class="['py-3', 'px-4', 'text-left', { 'hidden sm:table-cell': col.isResponsiveHidden }]">
                            {{ col.label }}
                        </th>
                    </tr>
                </thead>

                <!-- 表格内容 -->
                <tbody>
                    <tr v-for="(row, rowIndex) in data" :key="rowIndex"
                        class="border-b border-gray-200 hover:bg-indigo-50/50 transition duration-150">
                        <td v-for="col in columns" :key="col.key"
                            :class="['py-3', 'px-4', { 'hidden sm:table-cell': col.isResponsiveHidden }]">
                            <div class="font-medium text-gray-800 whitespace-nowrap">

                                <!-- ⚠️ 针对操作列和状态列的特殊渲染 -->
                                <template v-if="col.key === 'actions'">
                                    <div class="space-x-2 flex">
                                        <button class="btn btn-warning btn-xs text-white"
                                            @click="handleAction('edit', row)">
                                            <!-- 编辑图标 -->
                                            <svg class="w-3 h-3" xmlns="http://www.w3.org/2000/svg" width="24"
                                                height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                                                strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                                <path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z" />
                                            </svg> 编辑
                                        </button>
                                        <button class="btn btn-error btn-xs text-white"
                                            @click="handleAction('delete', row)">
                                            <!-- 删除图标 -->
                                            <svg class="w-3 h-3" xmlns="http://www.w3.org/2000/svg" width="24"
                                                height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                                                strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                                <path d="M3 6h18" />
                                                <path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6" />
                                                <path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
                                            </svg> 删除
                                        </button>
                                    </div>
                                </template>
                                <template v-else-if="col.key === 'status'">
                                    <div
                                        :class="['badge', row.status === '正常' ? 'badge-success' : 'badge-error', 'text-white']">
                                        {{ row.status }}
                                    </div>
                                </template>
                                <template v-else>
                                    <!-- 显示原始值 -->
                                    {{ row[col.key] }}
                                </template>
                            </div>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <!-- 加载/空数据状态 -->
        <div v-if="isLoading" class="flex justify-center items-center py-10">
            <!-- 加载图标 -->
            <svg class="w-8 h-8 animate-spin text-primary" xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round"
                strokeLinejoin="round">
                <path d="M21 12a9 9 0 1 1-6.219-8.56" />
            </svg>
            <span class="ml-2 text-primary font-semibold">数据加载中...</span>
        </div>
        <div v-if="!isLoading && data.length === 0" class="text-center py-10 text-gray-500">
            <p class="text-lg font-medium">暂无数据</p>
            <p class="text-sm mt-1">请尝试调整搜索条件或添加新数据</p>
        </div>
    </div>
</template>

<script setup lang="ts">
// --- TypeScript 接口定义 (为保证组件独立性，重复定义) ---
interface Column {
    key: string;
    label: string;
    isResponsiveHidden?: boolean;
}

interface MockRow {
    id: number;
    name: string;
    role: string;
    status: '正常' | '禁用';
    phone: string;
    entryDate: string;
    [key: string]: any;
}

interface TableContentProps {
    data: MockRow[];
    columns: Column[];
    isLoading?: boolean;
}

const props = withDefaults(defineProps<TableContentProps>(), {
    data: () => [],
    columns: () => [],
    isLoading: false,
});

const emit = defineEmits<{
    (e: 'action', type: 'edit' | 'delete', row: MockRow): void // 用于行操作按钮点击
}>()

const handleAction = (type: 'edit' | 'delete', row: MockRow) => {
    emit('action', type, row);
}
</script>