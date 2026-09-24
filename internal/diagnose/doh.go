package diagnose

import (
	"context"
	"fmt"
	"time"

	"github.com/shahradelahi/go-doh-client"
)

func DiagnoseDoH(domain string) ([]string, error) {
	client := doh.New()
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp4, err := client.Query(ctx, doh.Domain(domain), doh.TypeA)
	if err != nil {
		return nil, fmt.Errorf("DoH запрос не удался: %w", err)
	}

	var ips4 []string
	for _, answer := range resp4.Answer {
		switch answer.Type {
		case 1: //IPv4
			ips4 = append(ips4, answer.Data)
		}

	}

	resp6, err := client.Query(ctx, doh.Domain(domain), doh.TypeAAAA)
	if err != nil {
		return nil, fmt.Errorf("DoH запрос не удался: %w", err)
	}
	var ips6 []string
	for _, answer := range resp6.Answer {
		switch answer.Type {
		case 28: //IPv6
			ips6 = append(ips6, answer.Data)
		}

	}
	all := append(ips4, ips6...)
	return all, nil
}
