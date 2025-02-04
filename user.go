package main

import (
	"net"
	"strings"
)

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

// user online
func (this *User) Online() {
	//user online, add to OnlineMap
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	//put users' info into this.Message
	this.server.UserInfoEnQueue(this, "online")
}

// user offline
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

// user deal with message
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":online\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//rename
		newName := strings.Split(msg, "|")[1]
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("this username is already used!\n")
		} else {
			this.server.mapLock.Lock()
			this.server.OnlineMap[newName] = this
			delete(this.server.OnlineMap, this.Name)
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMsg("rename success: " + this.Name + "\n")
		}
	} else if len(msg) > 3 && msg[:3] == "to|" {
		//send private message
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMsg("message format error! format: to|someone|msg\n")
			return
		}

		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("user does not exist!\n")
			return
		}

		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMsg("message format error! format: to|someone|msg\n")
			return
		}

		remoteUser.SendMsg("(Private chat)" + this.Name + ": " + content + "\n")

	} else {
		this.server.UserInfoEnQueue(this, msg)
	}

}
