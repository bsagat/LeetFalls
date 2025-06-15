package handlers

import (
	"errors"
	"html/template"
	"leetFalls/internal/domain"
	"log/slog"
	"net/http"
)

func ErrorPage(w http.ResponseWriter, message error, code int) error {
	if message == nil {
		return errors.New("error message is empty")
	}
	data := struct {
		Code    int
		Message string
	}{
		Code:    code,
		Message: message.Error(),
	}

	w.WriteHeader(code)
	path := domain.Config.TemplatesPath

	temp, err := template.ParseFiles(path + "/error.html")
	if err != nil {
		slog.Error("Failed to Parse template file: ", "error", err.Error())
		return err
	}

	if err := temp.Execute(w, data); err != nil {
		slog.Error("Failed to Execute template: ", "error", err.Error())
		return err
	}

	return nil
}
