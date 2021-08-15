package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"thxy/logger"
	"thxy/model"
	"thxy/setting"
	"thxy/types"
)

func Login(c *gin.Context) {

	a := types.Response1{
		AccessToken: "add",
	}
	c.JSON(200, a)

}
func FileDownload(c *gin.Context) {
	//r := new(types.DownloadReqeust)
	//c.Bind(r)

	fileName := c.Query("fileName")


	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))//fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
	c.Writer.Header().Add("Content-Type", "application/octet-stream")


	storePath := setting.TomlConfig.Test.FilStore.FileStorePath

	c.File(storePath + fileName )

	logger.Info.Println(" 下载文件")
}
// 单文件上传
func FileUpload(context *gin.Context) {
	file, err := context.FormFile("file")
	if err != nil {
		logger.Info.Println("ERROR: upload file failed. ", err)
		context.JSON(http.StatusInternalServerError, gin.H{
			"msg": fmt.Sprintf("ERROR: upload file failed. %s", err),
		})
	}
	dst := fmt.Sprintf(`./`+file.Filename)
	// 保存文件至指定路径
	err = context.SaveUploadedFile(file, dst)
	if err != nil {
		logger.Info.Println("ERROR: save file failed. ", err)
		context.JSON(http.StatusInternalServerError, gin.H{
			"msg": fmt.Sprintf("ERROR: save file failed. %s", err),
		})
	}
	context.JSON(http.StatusOK, gin.H{
		"msg": "file upload succ.",
		"filepath": dst,
	})
}

func FindCourseFileByCourseId(c *gin.Context) {
	a := new(types.CourseFileReqeust)
	c.Bind(a)

	cc, err := model.FindCourseFileByCourseId(a.Id)

	if err != nil {
		c.JSON(501, err)
	}

	c.JSON(200, cc)
}

func FindCourseFileById(c *gin.Context) {
	a := new(types.CourseFileReqeust)
	c.Bind(a)

	cc, err := model.FindCourseFileById(a.Id)

	if err != nil {
		c.JSON(501, err)
	}

	c.JSON(200, cc)
}


func GetCourseTypes(c *gin.Context) {
	a := new(types.CourseFileReqeust)
	c.Bind(a)

	cc, err := model.FindAllCourseTypes()

	if err != nil {
		c.JSON(501, err)
	}

	c.JSON(200, cc)
}



func FindCourseByTypeId(c *gin.Context) {
	a := new(types.CourseFileReqeust)
	c.Bind(a)

	cc, err := model.FindCourseByTypeId(a.Id)

	if err != nil {
		c.JSON(501, err)
	}

	c.JSON(200, cc)
}
// 多文件上传
func MultiUpload(context *gin.Context) {
	// 多文件上传需要先解析form表单数据
	form, err := context.MultipartForm()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"msg": fmt.Sprintf("ERROR: parse form failed. %s", err),
		})
	}
	// 多个文件上传，要用同一个key
	files := form.File["files"]
	for _, file := range files {
		dst := fmt.Sprint(`./`+file.Filename)
		// 保存文件至指定路径
		err = context.SaveUploadedFile(file, dst)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"msg": fmt.Sprintf("ERROR: save file failed. %s", err),
			})
		}
	}
	context.JSON(http.StatusOK, gin.H{
		"msg": "upload file succ.",
		"filepath": `./`,
	})
}