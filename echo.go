package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// handlerV1Echo handles echo request
// This accepts GET, POST, PUT, DELETE and reflect body
func handleEcho(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// params: callback
	echoHeaders := r.URL.Query().Get("echo-hdr")
	dummyStatus := r.URL.Query().Get("dummy-status")

	var convStatus int64
	var err error

	// get status from dummyStatus query
	if dummyStatus == "" {
		convStatus = 200
	} else {
		convStatus, err = strconv.ParseInt(dummyStatus, 0, 16)
		if err != nil {
			log.Warningf("converting %s, but error %s", dummyStatus, err.Error())

			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			return
		}
	}
	contentType := r.Header.Get("Content-Type")
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(int(convStatus))

	// compose headers
	hdrs := strings.Split(echoHeaders, ",")
	for _, hdr := range hdrs {
		v := r.Header.Get(hdr)
		if v != "" {
			w.Header().Set(hdr, v)
		}
	}

	if r.Body != nil {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorNotFound)
			return
		}
		w.Write(body)
	}
	return
}
