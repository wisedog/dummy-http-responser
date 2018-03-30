package main

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/globalsign/mgo/bson"
)

var (
	collectionDummy = "dummy"
)

type requestModel struct {
	Content     string            `json:"content"`      // body to response
	Charset     string            `json:"charset"`      // charset
	ContentType string            `json:"content_type"` // http 'Content-Type'
	Status      int               `json:"status"`       // http status
	Headers     map[string]string `json:"headers"`
}

// validate requestModel. do not trust any input
func (m *requestModel) validate() error {
	var err error
	if m.Status == 0 {
		err = errors.New("status is not set")
	} else if m.ContentType == "" {
		err = errors.New("content type is empty")
	} else if m.Charset == "" {
		err = errors.New("charset is empty")
	}
	return err
}

// dummyModel is a model for manipulating databases' data
type dummyModel struct {
	ID          bson.ObjectId `bson:"_id"`
	Version     string        // API version. v1, v2 ... vn
	Content     string        // body to response
	Charset     string        // charset
	ContentType string        // http 'Content-Type'
	Headers     string        // stringify JSON
	Status      int           // http status
	CreatedAt   time.Time     // Time to created this record
}

func (d *dummyModel) updateWithRequestData(m *requestModel) error {
	d.Content = m.Content
	d.Charset = m.Charset
	d.ContentType = m.ContentType
	d.Status = m.Status
	d.CreatedAt = time.Now()
	d.Version = apiVersion
	// convert map to JSON
	jsonBytes, err := json.Marshal(m.Headers)
	if err != nil {
		return err
	}
	// stringify it
	d.Headers = string(jsonBytes)
	return nil
}
