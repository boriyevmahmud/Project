package model



type JwtRequestModel struct{
	Token string `string:"token"`
}

type ResponseError struct{
	Error interface{} `json:"error"`
}

// ServerError ...
type ServerError struct{
	Status string `json:"status"`
	Message string`json:"message"`
}