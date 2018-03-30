package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/julienschmidt/httprouter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Echo test", func() {
	Context("Test GET ", func() {
		It("should response to GET request without body with hdr application/json", func() {
			req, _ := http.NewRequest("GET", "/echo", nil)
			req.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r := createRoute()
			r.ServeHTTP(w, req)
			contentType := w.Header().Get("Content-Type")
			Expect(contentType).To(Equal("application/json"))
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(Equal(""))
		})

		It("should response to GET request with body and hdr application/xml", func() {
			buff := bytes.NewBufferString("blahblah")
			req, _ := http.NewRequest("GET", "/echo", buff)
			req.Header.Add("Content-Type", "application/xml")
			w := httptest.NewRecorder()
			r := createRoute()
			r.ServeHTTP(w, req)
			contentType := w.Header().Get("Content-Type")
			Expect(contentType).To(Equal("application/xml"))
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(Equal("blahblah"))
		})

		It("should response to POST request with body and a header", func() {
			buff := bytes.NewBufferString("blahblah")
			req, _ := http.NewRequest("GET", "/echo?echo-hdr=X-API-AUTH", buff)
			req.Header.Add("X-API-AUTH", "blahblah")
			req.Header.Add("Content-Type", "application/xml")
			w := httptest.NewRecorder()
			r := createRoute()
			r.ServeHTTP(w, req)
			contentType := w.Header().Get("Content-Type")
			customHdr := w.Header().Get("X-API-AUTH")
			Expect(contentType).To(Equal("application/xml"))
			Expect(customHdr).To(Equal("blahblah"))
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(Equal("blahblah"))
		})

		It("should response to POST request with body and multiple headers ", func() {
			buff := bytes.NewBufferString("blahblah")
			req, _ := http.NewRequest("GET", "/echo?echo-hdr=X-API-AUTH,X-AAA", buff)
			req.Header.Add("X-API-AUTH", "blahblah")
			req.Header.Add("X-AAA", "somebody")
			req.Header.Add("Content-Type", "application/xml")
			w := httptest.NewRecorder()
			r := createRoute()
			r.ServeHTTP(w, req)

			Expect("application/xml").To(Equal(w.Header().Get("Content-Type")))
			Expect("blahblah").To(Equal(w.Header().Get("X-API-AUTH")))
			Expect("somebody").To(Equal(w.Header().Get("X-AAA")))
			Expect(http.StatusOK).To(Equal(w.Code))
			Expect("blahblah").To(Equal(w.Body.String()))
		})
	})
	Context("test POST", func() {
		It("should response to POST request with body and hdr application/xml", func() {
			buff := bytes.NewBufferString("blahblah")
			r, _ := http.NewRequest("POST", "/echo", buff)
			r.Header.Add("Content-Type", "application/xml")
			w := httptest.NewRecorder()
			var k []httprouter.Param
			handleEcho(w, r, k)
			contentType := w.Header().Get("Content-Type")
			Expect(contentType).To(Equal("application/xml"))
			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(w.Body.String()).To(Equal("blahblah"))
		})

		It("should response 400 when dummy-status param is 400", func() {
			req, _ := http.NewRequest("POST", "/echo?dummy-status=400", nil)
			req.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r := createRoute()
			r.ServeHTTP(w, req)
			contentType := w.Header().Get("Content-Type")
			Expect(contentType).To(Equal("application/json"))
			Expect(w.Code).To(Equal(http.StatusBadRequest))
		})
	})
})
