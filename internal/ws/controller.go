package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Controller struct {
	hub *Hub
}

func NewController(h *Hub) *Controller{
	return &Controller{
		hub: h,
	}
}

type CreateRoomReq struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

func (c *Controller) CreateRoom(ctx *gin.Context){
	var req CreateRoomReq
	if err := ctx.ShouldBindJSON(&req); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.hub.Rooms[req.ID] = &Room{
		ID: req.ID,
		Name: req.Name,
		Clients: make(map[string]*Client),
	}

	ctx.JSON(http.StatusOK, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool{
		return true
	},
}

func (c *Controller) JoinRoom(ctx *gin.Context){
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request,nil)
	if err != nil{
		ctx.JSON(http.StatusBadRequest,gin.H{"error": err.Error()})
		return
	}

	roomID := ctx.Param("roomId")
	clientID := ctx.Query("userId")
	username := ctx.Query("username")

	cl := &Client{
		Conn: conn,
		Message: make(chan *Message, 10),
		ID: clientID,
		RoomID: roomID,
		Username: username,
	}

	m := &Message{
		Content: "A new user has joined the room",
		RoomID: roomID,
		Username: username,
	}

	//register a new client through the register channel
	c.hub.Register <- cl
	//broadcast that message
	c.hub.Broadcast <- m

	go cl.writeMessage()
	cl.readMessage(c.hub)

}