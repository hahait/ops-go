## 说明如下：
1. 后端基于 gin 框架开发;
2. 前端基于 [vue-element-admin](https://github.com/PanJiaChen/vue-element-admin);
3. 仅实现用户相关功能: 用户的增删改查, 用户组的增删改查, 权限表自动生成及查;
4. 另外还有如下功能:
   - JWT 登陆认证; 修改了 gin-jwt 源码, 实现 urls 白名单(免 token 验证), 并提了 [issue](https://github.com/appleboy/gin-jwt/issues/253);
   - 权限校验-casbin; 自定义了 RoleManager 和 使用了 gorm adapter, 实现了类 django 风格的 RBAC 权限, 但资源对象是 URL 而不是模型对象;
   - 使用 gorm 操作 mysql;
   - 全局的错误处理;
   - 类 django 的代码文件组织
