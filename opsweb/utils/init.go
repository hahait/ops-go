package utils

func init(){
	//* 1. 初始化数据库连接
	dbInit()
	//* 2. 初始化模型字段验证
	modelValidatorInit()
	//* 3. 初始化权限验证
	permsInit()
}
