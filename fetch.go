// Package twjudicial 提供司法院裁判書開放 API 的 Go 語言存取函式庫。
//
// 主要功能：
//  1. 權限驗證（Auth）：取得 API token
//  2. 取得裁判書異動清單（GetJList）
//  3. 取得裁判書內容（GetJDoc）
//
// 範例：
//
//	package main
//	import (
//	    "fmt"
//	    "os"
//	    "twjudicial"
//	)
//	func main() {
//	    token, err := twjudicial.Auth(os.Getenv("JUDICIAL_USER"), os.Getenv("JUDICIAL_PASSWORD"))
//	    if err != nil {
//	        panic(err)
//	    }
//	    jlist, err := twjudicial.GetJList(token)
//	    if err != nil {
//	        panic(err)
//	    }
//	    fmt.Println("第一天異動清單：", jlist[0].Date, jlist[0].List)
//	    jdoc, err := twjudicial.GetJDoc(token, "TPBA,113,訴,501,20240808,1")
//	    if err != nil {
//	        panic(err)
//	    }
//	    fmt.Println("判決標題：", jdoc.JTitle)
//	}
//
// 詳細 API 文件請參閱 docs/api.md。
package twjudicial

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

const (
	apiAuth  = "https://data.judicial.gov.tw/jdg/api/Auth"
	apiJList = "https://data.judicial.gov.tw/jdg/api/JList"
	apiJDoc  = "https://data.judicial.gov.tw/jdg/api/JDoc"
)

// AuthRequest 用於登入取得 token
type AuthRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

// AuthResponse 回傳 token 或 error
type AuthResponse struct {
	Token string `json:"token"`
	Error string `json:"error"`
}

// JListRequest 用於取得異動清單
type JListRequest struct {
	Token string `json:"token"`
}

// JListResponseItem 單日異動清單
type JListResponseItem struct {
	Date string   `json:"DATE"`
	List []string `json:"LIST"`
}

// JDocRequest 用於取得裁判書內容
type JDocRequest struct {
	Token string `json:"token"`
	J     string `json:"j"`
}

// Attachment 裁判書附件
type Attachment struct {
	Title string `json:"TITLE"`
	URL   string `json:"URL"`
}

// JFullX 裁判書全文
type JFullX struct {
	JFullType    string `json:"JFULLTYPE"`
	JFullContent string `json:"JFULLCONTENT"`
	JFullPDF     string `json:"JFULLPDF"`
}

// JDocResponse 裁判書內容
type JDocResponse struct {
	Attachments []Attachment `json:"ATTACHMENTS"`
	JFullX      JFullX       `json:"JFULLX"`
	JID         string       `json:"JID"`
	JYear       string       `json:"JYEAR"`
	JCase       string       `json:"JCASE"`
	JNo         string       `json:"JNO"`
	JDate       string       `json:"JDATE"`
	JTitle      string       `json:"JTITLE"`
	Error       string       `json:"error"`
}

// Auth 取得 token
func Auth(user, password string) (string, error) {
	reqBody, _ := json.Marshal(AuthRequest{User: user, Password: password})
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(apiAuth, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result AuthResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if result.Token != "" {
		return result.Token, nil
	}
	return "", errors.New(result.Error)
}

// GetJList 取得裁判書異動清單
func GetJList(token string) ([]JListResponseItem, error) {
	reqBody, _ := json.Marshal(JListRequest{Token: token})
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(apiJList, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if bytes.Contains(body, []byte("驗證失敗")) {
		return nil, errors.New("驗證失敗")
	}
	var result []JListResponseItem
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetJDoc 取得裁判書內容
func GetJDoc(token, jid string) (*JDocResponse, error) {
	reqBody, _ := json.Marshal(JDocRequest{Token: token, J: jid})
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(apiJDoc, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result JDocResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if result.Error != "" {
		return nil, errors.New(result.Error)
	}
	return &result, nil
}
