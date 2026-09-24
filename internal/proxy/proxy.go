package proxy

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

func Start(ctx context.Context, addr string, resolver func(string) ([]string, error)) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	defer listener.Close()

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn, resolver)
	}
}

func handleConnection(clientConn net.Conn, resolver func(string) ([]string, error)) {
	defer clientConn.Close()

	reader := bufio.NewReader(clientConn)

	line, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	line = strings.TrimSpace(line)
	parts := strings.Split(line, " ")

	if len(parts) < 3 || parts[0] != "CONNECT" {
		_, _ = clientConn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}

	hostPort := strings.Split(parts[1], ":")
	if len(hostPort) != 2 {
		_, _ = clientConn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}
	domain := hostPort[0]
	port := hostPort[1]

	ips, err := resolver(domain)

	if err != nil || len(ips) == 0 {
		_, _ = clientConn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
		return
	}
	targetIP := ips[0]
	fmt.Println("Резолв:", targetIP)

	serverConn, err := net.DialTimeout("tcp", targetIP+":"+port, 10*time.Second)
	if err != nil {
		_, _ = clientConn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
		return
	}

	defer serverConn.Close()
	fmt.Println("Соединение с сервером установлено")

	_, _ = clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	go func() {
		_, _ = io.Copy(serverConn, clientConn)
	}()

	_, _ = io.Copy(clientConn, serverConn)
}
