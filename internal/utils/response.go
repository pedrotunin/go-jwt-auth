package utils

func GetSuccessResponse() {}

func GetErrorResponse(err error) map[string]string {
	return map[string]string{
		"error": err.Error(),
	}
}
