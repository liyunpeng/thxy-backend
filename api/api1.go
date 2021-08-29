package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
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

func GetLatest(c *gin.Context) {

	cc, err := model.FindCourseFileListLatest(10)

	if err != nil {
		c.JSON(501, err)
		return
	}

	type Resp struct {
		CourseFileList []*model.CourseFile `json:"courseFileList"`
	}

	ret1 := &Resp{
		CourseFileList: cc,
	}

	c.JSON(200, ret1)
}

func UpdateUserListenedFiles(c *gin.Context) {
	request := new(types.UserListenedFilesRequest)
	c.Bind(request)

	code := request.Code
	courseId := request.CourseId
	ret, err := model.FindUserListenedCourseByUserCodeAndCourseId(code, courseId)
	if err != nil {
		c.JSON(501, err)
		return
	}

	userListenedFiles := make([]*types.ListenedFile, 0)
	if len(ret) > 0 {


		//cMap := make(map[int]*model.UserListenedCourseFile, len(ret))
		//for _, v := range ret {
		//	if _, isExist := cMap[v.Id]; !isExist{
		//		cMap[v.Id] = v
		//	}
		//
		//}

		json.Unmarshal([]byte(ret[0].ListenedFiles), &userListenedFiles)
		found := false
		for k, v := range userListenedFiles {

			//if  _, isExxist := cMap[v.CourseFileId]; isExxist {
			//	cMap[v.CourseFileId]
			//
			//}

			if v.CourseFileId == request.ListenedFile.CourseFileId {
				userListenedFiles[k].ListenedPercent = request.ListenedFile.ListenedPercent
				found = true
				break
			}
		}

		if found == false {
			userListenedFiles = append(userListenedFiles, request.ListenedFile)
		}

		ulStr, _ := json.Marshal(userListenedFiles)
		err = model.UpdateUserListenedCourseByUserCodeAndCourseId(string(ulStr), code, courseId)
	} else {
		userListenedFiles = append(userListenedFiles, request.ListenedFile)

		ulStr, _ := json.Marshal(userListenedFiles)

		ulc := &model.UserListenedCourseFile{
			Code:          code,
			CourseId:      courseId,
			ListenedFiles: string(ulStr),
		}

		tx := model.GetDB().Begin()

		err = model.InsertUserListenedCourse(tx, ulc)

		if err != nil {
			tx.Rollback()
			c.JSON(500, nil)
		}

		tx.Commit()
	}

	c.JSON(200, nil)
	return
}

func GetConfig(c *gin.Context) {

	cc, err := model.FindConfig()

	if err != nil {
		c.JSON(501, err)
	}

	type Resp struct {
		Config *model.Config `json:"config"`
	}

	//ret1 := &Resp{
	//	Config: cc,
	//}

	c.JSON(200, cc)
}

func FindCourseFileByCourseId(c *gin.Context) {
	a := new(types.CourseFileReqeust)
	c.Bind(a)

	cc, err := model.FindCourseFileByCourseId(a.Id)

	if err != nil {
		c.JSON(501, err)
		return
	}

	c.JSON(200, cc)
}

func FindCourseFileByCourseIdOk(c *gin.Context) {
	a := new(types.CourseFileReqeustOkhttp)
	c.Bind(a)

	courseId := c.Request.PostForm["course_id"]
	userCode := c.Request.PostForm["user_code"]
	if courseId == nil {
		c.JSON(501, "c.Request.PostForm[\"id\"] 为空")
	}

	courseIdInt, _ := strconv.Atoi(courseId[0])

	ulc, err := model.FindUserListenedCourseByUserCodeAndCourseId(userCode[0], courseIdInt)
	if err != nil {
		c.JSON(501, err)
		return
	}

	userListenedFiles := make([]*types.ListenedFile, 0)

	err = json.Unmarshal([]byte(ulc[0].ListenedFiles), &userListenedFiles)
	if err != nil {
		c.JSON(501, err)
		return
	}

	courseFiles, err := model.FindCourseFileByCourseId(courseIdInt)
	if err != nil {
		c.JSON(501, err)
		return
	}

	courseFileMap := make(map[int]*model.CourseFile, len(courseFiles))
	for _, v := range courseFiles {
		if _, isExist := courseFileMap[v.Id]; !isExist {
			courseFileMap[v.Id] = v
		}
	}

	for _, v := range userListenedFiles {
		if _, isExist :=   courseFileMap[v.CourseFileId]; isExist {
			courseFileMap[v.CourseFileId].ListenedPercent = v.ListenedPercent
		}
	}

	cc := make([]*model.CourseFile, 0)
	for _, v :=  range courseFileMap{
		cc = append(cc, v)
	}

	sort.Slice(cc, func(i, j int) bool {
		if cc[i].Number < cc[j].Number{
			return true
		}else{
			return false
		}
	})

	logger.Info.Println(courseFiles)

	type Resp struct {
		CourseFileList []*model.CourseFile `json:"courseFileList"`
	}

	ret := &Resp{
		CourseFileList: cc,
	}

	c.JSON(200, ret)

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
	type Request struct {
		CourseId int `json:"courseId"`
	}
	r := new(Request)

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
	courseIdStr := context.Request.PostForm["courseId"][0]
	courseId, _ := strconv.Atoi(courseIdStr)

	tx := model.GetDB()
	tx.Begin()
	for _, filea := range files {
		file := filea[0]
		regExp := regexp.MustCompile("[0-9]+")

		titleArr := strings.Split(file.Filename, ".")
		title := titleArr[0]

		numberStr := regExp.FindAllString(title, -1)
		number, err := strconv.Atoi(numberStr[0])

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"msg": fmt.Sprintf("ERROR: Atoi failed. %s", err),
			})
		}
		courseFile := &model.CourseFile{
			CourseId: courseId,
			//Title:       title,
			Number: number,
			//Mp3Url:      setting.TomlConfig.Test.Server.FileDownload,
			Mp3FileName: file.Filename,
		}

		err = model.InsertCourseFile(tx, courseFile)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"msg": fmt.Sprintf("ERROR: InsertCourseFileinsert failed. %s", err),
			})
			tx.Rollback()
		}
	}
	tx.Commit()

	for _, filea := range files {
		file := filea[0]
		dst := fmt.Sprint(setting.TomlConfig.Test.FilStore.FileStorePath + file.Filename)
		logger.Debug.Println("dst: ", dst)
		err = context.SaveUploadedFile(file, dst)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"msg": fmt.Sprintf("ERROR: save file failed. %s", err),
			})
		}
	}

	context.JSON(http.StatusOK, gin.H{
		"msg":      "upload file succ.",
		"filepath": setting.TomlConfig.Test.FilStore.FileStorePath,
	})
}
