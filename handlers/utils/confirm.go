package utils

import (
	"embed"
	"errors"
	"html/template"
	"net/http"
)

const (
	confirmCheckbox = "confirm"
)

//go:embed confirm.html
var f embed.FS
var confirmTpl *template.Template = GetTemplate(f, "confirm.html")

type confirmModel struct {
	Action      string
	ButtonValue string
}

func RenderConfirmationPage(w http.ResponseWriter, action string, buttonValue string) error {
	err := confirmTpl.Execute(w, confirmModel{
		Action:      action,
		ButtonValue: buttonValue,
	})

	if err != nil {
		return errors.Join(errors.New("Cannot render confirmation page"), err)
	}

	return nil
}

// Requires the form to be parsed before this is called
func IsConfirmed(r *http.Request) bool {
	confirmed := r.Form.Get(confirmCheckbox)
	return confirmed == "confirmed"
}
