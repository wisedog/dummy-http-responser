package main

import (
	"log"
	"os"
	"testing"
	"time"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testData [5]dummyModel

var testDB *mgo.Database

var _ = BeforeSuite(func() {
	log.Print("Setup")

	mongoDBURI := os.Getenv("MONGODB_URI")
	if mongoDBURI == "" {
		panic("mongoDB URI is empty")
	}

	// create data
	session, err := mgo.Dial(mongoDBURI)
	if err != nil {
		log.Fatal(err)
	}

	databaseName := os.Getenv("MONGODB_DATABASE")
	testDB = session.DB(databaseName)

	for i := range testData {
		testData[i] = dummyModel{
			Content:     `{"test":"hello"}`,
			Status:      200,
			Headers:     "",
			Charset:     "utf-8",
			ContentType: "application/json",
			Version:     "v1",
			CreatedAt:   time.Now(),
		}
		testData[i].ID = bson.NewObjectId()
		db.C(collectionDummy).Insert(&testData[i])
		log.Printf("item : %+v", testData[i])
	}
})

var _ = AfterSuite(func() {
	log.Print("TearDown")

	for i := range testData {
		db.C(collectionDummy).Remove(&testData[i])
	}
})

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Dummy Http Responser Suite")
}
