package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ImmortaL-jsdev/netcheck/internal/diagnose"
	"github.com/ImmortaL-jsdev/netcheck/internal/proxy"
	"github.com/ImmortaL-jsdev/netcheck/internal/zapret"
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
	case "serve":
		addr := "127.0.0.1:8080"
		resolver := func(domain string) ([]string, error) {
			return diagnose.DiagnoseDoH(domain)
		}
		fmt.Println("Прокси запущен на", addr)
		if err := proxy.Start(context.Background(), addr, resolver); err != nil {
			fmt.Println("Ошибка прокси:", err)
			os.Exit(1)
		}
	case "zapret":
		if len(os.Args) < 3 {
			fmt.Println("usage: netcheck zapret <status>")
			os.Exit(1)
		}
		subcommand := os.Args[2]
		switch subcommand {
		case "status":
			z, err := zapret.Detect()
			if err != nil {
				fmt.Println("❌", err)
				os.Exit(1)
			}
			fmt.Println("✅ Zapret найден:", z.Path)
			fmt.Println("📦 Версия:", z.Version)
			if z.IsRunning() {
				fmt.Println("⚙️ Сервис: активен")
			} else {
				fmt.Println("⚙️ Сервис: остановлен")
			}

			if z.HasBinary("nfqws") {
				fmt.Println("📁 nfqws: ✅")
			} else {
				fmt.Println("📁 nfqws: ❌")
			}
		default:
			fmt.Println("unknown zapret subcommand:", subcommand)
			os.Exit(1)
		}
	default:
		fmt.Println("unknown command:", command)
		os.Exit(1)
	}
}
