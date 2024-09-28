package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	go user.ListenMessage()

	return user
}

func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

//user online
func (this *User) Online() {
	//user online, add to OnlineMap
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	//put users' info into this.Message
	this.server.UserInfoEnQueue(this, "online")
}

//user offline
func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	//put users' info into this.Message
	this.server.UserInfoEnQueue(this, "offline")
}

func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

//user deal with message
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":online...\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else {
		this.server.UserInfoEnQueue(this, msg)
	}

}
