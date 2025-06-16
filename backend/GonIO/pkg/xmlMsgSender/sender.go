package xmlsender

import (
	"GonIO/internal/domain"
	"encoding/xml"
	"errors"
	"net/http"
)

type Response struct {
	Message string `xml:"message"`
}

func SendMessage(w http.ResponseWriter, status int, message string) error {
	resp := Response{Message: message}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	if err := xml.NewEncoder(w).Encode(resp); err != nil {
		return err
	}
	return nil
}

func SendBucketList(w http.ResponseWriter, list any) error {
	converted, ok := list.([]domain.Bucket)
	if !ok {
		return errors.New("list value type assertion is failed")
	}

	bucketList := domain.BucketList{Buckets: converted}
	w.Header().Set("Content-Type", "application/xml")
	if err := xml.NewEncoder(w).Encode(bucketList); err != nil {
		return err
	}
	return nil
}
