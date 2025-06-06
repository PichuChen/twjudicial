## 司法院裁判書開放 API 文件重點整理

### 1. 權限驗證
- 路徑：`https://data.judicial.gov.tw/jdg/api/Auth`
- 方法：POST
- Content-Type：application/json
- 輸入（JSON）：
  ```json
  {
    "user": "帳號",
    "password": "密碼"
  }
  ```
- 成功回傳：
  ```json
  { "token": "xxxxxx" }
  ```
- 失敗回傳：
  ```json
  { "error": "驗證失敗" }
  ```
- Token 有效 6 小時。

---
### 2. 取得裁判書異動清單
- 路徑：`https://data.judicial.gov.tw/jdg/api/JList`
- 方法：POST
- Content-Type：application/json
- 輸入（JSON）：
  ```json
  { "token": "xxxxxx" }
  ```
- 成功回傳（陣列，每天一筆）：
  ```json
  [
    {
      "DATE": "2016-12-23",
      "LIST": [
        "CDEV,105,橋司附民移調,101,20161219,1",
        ...
      ]
    },
    ...
  ]
  ```
- Token 錯誤會回傳「驗證失敗」。

---
### 3. 取得裁判書內容
- 路徑：`https://data.judicial.gov.tw/jdg/api/JDoc`
- 方法：POST
- Content-Type：application/json
- 輸入（JSON）：
  ```json
  {
    "token": "xxxxxx",
    "j": "裁判書ID"
  }
  ```
- 成功回傳（部分欄位）：
  ```json
  {
    "ATTACHMENTS": [
      { "TITLE": "xxx.pdf", "URL": "..." }
    ],
    "JFULLX": {
      "JFULLTYPE": "text" 或 "file",
      "JFULLCONTENT": "全文內容",
      "JFULLPDF": "PDF下載連結或空值"
    },
    "JID": "...",
    "JYEAR": "...",
    "JCASE": "...",
    "JNO": "...",
    "JDATE": "...",
    "JTITLE": "..."
  }
  ```
- 若資料被移除或未公開，會回傳：
  ```json
  { "error": "查無資料，本裁判可能未公開或已從系統移除..." }
  ```

---
### 其他注意事項
- API 僅於每日 0:00~6:00 提供服務。
- 若 jid 重複，表示同一筆裁判書，應以最新內容為主。
- 若裁判書被移除，應刪除本地資料。

---

司法院裁判書開放 API 規格說明
113.10.23
壹、 前言
提供取得司法院裁判書異動清單及裁判書全文之 API。
貳、 API 內容
一、驗證權限
(一)使用於「司法院資料開放平臺」申請之帳號密碼，經驗證通過後，回傳一組token，
此 token 有效時間為驗證通過後 6 小時，逾時需重新驗證。
(二)使用以下二組 api 時，需帶入前揭 token 以驗證使用權限。
二、取得裁判書異動清單
依據當日日期提供 7 日前裁判書異動清單，例如於 2017/10/16 呼叫本 API，則回傳
2017/10/9 異動的裁判書 ID（jid）。
三、取得裁判書內容
依據所輸入的 jid，提供該筆裁判書內容。
參、 一般說明
一、本 API 採 RESTful 型態，資料為 json 格式。
二、裁判書有附件或裁判全文為 pdf 檔時，系統回應下載該檔案的網址。
三、因考慮本院網路頻寬及系統負擔，並避免損及多數一般使用者權益，本 API 提供服
務時間為每日凌晨 0 時至 6 時止，其餘時間恕不提供服務。
肆、 裁判書異動注意事項
因裁判書資料上傳後仍可能有所異動或移除，使用者應遵循以下異動規則：
一、裁判書 jid 重複者，即為同一筆裁判書資料，若使用者曾於先前取得該裁判書資料，
之後再次以同一 jid 取得者，表示改筆裁判內容可能有所異動，應將後取得之內容覆
蓋先前取得之內容。
二、若系統回傳「{"error":"查無資料，本裁判可能未公開或已從系統移除，若您曾經下
載過本裁判，亦請您將其移除！謝謝！"}」，表示該筆裁判書業經本院移除或不再公
開，使用者應將先前所取得的裁判書內容刪除，以免損害當事人隱私及權益。
2
伍、 API 規格說明
一、驗證權限
功能說明 驗證是否有讀取本 API 之權限。
服務路徑 https://data.judicial.gov.tw/jdg/api/Auth
HTTP Method POST Content-Type application/json
輸
入
說
明
輸入內容 於 request body 中以 json 格式帶入使用者於「司法院資料開放平臺」申請
之帳號及密碼，欄位如下：
1. user：字串，帳號。
2. password：字串，密碼。
輸入範例 {
"user": "...略...",
"password": "...略..."
}
輸
出
說
明
輸出內容 若帳號密碼驗證通過，則回傳一組 token，若未通過，則回傳「驗證失敗」。
輸出範例 1. 驗證通過
{
"token": "ddf8bb4f32f746bdb5510c1eed76db51"
}
2. 驗證未通過
{
"error": "驗證失敗"
}
備註
3
二、取得裁判書異動清單
功能說明 取得 7 日前裁判書異動清單
服務路徑 https://data.judicial.gov.tw/jdg/api/JList
HTTP Method POST Content-Type application/json
輸
入
說
明
輸入內容 於 request body 中以 json 格式帶入第一組 API 所取得的 token，欄位如下：
1. token：字串，由第一組 API 所取得的 token。
輸入範例 {
"token": "ddf8bb4f32f746bdb5510c1eed76db51"
}
輸
出
說
明
輸出內容 1 每日異動清單的陣列，每個陣列項目中包括：
1.1 DATE: 異動日期
1.2 LIST: 異動的裁判書 ID（jid）陣列
輸出範例 [
{
"DATE": "2016-12-23",
"LIST":
[
"CDEV,105,橋司附民移調,101,20161219,1",
"CDEV,105,橋司附民移調,95,20161219,1",
"CDEV,105,橋司附民移調,98,20161219,1",……..
]
},
{
"DATE": "2016-12-24",
"LIST":
[
"CYDM,105,原訴,12,20161214,1",
"CYDM,105,朴簡,498,20161216,1",
"CYDM,105,易,509,20161216,1",……..
]
}
]
備註 1.若 token 驗證有誤，系統將回傳「驗證失敗」。
4
三、取得裁判書內容
功能說明 依據所輸入的 jid，提供該筆裁判書內容
服務路徑 https://data.judicial.gov.tw/jdg/api/JDoc
HTTP Method POST Content-Type application/json
輸
入
說
明
輸入內容 於 request body 中以 json 格式帶入第一組 API 所取得的 token，以及所要
取得的裁判書 id，欄位如下：
1. token：字串，由第一組 API 所取得的 token。
2. j：字串，要取得的裁判書 id
輸入範例 {
"token": "ddf8bb4f32f746bdb5510c1eed76db51",
"j": " CHDM,105,交訴,51,20161216,1"
}
輸
出
說
明
輸出內容 1 ATTACHMENTS: 裁判書附檔（多組）
1.1 TITLE: 檔案標題
1.2 URL: 下載網址
2 JFULLX: 裁判書全文
2.1 JFULLTYPE: 全文型態，含 text 及 file
2.2 JFULLCONTENT: 全文內容，本欄位為全文文字內容。
2.3 JFULLPDF: 全文內容檔案連結，若全文型態是 text，本欄位為空值，若型態
是 file，則本欄位為檔案下載 url。
3 JID: 裁判書 ID，係由法院別+裁判類別,年度,字別,號次,裁判日期,檢查單號組成，
裁判類別有民事(V)、刑事(M)、行政(A)、懲戒(P)、憲法(C)。
4 JYEAR: 年度
5 JCASE: 字別
6 JNO: 號次
7 JDATE: 裁判日期
8 JTITLE: 裁判案由
輸出範例 範例 1：全文內容為 text
{
"ATTACHMENTS":
[
{"TITLE":"附表三.pdf",
"URL":"
https://data.judicial.gov.tw/jdg/api/JFile/CHDM/100%2c%e8%a8%b4%2c1552%2c
1020517%2c2%2cCHDM1025H015_003.pdf"},
{"TITLE":"附表一.pdf",
"URL":"
https://data.judicial.gov.tw/jdg/api/JFile/CHDM/100%2c%e8%a8%b4%2c1552%2c
1020517%2c2%2cCHDM1025H015_001.pdf"}
],
5
"JFULLX":
{
"JFULLTYPE":"text",
"JFULLCONTENT":"臺灣彰化地方法院刑事判決 100 年度訴字第
1552 號\r\n 公 訴 人 臺灣彰化地方法院檢察署檢察官\r\n 被 告 林
麗勤\r\n 林美慧\r\n 共 同\r\n 選任辯護人 陳鎮律師\r\n
許富雄律師\r\n 被 告 林子庭\r\n 江怡芳\r\n
賴玉釵\r\n 曹仲凱\r\n 楊洪桂蘭\r\n 林
鄭喜美\r\n 蔣文魁\r\n 上列被告等因偽造文書等案件...",
"JFULLPDF":""
},
"JID":"CHDM,100,訴,1552,20130517,2",
"JYEAR":"100",
"JCASE":"訴",
"JNO":"1552",
"JDATE":"20130517",
"JTITLE":"偽造文書等"
}
範例 2：全文內容為檔案
{
"ATTACHMENTS":
[
{"TITLE":"110 毒抗 1212","URL":"
https://law.judicial.gov.tw/FJUD/GetFile.ashx?id=110%2c%e6%af%92%e6%8a%97%
2c1212%2c20210831%2c1%2cTPHH1108V_04Zpdf_001.pdf&tablename=TPHM&fi
lename=S45C-921090309230.pdf"}
],
"JFULLX":
{
"JFULLTYPE":"file",
"JFULLCONTENT":" \r\n 臺灣高等法院刑事裁定\r\n110 年度毒抗字第 1212 號
\r\n\r\n 抗 告 人\r\n\r\n 即 被 告 游嘉宏\r\n 上列抗告人因毒品危害防制條
例案件，不服臺灣基隆地方法院中華民國 110 年 6 月 15 日裁定（110 年度毒聲字
第 233 號），提起抗告，本院裁定如下：\r\n 主 文\r\n 抗告駁回。
\r\n 理 由\r\n 一、原裁定意旨如附件。\r\n 二、抗告意旨略以：抗告人即被
告游嘉宏（下稱抗告人）已另完成戒癮治療課程，又係因配合警方偵辦販毒案件，
為取信同夥而吸食毒品，且被告為專職之抓漏師，有正當職業，並須扶養兩名小
孩，以及有 200 多位客戶需抗告人提供保固服務，為此，請求撤銷原裁定，准予
易科罰金或勞動服務等語。\r\n\r\n 三、按施用第一、二級毒品為犯罪行為，毒品
危害防制條例第 10 條設有處罰規定。惟基於刑事政策，對合於一定條件之施用
6
者，依同條例第 20 條之規定，施以觀察、勒戒及強制戒治之保安處分。於民國
109 年 1 月 15 日修正公布、同年 7 月 15 日施行之毒品危害防制條例第 20 條第 3
項及第 23 條第 2 項，將原條文之再犯期間由「5 年」改為「3 年」，而第 20 條第
3 項規定中所謂「3 年後再犯」，只要本次再犯（不論修正施行前、後）距最近 1
次觀察、勒戒或強制戒治執行完畢釋放，已逾 3 年者，即該當之，不因其間有無
犯第 10 條之罪經起訴、判刑...",
"JFULLPDF": "https://data.judicial.gov.tw/jdg/api/JDocFile/TPHM /110%2c 毒抗
%2c1212%2c20210831%2c1.pdf"
},
"JID":"TPHM,110,毒抗,1212,20210831,1",
"JYEAR":"110",
"JCASE":"毒抗",
"JNO":"1212",
"JDATE":"20210831",
"JTITLE":"毒品危害防制條例"
}
備註 1 jid 為裁判書的 pkey，若 jid 相同，代表是同一筆裁判書，若該筆裁判
jid 出現在不同日期異動清單中，代表其內容可能有異動，應該以後日
期的內容覆蓋前日期的內容。
2 裁判書上傳後也可能變更為不可公開或從系統移除，這時系統會回應
一個 error message:
{"error":"查無資料，本裁判可能未公開或已從系統移除，若您曾經下載過
本裁判，亦請您將其移除！謝謝！"}
這種情況下本 API 使用者應該將先前已取得的裁判書刪除，以免損及
當事人隱私或權益。
3 若 token 驗證有誤，系統將回傳「驗證失敗」。