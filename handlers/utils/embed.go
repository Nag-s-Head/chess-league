package utils

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"runtime"
)

func GetTemplate(f embed.FS, name string) *template.Template {
	data, err := f.ReadFile(name)
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		caller := fmt.Sprintf("%s:%d", file, line)
		slog.Error("Cannot read embedded file", "name", name, "caller", caller)
		panic("Cannot read embedded file")
	}

	tpl, err := template.New("tpl").Parse(string(data))
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		caller := fmt.Sprintf("%s:%d", file, line)
		slog.Error("Cannot parse template", "name", name, "caller", caller)
		panic("Cannot parse template")
	}

	return tpl
}
