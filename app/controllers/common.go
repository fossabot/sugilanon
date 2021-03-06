package controllers

import (
	"os"
	"runtime"
	"strings"

	"github.com/XanderDwyl/sugilanon/app/libs/tmplname"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// OutputJSON ...
func OutputJSON(c *gin.Context, status, msg string) {
	c.JSON(200, gin.H{
		"status":  status,
		"message": msg,
	})
}

// OutputDataJSON ...
func OutputDataJSON(c *gin.Context, status, msg string, data gin.H) {
	c.JSON(200, gin.H{
		"status":  status,
		"message": msg,
		"data":    data,
	})
}

// RenderHTML ...
func RenderHTML(c *gin.Context, data gin.H) {
	pc, _, _, _ := runtime.Caller(1)
	callerName := runtime.FuncForPC(pc).Name()

	for strings.Contains(callerName, ".") {
		a := strings.Index(callerName, ".")
		callerName = callerName[a+1:]
	}

	tmpl := tmplname.Convert(callerName)
	_, err := os.Stat("app/views/" + tmpl + ".tmpl")
	if err != nil {
		c.String(200, "%s not found", "app/views/"+tmpl+".tmpl")

		return
	}

	RenderTemplate(c, tmpl, data, 200)
}

// RenderTemplate ...
func RenderTemplate(c *gin.Context, tmpl string, data gin.H, statusCode int) {
	data["fb_app_id"] = os.Getenv("FB_APP_ID")
	data["is_login"] = IsLogin(c)
	data["app_user_id"] = GetAppUserId(c)
	data["is_verified"] = IsVerified(c)
	data["facebook_id"] = GetFacebookId(c)
	data["username"] = GetAppUsername(c)

	if _, ok := data["page"]; !ok {
		data["page"] = "SUGILANON"
	}

	c.HTML(statusCode, tmpl, data)
}

func IsLogin(c *gin.Context) bool {
	isLogin, ok := c.Get("is_login")
	if ok && isLogin.(bool) {
		return true
	}

	session := sessions.Default(c)
	flag := session.Get("is_login")
	if flag != nil {
		val, ok := flag.(int)
		if ok && val == 1 {
			return true
		}
	}

	return false
}

func IsVerified(c *gin.Context) bool {
	session := sessions.Default(c)
	flag := session.Get("is_verified")
	if flag != nil {
		val, ok := flag.(bool)
		if ok {
			return val
		}
	}

	return false
}

func GetFacebookId(c *gin.Context) string {
	session := sessions.Default(c)
	facebookId := session.Get("facebook_id")
	if facebookId != nil {
		val, ok := facebookId.(string)
		if ok {
			return val
		}
	}

	return ""
}

func GetAppUserId(c *gin.Context) int64 {
	session := sessions.Default(c)
	appUserId := session.Get("app_user_id")
	if appUserId != nil {
		val, ok := appUserId.(int64)
		if ok {
			return val
		}
	}

	return 0
}

func GetAppUsername(c *gin.Context) string {
	session := sessions.Default(c)
	appUsername := session.Get("username")
	if appUsername != nil {
		val, ok := appUsername.(string)
		if ok {
			return val
		}
	}

	return ""
}

func SetAuth(c *gin.Context, user User) {
	session := sessions.Default(c)

	session.Set("is_login", 1)
	session.Set("app_user_id", user.AppUserId)
	session.Set("is_verified", user.IsVerified)
	session.Set("facebook_id", user.FacebookId)
	session.Set("username", user.Username)
	session.Set("name", user.Name)
	session.Set("email", user.Email)
	session.Save()
}

func ClearAuth(c *gin.Context) {
	session := sessions.Default(c)

	session.Delete("is_login")
	session.Delete("app_user_id")
	session.Delete("is_verified")
	session.Delete("facebook_id")
	session.Delete("username")
	session.Delete("name")
	session.Delete("email")
	session.Save()
}

func NoRoute(c *gin.Context) {
	c.Redirect(302, "/")
	c.Abort()

	return
}
