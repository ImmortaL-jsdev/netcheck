package proxy

import (
	"crypto/tls"
	"net"
	"testing"
)

func TestFindSni(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	go func() {
		tlsConn := tls.Client(client, &tls.Config{
			ServerName: "rutracker.org",
			MinVersion: tls.VersionTLS12,
		})
		_ = tlsConn.Handshake()
	}()

	buf := make([]byte, 4096)
	n, _ := server.Read(buf)
	data := buf[:n]

	offset, length, ok := findSNI(data)

	if !ok {
		t.Fatal("SNI не найден")
	}

	domain := string(data[offset : offset+length])
	if domain != "rutracker.org" {
		t.Errorf("ожидался rutracker.org, получен %q", domain)
	}
}
