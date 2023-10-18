package response

import (
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type Success struct {
	Status  string      `json:"status"`
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Response(context *gin.Context, statusCode int, data interface{}) {
	context.JSON(statusCode, data)
}

func ShowResponse(message string, statusCode int64, status string, data interface{}, context *gin.Context) {
	context.Writer.Header().Set("Content-Type", "application/json")
	context.Writer.WriteHeader(int(statusCode))
	response := Success{
		Status:  status,
		Code:    statusCode,
		Message: message,
		Data:    data,
	}

	Response(context, int(statusCode), response)
}

type SocketResp struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SocketResponse(message string, statusCode int64, status string, data interface{}, eventName string, s socketio.Conn) {
	socketResponse := Success{
		Status:  status,
		Code:    statusCode,
		Message: message,
		Data:    data,
	}

	s.Emit(eventName, socketResponse)

}
