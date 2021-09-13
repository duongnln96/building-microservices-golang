package helper

type Response struct {
	Status string      `json:"status"`
	Errors string      `json:"error"`
	Data   interface{} `json:"data"`
}

func BuildResponse(status string, data interface{}) Response {
	return Response{
		Status: status,
		Errors: "",
		Data:   data,
	}
}

func BuildErrorResponse(status string, err string, data interface{}) Response {
	return Response{
		Status: status,
		Errors: err,
		Data:   data,
	}
}
