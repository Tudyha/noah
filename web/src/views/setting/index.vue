<template>
  <el-form ref="form" label-width="120px">
    <el-form-item label="待定">
      <el-input v-model="form.name" />
    </el-form-item>
  </el-form>
</template>
<script>
export default {
  components: {
  },
  data() {
    return {
      form: {
        name: ''
      }
    }
  },
  computed: {
  },
  created() {
  },
  methods: {
  }
}
</script>
