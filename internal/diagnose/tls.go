package diagnose

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	TLSStatusOK           = "ok"
	TLSStatusTimeout      = "timeout"
	TLSStatusReset        = "reset"
	TLSStatusEOF          = "eof"
	TLSStatusCert         = "certificate"
	TLSStatusTCPFailed    = "tcp_failed"
	TLSStatusOther        = "other"
	TLSStatusNetworkError = "network_error"
)

func DiagnoseTLS(domain string) string {
	conn, err := net.DialTimeout("tcp", domain+":443", 5*time.Second)
	if err != nil {
		fmt.Printf("❌ TCP: не удалось подключиться:%v\n", err)
		return TLSStatusTCPFailed
	}
	defer conn.Close()

	config := &tls.Config{
		ServerName: domain,
		MinVersion: tls.VersionTLS12,
	}
	tlsConn := tls.Client(conn, config)

	_ = tlsConn.SetDeadline(time.Now().Add(5 * time.Second))

	if err := tlsConn.Handshake(); err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "timeout") {
			fmt.Println("⚠️ DPI: соединение зависло (timeout).")
			return TLSStatusTimeout
		}
		if strings.Contains(errStr, "reset") {
			fmt.Println("⚠️ DPI: соединение оборвано (SNI-фильтрация).")
			return TLSStatusReset
		}
		if strings.Contains(errStr, "EOF") {
			fmt.Println("⚠️ DPI: соединение закрыто после ClientHello.")
			return TLSStatusEOF
		}
		if strings.Contains(errStr, "certificate") {
			fmt.Println("❌ Проблема с сертификатом (не блокировка, но соединение небезопасно).")
			return TLSStatusCert
		}
		if strings.Contains(errStr, "no route to host") || strings.Contains(errStr, "network is unreachable") {
			fmt.Println("❌ Хост недоступен. Проверьте IP-адрес или сеть.")
			return TLSStatusNetworkError
		}
		fmt.Printf("❌ TLS ошибка: %v\n", err)
		return TLSStatusOther
	}
	fmt.Println("✅ TLS: рукопожатие успешно")
	return TLSStatusOK
}
