package twjudicial

import (
	"os"
	"testing"
)

func TestAuth(t *testing.T) {
	user := os.Getenv("JUDICIAL_USER")
	password := os.Getenv("JUDICIAL_PASSWORD")
	if user == "" || password == "" {
		t.Skip("需設定 JUDICIAL_USER 與 JUDICIAL_PASSWORD 環境變數")
	}
	token, err := Auth(user, password)
	if err != nil {
		t.Fatalf("Auth failed: %v", err)
	}
	if token == "" {
		t.Fatal("token is empty")
	}
}

func TestJListAndJDoc(t *testing.T) {
	user := os.Getenv("JUDICIAL_USER")
	password := os.Getenv("JUDICIAL_PASSWORD")
	if user == "" || password == "" {
		t.Skip("需設定 JUDICIAL_USER 與 JUDICIAL_PASSWORD 環境變數")
	}
	token, err := Auth(user, password)
	if err != nil {
		t.Fatalf("Auth failed: %v", err)
	}
	jlist, err := GetJList(token)
	if err != nil {
		t.Fatalf("GetJList failed: %v", err)
	}
	if len(jlist) == 0 || len(jlist[0].List) == 0 {
		t.Skip("JList 無資料，無法測試 JDoc")
	}
	// 測試用固定 JID
	jid := "TPBA,113,訴,501,20240808,1"
	jdoc, err := GetJDoc(token, jid)
	if err != nil {
		t.Fatalf("GetJDoc failed for jid=%s: %v", jid, err)
	}
	if jdoc.JID == "" {
		t.Errorf("JDocResponse JID is empty for jid=%s", jid)
	}
}
