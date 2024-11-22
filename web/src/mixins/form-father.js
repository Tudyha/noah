export default {
  data() {
    return {
      editOj: {},
      addOj: {},
      isAddShow: false,
      isUpdateShow: false
    }
  },
  methods: {
    hide(oj) {
      if (oj.name === 'add') {
        this.isAddShow = false
      }
      if (oj.name === 'edit') {
        this.isEditShow = false
      }
      if (oj.name === 'check') {
        this.isCheckShow = false
      }
      if (oj.name === 'update') {
        this.isUpdateShow = false
      }
      if (oj.name === 'ban') {
        this.isBanShow = false
      }
      if (oj.name === 'video') {
        this.isVideoShow = false
      }
      if (oj.name === 'label') {
        this.isLabelShow = false
      }
      if (oj && oj.isFresh) {
        this.fetchList()
      }
    },
    showAddDialog() {
      this.isAddShow = true
    },
    showCheckDialog() {
      console.log('check')
      this.isCheckShow = true
    },
    showEditDialog(row) {
      this.editOj = { ...row }
      this.isEditShow = true
    }
  }
}

