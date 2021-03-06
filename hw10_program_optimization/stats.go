package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var user User
	var err error
	result := make(DomainStat)

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if !json.Valid(scanner.Bytes()) {
			return nil, fmt.Errorf("invalid json")
		}
		if strings.Contains(scanner.Text(), "."+domain) {
			if err = json.Unmarshal(scanner.Bytes(), &user); err != nil {
				return nil, fmt.Errorf("get users error: %w", err)
			}
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}
