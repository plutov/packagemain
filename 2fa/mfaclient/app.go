package main

import (
	"context"
)

type App struct {
	ctx     context.Context
	storage *Storage
}

type AccountWithCode struct {
	Account
	Code          string `json:"code"`
	TimeRemaining int    `json:"time_remaining"`
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.storage = &Storage{
		Filepath: "./accounts.json",
	}
}

func (a *App) GetAccounts() ([]AccountWithCode, error) {
	accounts, err := a.storage.LoadAccounts()
	if err != nil {
		return nil, err
	}

	result := make([]AccountWithCode, len(accounts))
	for i, acc := range accounts {
		code, timeRemaining, err := GenerateTotp(acc.Secret)
		if err != nil {
			return []AccountWithCode{}, err
		}

		result[i] = AccountWithCode{
			Account:       acc,
			Code:          code,
			TimeRemaining: timeRemaining,
		}
	}
	return result, nil
}

func (a *App) AddAccount(issuer, secret string) error {
	newAccount := Account{
		Issuer: issuer,
		Secret: secret,
	}

	return a.storage.AddAccount(newAccount)
}

func (a *App) DeleteAccount(id string) error {
	return a.storage.DeleteAccount(id)
}
