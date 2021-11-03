package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tosone/minimp3"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"thxy/logger"
	"thxy/model"
	"thxy/redisclient"
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

func UpdatePwd(c *gin.Context) {
	type A struct {
		Pwd string `json:"pwd"`
	}

	a := new(A)
	c.Bind(a)

	err := model.UpdatePwd(a.Pwd)

	if err != nil {
		JSONError(c, err.Error(), nil)
		return
	}

	JSON(c, "ok", nil)

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

func ApkUpload(c *gin.Context) {
	//r := new(types.DownloadRequest)
	//c.Bind(r)

	// 组合成文件路径， 如：./data/37/mp3/music.mp3
	//courseId := c.Query("course_id")
	//fileType := c.Query("file_type") // fileType 为img， 或 mp3
	//fileName := c.Query("file_name")
	fileName := "app-debug.apk"

	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName)) //fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
	c.Writer.Header().Add("Content-Type", "application/octet-stream")

	storePath := setting.TomlConfig.Test.FileStore.FileStorePath
	filePath := storePath + fileName

	logger.Info.Println(" 下载文件路径=", filePath)
	c.File(filePath)

}

func FileDownloadV1(c *gin.Context) {
	response, err := http.Get("https://raw.githubusercontent.com/gin-gonic/logo/master/color.png")
	if err != nil || response.StatusCode != http.StatusOK {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	reader := response.Body
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")

	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="gopher.png"`,
	}


	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)

}

func LogDownload(c *gin.Context) {
	//r := new(types.DownloadRequest)
	//c.Bind(r)

	// 组合成文件路径， 如：./data/37/mp3/music.mp3
	fileName := c.Query("file_name")

	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName)) //fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
	//c.Writer.Header().Add("Content-Type", "application/octet-stream")
	//w := c.Writer
	c.Writer.Header().Set("Accept-Ranges", "bytes")


	//storePath := setting.TomlConfig.Test.FileStore.FileStorePath

	filePath := "./" + fileName

	logger.Info.Println(" log 下载文件路径=", filePath)

	c.File(filePath)
}


func FileDownload(c *gin.Context) {
	//r := new(types.DownloadRequest)
	//c.Bind(r)

	// 组合成文件路径， 如：./data/37/mp3/music.mp3
	courseId := c.Query("course_id")
	fileType := c.Query("file_type") // fileType 为img， 或 mp3
	fileName := c.Query("file_name")

	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName)) //fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
	//c.Writer.Header().Add("Content-Type", "application/octet-stream")
	//w := c.Writer
	c.Writer.Header().Set("Accept-Ranges", "bytes")

	storePath := setting.TomlConfig.Test.FileStore.FileStorePath

	filePath := storePath + courseId + "/" + fileType + "/" + fileName

	logger.Info.Println(" 下载文件路径=", filePath)

	//c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	//c.Writer.Header().Set("Content-Type", "application/zip")

	c.File(filePath)


	if false {
		key := filePath
		data, err := redisclient.GetFileData(key)
		if err != nil {
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				logger.Error.Println(" 下载错误：", err)
			} else {
				err = redisclient.SetFileData(key, data)
				if err != nil {
					logger.Error.Println(" redis 写错误 ：", err)
				}
				logger.Info.Println("  从硬盘 读到 ", fileType, "文件数据 ")
				c.File(filePath)

			}
			return
		} else {
			logger.Info.Println(" 从redis 读到 ", fileType, "文件数据, key=", key, ", data长度=", len(data))

			if c.Writer.Header().Get("Content-Encoding") == "" {
				c.Writer.Header().Set("Content-Length", strconv.FormatInt(int64(len(data)), 10))
			}
			//c.Data(http.StatusOK, "application/octet-stream" , data)
			c.Data(http.StatusOK, "audio", data)
		}
	}
}

func GetLogList(c *gin.Context) {
	l, _ := model.GetLogList()
	JSON(c, "ok", l)
}

// 单文件上传, 用于客户端log 上传
func FileUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		logger.Info.Println("ERROR: upload file failed. ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": fmt.Sprintf("ERROR: upload file failed. %s", err),
		})
	}

	dst := fmt.Sprintf(`./` + file.Filename)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		logger.Info.Println("ERROR: save file failed. ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": fmt.Sprintf("ERROR: save file failed. %s", err),
		})
	}

	systemVersion, _ := c.GetPostForm("system_version")
	brand, _ := c.GetPostForm("brand")
	modelVersion, _ := c.GetPostForm("model")
	appVersion, _ := c.GetPostForm("app_version")

	log := &model.Log{
		FileName:      file.Filename,
		SystemVersion: systemVersion,
		Brand:         brand,
		AppVersion:    appVersion,
		ModelVersion:  modelVersion,
	}

	model.InsertLog(log)

	logger.Info.Println(" FileUpload 文件上传 ")

	c.JSON(http.StatusOK, gin.H{
		"msg":      "file upload succ.",
		"filepath": dst,
	})
}

func GetLatestCourseFile(c *gin.Context) {

	latestCourseFile, err := model.FindCourseFileListLatest(10)
	if err != nil {
		c.JSON(501, err)
		return
	}

	type Resp struct {
		CourseFileList []*model.CourseFile `json:"courseFileList"`
	}

	resp := &Resp{
		CourseFileList: latestCourseFile,
	}

	c.JSON(200, resp)
}

func UpdateUserListenedFiles(c *gin.Context) {
	request := new(types.UserListenedFilesRequest)
	c.Bind(request)

	code := request.Code
	courseId := request.CourseId
	listenedCourseByUserCodeAndCourseId, err := model.FindUserListenedCourseByUserCodeAndCourseId(code, courseId)
	if err != nil {
		JSONError(c, "查找已听记录错误 ", err)
		return
	}

	lastListenedFileId := request.LastListenedFileId
	userListenedFiles := make([]*types.ListenedFile, 0)
	if len(listenedCourseByUserCodeAndCourseId) > 0 {
		//cMap := make(map[int]*model.UserListenedCourseFile, len(listenedCourseByUserCodeAndCourseId))
		//for _, v := range listenedCourseByUserCodeAndCourseId {
		//	if _, isExist := cMap[v.Id]; !isExist{
		//		cMap[v.Id] = v
		//	}
		//
		//}

		json.Unmarshal([]byte(listenedCourseByUserCodeAndCourseId[0].ListenedFiles), &userListenedFiles)
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

		var ulStr []byte
		ulStr, err = json.Marshal(userListenedFiles)
		if err != nil {
			JSONError(c, "Marshal 序列化错误 ", err)
			return
		}
		err = model.UpdateUserListenedCourseByUserCodeAndCourseId(string(ulStr), code, courseId, lastListenedFileId)
	} else {
		userListenedFiles = append(userListenedFiles, request.ListenedFile)

		ulfStr, err := json.Marshal(userListenedFiles)
		if err != nil {
			JSONError(c, "Marshal 序列化错误 ", err)
			return
		}

		// 那些没报上名
		ulcf := &model.UserListenedCourseFile{
			Code:                     code,
			CourseId:                 courseId,
			ListenedFiles:            string(ulfStr),
			LastListenedCourseFileId: lastListenedFileId,
		}

		tx := model.GetDB().Begin()

		err = model.InsertUserListenedCourse(tx, ulcf)

		if err != nil {
			tx.Rollback()
			JSONError(c, "创建 user_listened_course 表错误 ", err)
			return
		}

		tx.Commit()
	}

	c.JSON(200, nil)
	return
}

func GetConfig(c *gin.Context) {
	config, err := model.FindConfig()
	if err != nil {
		c.JSON(501, err)
	}

	type Resp struct {
		Config *model.Config `json:"config"`
	}

	c.JSON(200, config)
}

func FindCourseFileByCourseId(c *gin.Context) {
	a := new(types.CourseTypeRequest)
	c.Bind(a)

	cc, err := model.FindCourseFileByCourseId(a.Id)

	if err != nil {
		c.JSON(501, err)
		return
	}

	c.JSON(200, cc)
}

func FindCourseFileByCourseIdOkhttpV1(c *gin.Context) {
	a := new(types.CourseFileRequestOkhttp)
	c.Bind(a)
	courseId := c.Request.PostForm["course_id"]
	if courseId == nil {
		c.JSON(501, "c.Request.PostForm[\"id\"] 为空")
		return
	}

	courseIdInt, _ := strconv.Atoi(courseId[0])

	courseFiles, err := model.FindCourseFileByCourseId(courseIdInt)
	if err != nil {
		logger.Warning.Println("FindCourseFileByCourseId err =  ", err)
		c.JSON(501, err)
		return
	}

	type Resp struct {
		CourseFileList []*model.CourseFile `json:"courseFileList"`
	}

	resp := &Resp{
		CourseFileList: courseFiles,
	}

	c.JSON(200, resp)
	return
}

func FindUserListenedFilesByCodeAndCourseId(c *gin.Context) {
	a := new(types.CourseFileRequestOkhttp)
	c.Bind(a)

	userCode := c.Request.PostForm["user_code"]
	courseId := c.Request.PostForm["course_id"]
	if courseId == nil {
		c.JSON(501, "c.Request.PostForm[\"id\"] 为空")
		return
	}

	courseIdInt, _ := strconv.Atoi(courseId[0])
	ulc, err := model.FindUserListenedCourseByUserCodeAndCourseId(userCode[0], courseIdInt)
	if err != nil {
		logger.Warning.Println("  FindUserListenedCourseByUserCodeAndCourseId err=", err)
		c.JSON(501, err)
		return
	}

	c.JSON(200, ulc)
	return
}

func FindCourseFileByCourseIdOk(c *gin.Context) {
	a := new(types.CourseFileRequestOkhttp)
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
	a := new(types.CourseTypeRequest)
	c.Bind(a)

	cc, err := model.FindCourseFileById(a.Id)

	if err != nil {
		c.JSON(501, err)
	}

	c.JSON(200, cc)
}

func GetCourseTypes(c *gin.Context) {
	a := new(types.CourseTypeRequest)
	c.Bind(a)

	courseTypes, err := model.FindAllCourseTypes()

	if err != nil {
		c.JSON(501, err)
	}

	c.JSON(200, courseTypes)
}

func UpdateCourseType(c *gin.Context) {
	r := new(types.CourseTypeRequest)
	c.Bind(r)

	err := model.UpdateCourseTypeById(r.Name, r.Id)
	if err != nil {
		logger.Error.Println(" 更新课程类型失败 , err=", err)
		JSONError(c, "UpdateCourseType err= "+err.Error(), nil)
		return
	}

	JSON(c, "ok", nil)
	return
}

func AddCourseType(c *gin.Context) {
	r := new(types.CourseTypeRequest)
	c.Bind(r)

	courseType := &model.CourseType{
		Name: r.Name,
	}
	err := model.AddCourseType(courseType)
	if err != nil {
		logger.Error.Println(" 添加课程类型失败 , err=", err)
		JSONError(c, "AddCourseType err="+err.Error(), err)
		return
	}

	JSON(c, "ok", nil)
}

func GetCourseTypesOk(c *gin.Context) {
	a := new(types.CourseTypeRequest)
	c.Bind(a)

	courseTypes, err := model.FindAllCourseTypes()

	if err != nil {
		c.JSON(501, err)
	}

	type Resp struct {
		CouseTypes []*model.CourseType `json:"coursetypes"`
	}

	res := &Resp{
		CouseTypes: courseTypes,
	}
	c.JSON(200, res)
}

func AdminGetAllCourseType(c *gin.Context) {
	courseGroup, err := model.FindAllCourseTypes()
	if err != nil {
		logger.Warning.Println(" admin 选择框获取课程类型错误, err = ", err)
		c.JSON(501, err)
		return
	}

	optionItems := make([]types.OptionItem, 0, len(courseGroup))
	for _, v := range courseGroup {
		optionItem := types.OptionItem{
			Value: v.Id,
			Label: v.Name,
		}
		optionItems = append(optionItems, optionItem)
	}

	c.JSON(200, optionItems)

	return
}
func AdminGetAllCourseIds(c *gin.Context) {
	request := new(types.CourseTypeRequest)
	c.Bind(request)

	courseGroup, err := model.GetAllCourseGroup()
	if err != nil {
		logger.Warning.Println(" admin 选择框获取课程类型错误, err = ", err)
		c.JSON(501, err)
		return
	}

	optionItems := make([]types.OptionItem, 0, len(courseGroup))
	optionItemsMap := make(map[int]string, len(optionItems))
	for _, courseGroupItem := range courseGroup {
		optionItemsMap[courseGroupItem.TypeId] = courseGroupItem.Name
	}
	allCourseIds, err := model.GetAllCourseIds()

	for _, v := range courseGroup {
		optionItem := types.OptionItem{
			Value: v.TypeId,
			Label: v.Name,
		}
		optionItems = append(optionItems, optionItem)
	}

	for _, courseItem := range allCourseIds {
		children := types.Children{
			Value: courseItem.Id,
			Label: courseItem.Title,
		}

		for index, optionItem := range optionItems {
			if optionItem.Value == courseItem.TypeId {
				optionItems[index].Children = append(optionItems[index].Children, children)
				break
			}
		}
	}
	if err != nil {
		logger.Warning.Println(" admin 选择框获取课程类型错误, err = ", err)
		c.JSON(501, err)
		return
	}

	JSON(c, "ok", optionItems)
	//c.JSON(200, optionItems)
	return
}

func FindCourseByTypeId(c *gin.Context) {
	request := new(types.CourseTypeRequest)
	c.Bind(request)

	courseByTypeId, err := model.FindCourseByTypeId(request.Id)
	if err != nil {
		logger.Warning.Println(" admin FindCourseByTypeId err = ", err)
		c.JSON(501, err)
		return
	}

	c.JSON(200, courseByTypeId)
}

func GetCourseById(c *gin.Context) {
	type RequestA struct {
		Id int `json:"id"`
	}
	a := new(RequestA)
	c.Bind(a)

	//if reauestId == nil {
	//	c.JSON(501, "c.Request.PostForm[\"id\"] 为空")
	//	return
	//}
	//id, _ := strconv.Atoi(reauestId[0])

	id := a.Id
	course, err := model.FindCourseById(id)

	if err != nil {
		c.JSON(501, err)
		JSONError(c, "FindCourseById err"+err.Error(), err)
		return
	}

	JSON(c, "ok", course)

}
func FindCourseByTypeIdOkhttp(c *gin.Context) {
	a := new(types.CourseFileRequestOkhttp)
	c.Bind(a)

	reauestId := c.Request.PostForm["id"]
	if reauestId == nil {
		c.JSON(501, "c.Request.PostForm[\"id\"] 为空")
		return
	}
	id, _ := strconv.Atoi(reauestId[0])
	cc, err := model.FindCourseByTypeId(id)

	if err != nil {
		c.JSON(501, err)
		return
	}

	type Resp struct {
		Course []*model.Course `json:"courseList"`
	}

	ret := &Resp{
		Course: cc,
	}
	c.JSON(200, ret)

}

func AddCourse(c *gin.Context) {
	r := new(types.CourseRequest)
	c.Bind(r)
	cs := &model.Course{
		Title:  r.Title,
		TypeId: r.TypeId,
	}
	err := model.InsertCourse(cs)
	if err != nil {
		JSONError(c, "AddCourse err= "+err.Error(), nil)
		return
	}

	JSON(c, "ok", nil)

	return
}

func UpdateCourse(c *gin.Context) {
	r := new(types.CourseRequest)
	c.Bind(r)
	course := &model.Course{
		Title:        r.Title,
		StorePath:    r.StorePath,
		ImgFileName:  r.ImgSrc,
		Introduction: r.Introduction,
	}
	err := model.UpdateCourse(course.Title, course.Introduction, r.Id)
	if err != nil {
		JSONError(c, "AddCourse err= "+err.Error(), nil)
		return
	}

	JSON(c, "ok", nil)
	return
}

func DeleteCourse(c *gin.Context) {
	r := new(types.CourseRequest)
	c.Bind(r)

	err := model.DeleteCourse(r.Id)

	if err != nil {
		JSONError(c, "DeleteCourse err= "+err.Error(), nil)
		return
	}

	JSON(c, "ok", nil)
	return
}

func DeleteCourseType(c *gin.Context) {
	r := new(types.CourseRequest)
	c.Bind(r)

	err := model.DeleteCourseType(r.Id)
	if err != nil {
		JSONError(c, "DeleteCourseType err= "+err.Error(), nil)
		return
	}
	JSON(c, "ok", nil)
	return
}

func MultiUpload(c *gin.Context) {
	r := new(types.MultiUploadRequest)
	c.Bind(r)

	logger.Info.Println("r:", r)
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
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
	courseIdStr := c.Request.PostForm["courseId"][0]
	courseId, _ := strconv.Atoi(courseIdStr)

	durationStr := c.Request.PostForm["duration"][0]
	durationInt, _ := strconv.Atoi(durationStr)
	courseFiles := make([]*model.CourseFile, 0)

	logger.Info.Println("时长： ", durationInt)

	for _, fileArr := range files {
		file := fileArr[0]
		storePath := setting.TomlConfig.Test.FileStore.FileStorePath
		dst := storePath + strconv.Itoa(courseId) + "/mp3" + "/" + file.Filename
		logger.Debug.Println("dst: ", dst)
		err = c.SaveUploadedFile(file, dst)
		if err != nil {
			logger.Error.Println(" 上传mp3文件出错： err= ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": fmt.Sprintf("ERROR: save file failed. %s", err),
			})
			return
		}

		//f, _  := file.Open()
		//
		//ss := make([]byte, 1024*1024)
		////a, _ := f.Read(ss)
		//f.Read(ss)
		//offset := 0
		//var s1 []byte
		//for {
		//	len1, _ := f.ReadAt(ss, int64(offset))
		//	offset = offset + len1
		//	if len1 == 0 {
		//		break
		//	}
		//	s1 = append(s1, ss...)
		//	//fmt.Println( string(ss))
		//}
		//f.Close()
		//durationMp3,  _ := GetMP3PlayDuration(s1)
		//logger.Info.Println("filename=",  file.Filename, ",  duration=", durationMp3)

		regExp := regexp.MustCompile("[0-9]+")

		titleArr := strings.Split(file.Filename, ".")
		title := titleArr[0]

		numberStr := regExp.FindAllString(title, -1)
		number, err := strconv.Atoi(numberStr[0])

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": fmt.Sprintf("ERROR: Atoi failed. %s", err),
			})
		}

		courseFile := &model.CourseFile{
			CourseId:    courseId,
			Number:      number,
			Duration:    utils.GetTimeStrFromSecond(durationInt),
			Mp3FileName: file.Filename,
		}
		courseFiles = append(courseFiles, courseFile)
	}

	db := model.GetDB()
	tx := db.Begin()
	// 文件 上传完成后，再保存到数据库
	for _, v := range courseFiles {
		err = model.InsertCourseFile(tx, v)
		if err != nil {
			logger.Error.Println("InsertCourseFile err=", err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": fmt.Sprintf("ERROR: InsertCourseFileinsert failed. %s", err),
			})
			return
		}
	}

	tx.Commit()

	fileCount, err := model.FindCourseFileCountByCourseId(courseId)
	if err != nil {
		JSONError(c,  " 查询出错："+ err.Error(), nil)
		return
	}
	err = model.UpdateCourseFileCount(courseId, fileCount)
	if err != nil {
		JSONError(c,  " course表更新文件数目出错："+ err.Error(), nil)
		return
	}else{
		c.JSON(http.StatusOK, gin.H{
			"msg":      "upload file success",
			"filepath": setting.TomlConfig.Test.FileStore.FileStorePath,
		})
	}

	return
}

func CoursePictureUpload(c *gin.Context) {
	r := new(types.MultiUploadRequest)
	c.Bind(r)

	logger.Info.Println("r:", r)
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
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
	courseTitle := c.Request.PostForm["course_title"][0]
	//courseId, _ := strconv.Atoi(courseIdStr)

	durationStr := c.Request.PostForm["type_id"][0]
	durationInt, _ := strconv.Atoi(durationStr)
	courses := make([]*model.Course, 0)
	for _, fileArr := range files {
		file := fileArr[0]

		//regExp := regexp.MustCompile("[0-9]+")
		//
		//titleArr := strings.Split(file.Filename, ".")
		//title := titleArr[0]

		//numberStr := regExp.FindAllString(title, -1)
		//number, err := strconv.Atoi(numberStr[0])

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": fmt.Sprintf("ERROR: Atoi failed. %s", err),
			})
		}

		course := &model.Course{
			TypeId:      durationInt,
			Title:       courseTitle,
			ImgFileName: file.Filename,
		}
		courses = append(courses, course)
	}

	db := model.GetDB()
	tx := db.Begin()

	var imgDir string
	var mp3Dir string
	for _, v := range courses {
		err = model.InsertCourseT(tx, v)
		if err != nil {
			logger.Error.Println("InsertCourseFile err=", err)
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": fmt.Sprintf("ERROR: InsertCourseFileinsert failed. %s", err),
			})
			return
		}

		imgDir = utils.GetDir(v.Id, "img")
		mp3Dir = utils.GetDir(v.Id, "mp3")

		logger.Info.Println(v.Title, " 创建img目录：", imgDir)
		logger.Info.Println(v.Title, " 创建mp3目录：", mp3Dir)

		err = os.MkdirAll(imgDir, os.ModePerm)
		if err != nil {
			tx.Rollback()
		}
		err = os.MkdirAll(mp3Dir, os.ModePerm)
	}

	tx.Commit()

	for _, filea := range files {
		file := filea[0]
		dst := imgDir + file.Filename
		logger.Debug.Println("dst: ", dst)
		err = c.SaveUploadedFile(file, dst)
		if err != nil {
			JSON(c, " SaveUploadedFile err= "+err.Error(), nil)
			logger.Error.Println(" 创建失败： SaveUploadedFile err=  ", err)
			return
		}
	}

	JSON(c, "ok", nil)
	return
}
