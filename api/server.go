package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db "server/db/sqlc"
	"server/token"
	"server/util"
	"server/ws"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store, wsHandler *ws.Handler) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker([]byte(config.TokenSymmetric))
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	router := gin.Default()
	router.POST("/user", server.CreateNewUser)

	router.POST("/ws/createRoom", wsHandler.CreateRoom)
	router.GET("/ws/joinRoom/:roomId", wsHandler.JoinRoom)
	router.GET("/ws/rooms", wsHandler.JoinRoom)
	router.GET("/ws/getRooms", wsHandler.GetRooms)
	router.GET("/ws/getClients/:roomId", wsHandler.GetClients)

	//router.GET("/get_user/:id", server.GetUser)
	//router.POST("/user/login", server.LoginUser)

	//authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	server.router = router
	return server, nil
}

func (server *Server) RunHTTPServer(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
