import { getToken } from '@/utils/auth'

export default {
  data() {
    return {
      uploadUrl: `${process.env.VUE_APP_BASE_API}/upload`,
      headers: getGatewayHeaders(),
      imgKey: '',
      baseUrl: `${process.env.VUE_APP_FILE_BASE_URL}/`,
      limitNumber: 1
    }
  },
  methods: {
    beforeImageUpload(file) {
      const isImg = file.type === 'image/jpeg' || file.type === 'image/png'
      if (!isImg) {
        this.$message.error('只能上传png或jpg格式')
        return false
      }
      return true
    },
    handleImageSuccess(response) {
      if (response.code === 0) {
        this._setImg(response.data)
      }
    },
    _setImg(value) {
      if (!this.imgKey) {
        throw new Error('请指定imgKey')
      }
      const keys = this.imgKey.split('.')
      const len = keys.length
      let p = this
      for (let i = 0; i < len; i++) {
        if (i === len - 1) {
          p[keys[i]] = value
        } else {
          p = p[keys[i]]
        }
      }
    }
  }
}
