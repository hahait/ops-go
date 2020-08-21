package middlewares

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"ops.was.ink/opsweb/account"
	"ops.was.ink/opsweb/utils"
	"time"
)

var (
	identity_key = "username"
	toek_secret = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
)

type Login struct {
	UserName string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// 用于在 jwt 的 payload 中添加的自定义信息
type UserInfo struct {
	Username  string
}

// 用于验证用户密码
func AuthenticatorFunc(c *gin.Context) (interface{}, error) {
	var (
		lg Login
		u account.User
	)
	if err := c.ShouldBindBodyWith(&lg, binding.JSON); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	// 从数据库中查找该用户信息
	if err := utils.Db.Select("name, password").Where("name = ?", lg.UserName).Find(&u).Error; err != nil {
		return nil, jwt.ErrFailedAuthentication
	}
	// 检查密码是否与数据库中一致
	if ok := utils.CheckPassword(lg.Password, u.Password); !ok {
		return nil, jwt.ErrFailedAuthentication
	}

	c.Set("username", lg.UserName)

	return &UserInfo{ Username: lg.UserName}, nil
}

// 加载自定义的用户信息到 jwt 的 payload 中；在用户的 RBAC 权限验证中，可以将用户的权限放到这里，以免每次都要查数据库
func payloadFunc(data interface{}) jwt.MapClaims {
	if u, ok := data.(*UserInfo); ok {
		return jwt.MapClaims{
			identity_key: u.Username,
			// 权限相关可以放到这里
		}
	}
	return jwt.MapClaims{}
}

// 认证失败（用户验证失败, token 过期或无效）时就会调用这个函数返回 http response
func unauthorizedFunc(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code": 1,
		"msg": "token 认证失败",
		"emsg": message,
	})
}

// 登陆成功后，更新该用户在数据库中的 last_login 字段，因此需要自定义 LoginResponse() 函数
// 之所以不在 Authenticator() 函数中更新，是因为 Authenticator() 仅仅是用来验证用户密码，还没有生成 token, 不能证明用户已经登陆成功；
// 另外要注意，由于在 Authenticator() 中已经绑定了一次 request body，因此这里不能再次绑定，所以要修改这 2 个函数中的 ShouldBindJSON() 为 ShouldBindBodyWith()
func loginResponseFunc(c *gin.Context, code int, token string, expire time.Time) {
	var (
		lg Login
		u account.User
	)
	c.ShouldBindBodyWith(&lg, binding.JSON)
	l_time := time.Now()
	utils.ErrorHandler(utils.Db.Model(&u).Where("name = ?", lg.UserName).Updates(account.User{LastLogin: &l_time}).Error, 1, "更新登陆时间失败")
	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"token":  token,
		"expire": expire.Format(time.RFC3339),
	})
}

// 用于在每个 http 请求验证完 token 后，将用户名注入到 context 中, 以便后续的请求处理逻辑能够知道这是哪个用户的行为
func identityHandleFunc(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	c.Set("username", claims[identity_key])
	return claims[identity_key]
}

func JwtAuthTokenInit() *jwt.GinJWTMiddleware {
	gjm := jwt.GinJWTMiddleware{
		Realm:       "jwt auth token",
		Key:         []byte(toek_secret),
		Timeout:     time.Hour * 24,
		MaxRefresh:  time.Hour,
		IdentityKey: identity_key,
		Authenticator: AuthenticatorFunc,
		PayloadFunc: payloadFunc,
		Unauthorized: unauthorizedFunc,
		LoginResponse: loginResponseFunc,
		IdentityHandler: identityHandleFunc,
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc: time.Now,
		WhiteUrlList: []string{"/login"},
	}

	j, err := jwt.New(&gjm)
	if err != nil {
		panic(utils.Errors{
			Code: 1,
			Msg:  "实例化 jwt 对象失败",
			Errmsg: err.Error(),
		})
		return nil
	}

	return j
}


