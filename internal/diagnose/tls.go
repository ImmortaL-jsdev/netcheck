package diagnose

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"
)

func DiagnoseTLS(domain string) {
	conn, err := net.DialTimeout("tcp", domain+":443", 5*time.Second)
	if err != nil {
		fmt.Printf("❌ TCP: не удалось подключиться:%v\n", err)
		return
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
			return
		}
		if strings.Contains(errStr, "reset") {
			fmt.Println("⚠️ DPI: соединение оборвано (SNI-фильтрация).")
			return
		}
		if strings.Contains(errStr, "EOF") {
			fmt.Println("⚠️ DPI: соединение закрыто после ClientHello.")
			return
		}
		if strings.Contains(errStr, "certificate") {
			fmt.Println("❌ Проблема с сертификатом (не блокировка, но соединение небезопасно).")
			return
		}
		if strings.Contains(errStr, "no route to host") || strings.Contains(errStr, "network is unreachable") {
			fmt.Println("❌ Хост недоступен. Проверьте IP-адрес или сеть.")
			return
		}
		fmt.Printf("❌ TLS ошибка: %v\n", err)
		return
	}
	fmt.Println("✅ TLS: рукопожатие успешно")
}
