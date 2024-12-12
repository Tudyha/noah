<template>
  <el-form ref="form" :model="form" :rules="rules" label-position="right" label-width="120px"
    style="margin-left:50px;overflow: scroll;">

    <el-form-item label="系统类型：" prop="goos">
      <el-select v-model="form.goos" placeholder="请选择">
        <el-option v-for="item in map.osTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
    </el-form-item>
    <el-form-item label="系统架构：" prop="goarch">
      <el-select v-model="form.goarch" placeholder="请选择">
        <el-option v-for="item in map.goarchOptions" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
    </el-form-item>
    <el-form-item label="数据加密方式：" prop="compress">
      <el-select v-model="form.compress" placeholder="请选择">
        <el-option v-for="item in map.compressOptions" :key="item.value" :label="item.label" :value="+item.value" />
      </el-select>
    </el-form-item>
    <el-form-item label="文件名称：" prop="filename">
      <el-input v-model="form.filename" placeholder="" style="width: 300px" />
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
  name: 'ClientBuild',
  components: {},
  mixins: [],
  data() {
    return {
      map,
      form: { // 表单数据
        goos: null,
        goarch: null,
        compress: 0,
        filename: "",
      },
      rules: {
        goos: [{ required: true, message: '必填', trigger: 'change' }],
        goarch: [{ required: true, message: '必填', trigger: 'change' }],
        compress: [{ required: true, message: '必填', trigger: 'change' }],

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
      this.$refs.form.validate(async (valid) => {
        if (!valid) {
          return
        }
        this.loading = true
        this.downloadUrl = `/api/client/build?goos=${this.form.goos}&goarch=${this.form.goarch}&compress=${this.form.compress}&filename=${this.form.filename}&token=${getToken()}`
        // 生成a标签，并点击
        const a = document.createElement('a')
        a.href = this.downloadUrl
        a.click()
        this.loading = false
      })
    },
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

.c {
  display: flex;
  align-items: center;
}

.minWealthLevel-input {
  ::v-deep .el-form-item__content {
    margin-left: 0px !important;
  }
}
</style>
