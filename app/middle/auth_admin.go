package middle

import (
	"Kronos/library/apgs"
	"Kronos/library/casbin_helper"
	"Kronos/library/session"
	"encoding/json"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// AuthAdmin 中间件
func AuthAdmin(enforcer *casbin.SyncedEnforcer, nocheck ...casbin_helper.DontCheckFunc) gin.HandlerFunc {

	return func(c *gin.Context) {
		if casbin_helper.DontCheck(c, nocheck...) {
			c.Next()
			return
		}
		fmt.Println("123123123")
		// Session 判断权限
		get := session.GetSession(c, session.UserKey)
		var idss map[string]interface{}
		canGet := json.Unmarshal(get.([]byte), &idss)
		if canGet != nil {
			ginview.HTML(c, http.StatusUnauthorized, "err/401", apgs.NewApiRedirect(200, "无权限访问该内容", "/admin/login"))
		}
		// 超级管理员不验证权限
		if idss["IsSuper"] == 1 {
			c.Next()
			return
		}

		userId := strconv.Itoa(int(idss["id"].(float64)))
		p := strings.ToLower(c.Request.URL.Path)
		m := strings.ToLower(c.Request.Method)

		var b bool
		var err error
		if b, err = enforcer.Enforce(userId, p, m); err != nil {
			// TODO 判断是是否为调试模式
			// TODO 调试模式下 判断 异步，同步 返回 JSON HTML
			//c.JSON(403, helpers.NewApiReturn(401, err.Error(), b))
			//c.AbortWithStatus(403)
			ginview.HTML(c, http.StatusForbidden, "err/403", apgs.NewApiReturn(403, err.Error(), nil))
			c.Abort()
			return
		}
		if !b {
			//c.JSON(401, helpers.NewApiReturn(401, "权限验证失败", b))
			//c.Abort()
			//fmt.Println("Check:" + strconv.FormatBool(b))
			//c.Redirect(302, "/admin/login")
			ginview.HTML(c, http.StatusUnauthorized, "err/401", apgs.NewApiRedirect(200, "无权限访问该内容", "/admin/login"))
			c.Abort()
			return
		}
		c.Next()

	}
}
