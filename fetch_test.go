package twjudicial

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

// skipIfNotServiceTime skips the test when current time is outside 00:00~06:00.
func skipIfNotServiceTime(t *testing.T) {
	t.Helper()
	now := time.Now()
	if h := now.Hour(); h >= 6 {
		t.Skipf("目前非服務時間（00:00~06:00），現在時間：%02d:%02d", h, now.Minute())
	}
}

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
	skipIfNotServiceTime(t)
}

func TestAuth(t *testing.T) {
	skipIfNotServiceTime(t)
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
	skipIfNotServiceTime(t)
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
	skipIfNotServiceTime(t)
	user := "not_exist_user"
	password := "wrong_password"
	token, err := Auth(user, password)
	if err == nil {
		t.Fatalf("Auth 應該失敗，但卻成功，token=%s", token)
	}
}

// roundTripFunc allows custom HTTP responses in tests.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// TestNonOKStatus ensures functions return errors on non-200 responses.
func TestNonOKStatus(t *testing.T) {
	rt := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(bytes.NewBuffer(nil)),
			Header:     make(http.Header),
		}, nil
	})
	origClient := httpClient
	httpClient = &http.Client{Transport: rt}
	defer func() { httpClient = origClient }()
		t.Fatalf("Auth should fail on non-200 status, got %v", err)
	}
	if _, err := GetJListWithClient(client, "token"); err == nil || !strings.Contains(err.Error(), "status") {
		t.Fatalf("GetJList should fail on non-200 status, got %v", err)
	}
	if _, err := GetJDocWithClient(client, "token", "jid"); err == nil || !strings.Contains(err.Error(), "status") {
		t.Fatalf("GetJDoc should fail on non-200 status, got %v", err)
	}
}
