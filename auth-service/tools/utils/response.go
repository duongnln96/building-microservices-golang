package utils

type ResponseData struct {
	Status string      `json:"status"`
	Error  string      `json:"error"`
	Data   interface{} `json:"data"`
}

func BuildResponse(status string, data interface{}) ResponseData {
	return ResponseData{
		Status: status,
		Error:  "null",
		Data:   data,
	}
}

func BuildErrorResponse(status string, err string) ResponseData {
	return ResponseData{
		Status: status,
		Error:  err,
		Data:   "null",
	}
}
