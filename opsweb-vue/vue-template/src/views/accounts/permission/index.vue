<template>
  <div class="permission-container">
    <!-- 添加/搜索 -->
    <el-form :inline="true" class="permission-form-inline" size="small">
      <el-form-item>
        <el-input v-model="params.path" placeholder="请求路径" />
      </el-form-item>
      <el-form-item>
        <el-input v-model="params.method" placeholder="请求方法" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="searchSubmit">搜索</el-button>
      </el-form-item>
    </el-form>
    <!-- 用户组列表 -->
    <permission-list
      v-loading="listLoading"
      :permission-list="permissionList"
      :span-arr="spanArr"
      element-loading-text="拼命加载中"
      element-loading-spinner="el-icon-loading"
      element-loading-background="#f3f7f3cc"
      @commitChangePermissionName="commitChangePermissionName"/>
    <!-- 分页 -->
    <my-page :total="total" :page_size="params.page_size" :current_page="params.page" @CurrentChange="ChangeCurrentPage" @SizeChange="ChangePageSize"/>
  </div>
</template>

<script>
import { getPermission, editPermission } from '@/api/accounts/permission'
import MyPage from '@/components/Pagination'
import PermissionList from './permission_list'

export default {
  name: 'Permission',
  components: {
    MyPage,
    PermissionList
  },
  data: function() {
    return {
      permissionList: [],
      listLoading: true,
      total: 0,
      spanArr: [],
      pos: 0,
      permissionDialogVisible: false,
      params: {
        page: 1,
        page_size: 10,
        path: '',
        method: ''
      }
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    // 获取用权限列表
    fetchData: function() {
      this.listLoading = true
      getPermission(this.params).then(res => {
        this.total = res.total
        // 为了实现在表格中直接修改权限描述信息
        res.results.forEach(per => {
          this.$set(per, 'changePermissionNameIf', true)
        })
        this.permissionList = res.results
        this.getSpanArr(res.results)
        this.listLoading = false
      })
    },
    // 修改权限描述提交
    commitChangePermissionName: function(data) {
      editPermission(data).then(
        () => {
          this.$message({
            message: 'Good,修改权限备注成功',
            type: 'success'
          })
          this.fetchData()
        },
        (error) => {
          var errmsg = ''
          for (var v in error.response.data) {
            errmsg += error.response.data[v]
          }
          this.$message({
            message: 'Oh,修改权限备注失败,错误信息: ' + errmsg,
            type: 'error'
          })
        }
      )
    },
    // 合并行数
    getSpanArr: function(data) {
      this.spanArr = []
      for (var i = 0; i < data.length; i++) {
        if (i === 0) {
          this.spanArr.push(1)
          this.pos = 0
        } else {
          // 判断当前元素与上一个元素是否相同
          if (data[i].path === data[i - 1].path) {
            this.spanArr[this.pos] += 1
            this.spanArr.push(0)
          } else {
            this.spanArr.push(1)
            this.pos = i
          }
        }
      }
    },
    // 搜索提交
    searchSubmit: function() {
      this.params.page = 1
      this.params.page_size = 10
      this.fetchData()
    },
    // 分页改变页码
    ChangeCurrentPage: function(val) {
      this.params.page = val
      this.fetchData()
    },
    // 分页改变展示的数据行数
    ChangePageSize: function(val) {
      this.params.page_size = val
      this.fetchData()
    }
  }
}
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
  .permission-container {
    margin: 30px 10px 10px;
    .permission-form-inline {
      margin: 10px 10px 0px;
      text-align: right;
    }
    .permissionDialogForm {
      width: 90%;
    }
  }
</style>
