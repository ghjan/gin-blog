package util

import (
	"github.com/unknwon/com"
	"github.com/gin-gonic/gin"

	"github.com/ghjan/gin-blog/pkg/setting"
)

// GetPage get page parameters
func GetPage(c *gin.Context) int {
	result := 0
	page := com.StrTo(c.Query("page")).MustInt()
	if page > 0 {
		result = (page - 1) * setting.AppSetting.PageSize
	}

	return result
}
