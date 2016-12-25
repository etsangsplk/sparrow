package oauth

import (
	"fmt"
	logger "github.com/juju/loggo"
	"os"
	"testing"
)

var dbFilePath = "/tmp/oauth_silo_test.db"
var osl *OauthSilo

func TestMain(m *testing.M) {
	logger.ConfigureLoggers("<root>=warn;sparrow.oauth=debug")

	initSilo()

	// now run the tests
	m.Run()

	// cleanup
	os.Remove(dbFilePath)
}

func initSilo() {
	if osl != nil {
		osl.Close()
	}

	os.Remove(dbFilePath)

	var err error
	osl, err = Open(dbFilePath)

	if err != nil {
		fmt.Println("Failed to open oauth silo\n", err)
		os.Exit(1)
	}
}

func TestAddClient(t *testing.T) {
	cl := NewClient()

	osl.AddClient(cl)

	loaded, _ := osl.GetClient(cl.Id)
	if loaded == nil {
		t.Errorf("Failed to retrieve the client %s", cl.Id)
	}

	/*
		tx, _ := osl.db.Begin(false)
		bc := tx.Bucket(BUC_OAUTH_CLIENTS)
		cursor := bc.Cursor()

		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			fmt.Println(string(k))
		}
	*/
}
