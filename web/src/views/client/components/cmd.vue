<template>
  <common-dialog :title="title" :visible="visible" width="60%" @closed="handleClose()">
    <el-input v-model="result" type="textarea" :rows="15" disabled placeholder="执行结果" />
    <el-input v-model="cmd" placeholder="请输入命令">
      <el-button slot="append" icon="el-icon-position" @click="handleCmd()" />
    </el-input>
  </common-dialog>
</template>

<script>
import { cmd } from '@/api/client'
import formMixin from '@/mixins/form-children'
import CommonDialog from '@/components/CommonDialog/index.vue'

export default {
  name: 'Cmd',
  mixins: [formMixin],
  components: {
    CommonDialog
  },
  props: {
    clientId: {
      type: Number,
      default: null
    },
    visible: {
      type: Boolean,
      default: false
    },
    title: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      cmd: '',
      result: ''
    }
  },
  mounted() {

  },
  methods: {
    handleCmd() {
      cmd({ id: this.clientId, command: this.cmd }).then(res => {
        this.result = res.data
      })
    },
    handleClose() {
      this.cmd = ''
      this.result = ''
      this.$emit('hide')
    }
  }
}
</script>

<style lang="scss" scoped>
html,
body,
#app {
    height: 100%;
}
</style>
