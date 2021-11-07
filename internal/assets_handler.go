package internal

import (
	"net/http"
)

var assets = LoadNestedTemplates("template/assets.htm")

func AssetsHandler(w http.ResponseWriter, r *http.Request) {
	Render(w, assets, nil)
}
