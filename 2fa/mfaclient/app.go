package main

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"
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
	a.storage, _ = NewStorage("./accounts.data")
}

func (a *App) GetAccounts() ([]AccountWithCode, error) {
	accounts, err := a.storage.LoadAccounts()
	if err != nil {
		return nil, err
	}

	result := make([]AccountWithCode, len(accounts))
	for i, acc := range accounts {
		code, timeLeft, _ := GenerateTOTP(acc.Secret)
		result[i] = AccountWithCode{
			Account:       acc,
			Code:          code,
			TimeRemaining: timeLeft,
		}
	}
	return result, nil
}

func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (a *App) AddAccount(issuer, label, secret string) error {
	cleanSecret := strings.TrimRight(strings.ToUpper(secret), "=")
	_, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(cleanSecret)
	if err != nil {
		return fmt.Errorf("unable to decode the secret")
	}

	accounts, _ := a.storage.LoadAccounts()

	newAccount := Account{
		ID:      generateID(),
		Issuer:  issuer,
		Label:   label,
		Secret:  secret,
		AddedAt: time.Now().Unix(),
	}

	accounts = append(accounts, newAccount)
	return a.storage.SaveAccounts(accounts)
}

func (a *App) ParseOTPAuthURI(uri string) (issuer, label, secret string, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", "", "", err
	}

	if u.Scheme != "otpauth" || u.Host != "totp" {
		return "", "", "", fmt.Errorf("invalid OTP URI")
	}

	label = strings.TrimPrefix(u.Path, "/")
	secret = u.Query().Get("secret")
	issuer = u.Query().Get("issuer")

	return issuer, label, secret, nil
}

func (a *App) AddAccountFromURI(uri string) error {
	issuer, label, secret, err := a.ParseOTPAuthURI(uri)
	if err != nil {
		return err
	}
	return a.AddAccount(issuer, label, secret)
}

func (a *App) DeleteAccount(id string) error {
	accounts, _ := a.storage.LoadAccounts()

	filtered := []Account{}
	for _, acc := range accounts {
		if acc.ID != id {
			filtered = append(filtered, acc)
		}
	}

	return a.storage.SaveAccounts(filtered)
}
