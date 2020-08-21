package utils

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"os"
	"strings"
)

var (
	Efc *casbin.Enforcer
	Adp *gormadapter.Adapter
	pwd, _ = os.Getwd()
)

// 自定义 matcher 函数; 主要用于获取不带查询字符传的 url
func PathMatcher(rpath string, dpath string) bool {
	path := strings.Split(rpath, "?")[0]
	return util.KeyMatch2(path, dpath)
}

func PathMatcherFunc(args ...interface{}) (interface{}, error) {
	rpath := args[0].(string)
	dpath := args[1].(string)

	//fmt.Println("我获取的 rpath 和 dpath: ", rpath, dpath)
	return PathMatcher(rpath, dpath), nil
}

func permsInit() {
	if Adp, err = gormadapter.NewAdapterByDB(Db); err != nil {
		panic(fmt.Sprintf("初始一个 mysql adapter 失败, 错误信息: %s", err.Error()))
	}
	RM := NewRoleManager(10)
	if Efc, err = casbin.NewEnforcer(pwd + "/config/casbin_model.conf", Adp); err != nil {
		panic(fmt.Sprintf("初始一个 enforcer 对象失败, 错误信息: %s", err.Error()))
	}
	Efc.SetRoleManager(RM)
	Efc.EnableAutoSave(true)
	Efc.AddFunction("PathMatcher", PathMatcherFunc)
	Efc.BuildRoleLinks()
}
