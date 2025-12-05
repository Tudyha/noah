<template>
<div v-if="pagination && pagination.total > 0" class="flex flex-col sm:flex-row justify-between items-center mt-6 p-4 bg-white rounded-xl shadow-lg">
    <!-- 页面信息 -->
    <div class="text-sm text-gray-600 mb-4 sm:mb-0">
        总共 {{ pagination.total }} 条数据，当前第 {{ pagination.currentPage }} / {{ totalPages }} 页
    </div>

    <!-- 分页按钮 -->
    <div class="flex items-center space-x-4">
         <!-- 每页大小选择 -->
        <select
            class="select select-bordered select-sm w-auto focus:border-indigo-500 transition duration-150"
            :value="pagination.pageSize"
            @change="(e: { target: HTMLSelectElement; }) => pagination.onPageSizeChange(Number((e.target as HTMLSelectElement).value))"
        >
            <option value="10">10 / 页</option>
            <option value="20">20 / 页</option>
            <option value="50">50 / 页</option>
        </select>

        <div class="join shadow-md">
            <button
                class="join-item btn btn-sm btn-outline btn-info transition duration-150"
                @click="pagination.onPageChange(pagination.currentPage - 1)"
                :disabled="pagination.currentPage <= 1"
            >
                <!-- 左箭头 -->
                <svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M15 18l-6-6 6-6"/></svg>
            </button>
            
            <!-- 分页数字按钮 -->
            <button
                v-for="(page, index) in renderPaginationButtons.pages"
                :key="index"
                :class="['join-item', 'btn', 'btn-sm', 'transition', 'duration-150', { 'btn-primary shadow-lg shadow-indigo-500/50': page === pagination.currentPage, 'btn-ghost': page !== pagination.currentPage, 'btn-disabled pointer-events-none': page === '...' }]"
                @click="pagination.onPageChange(page as number)"
                :disabled="page === pagination.currentPage || page === '...'"
            >
                {{ page }}
            </button>

            <button
                class="join-item btn btn-sm btn-outline btn-info transition duration-150"
                @click="pagination.onPageChange(pagination.currentPage + 1)"
                :disabled="pagination.currentPage >= totalPages"
            >
                <!-- 右箭头 -->
                <svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M9 18l6-6-6-6"/></svg>
            </button>
        </div>
    </div>
</div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

// --- TypeScript 接口定义 (为保证组件独立性，重复定义) ---
interface Pagination {
  currentPage: number;
  pageSize: number;
  total: number;
  onPageChange: (page: number) => void;
  onPageSizeChange: (size: number) => void;
}

interface TablePaginationProps {
  pagination: Pagination;
}

const props = defineProps<TablePaginationProps>();

// --- 分页计算 ---

const totalPages = computed(() => {
    if (!props.pagination || props.pagination.total === 0 || props.pagination.pageSize === 0) return 1;
    return Math.ceil(props.pagination.total / props.pagination.pageSize);
});

// 计算分页按钮数组
const renderPaginationButtons = computed(() => {
    const pages: (number | '...')[] = [];
    const { currentPage } = props.pagination;
    const maxPages = totalPages.value;

    if (maxPages === 0) return { pages: [] };

    pages.push(1);

    if (currentPage > 3) pages.push('...');

    for (let i = Math.max(2, currentPage - 1); i <= Math.min(maxPages - 1, currentPage + 1); i++) {
        if (i > 1 && i < maxPages) {
            pages.push(i);
        }
    }

    if (currentPage < maxPages - 2 && maxPages > 1) pages.push('...');

    if (maxPages > 1 && !pages.includes(maxPages)) pages.push(maxPages);
    
    // 清理重复的页码和相邻的省略号
    const finalPages: (number | '...')[] = [];
    pages.forEach(p => {
        if (p === '...' && finalPages[finalPages.length - 1] === '...') return;
        if (typeof p === 'number' && finalPages.includes(p)) return;
        finalPages.push(p);
    });

    return { pages: finalPages };
});
</script>