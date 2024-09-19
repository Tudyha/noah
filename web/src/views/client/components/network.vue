<template>
      <el-table :data="list" height="70vh">
        <el-table-column prop="pid" label="进程id" width="70px"></el-table-column>
        <el-table-column prop="localaddr" label="Local Address" :formatter="(row, column, cellValue, index) => cellValue.ip + ':' + cellValue.port"></el-table-column>
        <el-table-column prop="remoteaddr" label="Remote Address" :formatter="(row, column, cellValue, index) => cellValue.ip + ':' + cellValue.port"></el-table-column>
        <el-table-column prop="status" label="状态"></el-table-column>
      </el-table>
</template>
<script>
import { fetchNetworkList } from '@/api/client'

export default {
  props: {
    id: {
      type: Number,
      required: true
    }
  },
  data() {
    return {
      list: [],
    };
  },
  methods: {
    fetchList() {
      fetchNetworkList(this.id)
        .then(response => {
          this.list = response.data;
        })
        .catch(error => {
          console.error('Error fetching process list:', error);
        });
    },
    refresh() {
      this.fetchList();
    },
  },
  components: {
  },
  created() {
    this.fetchList();
  },
}
</script>
<style scoped>
</style>
