package twjudicial

import (
	"os"
	"testing"
	"time"
)

// 測試未設定環境變數時，測試應該被 Skip
func TestEnvVarsRequired(t *testing.T) {
	origUser := os.Getenv("JUDICIAL_USER")
	origPass := os.Getenv("JUDICIAL_PASSWORD")
	os.Unsetenv("JUDICIAL_USER")
	os.Unsetenv("JUDICIAL_PASSWORD")
	defer func() {
		os.Setenv("JUDICIAL_USER", origUser)
		os.Setenv("JUDICIAL_PASSWORD", origPass)
	}()

	user := os.Getenv("JUDICIAL_USER")
	password := os.Getenv("JUDICIAL_PASSWORD")
	if user != "" || password != "" {
		t.Fatal("JUDICIAL_USER 與 JUDICIAL_PASSWORD 應為空")
	}
	// 呼叫 Auth 前檢查
	if user == "" || password == "" {
		t.Skip("需設定 JUDICIAL_USER 與 JUDICIAL_PASSWORD 環境變數")
	}
}

// 檢查服務時間（司法院 API 僅於每日 00:00~06:00 提供服務）
func TestServiceHours(t *testing.T) {
	now := time.Now()
	hour := now.Hour()
	if hour >= 6 {
		t.Skipf("目前非服務時間（00:00~06:00），現在時間：%02d:%02d", hour, now.Minute())
	}
}

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

// 測試 Auth 錯誤帳密
func TestAuthWrongPassword(t *testing.T) {
	user := "not_exist_user"
	password := "wrong_password"
	token, err := Auth(user, password)
	if err == nil {
		t.Fatalf("Auth 應該失敗，但卻成功，token=%s", token)
	}
}
