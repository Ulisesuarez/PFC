package model

import "time"

type ShareRQ struct {
	Account    int        `json:"account"`
	Group      int        `json:"group"`
	Guests     string     `json:"guest"`
	ExpiredAt  *time.Time `json:"expired_at"`
	Permission string     `json:"permission"`
	Message    string     `json:"message,omitempty"`
}

type TransferRQ struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type PublicURLRQ struct {
	Namespace string `json:"namespace"`
}
