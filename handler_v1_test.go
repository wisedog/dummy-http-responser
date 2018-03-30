package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/globalsign/mgo/bson"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("Handler V1", func() {
	Context("with valid hexify url", func() {
		It("should response 200 when existing url ", func() {
			log.Infof("id: %s", testData[0].ID.Hex())
			req, _ := http.NewRequest("GET", "/v1/"+testData[0].ID.Hex(), nil)
			req.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r := createRoute()
			r.ServeHTTP(w, req)
			contentType := w.Header().Get("Content-Type")

			Expect(contentType).To(Equal("application/json; charset=utf-8"))
			Expect(w.Code).To(Equal(http.StatusOK))
		})

		It("should response 404 when not exist url", func() {
			req, _ := http.NewRequest("GET", "/v1/152ab3829d918f9", nil)
			req.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r := createRoute()
			r.ServeHTTP(w, req)
			contentType := w.Header().Get("Content-Type")
			Expect(contentType).To(Equal("application/json; charset=utf-8"))
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})

	})

	Context("with invalid hexify", func() {
		It("should response 400", func() {
			req, _ := http.NewRequest("GET", "/v1/ksjnfkwjenfkjwen", nil)
			req.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r := createRoute()
			r.ServeHTTP(w, req)
			contentType := w.Header().Get("Content-Type")
			Expect("application/json; charset=utf-8").To(Equal(contentType))
			Expect(http.StatusBadRequest).To(Equal(w.Code))
		})
	})

	Context("with creating dummy url", func() {
		It("should response 400 with no status value", func() {
			dat := `{
				"content":      "blahblah",
				"content_type": "application/json",
				"charset":      "utf-8",
				"headers":      {}
			}`

			buf := bytes.NewBufferString(dat)

			req, _ := http.NewRequest("POST", "/create", buf)
			req.Header.Add("Content-Type", "application/json; charset=utf-8")
			w := httptest.NewRecorder()
			r := createRoute()
			r.ServeHTTP(w, req)
			contentType := w.Header().Get("Content-Type")
			Expect(http.StatusBadRequest).To(Equal(w.Code))
			Expect("application/json; charset=utf-8").To(Equal(contentType))
		})

		It("should response 400 with no charset value", func() {
			dat := `{
				"content":      "blahblah",
				"content_type": "application/json",
				"status": 200,
				"headers":      {}
			}`

			buf := bytes.NewBufferString(dat)

			req, _ := http.NewRequest("POST", "/create", buf)
			req.Header.Add("Content-Type", "application/json; charset=utf-8")
			w := httptest.NewRecorder()
			r := createRoute()
			r.ServeHTTP(w, req)
			contentType := w.Header().Get("Content-Type")
			Expect(http.StatusBadRequest).To(Equal(w.Code))
			Expect("application/json; charset=utf-8").To(Equal(contentType))

			byt := w.Body.Bytes()
			w.Body.Read(byt)
			var resp map[string]interface{}
			if err := json.Unmarshal(byt, &resp); err != nil {
				panic("")
			}
			Expect("InvalidData").To(Equal(resp["error"].(string)))
			Expect("charset is empty").To(Equal(resp["error_msg"].(string)))
		})

		It("should response 200", func() {
			dat := `{
				"content":      "blahblah",
				"content_type": "application/json",
				"status": 200,
				"charset":      "utf-8",
				"headers":      {}
			}`

			buf := bytes.NewBufferString(dat)
			req, _ := http.NewRequest("POST", "/create", buf)
			req.Header.Add("Content-Type", "application/json; charset=utf-8")
			w := httptest.NewRecorder()
			r := createRoute()
			r.ServeHTTP(w, req)
			contentType := w.Header().Get("Content-Type")
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(contentType).To(Equal("application/json; charset=utf-8"))

			byt := w.Body.Bytes()
			w.Body.Read(byt)
			var resp map[string]interface{}
			if err := json.Unmarshal(byt, &resp); err != nil {
				panic("")
			}
			Expect("").NotTo(Equal(resp["id"].(string)))
			Expect(strings.HasPrefix(resp["url"].(string), "https://httpdummyresponser.herokuapp.com")).To(Equal(true))

			// should delete test data
			if err := testDB.C(collectionDummy).RemoveId(bson.ObjectIdHex(resp["id"].(string))); err != nil {
				panic(err.Error())
			}
		})
		It("should response 200 with custom headers", func() {
			dat := `{
				"content":      "blahblah",
				"content_type": "application/json",
				"status": 200,
				"charset":      "utf-8",
				"headers":      {"X-TEST-XXX": "ASDF", "BLAHBLAH": "AAA"}
			}`

			buf := bytes.NewBufferString(dat)
			req, _ := http.NewRequest("POST", "/create", buf)
			req.Header.Add("Content-Type", "application/json; charset=utf-8")
			w := httptest.NewRecorder()
			r := createRoute()
			r.ServeHTTP(w, req)
			contentType := w.Header().Get("Content-Type")
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(contentType).To(Equal("application/json; charset=utf-8"))

			byt := w.Body.Bytes()
			w.Body.Read(byt)
			var resp map[string]interface{}
			if err := json.Unmarshal(byt, &resp); err != nil {
				panic("")
			}
			log.Info(string(byt))
			Expect("").NotTo(Equal(resp["id"].(string)))
			Expect(true).To(Equal(strings.HasPrefix(resp["url"].(string), "https://httpdummyresponser.herokuapp.com")))

			// should delete test data
			if err := testDB.C(collectionDummy).RemoveId(bson.ObjectIdHex(resp["id"].(string))); err != nil {
				panic(err.Error())
			}
		})
	})

	Context("with created dummy", func() {
		It("should response 200 ", func() {
			// create a record to test it
			dat := `{
				"content":      "{\"blahblah\": \"aaa\"}",
				"content_type": "application/json",
				"status": 200,
				"charset":      "utf-8",
				"headers":      {}
			}`

			buf := bytes.NewBufferString(dat)
			req, _ := http.NewRequest("POST", "/create", buf)
			req.Header.Add("Content-Type", "application/json; charset=utf-8")
			w := httptest.NewRecorder()
			r := createRoute()
			r.ServeHTTP(w, req)
			contentType := w.Header().Get("Content-Type")
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(contentType).To(Equal("application/json; charset=utf-8"))

			byt := w.Body.Bytes()
			w.Body.Read(byt)
			var resp map[string]interface{}
			if err := json.Unmarshal(byt, &resp); err != nil {
				panic("")
			}
			Expect("").NotTo(Equal(resp["id"].(string)))
			Expect(strings.HasPrefix(resp["url"].(string), "https://httpdummyresponser.herokuapp.com")).To(Equal(true))

			req, _ = http.NewRequest("GET", "/"+apiVersion+"/"+resp["id"].(string), nil)
			req.Header.Add("Content-Type", "application/json")
			w = httptest.NewRecorder()
			r = createRoute()
			r.ServeHTTP(w, req)
			contentType = w.Header().Get("Content-Type")
			Expect("application/json; charset=utf-8").To(Equal(contentType))
			Expect(w.Code).To(Equal(http.StatusOK))

			byt = w.Body.Bytes()
			w.Body.Read(byt)
			var resp1 map[string]string
			if err := json.Unmarshal(byt, &resp1); err != nil {
				panic("")
			}
			Expect(map[string]string{"blahblah": "aaa"}).To(Equal(resp1))

			// should delete test data
			if err := testDB.C(collectionDummy).RemoveId(bson.ObjectIdHex(resp["id"].(string))); err != nil {
				panic(err.Error())
			}

		})
	})
})
