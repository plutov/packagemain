package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
)

type Account struct {
	Id     string `json:"id"`
	Issuer string `json:"issuer"`
	Secret string `json:"secret"`
}

type Storage struct {
	Filepath string
}

func (s *Storage) LoadAccounts() ([]Account, error) {
	data, err := os.ReadFile(s.Filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Account{}, nil
		}

		return nil, err
	}

	var accounts []Account
	json.Unmarshal(data, &accounts)
	return accounts, nil
}

func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (s *Storage) AddAccount(a Account) error {
	accounts, err := s.LoadAccounts()
	if err != nil {
		return err
	}

	a.Id = generateID()
	accounts = append(accounts, a)
	return s.save(accounts)
}

func (s *Storage) DeleteAccount(id string) error {
	accounts, err := s.LoadAccounts()
	if err != nil {
		return err
	}

	filtered := []Account{}
	for _, acc := range accounts {
		if acc.Id != id {
			filtered = append(filtered, acc)
		}
	}

	return s.save(filtered)
}

func (s *Storage) save(accounts []Account) error {
	data, err := json.Marshal(accounts)
	if err != nil {
		return err
	}

	return os.WriteFile(s.Filepath, data, 0o600)
}
