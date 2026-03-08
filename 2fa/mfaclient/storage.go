package main

import (
	"encoding/json"
	"os"
)

type Account struct {
	ID      string `json:"id"`
	Issuer  string `json:"issuer"`
	Label   string `json:"label"`
	Secret  string `json:"secret"`
	AddedAt int64  `json:"added_at"`
}
type Storage struct {
	filepath string
}

func NewStorage(filepath string) (*Storage, error) {
	return &Storage{filepath: filepath}, nil
}

func (s *Storage) SaveAccounts(accounts []Account) error {
	data, err := json.Marshal(accounts)
	if err != nil {
		return err
	}

	return os.WriteFile(s.filepath, data, 0600)
}

func (s *Storage) LoadAccounts() ([]Account, error) {
	data, err := os.ReadFile(s.filepath)
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
