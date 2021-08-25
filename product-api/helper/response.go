package helper

type Response struct {
	Status string      `json:"status"`
	Errors error       `json:"error"`
	Data   interface{} `json:"data"`
}

func BuildResponse(status string, data interface{}) Response {
	return Response{
		Status: status,
		Errors: nil,
		Data:   data,
	}
}

func BuildErrorResponse(status string, err error, data interface{}) Response {
	return Response{
		Status: status,
		Errors: err,
		Data:   data,
	}
}
