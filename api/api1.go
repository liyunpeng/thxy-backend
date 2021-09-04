package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tosone/minimp3"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"thxy/logger"
	"thxy/model"
	"thxy/setting"
	"thxy/types"
	"thxy/utils"
)

func Login(c *gin.Context) {

	a := types.Response1{
		AccessToken: "add",
	}
	c.JSON(200, a)

}

func GetMP3PlayDuration(mp3Data []byte) (seconds int, err error) {
	dec, _, err := minimp3.DecodeFull(mp3Data)
	if err != nil {
		return 0, err
	}
	// 音乐时长 = (文件大小(byte) - 128(ID3信息)) * 8(to bit) / (码率(kbps b:bit) * 1000)(kilo bit to bit)
	seconds = (len(mp3Data) - 128) * 8 / (dec.Kbps * 1000)
	return seconds, nil
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

	lfcfid := request.ListenedFile.CourseFileId
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
			if v.CourseFileId == request.ListenedFile.CourseFileId {
				userListenedFiles[k].ListenedPercent = request.ListenedFile.ListenedPercent
				userListenedFiles[k].Position = request.ListenedFile.Position
				found = true
				break
			}
		}

		if found == false {
			userListenedFiles = append(userListenedFiles, request.ListenedFile)
		}

		ulStr, _ := json.Marshal(userListenedFiles)
		err = model.UpdateUserListenedCourseByUserCodeAndCourseId(string(ulStr), code, courseId, lfcfid)
	} else {
		userListenedFiles = append(userListenedFiles, request.ListenedFile)

		ulfStr, _ := json.Marshal(userListenedFiles)

		ulcf := &model.UserListenedCourseFile{
			Code:                     code,
			CourseId:                 courseId,
			ListenedFiles:            string(ulfStr),
			LastListenedCourseFileId: lfcfid,
		}

		tx := model.GetDB().Begin()

		err = model.InsertUserListenedCourse(tx, ulcf)

		if err != nil {
			tx.Rollback()
			c.JSON(500, nil)
			return
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

	courseFiles, err := model.FindCourseFileByCourseId(courseIdInt)
	if err != nil {
		logger.Warning.Println("FindCourseFileByCourseId err =  ", err)
		c.JSON(501, err)
		return
	}

	ulc, err := model.FindUserListenedCourseByUserCodeAndCourseId(userCode[0], courseIdInt)
	if err != nil {
		logger.Warning.Println("  FindUserListenedCourseByUserCodeAndCourseId err=", err)
		c.JSON(501, err)
		return
	}

	type Resp struct {
		CourseFileList           []*model.CourseFile `json:"courseFileList"`
		LastListenedCourseFileId int                 `json:"last_listened_course_file_id"`
	}

	if len(ulc) <= 0 {
		ret := &Resp{
			CourseFileList:           courseFiles,
			LastListenedCourseFileId: -2,
		}
		c.JSON(200, ret)
		return
	}

	LastListenedCourseFileId := ulc[0].LastListenedCourseFileId

	userListenedFiles := make([]*types.ListenedFile, 0)
	err = json.Unmarshal([]byte(ulc[0].ListenedFiles), &userListenedFiles)
	if err != nil {
		logger.Warning.Println(" ListenedFiles  Unmarshal err =  ", err)
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
		if _, isExist := courseFileMap[v.CourseFileId]; isExist {
			courseFileMap[v.CourseFileId].ListenedPercent = v.ListenedPercent
			courseFileMap[v.CourseFileId].ListenedPosition = v.Position
		}
	}

	courseFileList := make([]*model.CourseFile, 0)
	for _, v := range courseFileMap {
		courseFileList = append(courseFileList, v)
	}

	sort.Slice(courseFileList, func(i, j int) bool {
		if courseFileList[i].Number < courseFileList[j].Number {
			return true
		} else {
			return false
		}
	})

	logger.Info.Println(courseFiles)

	ret := &Resp{
		CourseFileList:           courseFileList,
		LastListenedCourseFileId: LastListenedCourseFileId,
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

	durationStr := context.Request.PostForm["duration"][0]
	durationInt, _ := strconv.Atoi(durationStr)

	cfs := make([]*model.CourseFile, 0)
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
			CourseId:    courseId,
			Number:      number,
			Duration:    utils.GetTimeStrFromSecond(durationInt),
			Mp3FileName: file.Filename,
		}
		cfs = append(cfs, courseFile)
	}

	db := model.GetDB()
	tx := db.Begin()
	for _, v := range cfs {
		err = model.InsertCourseFile(tx, v)
		if err != nil {
			logger.Error.Println("InsertCourseFile err=", err)
			tx.Rollback()
			context.JSON(http.StatusInternalServerError, gin.H{
				"msg": fmt.Sprintf("ERROR: InsertCourseFileinsert failed. %s", err),
			})
			return
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
	return
}
