package server

import (
	"sync"

	"duel-masters/db"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var sockets = make(map[*Socket]Hub)
var socketsMutex = sync.Mutex{}

// Socket links a ws connection to a user id and handles safe reading and writing of data
type Socket struct {
	conn   *websocket.Conn
	User   db.User
	hub    Hub
	ready  bool
	mutex  *sync.Mutex
	closed bool
}

// NewSocket creates and returns a new Socket instance
func NewSocket(c *websocket.Conn, hub Hub) *Socket {

	s := &Socket{
		conn:   c,
		hub:    hub,
		ready:  false,
		mutex:  &sync.Mutex{},
		closed: false,
	}

	socketsMutex.Lock()
	sockets[s] = hub
	socketsMutex.Unlock()

	return s

}

// Listen sets up reader and writer for the socket
func (s *Socket) Listen() {

	defer s.Close()

	for {

		_, message, err := s.conn.ReadMessage()

		if err != nil {
			//logrus.Debug(err)
			return
		}

		if !s.ready {

			// Look for authorization token as the first message
			u, err := db.GetUserForToken(string(message))

			if err != nil {
				continue
			}

			s.User = u
			s.ready = true

			s.Send(Message{Header: "hello"})

			continue

		}

		go s.hub.Parse(s, message)

	}

}

// Send sends a struct v to the client
func (s *Socket) Send(v interface{}) {

	if s.closed {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("Recovered from panic in socket Send. %v", r)
			return
		}
	}()

	s.mutex.Lock()
	if err := s.conn.WriteJSON(v); err != nil {
		logrus.Debug(err)
	}
	s.mutex.Unlock()

}

// Close closes the client connection
func (s *Socket) Close() {

	if s.closed {
		return
	}

	delete(sockets, s)

	if s == nil || s.conn == nil {
		return
	}

	socketsMutex.Lock()

	defer socketsMutex.Unlock()

	s.conn.Close()

	s.closed = true

	logrus.Debug("Closed a connection")

}