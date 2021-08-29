package types

type Response1 struct {
	AccessToken string `json:"access_token"`
}

type CourseFileReqeust struct {
	Id int `json:"id"`
}

type CourseFileReqeustOkhttp struct {
	Id string `json:"id"`
}

type DownloadReqeust struct {
	Id       int    `json:"id"`
	FileName string `json:"file_name"`
}

type ListenedFile struct {
	CourseFileId    int     `json:"cfi"`  // 为了节约数据库存储空间
	ListenedPercent int `json:"pc"`
}

type UserListenedFilesRequest struct {
	Code         string       `json:"code"`
	CourseId     int          `json:"course_id"`
	ListenedFile *ListenedFile `json:"listened_file"`
}
