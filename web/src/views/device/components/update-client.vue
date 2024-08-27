<template>
  <el-dialog :visible="visible" width="30%" @close="handleClose()">
    <el-form
      ref="form"
      :model="form"
      :rules="rules"
    >
      <el-form-item label="服务器地址：" prop="serverAddr">
        <el-input v-model="form.serverAddr" placeholder="" style="width: 300px" />
      </el-form-item>
      <el-form-item label="服务器端口：" prop="port">
        <el-input v-model="form.port" placeholder="" style="width: 100px" />
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button type="primary" @click="handleSubmit()">确 定</el-button>
    </div>
  </el-dialog>
</template>

<script>
import {update} from '@/api/client'
import formMixin from '@/mixins/form-children'

export default {
  name: 'ClientUpdate',
  mixins: [formMixin],
  props: {
    clientId: {
      type: Number,
      default: null
    },
    visible: {
      type: Boolean,
      default: false
    },
    osType: {
      type: Number,
      default: null
    }
  },
  data() {
    return {
      form: { // 表单数据
        serverAddr: null,
        port: null,
        filename: null,
        osType: null
      },
      rules: {
        serverAddr: [{ required: true, message: '必填', trigger: 'change' }],
        port: [{ required: true, message: '必填', trigger: 'change' }]
      }
    }
  },
  mounted() {

  },
  methods: {
    handleSubmit() {
      this.form.osType = this.osType
      this.$refs.form.validate(async(valid) => {
        if (!valid) {
          return
        }
        this.form.id = this.clientId
        update(this.form).then(res => {
          if (res.code === 0) {
            this.$message({
              message: '更新成功',
              type: 'success'
            })
            this.handleClose()
          } else {
            this.$message({
              message: res.msg,
              type: 'error'
            })
          }
        })
      })

    },
    handleClose() {
      this.form = {
        serverAddr: null,
        port: null,
        filename: null,
        osType: null
      }
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
