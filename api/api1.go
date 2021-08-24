package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName)) //fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
	c.Writer.Header().Add("Content-Type", "application/octet-stream")

	storePath := setting.TomlConfig.Test.FilStore.FileStorePath

	c.File(storePath + fileName)

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
	dst := fmt.Sprintf(`./` + file.Filename)
	// 保存文件至指定路径
	err = context.SaveUploadedFile(file, dst)
	if err != nil {
		logger.Info.Println("ERROR: save file failed. ", err)
		context.JSON(http.StatusInternalServerError, gin.H{
			"msg": fmt.Sprintf("ERROR: save file failed. %s", err),
		})
	}
	context.JSON(http.StatusOK, gin.H{
		"msg":      "file upload succ.",
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

func FindCourseFileByCourseIdOk(c *gin.Context) {
	a := new(types.CourseFileReqeustOkhttp)
	c.Bind(a)

	//
	//a.Id = "1"
	//id, _ := strconv.Atoi(a.Id)
	//
	//a := new(types.CourseFileReqeustOkhttp)
	//c.Bind(a)
	reauestId := c.Request.PostForm["course_id"]
	if reauestId == nil {
		c.JSON(501, "c.Request.PostForm[\"id\"] 为空")
	}

	id1, _ := strconv.Atoi(reauestId[0])
	//id = "1"
	cc, err := model.FindCourseFileByCourseId(id1)

	if err != nil {
		c.JSON(501, err)
	}

	logger.Info.Println(cc)

	//type Song struct {
	//	Songname string `json:"songname"`
	//	Artistname string `json:"artistname"`
	//	Songid	string 	`json:"songid"`
	//}
	//
	//a1 := make([]Song, 1, 1)
	//a1[0] = Song{
	//	Songid: "1",
	//	Songname: "name11111111",
	//	Artistname: "art",
	//
	//}


	type Resp struct {
		CourseFileList []*model.CourseFile	`json:"courseFileList"`
	}

	ret1 := &Resp{
		CourseFileList: cc,
	}


	c.JSON(200, ret1)
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

func GetCourseTypesOk(c *gin.Context) {
	a := new(types.CourseFileReqeust)
	c.Bind(a)

	cc, err := model.FindAllCourseTypes()

	if err != nil {
		c.JSON(501, err)
	}

	type Resp struct {
		CouseTypes []*model.CourseType `json:"coursetypes"`
	}

	res := &Resp{
		CouseTypes: cc,

	}
	c.JSON(200, res)
}

func GetAllCourseIds(c *gin.Context) {
	a := new(types.CourseFileReqeust)
	c.Bind(a)

	type AA struct {
		Value int    `json:"value"`
		Label string `json:"label"`
	}

	type CC struct {
		Value    int    `json:"value"`
		Label    string `json:"label"`
		Children []AA   `json:"children"`
	}

	courseGroup, err := model.GetAllCourseGroup()

	add := make([]CC, 0, len(courseGroup))

	aMap := make(map[int]string, len(add))
	for _, v2 := range courseGroup {
		aMap[v2.TypeId] = v2.Name
	}

	cc, err := model.GetAllCourseIds()

	for _, v := range courseGroup {
		cccc := CC{
			Value: v.TypeId,
			Label: v.Name,
		}
		add = append(add, cccc)
	}

	for _, s2 := range cc {
		aacc := AA{
			Value: s2.Id,
			Label: s2.Title,
		}
		//aacc.Value = s2.Id
		//aacc.Label = s2.Title

		for k, v3 := range add {
			if v3.Value == s2.TypeId {
				add[k].Children = append(add[k].Children, aacc)
				break
			}
		}
	}

	//for k, v := range cc {
	//	v.TypeId
	//}
	if err != nil {
		c.JSON(501, err)
	}

	c.JSON(200, add)
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

func FindCourseByTypeIdOkhttp(c *gin.Context) {
	a := new(types.CourseFileReqeustOkhttp)
	c.Bind(a)


	reauestId := c.Request.PostForm["id"]
	if reauestId == nil {
		c.JSON(501, "c.Request.PostForm[\"id\"] 为空")
	}
	id, _ := strconv.Atoi(reauestId[0])
	cc, err := model.FindCourseByTypeId(id)

	if err != nil {
		c.JSON(501, err)
	}

	type Resp struct {
		Course []*model.Course `json:"courseList"`
	}

	ret := &Resp{
		Course: cc,
	}
	c.JSON(200, ret)
}

func MultiUpload(context *gin.Context) {
	type AAA struct {
		CourseId  int `json:"courseId"`
	}
	r := new(AAA)

	context.Bind(r)

	logger.Info.Println("r:", r)

	form, err := context.MultipartForm()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"msg": fmt.Sprintf("ERROR: parse form failed. %s", err),
		})
		return
	}
	// 多个文件上传，要用同一个key
	//files := form.File["files"]
	//for k, v := range form {
	//	fmt.Println("key is: ", k)
	//	fmt.Println("val is: ", v)
	//}

	files := form.File
	for _, filea := range files {
		file := filea[0]
		dst := fmt.Sprint(setting.TomlConfig.Test.FilStore.FileStorePath + file.Filename)
		// 保存文件至指定路径
		//file.
		err = context.SaveUploadedFile(file, dst)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"msg": fmt.Sprintf("ERROR: save file failed. %s", err),
			})
		}
	}
	context.JSON(http.StatusOK, gin.H{
		"msg":      "upload file succ.",
		"filepath": `./`,
	})
}
