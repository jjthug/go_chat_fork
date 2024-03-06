package ws

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

type CreateRoomRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	c.JSON(http.StatusOK, req)
}

var upgrader = websocket.Upgrader{
	//HandshakeTimeout: 0,
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	//WriteBufferPool:  nil,
	//Subprotocols:     nil,
	//Error:            nil,
	CheckOrigin: func(r *http.Request) bool {
		//origin := r.Header.Get("Origin")
		//return origin == "http://localhost:3000"
		return true
	},
	//EnableCompression: false,
}

func (h *Handler) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	roomId := c.Param("roomId")
	clientId := c.Query("userId")
	username := c.Query("username")

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientId,
		RoomID:   roomId,
		Username: username,
	}

	m := &Message{
		Content:  "New user joined",
		RoomID:   roomId,
		Username: username,
	}

	//Register a new client through register channel
	h.hub.Register <- cl

	// TODO handle from redis queue to worker
	//Broadcast that message
	//h.hub.Broadcast <- m
	err = h.hub.Distributor.DistributeTaskSendMessage(c, &PayloadSendMessage{Message: *m})
	if err != nil {
		fmt.Printf("error sending message to queue: %v", err)
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	go cl.writeMessage()
	cl.readMessage(h.hub)
}

type RoomsResp struct {
	ID    string `json:"id"`
	Names string `json:"name"`
}

func (h *Handler) GetRooms(c *gin.Context) {
	rooms := make([]RoomsResp, 0)
	for _, room := range h.hub.Rooms {
		rooms = append(rooms, RoomsResp{
			ID:    room.ID,
			Names: room.Name,
		})
	}
	c.JSON(http.StatusOK, rooms)
}

type ClientsResp struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) GetClients(c *gin.Context) {
	var clients []ClientsResp
	roomId := c.Param("roomId")

	if _, ok := h.hub.Rooms[roomId]; !ok {
		clients = make([]ClientsResp, 0)
		c.JSON(http.StatusOK, clients)
	}

	for _, client := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, ClientsResp{
			ID:       client.ID,
			Username: client.Username,
		})
	}

	c.JSON(http.StatusOK, clients)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
