package xmlsender

import (
	"GonIO/internal/domain"
	"encoding/xml"
	"net/http"
)

type Response struct {
	Message string `xml:"message"`
}

// Sends message in XML format
func SendMessage(w http.ResponseWriter, status int, message string) error {
	resp := Response{Message: message}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	if err := xml.NewEncoder(w).Encode(resp); err != nil {
		return err
	}
	return nil
}

func SendBucketList(w http.ResponseWriter, list []domain.Bucket) error {
	bucketList := domain.BucketList{Buckets: list}
	w.Header().Set("Content-Type", "application/xml")
	if err := xml.NewEncoder(w).Encode(bucketList); err != nil {
		return err
	}
	return nil
}

func SendObjectList(w http.ResponseWriter, list []domain.Object) error {
	objectList := domain.ObjectsList{Objects: list}
	w.Header().Set("Content-Type", "application/xml")
	if err := xml.NewEncoder(w).Encode(objectList); err != nil {
		return err
	}
	return nil
}
