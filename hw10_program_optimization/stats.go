package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
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
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		var user User
		line := scanner.Text()
		if err := json.Unmarshal([]byte(line), &user); err != nil {
			return nil, fmt.Errorf("error unmarshalling user: %w", err)
		}

		emailParts := strings.Split(user.Email, "@")
		if len(emailParts) != 2 {
			continue
		}
		emailDomain := strings.ToLower(emailParts[1])

		if strings.HasSuffix(emailDomain, "."+domain) {
			primaryDomain := strings.TrimSuffix(emailDomain, "."+domain)
			result[primaryDomain]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	return result, nil
}

// func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
// 	u, err := getUsers(r)
// 	if err != nil {
// 		return nil, fmt.Errorf("get users error: %w", err)
// 	}
// 	return countDomains(u, domain)
// }

// type users [100_000]User

// func getUsers(r io.Reader) (result users, err error) {
// 	content, err := io.ReadAll(r)
// 	if err != nil {
// 		return
// 	}

// 	lines := strings.Split(string(content), "\n")
// 	for i, line := range lines {
// 		var user User
// 		if err = json.Unmarshal([]byte(line), &user); err != nil {
// 			return
// 		}
// 		result[i] = user
// 	}
// 	return
// }

// func countDomains(u users, domain string) (DomainStat, error) {
// 	result := make(DomainStat)

// 	for _, user := range u {
// 		matched, err := regexp.Match("\\."+domain, []byte(user.Email))
// 		if err != nil {
// 			return nil, err
// 		}

// 		if matched {
// 			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
// 			num++
// 			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
// 		}
// 	}
// 	return result, nil
// }
