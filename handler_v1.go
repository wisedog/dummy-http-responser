package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// handler for /v1/:id
// If the user specifies 'dummy-status', the status is overrided
func handleV1Custom(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO jsonp will be supported
	// param : dummy-status
	dummyStatus := r.URL.Query().Get("dummy-status")
	dummyID := ps.ByName("id")
	if ok := bson.IsObjectIdHex(dummyID); !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(errorInvalidID)
		return
	}

	var convStatus int64
	var err error
	if dummyStatus == "" {
		convStatus = 200
	} else {
		convStatus, err = strconv.ParseInt(dummyStatus, 0, 16)
		if err != nil {
			log.Warningf("converting %s, but error %s", dummyStatus, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(errorInvalidStatus)
			return
		}
	}

	// get data from db
	var dummyOne dummyModel
	if err := db.C(collectionDummy).FindId(bson.ObjectIdHex(dummyID)).One(&dummyOne); err != nil {
		if err == mgo.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(errorNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.WithField("error_msg", err.Error())
		return
	}

	if dummyOne.Headers != "" {
		byt := []byte(dummyOne.Headers)
		var dat map[string]string
		if err := json.Unmarshal(byt, &dat); err != nil {
			log.Errorf("JSON marshal error: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// traverse it
		for k, v := range dat {
			w.Header().Set(k, v)
		}
	}

	// set content type and charset
	w.Header().Set("Content-Type", dummyOne.ContentType+"; charset="+dummyOne.Charset)

	if dummyStatus == "" {
		w.WriteHeader(dummyOne.Status)
	} else {
		w.WriteHeader(int(convStatus))
	}

	w.Write([]byte(dummyOne.Content))
}

// content-type and charset is defined
// only JSON body is accepted
func handleV1CreateDummy(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Read body, custom headers, status, content-type from request
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var reqModel requestModel
	if err := json.NewDecoder(r.Body).Decode(&reqModel); err != nil {
		log.Warningf("fail to parse json %s ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{"InvalidJSON", "fail to parse JSON"})
		return
	}

	if err := reqModel.validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{"InvalidData", err.Error()})
		return
	}

	var dummyToSave dummyModel
	dummyToSave.updateWithRequestData(&reqModel)

	// save it to db
	dummyToSave.ID = bson.NewObjectId()
	if err := db.C(collectionDummy).Insert(&dummyToSave); err != nil {
		log.Error("error on saving an entity", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)

	type successResponse struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}

	json.NewEncoder(w).Encode(
		successResponse{
			ID:  dummyToSave.ID.Hex(),
			URL: "https://httpdummyresponser.herokuapp.com/" + apiVersion + "/" + dummyToSave.ID.Hex(),
		})
}
