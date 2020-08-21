<template>
  <div class="permission-list">
    <!-- 权限列表 -->
    <el-table
      :data="permissionList"
      :span-method="objectSpanMethod"
      stripe
      border
      size="mini"
      highlight-current-row
      style="min-width: 100%;">
      <el-table-column
        label="请求路径"
        align="center">
        <template slot-scope="scope">
          <span style="margin-left: 10px">{{ scope.row.path }}</span>
        </template>
      </el-table-column>
      <el-table-column
        label="请求方法"
        align="center">
        <template slot-scope="scope">
          <span style="margin-left: 10px">{{ scope.row.method }}</span>
        </template>
      </el-table-column>
      <el-table-column
        label="描述"
        align="center">
        <template slot-scope="scope">
          <span v-if="scope.row.changePermissionNameIf" style="margin-left: 10px">{{ scope.row.describe }}</span>
          <el-input v-if="!scope.row.changePermissionNameIf" v-model="scope.row.describe" size="small"/>
        </template>
      </el-table-column>
      <el-table-column
        label="管理权限"
        align="center">
        <template slot-scope="scope">
          <el-button v-if="scope.row.changePermissionNameIf" type="primary" size="mini" plain @click="changePermissionNameBtn(scope.row)">
            修改描述
          </el-button>
          <el-button v-if="!scope.row.changePermissionNameIf" type="success" size="mini" plain @click="commitChangePermissionNameBtn(scope.row)">
            确定
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>

export default {
  name: 'PermissionList',
  props: {
    'permissionList': {
      type: Array,
      required: true
    },
    'spanArr': {
      required: true,
      type: Array,
      default: []
    }
  },
  data: function() {
    return {
      row: {}
    }
  },
  methods: {
    // 修改权限名称按钮
    changePermissionNameBtn: function(row) {
      row.changePermissionNameIf = false
    },
    // 修改权限名称的表单提交
    commitChangePermissionNameBtn: function(row) {
      const data = { id: row.id, describe: row.describe }
      this.$emit('commitChangePermissionName', data)
      row.changePermissionNameIf = true
    },
    // 合并行的自定义方法
    objectSpanMethod: function({ row, column, rowIndex, columnIndex }) {
      if (columnIndex === 0) {
        const _row = this.spanArr[rowIndex]
        const _col = _row > 0 ? 1 : 0
        return {
          rowspan: _row,
          colspan: _col
        }
      }
    }
  }
}
</script>

<style rel="stylesheet/scss" lang="scss">
  .permission-list {
    margin: 0px 10px 10px;
    .el-table__body-wrapper .el-table__body {
      min-width: 100%;
    }
    .el-table__header-wrapper .el-table__header {
      min-width: 100%;
    }
  }
</style>
