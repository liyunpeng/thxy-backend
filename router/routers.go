package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"thxy/api"
	"thxy/logger"
	"thxy/middleware/cors"
	"thxy/setting"
)

func InitRouter() *gin.Engine {
	conf := setting.TomlConfig
	logger.Info.Println(conf)

	//paths := conf.Paths
	//if conf.App.Runmode == settings.RunmodeProd {
	//	fmt.Printf("app runmode: %s %s", conf.App.Runmode, settings.RunmodeProd)
	//	gin.SetMode(gin.ReleaseMode)
	//}
	r := gin.New()
	//r.Use(gin.Logger())
	//r.Use(gin.Recovery())
	//cors.
	r.Use(cors.AddCorsHeaders())
	//r.Use(auth.ImgCheck())
	//{
	//	//r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	//	//r.POST("/upload", api.UploadImage)
	//	r.Static(settings.StaticHead, paths.Head)
	//	r.Static(settings.StaticAuth, paths.Authentication)
	//	r.Static(settings.StaticProductImg, paths.ProductImg)
	//}

	baseGroup := r.Group("/")
	{
		baseGroup.POST("/version", func(c *gin.Context) {
			//msg := mes.ByCtx(c, mes.SearchSuccess)
			//data := settings.AppConfig.App.Version
			//data := settings.Version
			fmt.Println("11111111")
			c.JSON(200, "v1")
		})

		//baseGroup.POST("/fileItem", api.PaymentOrders)
		adminGroup := baseGroup.Group("/api")
		//adminGroup.Use(session.CheckAdminSession())
		//adminGroup.Use(admin.AdminAccessRightFilter())
		{
			adminGroup.POST("/login", api.Login)
			adminGroup.POST("/fileUpload", api.FileUpload)
			adminGroup.GET("/fileDownload", api.FileDownload)
			adminGroup.POST("/multiUpload", api.MultiUpload)
			adminGroup.POST("/getCourseTypes", api.GetCourseTypes)
			adminGroup.POST("/findCourseByTypeId", api.FindCourseByTypeId)
			adminGroup.POST("/findCourseFileById", api.FindCourseFileById)
			adminGroup.POST("/findCourseFileByCourseId", api.FindCourseFileByCourseId)
			//adminGroup.POST("/fileItem", api.PaymentOrders)

		}
	}

	return r
}
