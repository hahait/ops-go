package account

import (
	"fmt"
	"ops.was.ink/opsweb/utils"
)

func init() {
	//defer utils.Db.Close()
	// 模型初始化
	err := utils.Db.AutoMigrate(&User{}, &Group{}, &UserGroups{}, &Permissions{})
	if err != nil {
		fmt.Println("自动迁移 model 到数据库失败, 错误信息: ", err.Error())
	}

}