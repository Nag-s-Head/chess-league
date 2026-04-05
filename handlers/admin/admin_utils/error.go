package adminutils

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"runtime"

	"github.com/Nag-s-Head/chess-league/handlers/utils"
)

var errorTpl *template.Template = utils.GetTemplate(f, "error.html")

type ErrorModel struct {
	Error error
}

func RenderError(w http.ResponseWriter, inErr error) {
	mdl := ErrorModel{Error: inErr}
	err := errorTpl.Execute(w, mdl)
	if err != nil {
		slog.Error("Could not render error page from template")
		w.Write([]byte("An error has occurred showing you the error"))
	}

	_, file, line, _ := runtime.Caller(1)
	slog.Error("An API error has occurred", "err", inErr, "caller", fmt.Sprintf("%s:%d", file, line))
}
