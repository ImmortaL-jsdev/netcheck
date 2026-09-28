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

	for {
		h, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		if strings.TrimSpace(h) == "" {
			break
		}
	}

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

	// header: [CT=0x16][Ver:2][Len:2]
	header := make([]byte, 5)
	if _, err := io.ReadFull(reader, header); err != nil {
		fmt.Println("Ошибка чтения заголовка:", err)
		return
	}

	recordLen := int(header[3])<<8 | int(header[4]) // длина payload

	// payload: [HandshakeHdr:4][ClientHello Body]
	payload := make([]byte, recordLen)
	_ = clientConn.SetReadDeadline(time.Now().Add(5 * time.Second))

	total := 0
	for total < recordLen {
		n, err := reader.Read(payload[total:])
		if err != nil {
			fmt.Printf("Ошибка чтения payload (%d%d): %v\n", total, recordLen, err)
			return
		}
		total += n
	}

	data := append(header, payload...) // полный ClientHello
	fmt.Println("Прочитано байт:", len(data))

	offset, length, ok := findSNI(data)
	if ok {
		domain := string(data[offset : offset+length])
		fmt.Println("SNI найден:", domain)

	} else {
		fmt.Println("SNI не найден")
	}

	if tcpConn, ok := serverConn.(*net.TCPConn); ok {
		_ = tcpConn.SetNoDelay(true)
	}

	chunkSize := 5
	if tcpConn, ok := serverConn.(*net.TCPConn); ok {
		_ = tcpConn.SetNoDelay(true)
	}

	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		_, _ = serverConn.Write(data[i:end])
		time.Sleep(20 * time.Millisecond)
	}

	fmt.Printf("Multisplit: %d фрагментов\n", len(data)/chunkSize+1)

	go func() {
		_, _ = io.Copy(serverConn, reader)
	}()

	_, _ = io.Copy(clientConn, serverConn)
}
