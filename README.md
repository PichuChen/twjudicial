# Judicial Open Data Library

This Go library provides functionality to fetch and parse open judicial data from the Judicial Yuan of Taiwan. It allows users to easily retrieve and work with judicial decisions.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

## Installation

To install the library, use the following command:

```
go get github.com/PichuChen/twjudicial
```

Replace `PichuChen` with your GitHub username or the appropriate path where the library is hosted.

## Usage

> **注意：司法院 API 僅於每日凌晨 0:00~6:00 提供服務，其餘時間無法存取。**
>
> 請先於系統環境變數設定 JUDICIAL_USER 與 JUDICIAL_PASSWORD，否則無法通過驗證。
>
> 設定方式範例（Linux/macOS）：
> ```bash
> export JUDICIAL_USER=你的帳號
> export JUDICIAL_PASSWORD=你的密碼
> ```

To use the library, import it in your Go application:

```go
import "github.com/PichuChen/twjudicial"
```

### Example: Authenticate and Fetch Judicial Data

```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/PichuChen/twjudicial"
)

func main() {
    // 1. 取得 token
    token, err := twjudicial.Auth(os.Getenv("JUDICIAL_USER"), os.Getenv("JUDICIAL_PASSWORD"))
    if err != nil {
        log.Fatal("Auth failed:", err)
    }
    // 2. 取得異動清單
    jlist, err := twjudicial.GetJList(token)
    if err != nil {
        log.Fatal("GetJList failed:", err)
    }
    fmt.Println("第一天異動清單：", jlist[0].Date, jlist[0].List)
    // 3. 取得特定裁判書內容
    jdoc, err := twjudicial.GetJDoc(token, "TPBA,113,訴,501,20240808,1")
    if err != nil {
        log.Fatal("GetJDoc failed:", err)
    }
    fmt.Println("判決標題：", jdoc.JTitle)
}
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
