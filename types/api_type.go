package types

type Response1 struct {
	AccessToken string `json:"access_token"`
}


type CourseFileReqeust struct {
	Id int `json:"id"`
}

type DownloadReqeust struct {
	Id int `json:"id"`
	FileName string `json:"file_name"`
}