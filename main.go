package main

import (
	"crypto/sha1"
	"encoding/base64"
	"log"
	"net"
	"strings"
	"sync"
)

const PORT string = ":5900"
const WS_UUID string = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

var (
	clients []net.Conn
	mutex   sync.Mutex
)

func main() {
	server, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Printf("Could not start server: %s\n", err)
	}

	log.Printf("Listening for TCP on %s\n", PORT)

	for {
		connection, err := server.Accept()
		if err != nil {
			log.Printf("Could not accept new TCP connection: %s\n", err)
			continue
		}

		go notifyConnected(clients, "New brother joined!")

		go handleConnection(clients, connection)
	}
}

func handleConnection(clients []net.Conn, client net.Conn) {
	buffer := make([]byte, 1024)

	_, err := client.Read(buffer)
	if err != nil {
		log.Printf("Error reading handshake request: %s", err)
	}

	valid, nonce := validateWebSocketHandshake(buffer)
	if !valid {
		log.Printf("Error or invalid WebSocket handshake: %s", err)

		return
	}

	_, err = acceptConnection(client, nonce)
	if err != nil {
		log.Printf("Could not accept connection: %s", err)

		client.Close()
	}

	mutex.Lock()
	clients = append(clients, client)
	mutex.Unlock()

	go keepAlive(client)
}

func keepAlive(client net.Conn) {
	//    reader := bufio.NewReader()
	//
	// for {
	//
	// 	client.Read(buffer)
	//
	// 	log.Printf("Received payload: %s", buffer)
	// }
}

func validateWebSocketHandshake(request []byte) (bool, string) {
	requestAsString := string(request)
	headers := strings.Split(requestAsString, "\r\n")

	log.Printf("Request:\r\n%s", requestAsString)

	for _, headerAndValue := range headers {
		headerParts := strings.SplitN(headerAndValue, ":", 2)

		if len(headerParts) < 2 {
			continue
		}

		if strings.ToLower(headerParts[0]) == "sec-websocket-key" {
			value := strings.TrimSpace(headerParts[1])

			base64EncodedNonce := generateWebSocketNonce(value)

			log.Printf("Base64EncodedNonce: %s", base64EncodedNonce)

			return true, base64EncodedNonce
		}

	}

	return false, ""
}

func generateWebSocketNonce(key string) string {
	keyAndUUID := key + WS_UUID

	h := sha1.New()

	h.Write([]byte(keyAndUUID))

	sha := h.Sum(nil)

	return base64.StdEncoding.EncodeToString(sha)
}

func acceptConnection(client net.Conn, nonce string) (int, error) {
	acceptHeaders := []string{
		"HTTP/1.1 101 Switching Protocols",
		"Upgrade: websocket",
		"Connection: upgrade",
		"Sec-Websocket-Accept: " + nonce,
		"",
		"",
	}

	return client.Write([]byte(strings.Join(acceptHeaders, "\r\n")))
}

func notifyConnected(clients []net.Conn, message string) {
	for _, client := range clients {
		sendMessage(client, message)
	}
}

func sendMessage(conn net.Conn, message string) {
	conn.Write([]byte(message))
}
