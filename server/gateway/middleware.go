package gateway

import (
	"main/server/db"
	"main/server/response"
	"main/server/services/token"
	"main/server/utils"

	"fmt"

	"github.com/gin-gonic/gin"
)

func AdminAuthorization(ctx *gin.Context) {

	fmt.Println("inside middleware")
	tokenString := ctx.Request.Header.Get(utils.Authorization)
	var exists bool
	//first check if the session is valid or not
	query := "SELECT EXISTS(SELECT 1 FROM sessions WHERE token=?)"
	err := db.QueryExecutor(query, &exists, tokenString)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_INTERNAL_SERVER_ERROR, utils.FAILURE, nil, ctx)
		ctx.Abort()
		return
	}
	if !exists {
		response.ShowResponse("Invalid session", utils.HTTP_FORBIDDEN, utils.FAILURE, nil, ctx)
		ctx.Abort()
		return
	}

	claims, err := token.DecodeToken(tokenString)
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		ctx.Abort()
		return
	}
	err = claims.Valid()
	if err != nil {
		response.ShowResponse(err.Error(), utils.HTTP_UNAUTHORIZED, utils.FAILURE, nil, ctx)
		ctx.Abort()
		return
	}
	if claims.Role == "admin" || claims.Role == "player" {
		ctx.Set("role", claims.Role)
		ctx.Set(utils.PLAYERID, claims.Id)
		ctx.Next()
	} else {
		response.ShowResponse(utils.ACCESS_DENIED, utils.HTTP_FORBIDDEN, utils.FAILURE, nil, ctx)
		ctx.Abort()
		return
	}
	//set the token details into context for further processing in handler function
	ctx.Next()

}

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(int(utils.HTTP_NO_CONTENT))
			return
		}

		ctx.Next()
	}
}
