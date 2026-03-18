package qrcode

import (
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	dbmodel "github.com/Nag-s-Head/chess-league/db/model"
	submitgame "github.com/Nag-s-Head/chess-league/handlers/submit_game"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

type bufferWriteCloser struct {
	*bytes.Buffer
}

// Close satisfies the io.Closer interface
func (b *bufferWriteCloser) Close() error {
	return nil
}

type Model struct {
	QrCodeB64 string
	Url       string
}

func getModel() Model {
	url := fmt.Sprintf("%s%s?%s=%s",
		os.Getenv("APP_BASE_URL"),
		submitgame.BasePath,
		submitgame.MagicNumberParam,
		os.Getenv("MAGIC_NUMBER"))

	qrc, err := qrcode.NewWith(url, qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionMedium))
	if err != nil {
		panic(fmt.Sprintf("Could not generate QR code: %s", err))
	}

	options := []standard.ImageOption{
		standard.WithBgColorRGBHex("#ffffff"),
		standard.WithFgColorRGBHex("#000000"),
		standard.WithHalftone("knight.png"),
	}

	buffer := bytes.NewBuffer(make([]byte, 0))
	w := &bufferWriteCloser{buffer}
	writer := standard.NewWithWriter(w, options...)

	err = qrc.Save(writer)
	if err != nil {
		panic(fmt.Sprintf("Could not generate QR code: %s", err))
	}

	return Model{
		QrCodeB64: base64.StdEncoding.EncodeToString(buffer.Bytes()),
		Url:       url,
	}
}

var model Model = getModel()

//go:embed qr_code.html
var f embed.FS
var indexTpl *template.Template = utils.GetTemplate(f, "qr_code.html")

func Render(user *dbmodel.AdminUser) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		err := indexTpl.Execute(&buf, model)
		if err != nil {
			slog.Error("Could not execute template", "err", err)
			w.Write([]byte("Could not execute template"))
			return
		}

		w.Write(buf.Bytes())
	}
}
