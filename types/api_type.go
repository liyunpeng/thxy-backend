package types

type Response1 struct {
	AccessToken string `json:"access_token"`
}

type CourseTypeRequest struct {
	Id   int    `json:"id"` // id 为course type id
	Name string `json:"name"`
}

type CourseRequest struct {
	Id        int    `json:"id"`
	TypeId    int    `json:"type_id"`
	Title     string `json:"title"`
	StorePath string `json:"store_path"`
	ImgSrc    string `json:"img_src"`
}

type CourseFileRequestOkhttp struct {
	Id string `json:"id"`
}

type DownloadRequest struct {
	Id       int    `json:"id"`
	FileName string `json:"file_name"`
}

type ListenedFile struct {
	CourseFileId    int `json:"cfi"` // 为了节约数据库存储空间
	ListenedPercent int `json:"pc"`
	Position        int `json:"pos"`
}

type UserListenedFilesRequest struct {
	Code         string        `json:"code"`
	CourseId     int           `json:"course_id"`
	ListenedFile *ListenedFile `json:"listened_file"`
}

type CommonRes struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

type CommonList struct {
	Total                int         `json:"total"`
	ServiceSignedCount   int         `json:"service_signed_count"`
	ServiceUnsignedCount int         `json:"service_unsigned_count"`
	ServiceNotShowCount  int         `json:"service_not_show_count"`
	ServiceShowCount     int         `json:"service_show_count"`
	List                 interface{} `json:"list"`
}

type MultiUploadRequest struct {
	CourseId int `json:"courseId"`
}
