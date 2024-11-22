<template>
  <!-- 通用表格 -->
  <div :class="classObj">
    <div class="common-table__content">
      <lb-table
        ref="commonTable"
        border
        tooltip-effect="dark"
        v-bind="$attrs"
        :data="data || []"
        :column="lbColumns"
        :pagination="false"
        :row-height="rowHeight"
        height="100%"
        align="center"
        header-row-class-name="common-table-header"
        v-on="$listeners"
      />
    </div>
    <div v-if="hasPagination" class="common-table__pagination" style="text-align: right;">
      <component
        :is="paginationType"
        v-bind="$attrs"
        :page="page"
        :limit="limit"
        v-on="$listeners"
      />
    </div>
  </div>
</template>
<script>
import LbTable from '@/components/LbTable'
import Pagination from '@/components/Pagination'
import NotTotalPagination from '@/components/NotTotalPagination'
import { formatterField } from '@/utils'

// const column = {
//   type: String, // 如：selection,
//   render: Function, // render函数
//   prop: String, // 属性key，唯一
//   propType: String, // 指定属性类型，根据属性类型解析属性值，不指定直接返回对应属性值
//   map: Object, // propType 是 map时必填，对应属性值的显示值，可以操作formatterField函数
//   yes: String, // 可以指定 propType 是 boolean时，属性值为true时的显示值
//   no: String, // 可以指定 propType 是 boolean时，属性值为false时的显示值
//   formatter: Function, // 自己解析属性值,
//   click: Function, // 点击属性值时的回调函数
//   opts: [{// 多个操作，常用于操作列表
//     // 请参考el-button 属性
//     label: String, // 按钮label
//     type: String, // 按钮 type
//     click: Function
//   }],
// }

export default {
  name: 'CommanTable',
  components: { Pagination, NotTotalPagination, LbTable },
  props: {
    data: { type: Array, default: () => ([]) },
    columns: { type: Array, default: () => {} },
    hasPagination: { type: Boolean, default: true },
    isFlexLayout: { type: Boolean, default: true }, // 默认flex 布局，需要指定外层容器高度；否则自动乘客
    page: { type: Number, default: 1 },
    limit: { type: Number, default: 100 },
    dataChange: { type: Function, default: null }, // 指定data改变时触发调用的函数,
    rowHeight: { type: Number, default: 32 }, // 指定每一行高度，如果是0代表不限制高度,
    paginationType: { type: String, default: 'pagination' } // 置顶分页组件类型
  },
  computed: {
    lbColumns() {
      return this.columns.filter(c => {
        // 展示条件
        // 1. 指定hidden为false
        // 2. 指定permission且有权限
        // 3. 有opts属性(一般是操作列)，没有指定permission属性，检查opts 子操作是否一个有权限或者不做权限校验的
        let isShow = !c.hidden && (!c.permission)
        if (isShow && !c.permission && c.opts) {
          isShow = c.opts.some(o => !o.permission)
        }
        return isShow
      }).map(this.formatterColumn)
    },
    classObj() {
      return {
        'common-table': true,
        'is-flex-layout': this.isFlexLayout
      }
    },
    showOverflowTooltip() {
      return this.rowHeight > 0
    },
    elTable() {
      return this.$refs.commonTable && this.$refs.commonTable.$refs.elTable
    }
  },
  watch: {
    data(newData, oldData) {
      if (newData !== oldData) {
        this.dataChange ? this.dataChange(newData) : this.$refs.commonTable.clearSelection()
      }
    }
  },
  methods: {
    formatterColumn(c) {
      let overflowHiddenStyle = ''
      if (c.hasOwnProperty('overflowHidden')) {
        overflowHiddenStyle = 'text-overflow: -o-ellipsis-lastline;overflow: hidden; text-overflow: ellipsis;display: -webkit-box;-webkit-line-clamp: 2;line-clamp: 2;-webkit-box-orient: vertical;'
      } else {
        overflowHiddenStyle = 'display: inline-block; width: 100%; overflow: hidden;text-overflow:ellipsis;white-space: nowrap;'
      }
      if (c.children) {
        return { ...c, children: c.children.filter(ch => !ch.hidden).map(ch => ({ align: 'center', ...this.formatterColumn(ch) })) }
      } else if (c.type || c.render) {
        return c
      } else if (c.click && (!c.clickPermission)) {
        const { formatter, ...rest } = c
        return {
          sortable: false,
          ...rest,
          showOverflowTooltip: false,
          render: (h, { row }) => {
            const content = formatterField({ formatter, ...rest }, row)
            return (
              <el-button
                type='text'
                style={`padding:0!important;line-height:1.2em;${overflowHiddenStyle}`}
                onClick={() => c.click(row)}
              >
                {
                  c.showOverflowTooltip ? (
                    <span title={content}>{content}</span>
                  ) : content
                }
              </el-button>
            )
          }
        }
      } else if (c.opts) {
        return {
          ...c,
          render: (h, { row }) => (
            <div class='common-table-opts'>
              {
                c.opts.map(({ permission, ...opt }) => {
                  if ((opt.hidden && opt.hidden(row))) {
                    return null
                  }
                  const props = { ...opt, ...(opt.props ? opt.props(row) : {}) }
                  return (
                    <el-button
                      key={props.label}
                      type={props.type}
                      disabled={props.disabled}
                      {...props}
                      onClick={() => opt.click(row)}
                    >
                      {props.label}
                    </el-button>
                  )
                })
              }
            </div>
          )
        }
      }
      const { formatter, ...rest } = c
      return {
        sortable: false,
        ...rest,
        showOverflowTooltip: false,
        render: (h, { row }) => {
          const content = formatterField({ formatter, ...c }, row)
          return <span title={c.showOverflowTooltip ? content : ''} style={overflowHiddenStyle}>{content}</span>
        }
      }
    },
    clearSelection() {
      this.$refs.commonTable.clearSelection()
    },
    toggleRowSelection(row, selected) {
      this.$refs.commonTable.toggleRowSelection(row, selected)
    },
    toggleRowSelectionByRowKey(value, selected) {
      this.$refs.commonTable.toggleRowSelectionByRowKey(value, selected)
    },
    setCurrentRow(row) {
      this.$refs.commonTable.setCurrentRow(row)
    },
    toggleRowExpansion(row, expanded) {
      this.$refs.commonTable.toggleRowExpansion(row, expanded)
    }
  }
}
</script>
<style lang="scss">
.common-table-opts {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  .el-button {
    padding: 4px 10px!important;
  }
}
.common-table.is-flex-layout {
  .lb-table,
  .lb-table > div {
    height: 100%;
  }
}
.common-table-row--no-indent {
  .el-table__indent,
  .el-table__placeholder {
    display: none;
  }
}
</style>
.<style lang="scss" scoped>
.common-table.is-flex-layout {
  height: 100%;
  display: flex;
  flex-direction: column;
  .common-table__content {
    flex: 1;
    overflow: hidden;
  }
}
</style>
