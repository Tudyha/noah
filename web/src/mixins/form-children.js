export default {
  data() {
    return {
    }
  },
  methods: {
    // 重置
    resetForm(formName) {
      this.$refs[formName].resetFields()
    },
    // 交互事件
    handleClose(name) {
      this.$emit('hide', { name: name })
    },
    handleCloseRefresh(name) {
      this.$emit('hide', { name: name, isFresh: true })
    },
    submit(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.dataApi()
        } else {
          return false
        }
      })
    },
    // 参数属性
    beforeGifUpload(file) {
      const isImg = file.type === 'image/gif'
      if (!isImg) {
        this.$message.error('只能上传gif或svga格式')
        return false
      }
      return true
    },
    handleGifSuccess(res, file) {
      if (!this.dataOj.fileUrl) {
        this.dataOj.fileUrl = ''
      }
      var url = res.data
      this.dataOj.fileUrl = url
      this.dataOj.fileName = file.name
    },
    handleSvgaSuccess(res, file) {
      if (!this.dataOj.fileUrl) {
        this.dataOj.fileUrl = ''
      }
      var url = res.data
      this.dataOj.fileUrl = url
      this.dataOj.fileName = file.name
    },
    handleImageSuccess(res, file, name) {
      if (!this.dataOj.imagesUrl) {
        this.dataOj.imagesUrl = ''
      }
      var url = res.data
      this.dataOj.imagesUrl = url
    },
    beforeSvgaUpload(file) {
      const isImg = file.name.indexOf('svga')
      if (isImg === -1) {
        this.$message.error('只能上传svga格式')
        return false
      }
      return true
    },
    filterObj(arr, obj) {
      if (!Array.isArray(arr)) return
      const newObj = {}
      for (const i in obj) {
        if (arr.includes(i)) {
          newObj[i] = obj[i]
        }
      }
      return newObj
    }
  }
}

