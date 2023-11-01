package gateway

import (
	"main/server/db"
	"main/server/response"
	"main/server/services/token"
	"main/server/utils"
	"strings"

	socketio "github.com/googollee/go-socket.io"
)

func SocketAuthMiddleware(next func(s socketio.Conn, reqData map[string]interface{})) func(socketio.Conn, map[string]interface{}) {
	return func(s socketio.Conn, reqData map[string]interface{}) {
		// fmt.Println("query is ", s.URL().RawQuery)
		tokenString := strings.Split(strings.Split(s.URL().RawQuery, "&")[0], "=")[1]
		// tokenString := s.RemoteHeader().Get(utils.Authorization)
		// fmt.Println("Token is", tokenString)

		var exists bool
		//first check if the session is valid or not
		query := "SELECT EXISTS(SELECT 1 FROM sessions WHERE token=?)"
		err := db.QueryExecutor(query, &exists, tokenString)
		if err != nil {
			response.SocketResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, "ack", s)

			return
		}
		if !exists {
			response.SocketResponse("Invalid session", utils.HTTP_FORBIDDEN, utils.FAILURE, nil, "ack", s)
			return
		}

		claims, err := token.DecodeToken(tokenString)
		if err != nil {
			response.SocketResponse(err.Error(), utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, "ack", s)
			return
		}
		err = claims.Valid()
		if err != nil {
			response.SocketResponse(err.Error(), utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, "ack", s)

			return
		}
		if claims.Role == "admin" || claims.Role == "player" {
			s.SetContext(claims.Id)
		} else {
			response.SocketResponse(utils.ACCESS_DENIED, utils.HTTP_FORBIDDEN, utils.FAILURE, nil, "ack", s)
			return
		}
		// If authentication is successful, call the next event handler
		next(s, reqData)
	}
}
