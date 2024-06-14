package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/goccy/go-json"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	domainStat := make(DomainStat)
	scanner := bufio.NewScanner(r)
	domainSuffix := "." + domain

	for scanner.Scan() {
		line := scanner.Bytes()
		email := extractEmail(line)

		if strings.HasSuffix(email, domainSuffix) {
			domainPart := strings.SplitN(email, "@", 2)[1]
			// domainName := strings.ToLower(strings.SplitN(domainPart, ".", 2)[0])
			domainName := strings.ToLower(domainPart)
			domainStat[domainName]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	return domainStat, nil
}

func extractEmail(line []byte) string {
	var d User
	err := json.Unmarshal(line, &d)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	return d.Email
}
