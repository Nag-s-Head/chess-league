package assets_test

import (
	"net/http"
	"testing"

	"github.com/Nag-s-Head/chess-league/handlers/assets"
)

func TestAssetsHandler(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	assets.Register(mux, []byte(".test {}"))
}
