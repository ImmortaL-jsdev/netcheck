package diagnose

import (
	"fmt"
	"net"
)

func DiagnoseDNS(domain string) {
	ips, err := net.LookupHost(domain)
	if err != nil {
		fmt.Printf("❌ Ошибка DNS: %v\n", err)
		return
	}

	fmt.Printf("✅ Резолвится в: %v\n", ips)

	badIps := map[string]bool{
		"0.0.0.0":   true,
		"127.0.0.1": true,
	}

	for _, ip := range ips {
		if badIps[ip] {
			fmt.Printf("⚠️ Подмена DNS: %s\n", ip)
		}
	}
}
