package router

import (
	"github.com/gin-gonic/gin"
	"thxy/api"
	"thxy/api/admin"
	"thxy/api/user"
	"thxy/logger"
	"thxy/middleware/cors"
	"thxy/middleware/session"
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

		adminGroup := baseGroup.Group("/adminApi")
		adminGroup.Use(session.CheckAdminSession())
		{
			// login
			adminGroup.POST("/login", api.Login)
			adminGroup.POST("/updatePwd", api.UpdatePwd)

			// log
			adminGroup.POST("/logList", api.GetLogList)
			adminGroup.GET("/logDownload", api.LogDownload)

			// upload
			adminGroup.POST("/multiUpload", api.MultiUpload)

			// course type
			adminGroup.POST("/getCourseTypes", api.GetCourseTypes)
			adminGroup.POST("/updateCourseType", api.UpdateCourseType)
			adminGroup.POST("/addCourseType", api.AddCourseType)
			adminGroup.POST("/deleteCourseType", api.DeleteCourseType)
			adminGroup.POST("/coursePictureUpload", api.CoursePictureUpload)

			// course
			adminGroup.POST("/adminGetAllCourseIds", api.AdminGetAllCourseIds)
			adminGroup.POST("/adminGetAllCourseType", api.AdminGetAllCourseType)
			adminGroup.POST("/addCourse", api.AddCourse)
			adminGroup.POST("/updateCourse", api.UpdateCourse)
			adminGroup.POST("/deleteCourse", api.DeleteCourse)

			adminGroup.POST("/findCourseByTypeId", api.FindCourseByTypeId)

			// course file
			adminGroup.POST("/findCourseFileByCourseId", api.FindCourseFileByCourseId)
		}

		userGroup := baseGroup.Group("/api")

		{
			userGroup.POST("/getConfig", api.GetConfig)

			// log
			userGroup.POST("/fileUpload", api.FileUpload)

			// download
			userGroup.GET("/fileDownload", api.FileDownload)
			userGroup.GET("/apkUpload", api.ApkUpload)

			// courseType
			userGroup.POST("/getCourseTypesOk", api.GetCourseTypesOk)

			// course
			userGroup.POST("/findCourseByTypeIdOk", api.FindCourseByTypeIdOkhttp)
			userGroup.POST("/getCourseById", api.GetCourseById)
			userGroup.POST("/findCourseByTypeIdAndUpdateVersion", api.FindCourseByTypeIdAndUpdateVersion)

			// courseFile
			userGroup.POST("/findCourseFileById", api.FindCourseFileById)
			userGroup.POST("/findCourseFileByCourseIdAndUpdateVersion", api.FindCourseFileByCourseIdAndUpdateVersion)
			userGroup.POST("/getLatest", api.GetLatestCourseFile)
			userGroup.POST("/updateUserListenedFiles", api.UpdateUserListenedFiles)
			userGroup.POST("/findCourseFileByCourseIdOk", api.FindCourseFileByCourseIdOk)
			userGroup.POST("/findCourseFileByCourseIdOkhttpV1", api.FindCourseFileByCourseIdOkhttpV1)
			userGroup.POST("/findUserListenedFilesByCodeAndCourseId", api.FindUserListenedFilesByCodeAndCourseId)

			// weixin
			userGroup.POST("/wxBind", user.WXBind)
			userGroup.POST("/wxLogin", user.WXLogin)
			userGroup.POST("/wxToken", user.WXToken)
		}
	}

	return r
}
