package testmode

import (
	"bytes"
	"embed"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/Nag-s-Head/chess-league/handlers/admin/auth"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
)

//go:embed login.html
var f embed.FS
var loginTpl *template.Template = utils.GetTemplate(f, "login.html")

func Login(w http.ResponseWriter, r *http.Request) (template.HTML, error) {
	var buf bytes.Buffer
	err := loginTpl.Execute(&buf, nil)
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

const (
	adminId    = "admin-id"
	adminName  = "admin-name"
	submitType = "submit-type"
)

func LoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		slog.Error("Could not parse form", "err", err)
		w.Write([]byte("Could not paarse form"))
		return
	}

	adminId := r.Form.Get(adminId)
	adminName := r.Form.Get(adminName)
	submitType := r.Form.Get(submitType)

	cookieValue := "invalid"
	if submitType == "valid" {
		if adminId == "" || adminName == "" {
			w.Write([]byte("Admin name and Admin ID are required for a valid user"))
			return
		}

		// TODO: change to a session token for the user
		cookieValue = "valid"
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.AuthCookie,
		Value:    cookieValue,
		HttpOnly: true,
		Secure:   true,
		MaxAge:   60 * 60,
	})

	w.Write([]byte(`
<script>
	console.log("Executing test mode redirect");
  window.location.href = "/admin";
</script>
`))
}
