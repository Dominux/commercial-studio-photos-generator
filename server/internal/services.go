package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
)

type GenerateImageSchema struct {
	Id      uuid.UUID `json:"id"`
	Product string    `json:"product"`
}

type GenerateResultSchema struct {
	Id    uuid.UUID `json:"id"`
	State string    `json:"state"`
}

func RunGenerating(product string, p *Producer) (uuid.UUID, error) {
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

func SaveResults(results chan []byte, redisClient *redis.Client, resultsExpire time.Duration) {
	for msg := range results {
		var result GenerateResultSchema
		if err := json.Unmarshal(msg, &result); err != nil {
			fmt.Printf("error on unmarshaling msg: %s; %s", msg, err)
			continue
		}

		ctx := context.Background()

		if err := redisClient.Set(ctx, result.Id.String(), nil, resultsExpire).Err(); err != nil {
			fmt.Printf("error on saving into redis: %s\n", err)
		}
	}
}

func CheckProgress(ctx context.Context, rawId string, redisClient *redis.Client) error {
	return redisClient.Get(ctx, rawId).Err()
}

func GetResult(ctx context.Context, rawId string, minioClient *minio.Client, minioBucket string, minioResultsPath string) ([]byte, error) {
	stripedId := strings.ReplaceAll(rawId, "-", "")
	path := fmt.Sprintf("%s/%s.png", minioResultsPath, stripedId)

	obj, err := minioClient.GetObject(ctx, minioBucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	buf, err := io.ReadAll(obj)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
