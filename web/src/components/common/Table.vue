<template>
    <div class="flex flex-col gap-4">
        <div class="navbar shadow-sm p-4 rounded-xl">
            <div class="flex-1">{{ props.title }}</div>
            <div class="flex-none flex flex-row gap-2">
                <slot name="actions" />
            </div>
        </div>

        <!-- 搜索栏 -->
        <TableSearch v-if="props.search && props.search.length > 0" :items="props.search" @search="handleSearch"
            @reset="handleReset" class="card shadow-sm p-4" />

        <!-- 2. 表格内容 -->
        <TableContent :data="tableData" :columns="tableColumns" :is-loading="isLoading" @action="handleRowAction" />

        <!-- 3. 分页器 -->
        <TablePagination :pagination="paginationProps" />
    </div>
</template>

<script setup lang="ts">
import TableSearch from './Search.vue';
import TableContent from './TableContent.vue';
import TablePagination from './Pagination.vue';
import type { TableProps } from '@/types';

const props = defineProps<TableProps>();

// --- TypeScript 接口定义 (主组件作为配置中心) ---


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


// --- 模拟数据和状态 ---

const mockData: MockRow[] = [
    { id: 1001, name: '张三', role: '系统管理员', status: '正常', phone: '138****0001', entryDate: '2023-10-01' },
    { id: 1002, name: '李四', role: '普通用户', status: '禁用', phone: '139****0002', entryDate: '2023-10-05' },
    { id: 1003, name: '王五', role: '财务经理', status: '正常', phone: '137****0003', entryDate: '2023-11-12' },
    { id: 1004, name: '赵六', role: '普通用户', status: '正常', phone: '136****0004', entryDate: '2023-12-20' },
    { id: 1005, name: '钱七', role: '系统管理员', status: '正常', phone: '135****0005', entryDate: '2024-01-01' },
    { id: 1006, name: '孙八', role: '普通用户', status: '正常', phone: '134****0006', entryDate: '2024-02-15' },
];

const tableData = ref<MockRow[]>([]);
const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(45); // 假设总共有45条数据
const isLoading = ref(false);
const currentSearchParams = ref<Record<string, any>>({});


// --- 配置 ---
const tableColumns: Column[] = [
    { key: 'id', label: 'ID', isResponsiveHidden: true },
    { key: 'name', label: '姓名' },
    { key: 'role', label: '角色' },
    { key: 'phone', label: '电话', isResponsiveHidden: true },
    { key: 'status', label: '状态' },
    { key: 'entryDate', label: '入职日期' },
    { key: 'actions', label: '操作' },
];

// --- 核心方法 ---

// 统一数据请求逻辑
const fetchData = (params: Record<string, any>) => {
    console.log('--- 模拟发送数据请求 ---');
    console.log('请求参数:', params);

    isLoading.value = true;

    setTimeout(() => {
        // 模拟后端过滤逻辑
        const filteredData = mockData.filter(item => {
            let match = true;
            if (params.name && !item.name.includes(params.name)) match = false;
            if (params.role && item.role !== params.role) match = false;
            if (params.status && item.status !== params.status) match = false;
            return match;
        });

        // 模拟分页逻辑 (根据过滤后的数据)
        const start = (params.currentPage - 1) * params.pageSize;
        const end = start + params.pageSize;

        tableData.value = filteredData.slice(start, end);
        total.value = 45; // 假设总数不变
        isLoading.value = false;
        console.log('数据更新完成。');
    }, 800);
};


// --- 事件处理 (由子组件触发，主组件处理) ---

// 1. 搜索处理 (来自 TableSearch)
const handleSearch = (form: Record<string, any>) => {
    currentPage.value = 1;
    currentSearchParams.value = form;
    fetchData({ ...form, currentPage: 1, pageSize: pageSize.value });
};

// 2. 重置处理 (来自 TableSearch)
const handleReset = () => {
    currentPage.value = 1;
    currentSearchParams.value = {};
    fetchData({ currentPage: 1, pageSize: pageSize.value });
};

// 3. 分页页码变化 (来自 TablePagination)
const handlePageChange = (page: number) => {
    const totalPages = Math.ceil(total.value / pageSize.value);
    if (page > 0 && page <= totalPages) {
        currentPage.value = page;
        fetchData({ ...currentSearchParams.value, currentPage: page, pageSize: pageSize.value });
    }
};

// 4. 每页大小变化 (来自 TablePagination)
const handlePageSizeChange = (size: number) => {
    pageSize.value = size;
    currentPage.value = 1;
    fetchData({ ...currentSearchParams.value, currentPage: 1, pageSize: size });
};

// 5. 行操作 (来自 TableContent)
const handleRowAction = (type: 'edit' | 'delete', row: MockRow) => {
    console.log(`${type === 'edit' ? '编辑' : '删除'}行:`, row.id);
};


// --- Computed Props 传递给子组件 ---

// 传递给 TablePagination 的 Pagination 对象
const paginationProps = computed(() => ({
    currentPage: currentPage.value,
    pageSize: pageSize.value,
    total: total.value,
    onPageChange: handlePageChange,
    onPageSizeChange: handlePageSizeChange,
}));

// --- 生命周期 ---
onMounted(() => {
    fetchData({ currentPage: currentPage.value, pageSize: pageSize.value });
});
</script>