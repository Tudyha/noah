<script setup lang="ts">
import type { SearchItem, TableColumn } from '@/types';

const tableData = ref([])
const columns: TableColumn[] = [
  { key: 'id', label: 'ID', width: '80px', sortable: true },
  { key: 'name', label: '用户名', sortable: true },
  { key: 'email', label: '邮箱', sortable: true },
  {
    key: 'role', label: '角色', render: (column, row, _) => {
      const value = row[column.key];
      return h('span', {
        class: ['badge', value === 'admin' ? 'badge-primary' : 'badge-ghost']
      }, value === 'admin' ? '管理员' : '普通用户')
    }
  },
  {
    key: 'status', label: '状态', render: (column, row, _) => {
      const value = row[column.key];
      return h('span', {
        class: ['badge', value === 'active' ? 'badge-success' : 'badge-error']
      }, value === 'active' ? '正常' : '禁用')
    }
  },
  { key: 'createTime', label: '创建时间', sortable: true },
  {
    key: 'action', label: '操作', width: '200px', render: () => {
      return h('div', { class: 'flex gap-2' }, [
        h('button', {
          class: 'btn btn-sm btn-primary',
          onClick: (e) => {
            e.stopPropagation()
          }
        }, '编辑'),
        h('button', {
          class: 'btn btn-sm btn-error',
          onClick: (e) => {
            e.stopPropagation()
          }
        }, '删除')
      ])
    }
  }
]
const searchFields: SearchItem[] = [{
  key: 'name',
  label: '用户名',
  type: 'input',
  placeholder: '请输入用户名'
},
{
  key: 'status',
  label: '状态',
  type: 'select',
  options: [
    { label: '全部', value: '' },
    { label: '正常', value: 'active' },
    { label: '禁用', value: 'inactive' }
  ]
},
{
  key: 'role',
  label: '角色',
  type: 'select',
  options: [
    { label: '全部', value: '' },
    { label: '管理员', value: 'admin' },
    { label: '普通用户', value: 'user' }
  ]
},
{
  key: 'createDate',
  label: '创建日期',
  type: 'date'
},
]

const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const mockUsers = [
  { id: 1, name: '张三', email: 'zhangsan@example.com', role: 'admin', status: 'active', createTime: '2024-01-15 10:30:00' },
  { id: 2, name: '李四', email: 'lisi@example.com', role: 'user', status: 'active', createTime: '2024-02-20 14:20:00' },
  { id: 3, name: '王五', email: 'wangwu@example.com', role: 'user', status: 'inactive', createTime: '2024-03-10 09:15:00' },
  { id: 4, name: '赵六', email: 'zhaoliu@example.com', role: 'admin', status: 'active', createTime: '2024-03-25 16:45:00' },
  { id: 5, name: '孙七', email: 'sunqi@example.com', role: 'user', status: 'active', createTime: '2024-04-01 11:00:00' },
]

const fetchData = async () => {
  loading.value = true

  try {
    // 模拟 API 请求
    await new Promise(resolve => setTimeout(resolve, 500))

    // 模拟数据过滤和分页
    let filteredData = [...mockUsers]

    total.value = filteredData.length

    // 分页
    const start = (currentPage.value - 1) * pageSize.value
    const end = start + pageSize.value
    tableData.value = filteredData.slice(start, end)

  } catch (error) {
    console.error('获取数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 初始化
onMounted(() => {
  fetchData()
})

const handleSearch = (params: any) => {
  console.log(params)
}
</script>

<template>
  <Table title="主机列表" :data="tableData" :columns="columns" :search-items="searchFields" :total="total"
    :is-loading="loading" :current-page="currentPage" :page-size="pageSize" @search="handleSearch" />
</template>