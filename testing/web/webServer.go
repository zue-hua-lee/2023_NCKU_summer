/**
  author: jerry
 */
package web

import (
	"fmt"
	"github.com/hyperledger/fabric/eep/web/controller"
)

func  WebStart(app *controller.Application) {
	fmt.Println("abcdefg")
	app.Testing()

	// 讀取新充電請求(csv)按時間排序存到陣列中
	offers := app.LoadAllOffer()
	var num int = 0
	for i := 1; i <= 288; i++ {
		for ; num < len(offers) && offers[num].ArrTime == i; num++ {
			// 新充電申請上鏈
			fmt.Printf("time: %d, num: %d, EV: %d\n", i, num, offers[num].CarNum)
			app.SetOffer(offers[num])

			// 呼叫最佳化函式1 (先慢1後快2) 
			// 輸入參數: 新車資訊、時段
			// 回傳參數: 每個廠回傳充電樁編號、承諾SoC、單價
			optionA := controller.Option{StationID: "A", ChargerID: 1, MaxSoC: 100, Price: 100,}
			optionB := controller.Option{StationID: "B", ChargerID: 1, MaxSoC: 100, Price: 50,}
			optionC := controller.Option{StationID: "C", ChargerID: 1, MaxSoC: 100, Price: 30,}
			options := []controller.Option{optionA, optionB, optionC}

			// 媒合階段
			// 呼叫函式2
			// 輸入參數: 車子編號、選擇的充電站的樁位
			// 回傳參數:

			option := app.ChooseOption(options)
			fmt.Printf("station: %s, price: %d\n", option.StationID, option.Price)
			app.SetMatch(option)
		}

		// 5分鐘最佳化排程
		// 呼叫最佳化函式3
		// 輸入參數: 時段
		// 回傳參數: 每個樁的功率、目前場內汽車的SoC

		// 將各個充電裝的功率上鏈
		// fmt.Printf("第%d區間上鍊開始\n",i)
		// for j := 1; j <= 12; j++{
		// 	app.SetPower("A", j, 0, 0, i)
		// }
		// for j := 1; j <= 6; j++{
		// 	app.SetPower("B", j, 30, 1, i)
		// }
		// for j := 1; j <= 20; j++{
		// 	app.SetPower("C", j, 40, 1, i)
		// }
		// fmt.Printf("第%d區間上鍊結束\n",i)
	}
	fmt.Println("\n\nOffer")
	app.ShowAllOffer()
	fmt.Println("\n\nMatch")
	app.ShowAllMatch()
	fmt.Println("\n\nPower")
	app.ShowAllPower()
}
