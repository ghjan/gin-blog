package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/ghjan/gin-blog/models"
	"github.com/ghjan/gin-blog/pkg/e"
	"github.com/ghjan/gin-blog/pkg/logging"
	"github.com/ghjan/gin-blog/pkg/setting"
	"github.com/ghjan/gin-blog/util"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

//获取单个文章
func GetArticle(c *gin.Context) {
	id, _ := com.StrTo(c.Param("id")).Int()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	var data interface{}
	if ! valid.HasErrors() {
		if models.ExistArticleByID(id) {
			data = models.GetArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

//获取多个文章
func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state, _ = com.StrTo(arg).Int()
		maps["state"] = state

		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId, _ = com.StrTo(arg).Int()
		maps["tag_id"] = tagId

		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	}

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		code = e.SUCCESS

		data["lists"] = models.GetArticles(util.GetPage(c), setting.PageSize, maps)
		data["total"] = models.GetArticleTotal(maps)

	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

//新增文章
func AddArticle(c *gin.Context) {
	tag_id, _ := c.GetPostForm("tag_id")
	tagId, _ := com.StrTo(tag_id).Int()
	title, _ := c.GetPostForm("title")
	desc, _ := c.GetPostForm("desc")
	content, _ := c.GetPostForm("content")
	createdBy, _ := c.GetPostForm("created_by")
	state_ := c.DefaultPostForm("state", "0")
	state, _ := com.StrTo(state_).Int()

	valid := validation.Validation{}
	valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("简述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS
	data := make(map[string]interface{})
	err_map := make(map[string]string)
	if ! valid.HasErrors() {
		if models.ExistTagByID(tagId) {
			data["tag_id"] = tagId
			data["title"] = title
			data["desc"] = desc
			data["content"] = content
			data["created_by"] = createdBy
			data["state"] = state

			models.AddArticle(data)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			err_map[err.Key] = err.Message
			logging.Info(err.Key, err.Message)
		}
	}
	result := gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	}
	if err_map != nil && len(err_map) > 0 {
		result["errors"] = err_map
	}
	c.JSON(http.StatusOK, result)
}

//修改文章
func EditArticle(c *gin.Context) {
	valid := validation.Validation{}

	id, _ := com.StrTo(c.Param("id")).Int()
	tag_id, _ := c.GetPostForm("tag_id")
	tagId, _ := com.StrTo(tag_id).Int()
	title, _ := c.GetPostForm("title")
	desc, _ := c.GetPostForm("desc")
	content, _ := c.GetPostForm("content")
	modifiedBy, _ := c.GetPostForm("modified_by")

	var state int = -1
	if arg, _ := c.GetPostForm("state"); arg != "" {
		state, _ = com.StrTo(arg).Int()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.MaxSize(title, 100, "title").Message("标题最长为100字符")
	valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")
	valid.MaxSize(content, 65535, "content").Message("内容最长为65535字符")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")

	code := e.INVALID_PARAMS
	data := make(map[string]interface{})
	err_map := make(map[string]string)
	if ! valid.HasErrors() {
		if models.ExistArticleByID(id) {
			if models.ExistTagByID(tagId) {
				if tagId > 0 {
					data["tag_id"] = tagId
				}
				if title != "" {
					data["title"] = title
				}
				if desc != "" {
					data["desc"] = desc
				}
				if content != "" {
					data["content"] = content
				}

				data["modified_by"] = modifiedBy

				models.EditArticle(id, data)
				code = e.SUCCESS
			} else {
				code = e.ERROR_NOT_EXIST_TAG
			}
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			err_map[err.Key] = err.Message
			logging.Info(err.Key, err.Message)
		}
	}

	result := gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	}
	if err_map != nil && len(err_map) > 0 {
		result["errors"] = err_map
	}
	c.JSON(http.StatusOK, result)
}

//删除文章
func DeleteArticle(c *gin.Context) {
	id, _ := com.StrTo(c.Param("id")).Int()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		if models.ExistArticleByID(id) {
			models.DeleteArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}
