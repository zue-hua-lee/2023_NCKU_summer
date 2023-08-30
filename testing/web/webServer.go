/*
*

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
			// 回傳參數: 每個廠回傳充電樁編號、承諾SoC、單價、佔位價
			var optionA controller.Option
			var optionB controller.Option
			var optionC controller.Option
			// A場
			cmd := exec.Command("python", "../COMMERCIAL/call_commercial_newEV.py", i, offers[num].CarNum, offers[num].ArrTime, offers[num].DepTime, offers[num].ArrSoC, offers[num].DepSoC, offers[num].Capacity, offers[num].Acdc, offers[num].Location_x, offers[num].Location_y)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("執行失敗：", err)
				return
			}
			matrixStr := string(output)
			matrixStr = strings.TrimSpace(matrixStr)
			matrixValues1 := strings.Split(matrixStr, " ")
			optionA.StationID = "A"
			for j, value := range matrixValues1 {
				checkvalue, err := strconv.Atoi(value)
				if err != nil {
					fmt.Printf("無法將字符轉換成整數: %v\n", err)
					continue
				}
				switch j {
				case 1:
					optionA.ChargerID := checkvalue
				case 2:
					optionA.MaxSoC := checkvalue
				case 3:
					optionA.PerPrice := checkvalue
				case 4:
					optionA.ParkPrice := checkvalue
				case 5:
					optionA.TolPrice := checkvalue
				}
			}
			// B場
			cmd := exec.Command("python", "../COMMERCIAL/call_commercial_newEV.py", i, offers[num].CarNum, offers[num].ArrTime, offers[num].DepTime, offers[num].ArrSoC, offers[num].DepSoC, offers[num].Capacity, offers[num].Acdc, offers[num].Location_x, offers[num].Location_y)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("執行失敗：", err)
				return
			}
			matrixStr := string(output)
			matrixStr = strings.TrimSpace(matrixStr)
			matrixValues1 := strings.Split(matrixStr, " ")
			optionB.StationID = "B"
			for j, value := range matrixValues1 {
				checkvalue, err := strconv.Atoi(value)
				if err != nil {
					fmt.Printf("無法將字符轉換成整數: %v\n", err)
					continue
				}
				switch j {
				case 1:
					optionB.ChargerID := checkvalue
				case 2:
					optionB.MaxSoC := checkvalue
				case 3:
					optionB.PerPrice := checkvalue
				case 4:
					optionB.ParkPrice := checkvalue
				case 5:
					optionB.TolPrice := checkvalue
				}
			}
			// C場
			cmd := exec.Command("python", "../COMMERCIAL/call_commercial_newEV.py", i, offers[num].CarNum, offers[num].ArrTime, offers[num].DepTime, offers[num].ArrSoC, offers[num].DepSoC, offers[num].Capacity, offers[num].Acdc, offers[num].Location_x, offers[num].Location_y)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("執行失敗：", err)
				return
			}
			matrixStr := string(output)
			matrixStr = strings.TrimSpace(matrixStr)
			matrixValues1 := strings.Split(matrixStr, " ")
			optionC.StationID = "C"
			for j, value := range matrixValues1 {
				checkvalue, err := strconv.Atoi(value)
				if err != nil {
					fmt.Printf("無法將字符轉換成整數: %v\n", err)
					continue
				}
				switch j {
				case 1:
					optionC.ChargerID := checkvalue
				case 2:
					optionC.MaxSoC := checkvalue
				case 3:
					optionC.PerPrice := checkvalue
				case 4:
					optionC.ParkPrice := checkvalue
				case 5:
					optionC.TolPrice := checkvalue
				}
			}
			options := []controller.Option{optionA, optionB, optionC}

			// 媒合階段
			// 呼叫函式2
			// 輸入參數: 時段、(1/-1)
			// 回傳參數:

			option := app.ChooseOption(options)
			fmt.Printf("station: %s, tolPrice: %d\n", option.StationID, option.TolPrice)
			app.SetMatch(option)
		}

		// 5分鐘最佳化排程
		// 呼叫最佳化函式3
		// 輸入參數: 時段、(1/-1之外)
		// 回傳參數: 每個樁的功率、目前場內汽車的SoC

		// 將各個充電樁的功率上鏈
		fmt.Printf("第%d區間上鍊開始\n",i)
		var powers []controller.Power
		for j := 1; j <= 12; j++{
			powers = append(powers, controller.Power{StationID: "A", ChargerID: j, Power: 0, State: 0, TimeStamp: i})
		}
		for j := 1; j <= 6; j++{
			powers = append(powers, controller.Power{StationID: "B", ChargerID: j, Power: 30, State: 1, TimeStamp: i})
		}
		for j := 1; j <= 20; j++{
			powers = append(powers, controller.Power{StationID: "C", ChargerID: j, Power: 40, State: 1, TimeStamp: i})
		}
		app.SetPower(powers)
		fmt.Printf("第%d區間上鍊結束\n",i)
	}
	fmt.Println("\n\nOffer")
	app.ShowAllOffer()
	fmt.Println("\n\nMatch")
	app.ShowAllMatch()
	fmt.Println("\n\nPower")
	app.ShowAllPower()
}
