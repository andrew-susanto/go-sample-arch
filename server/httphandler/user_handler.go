package httphandler

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (h Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	paramID, ok := h.router.GetRouteParams(r)["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(paramID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := h.account.GetUser(id)
	// user.ID = id

	encoded, _ := json.Marshal(user)

	w.WriteHeader(http.StatusOK)
	w.Write(encoded)
	return
}
