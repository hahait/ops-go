package account

import "github.com/gin-gonic/gin"

func Routers( e *gin.Engine) {
	userGroup := e.Group("/user")
	{
		userGroup.POST("", createUserHandler)
		userGroup.DELETE("", deleteUserHandler)
		userGroup.PUT("", updateUserHandler)
		userGroup.GET("", queryUserHandler)
		userGroup.GET("/info", queryUserInfoHandler)
		userGroup.PUT("/pwd", updateUserPasswordHandler)
		userGroup.PUT("/admin", updateUserAdminHandler)
		userGroup.PUT("/status", updateUserStatusHandler)
		ugGroup := userGroup.Group("/groups")
		{
			ugGroup.PUT("", updateUserGroupsHandler)
			ugGroup.GET("", queryGroupUsersHandler)
		}
	}
	gGroup := e.Group("/group")
	{
		gGroup.POST("", createGrouphandler)
		gGroup.DELETE("", deleteGroupHandler)
		gGroup.PUT("", updateGroupHandler)
		gGroup.GET("", queryGroupHandler)
		guGroup := gGroup.Group("/users")
		{
			guGroup.PUT("", updateGroupUsersHandler)
			guGroup.GET("", queryGroupUsersHandler)
		}
		puGroup := gGroup.Group("/perms")
		{
			puGroup.PUT("", updateGroupPermsHandlers)
		}

	}
	pGroup := e.Group("/permission")
	{
		pGroup.GET("", queryPermissionHandler)
		pGroup.PUT("", updatePermissionDescribeHandler)
	}
}
