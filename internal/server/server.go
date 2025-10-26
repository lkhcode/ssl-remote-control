package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/RoboCup-SSL/ssl-remote-control/internal/rcon"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Server struct {
	connections []*ServerConnection
	rconClient  *rcon.Client
}

type ServerConnection struct {
	conn       *websocket.Conn
	marshaler  protojson.MarshalOptions
	rconClient *rcon.Client
}

func NewServer(rconClient *rcon.Client) (s *Server) {
	s = new(Server)
	s.connections = []*ServerConnection{}
	s.rconClient = rconClient
	return
}

// newServerConnection creates a new connection between server and UI client
func newServerConnection(rconClient *rcon.Client, conn *websocket.Conn) (s *ServerConnection) {
	s = new(ServerConnection)
	s.conn = conn
	s.marshaler.EmitUnpopulated = true
	s.rconClient = rconClient
	return
}

// Publish publishes the state to all connected clients
func (a *Server) Publish(state *rcon.ControllerToRemoteControl) {
	for _, serverConn := range a.connections {
		serverConn.publishState(state)
	}
}

// PublishRefereeData publishes referee data to all connected clients
func (a *Server) PublishRefereeData(data []byte) {
	for _, serverConn := range a.connections {
		serverConn.publishRefereeData(data)
	}
}

// WsHandler handles incoming web socket connections
func (a *Server) WsHandler(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(*http.Request) bool { return true },
	}

	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	serverConn := newServerConnection(a.rconClient, conn)
	a.connections = append(a.connections, serverConn)
	defer a.disconnect(serverConn)
	log.Println("UI Client connected")

	a.listenForMessages(conn)
}

func (s *ServerConnection) publishState(state *rcon.ControllerToRemoteControl) {
	b, err := s.marshaler.Marshal(state)
	if err != nil {
		log.Println("Marshal error:", err)
	}

	err = s.conn.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		log.Println("Could not write message to api client:", err)
	}
}

func (s *ServerConnection) publishRefereeData(data []byte) {
	// Unmarshal the binary protobuf data
	referee := &rcon.Referee{}
	err := proto.Unmarshal(data, referee)
	if err != nil {
		log.Println("Could not unmarshal referee data:", err)
		return
	}

	// Marshal the referee data to JSON
	refereeJSON, err := protojson.Marshal(referee)
	if err != nil {
		log.Println("Could not marshal referee data to JSON:", err)
		return
	}

	// Create a wrapper object with type and payload
	wrapper := map[string]interface{}{
		"type":    "referee",
		"payload": json.RawMessage(refereeJSON),
	}

	// Marshal the wrapper to JSON
	messageJSON, err := json.Marshal(wrapper)
	if err != nil {
		log.Println("Could not marshal wrapper to JSON:", err)
		return
	}

	err = s.conn.WriteMessage(websocket.TextMessage, messageJSON)
	if err != nil {
		log.Println("Could not write referee data to api client:", err)
	}
}

func (a *Server) disconnect(conn *ServerConnection) {
	err := conn.conn.Close()
	if err != nil {
		log.Println("Could not disconnect from websocket conn: ", err)
	}
	for i, c := range a.connections {
		if c == conn {
			a.connections = append(a.connections[:i], a.connections[i+1:]...)
			break
		}
	}
	log.Println("UI Client disconnected")
}

func (a *Server) listenForMessages(conn *websocket.Conn) {
	for {
		messageType, b, err := conn.ReadMessage()
		if err != nil || messageType != websocket.TextMessage {
			log.Println("Could not read message: ", err)
			return
		}

		a.handleNewMessage(b)
	}
}

func (a *Server) handleNewMessage(b []byte) {
	in := &rcon.RemoteControlToController{}
	err := protojson.Unmarshal(b, in)
	if err != nil {
		log.Println("Could not read input:", string(b), err)
		return
	}

	reply, err := a.rconClient.SendRequest(in)
	if err != nil {
		log.Println("Failed to send request: ", err)
	}
	a.Publish(reply)
}
