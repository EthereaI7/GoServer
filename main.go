package main

func main() {
	server := NewServer("192.168.205.156", 8888)
	server.Start()
}
