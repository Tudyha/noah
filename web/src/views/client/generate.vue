<template>
  <el-form
    ref="form"
    :model="form"
    :rules="rules"
    label-position="right"
    label-width="120px"
    style="margin-left:50px;overflow: scroll;"
  >
    <el-form-item label="服务器地址：" prop="serverAddr">
      <el-input v-model="form.serverAddr" placeholder="" style="width: 300px" />
    </el-form-item>
    <el-form-item label="服务器端口：" prop="port">
      <el-input v-model="form.port" placeholder="" style="width: 100px" />
    </el-form-item>
    <el-form-item label="文件名称：" prop="filename">
      <el-input v-model="form.filename" placeholder="" style="width: 300px" />
    </el-form-item>
    <el-form-item label="系统类型：" prop="osType">
      <el-select v-model="form.osType" placeholder="请选择">
        <el-option
          v-for="item in map.osTypeOptions"
          :key="item.value"
          :label="item.label"
          :value="+item.value"
        />
      </el-select>
    </el-form-item>
<!--    <el-form-item>-->
<!--      点击下载:-->
<!--      <a :href="downloadUrl" download>-->
<!--        {{ downloadUrl }}-->
<!--      </a>-->
<!--    </el-form-item>-->
    <el-form-item>
      <div class="margin-top-20">
        <el-button type="primary" :loading="loading" @click="handleGenerate()">立即生成</el-button>
      </div>
    </el-form-item>

  </el-form>
</template>
<script>
import * as map from '@/map/client'
import { generate } from '@/api/client'
import { getToken } from '@/utils/auth'

export default {
  name: 'Client',
  components: { },
  mixins: [],
  data() {
    return {
      map,
      form: { // 表单数据
        serverAddr: null,
        port: null,
        filename: null,
        osType: null
      },
      rules: {
        serverAddr: [{ required: true, message: '必填', trigger: 'change' }],
        port: [{ required: true, message: '必填', trigger: 'change' }],
        osType: [{ required: true, message: '必填', trigger: 'change' }]
      },
      downloadUrl: '',
      loading: false
    }
  },
  watch: {},
  created() {
  },
  methods: {
    handleGenerate() {
      this.$refs.form.validate(async(valid) => {
        if (!valid) {
          return
        }
        this.loading = true
        await generate(this.form).then((res) => {
          this.loading = false
          if (res.code === 0) {
            window.open(process.env.VUE_APP_BASE_API + '/file/download/' + res.data + '?token=' + getToken(), '_self')
            // this.downloadUrl = process.env.VUE_APP_BASE_API + '/file/download/' + res.data + '?token=' + getToken()
          } else {
            this.$message.error(res.msg)
          }
        }).catch(() => {
          this.loading = false
        })
      })
    }
  }
}
</script>
<style lang="scss" scoped>
  .avatar-uploader .el-upload {
    border: 1px dashed #d9d9d9;
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
  }
  .avatar-uploader .el-upload:hover {
    border-color: #409EFF;
  }
  .avatar-uploader-icon {
    font-size: 28px;
    color: #8c939d;
    width: 178px;
    height: 178px;
    line-height: 178px;
    text-align: center;
  }
  .avatar {
    width: 178px;
    height: 178px;
    display: block;
  }
  .c{
    display: flex;
    align-items: center;
  }
  .minWealthLevel-input{
    ::v-deep .el-form-item__content{
      margin-left: 0px !important;
    }
  }
</style>
