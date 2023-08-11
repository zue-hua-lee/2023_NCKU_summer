/**
  author: jerry
 */
package web

import (
	// "net/http"
	"fmt"
	"github.com/hyperledger/fabric/eep/web/controller"
)

func  WebStart(app *controller.Application)  {
	fmt.Println("abcdefg")
	app.Testing()

	for i := 1; i <= 288; i++ {
        
		// 讀取新充電請求: csv
		// 上鏈

		// 呼叫最佳化函式1 (先慢1後快2) 
		// 輸入參數: 新車資訊、時段
		// 回傳參數: 每個廠回傳單價、承諾SoC、哪一隻充電樁

		// 媒合階段
		// 呼叫函式2
		// 輸入參數: 車子編號、選擇的充電站的樁位
		// 回傳參數:

		// 5分鐘最佳化排程
		// 呼叫最佳化函式3
		// 輸入參數: 時段
		// 回傳參數: 每個樁的功率、目前場內汽車的SoC

		// 將各個充電裝的功率上鏈

	}
}
