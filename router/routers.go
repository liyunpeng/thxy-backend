package router

import (
	"github.com/gin-gonic/gin"
	"thxy/api"
	"thxy/api/admin"
	"thxy/api/user"
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
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
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
			c.JSON(200, "v1")
		})

		baseGroup.POST("/adminLogin", admin.Login)
		adminGroup := baseGroup.Group("/api")
		//adminGroup.Use(session.CheckAdminSession())
		//adminGroup.Use(admin.AdminAccessRightFilter())
		{
			adminGroup.POST("/getConfig", api.GetConfig)

			// login
			adminGroup.POST("/login", api.Login)
			adminGroup.POST("/updatePwd", api.UpdatePwd)

			// log
			adminGroup.POST("/fileUpload", api.FileUpload)
			adminGroup.POST("/logList", api.GetLogList)
			adminGroup.GET("/logDownload", api.LogDownload)

			// upload
			adminGroup.POST("/multiUpload", api.MultiUpload)

			// download
			adminGroup.GET("/fileDownload", api.FileDownload)
			adminGroup.GET("/apkUpload", api.ApkUpload)

			// courseType
			adminGroup.POST("/getCourseTypes", api.GetCourseTypes)
			adminGroup.POST("/getCourseTypesOk", api.GetCourseTypesOk)
			adminGroup.POST("/findCourseByTypeId", api.FindCourseByTypeId)
			adminGroup.POST("/updateCourseType", api.UpdateCourseType)
			adminGroup.POST("/addCourseType", api.AddCourseType)
			adminGroup.POST("/deleteCourseType", api.DeleteCourseType)

			// course
			adminGroup.POST("/findCourseByTypeIdOk", api.FindCourseByTypeIdOkhttp)
			adminGroup.POST("/getCourseById", api.GetCourseById)
			adminGroup.POST("/adminGetAllCourseIds", api.AdminGetAllCourseIds)
			adminGroup.POST("/adminGetAllCourseType", api.AdminGetAllCourseType)
			adminGroup.POST("/addCourse", api.AddCourse)
			adminGroup.POST("/updateCourse", api.UpdateCourse)
			adminGroup.POST("/deleteCourse", api.DeleteCourse)
			adminGroup.POST("/coursePictureUpload", api.CoursePictureUpload)

			// courseFile
			adminGroup.POST("/findCourseFileById", api.FindCourseFileById)
			adminGroup.POST("/findCourseFileByCourseId", api.FindCourseFileByCourseId)
			adminGroup.POST("/findCourseFileByCourseIdAndUpdateVersion", api.FindCourseFileByCourseIdAndUpdateVersion)
			adminGroup.POST("/getLatest", api.GetLatestCourseFile)
			adminGroup.POST("/updateUserListenedFiles", api.UpdateUserListenedFiles)
			adminGroup.POST("/findCourseFileByCourseIdOk", api.FindCourseFileByCourseIdOk)
			adminGroup.POST("/findCourseFileByCourseIdOkhttpV1", api.FindCourseFileByCourseIdOkhttpV1)
			adminGroup.POST("/findUserListenedFilesByCodeAndCourseId", api.FindUserListenedFilesByCodeAndCourseId)

			// weixin
			adminGroup.POST("/wxBind", user.WXBind)
			adminGroup.POST("/wxLogin", user.WXLogin)
			adminGroup.POST("/wxToken", user.WXToken)
		}
	}

	return r
}