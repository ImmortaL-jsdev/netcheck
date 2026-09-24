package main

import (
	"fmt"
	"os"

	"github.com/ImmortaL-jsdev/netcheck/internal/diagnose"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("usage: netcheck <diagnose|fix|serve> [args]")
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "diagnose":
		if len(os.Args) < 3 {
			fmt.Println("usage: netcheck diagnose <domain>")
			os.Exit(1)
		}
		domain := os.Args[2]
		dnsResult := diagnose.DiagnoseDNS(domain)
		tlsResult := diagnose.DiagnoseTLS(domain)

		dohIPs, dohErr := diagnose.DiagnoseDoH(domain)
		if dohErr != nil {
			fmt.Printf("❌ DoH: %v\n", dohErr)
		} else {
			fmt.Printf("✅ DoH резолвит в: %v\n", dohIPs)
		}

		switch {
		case dnsResult.Hijacked:
			fmt.Println("📋 Итог: DNS-подмена. Рекомендуется DoH.")
		case tlsResult == diagnose.TLSStatusTimeout || tlsResult == diagnose.TLSStatusReset || tlsResult == diagnose.TLSStatusEOF:
			fmt.Println("📋 Итог: SNI-фильтрация. Рекомендуется фрагментация ClientHello.")
		case tlsResult == diagnose.TLSStatusOK:
			fmt.Println("📋 Итог: блокировки нет.")
		case tlsResult == diagnose.TLSStatusNetworkError:
			fmt.Println("📋 Итог: сетевая ошибка. Проверьте подключение.")
		default:
			fmt.Println("📋 Итог: неизвестная ситуация.")
		}
	case "fix":
		if len(os.Args) < 3 {
			fmt.Println("usage: netcheck fix <domain>")
			os.Exit(1)
		}
	default:
		fmt.Println("unknown command:", command)
		os.Exit(1)
	}
}
