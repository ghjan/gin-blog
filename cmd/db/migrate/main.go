package main

import (
	"github.com/ghjan/gin-blog/models"
	connection "github.com/ghjan/gin-blog/pkg/database"
	"github.com/ghjan/gin-blog/pkg/logging"
	"github.com/ghjan/gin-blog/pkg/setting"
	"github.com/ghjan/gin-blog/pkg/util"
)

func init() {
	setting.Setup()
	models.Setup()
	logging.Setup()
	util.Setup()
}

func main() {
	db := connection.GetInstance()
	defer connection.Close()

	db.AutoMigrate(&models.Article{}, &models.Auth{}, &models.Tag{})
	db.Model(&models.Article{}).AddForeignKey("tag_id", "blog_tags(id)", "SET NULL", "CASCADE")
}
