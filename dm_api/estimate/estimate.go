package api_estimate

import (
	x "dm_server/dm_xml"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write(x.ParseSmeta())
	w.WriteHeader(http.StatusOK)
}
