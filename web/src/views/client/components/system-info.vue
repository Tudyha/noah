<template>
  <div class="space-y-4">
    <div class="space-x-2">
      <input type="datetime-local" class="input w-40" v-model="filter.start"/>
      -
      <input type="datetime-local" class="input w-40" v-model="filter.end"/>
      <button class="btn btn-primary" @click="refresh">查询</button>
      <button class="btn btn-primary" @click="refresh">刷新</button>
    </div>
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <div>
        <Line v-if="!loading" :data="cpuData" :options="options" />
      </div>
      <div>
        <Line v-if="!loading" :data="memoryData" :options="options" />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import Line from '@/components/chart/line.vue'
import { getClientSystemInfo } from '@/api/client'
import type { ClientSystemInfoResponse } from '@/types'
import { formatBytesToGB } from '@/utils'

const props = defineProps<{
  id: string
}>()

const filter = ref({
  start: '2025-12-19T11:00',
  end: '2025-12-19T11:55'
})

const data = ref<ClientSystemInfoResponse[]>([])
const loading = ref(false)
const labels = computed(() => {
  return data.value.map(item => new Date(item.created_at).toLocaleTimeString())
})
const cpuData = computed(() => {
  return {
    labels: labels.value,
    datasets: [
      {
        label: 'CPU使用率(%)',
        data: data.value.map(item => item.cpu_percent),
        fill: false,
        tension: 0.1,
        borderColor: 'rgb(75, 192, 192)',
      }
    ]
  }
})

const memoryData = computed(() => {
  return {
    labels: labels.value,
    datasets: [
      {
        label: '内存使用率(%)',
        data: data.value.map(item => item.mem_used_percent),
        fill: false,
        tension: 0.1,
        borderColor: 'rgb(75, 192, 192)',
      },
      {
        label: '内存使用量(GB)',
        data: data.value.map(item => formatBytesToGB(item.mem_used)),
        fill: false,
        tension: 0.1,
        borderColor: 'rgb(255, 99, 132)',
      }
    ]
  }
})

const options = {
  responsive: true,
}

const refresh = async () => {
  loading.value = true
  const res = await getClientSystemInfo(Number(props.id))
  data.value = res
  loading.value = false
}

onMounted(async () => {
  refresh()
})

defineExpose({
  refresh
})
</script>
