<template>
  <CommonTable
    title="主机列表"
    :columns="columns"
    :data="hosts"
    :search="searchFields"
    row-key="id"
    show-pagination
    :page-size="3"
  >
    <template #actions>
      <button class="btn btn-primary btn-sm">
        <span class="text-xl">+</span> 新建主机
      </button>
      <button class="btn btn-primary btn-sm">
        <span class="text-xl">+</span> 新建主机
      </button>
    </template>

    <template #header>
      <th>
        <label>
          <input type="checkbox" class="checkbox" />
        </label>
      </th>
      <th>主机名称</th>
      <th>IP 地址</th>
      <th>区域</th>
      <th>状态</th>
      <th>创建时间</th>
      <th>操作</th>
    </template>

    <template #body="{ data }">
      <tr v-for="host in data" :key="host.id" class="hover">
        <th>
          <label>
            <input type="checkbox" class="checkbox" />
          </label>
        </th>
        <td>
          <div class="flex items-center gap-3">
            <div class="avatar placeholder">
              <div class="mask mask-squircle w-10 h-10 bg-neutral text-neutral-content">
                <span>H</span>
              </div>
            </div>
            <div>
              <div class="font-bold">{{ host.name }}</div>
              <div class="text-sm opacity-50">ID: {{ host.id }}</div>
            </div>
          </div>
        </td>
        <td class="font-mono text-xs">{{ host.ip }}</td>
        <td>{{ host.region }}</td>
        <td>
          <div class="badge badge-outline gap-2" :class="getStatusColor(host.status)">
            {{ host.status }}
          </div>
        </td>
        <td>{{ host.createTime }}</td>
        <th>
          <button class="btn btn-ghost btn-xs">详情</button>
          <button class="btn btn-ghost btn-xs text-error">重启</button>
        </th>
      </tr>
    </template>
  </CommonTable>
</template>

<script setup lang="ts">
import CommonTable from '@/components/common/Table.vue'
import type { TableColumn, SearchItem } from '@/types/common'

const columns: TableColumn[] = [
  { key: 'selection', label: '' },
  { key: 'name', label: '主机名称' },
  { key: 'ip', label: 'IP 地址' },
  { key: 'region', label: '区域' },
  { key: 'status', label: '状态' },
  { key: 'createTime', label: '创建时间' },
  { key: 'actions', label: '操作' }
]

const searchFields: SearchItem[] = [
  { key: 'name', label: '主机名称', type: 'input', placeholder: '请输入主机名称' },
  { key: 'ip', label: 'IP地址', type: 'input', placeholder: '请输入IP地址' },
  { key: 'region', label: '区域', type: 'select', options: [
    { label: '全部', value: '' },
    { label: 'CN-East', value: 'CN-East' },
    { label: 'CN-South', value: 'CN-South' },
    { label: 'CN-North', value: 'CN-North' },
    { label: 'US-West', value: 'US-West' }
  ]},
  { key: 'status', label: '状态', type: 'select', options: [
    { label: '全部', value: '' },
    { label: '运行中', value: 'Running' },
    { label: '已停止', value: 'Stopped' },
    { label: '维护中', value: 'Maintenance' }
  ]},
  { key: 'createTime', label: '创建时间', type: 'date-range' }
]

const hosts = ref([
  { id: 1, name: 'Web-Server-01', ip: '192.168.1.101', status: 'Running', region: 'CN-East', createTime: '2023-01-15' },
  { id: 2, name: 'DB-Master', ip: '192.168.1.102', status: 'Running', region: 'CN-East', createTime: '2023-02-20' },
  { id: 3, name: 'Cache-Redis', ip: '192.168.1.105', status: 'Stopped', region: 'CN-South', createTime: '2023-03-10' },
  { id: 4, name: 'Worker-Node-A', ip: '192.168.1.120', status: 'Running', region: 'US-West', createTime: '2023-04-05' },
  { id: 5, name: 'Backup-Svr', ip: '192.168.1.200', status: 'Maintenance', region: 'CN-North', createTime: '2023-05-12' },
])

const getStatusColor = (status: string) => {
  switch (status) {
    case 'Running': return 'badge-success'
    case 'Stopped': return 'badge-error'
    case 'Maintenance': return 'badge-warning'
    default: return 'badge-ghost'
  }
}
</script>