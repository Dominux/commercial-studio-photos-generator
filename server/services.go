package main

import (
	"encoding/json"

	"github.com/google/uuid"
)

type GenerateImageSchema struct {
	Id      uuid.UUID `json:"id"`
	Product string    `json:"product"`
}

func runGenerating(product string, p *Producer) (uuid.UUID, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return uuid.Nil, err
	}

	msg := GenerateImageSchema{Id: id, Product: product}
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return uuid.Nil, err
	}

	if err := p.Send(msgJson); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
