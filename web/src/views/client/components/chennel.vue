<template>
    <div>
        <el-container>
            <el-header>
                <el-button type="primary" @click="dialogVisible = true">新增</el-button>
            </el-header>
            <el-main>
                <el-table :data="list">
                    <el-table-column label="模式">
                        <template slot-scope="scope">
                            {{ m.channelType[scope.row.channelType] }}
                        </template>
                    </el-table-column>
                    <el-table-column label="端口">
                        <template slot-scope="scope">
                            {{ scope.row.serverPort }}
                        </template>
                    </el-table-column>
                    <el-table-column label="目标IP">
                        <template slot-scope="scope">
                            {{ scope.row.clientIp }}
                        </template>
                    </el-table-column>
                  <el-table-column label="目标端口">
                    <template slot-scope="scope">
                      {{ scope.row.clientPort }}
                    </template>
                  </el-table-column>
                </el-table>
            </el-main>


        </el-container>
        <el-dialog title="新增隧道" :visible.sync="dialogVisible" width="30%">
            <el-form :model="form">
                <el-form-item label="模式">
                    <el-select v-model="form.channelType" placeholder="请选择">
                        <el-option label="TCP" :value="2"></el-option>
                        <el-option label="UDP" :value="3"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="端口">
                    <el-input v-model.number="form.serverPort" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="目标ip">
                    <el-input v-model="form.clientIp" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="port">
                    <el-input v-model.number="form.clientPort" autocomplete="off"></el-input>
                </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click="dialogVisible = false">取 消</el-button>
                <el-button type="primary" @click="newChannel">确 定</el-button>
            </div>
        </el-dialog>
    </div>
</template>

<script>
import { newChannel, fetchList } from '@/api/channel'
import * as m from '@/map/channel'

export default {
    props: {
        id: {
            type: Number,
            required: true
        }
    },
    data() {
        return {
          m,
            list: [],
            dialogVisible: false,
            form: {
                channelType: null,
                serverPort: null,
                clientIp: null,
                clientPort: null
            },
        };
    },
    methods: {
      fetchList() {
        fetchList(this.id)
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
        newChannel() {
            this.form.id = this.id
            newChannel(this.form).then(res => {
                if (res.code === 0) {
                    this.$message({
                        message: '新增隧道成功',
                        type: 'success'
                    })
                    this.dialogVisible = false
                    this.fetchList()
                } else {
                    this.$message({
                        message: res.msg,
                        type: 'error'
                    })
                }
            })
        }
    },
    components: {
    },
    created() {
        this.fetchList();
    },
}
</script>
