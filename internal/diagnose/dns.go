package diagnose

import (
	"fmt"
	"net"
)

type DNSResult struct {
	IPs      []string
	Hijacked bool
}

func DiagnoseDNS(domain string) DNSResult {
	ips, err := net.LookupHost(domain)
	if err != nil {
		fmt.Printf("❌ Ошибка DNS: %v\n", err)
		return DNSResult{IPs: nil, Hijacked: false}
	}

	fmt.Printf("✅ Резолвится в: %v\n", ips)

	badIps := map[string]bool{
		"0.0.0.0":   true,
		"127.0.0.1": true,
	}
	hijacked := false

	for _, ip := range ips {
		if badIps[ip] {
			hijacked = true
			fmt.Printf("⚠️ Подмена DNS: %s\n", ip)
		}
		parsed := net.ParseIP(ip)
		if parsed != nil && parsed.IsPrivate() {
			hijacked = true
			fmt.Printf("⚠️ Приватный IP для публичного домена: %s (возможна подмена или локальная настройка)\n", ip)
		}
	}
	return DNSResult{IPs: ips, Hijacked: hijacked}
}
