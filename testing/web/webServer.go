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

			checknowtime := strconv.Itoa(i)
			checkcarnum := strconv.Itoa(offers[num].CarNum)
			checkarrtime := strconv.Itoa(offers[num].ArrTime)
			checkdeptime := strconv.Itoa(offers[num].DepTime)
			checkarrsoc := strconv.Itoa(offers[num].ArrSoC)
			checkdepsoc := strconv.Itoa(offers[num].DepSoC)
			checkcapacity := strconv.Itoa(offers[num].Capacity)
			checkacdc := strconv.Itoa(offers[num].Acdc)
			checklocation_x := strconv.Itoa(offers[num].Location_x)
			checklocation_y := strconv.Itoa(offers[num].Location_y)

			// A場最佳化函式1
			cmd := exec.Command("python", "../ROAD/call_road_newEV.py", checknowtime, checkcarnum, checkarrtime, checkdeptime, checkarrsoc, checkdepsoc, checkcapacity, checkacdc, checklocation_x, checklocation_y)
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
			// B場最佳化函式1
			cmd := exec.Command("python", "../FCS/call_FCS_newEV.py", checknowtime, checkcarnum, checkarrtime, checkdeptime, checkarrsoc, checkdepsoc, checkcapacity, checkacdc, checklocation_x, checklocation_y)
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
			// C場最佳化函式1
			cmd := exec.Command("python", "../COMMERCIAL/call_commercial_newEV.py", checknowtime, checkcarnum, checkarrtime, checkdeptime, checkarrsoc, checkdepsoc, checkcapacity, checkacdc, checklocation_x, checklocation_y)
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
					return
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
			app.SetMatch(option)

			fmt.Printf("station: %s, tolPrice: %d\n", option.StationID, option.TolPrice)

			switch option.StationID {
			case "A":
				choice := [3]int{1, -1, -1}
			case "B":
				choice := [3]int{-1, 1, -1}
			case "C":
				choice := [3]int{-1, -1, 1}
			}

			// A場最佳化函式2
			cmd := exec.Command("python", "../ROAD/call_road.py", strconv.Itoa(i), choice[0])
			_, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("執行失敗：", err)
				return
			}
			// B場最佳化函式2
			cmd := exec.Command("python", "../FCS/call_FCS.py", strconv.Itoa(i), choice[1])
			_, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("執行失敗：", err)
				return
			}
			// C場最佳化函式3
			cmd := exec.Command("python", "../COMMERCIAL/call_commercial.py", strconv.Itoa(i), choice[2])
			_, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("執行失敗：", err)
				return
			}
		}

		// 5分鐘最佳化排程
		// 呼叫最佳化函式3
		// 輸入參數: 時段、(1/-1之外)
		// 回傳參數: 每個樁的功率、目前場內汽車的SoC

		// A場最佳化函式3
		cmd := exec.Command("python", "../ROAD/call_road.py", strconv.Itoa(i), "2")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("執行失敗：", err)
			return
		}
		matrixStr := string(output)
		matrixStr = strings.TrimSpace(matrixStr)
		rows := strings.Split(matrixStr, "\n")
	
		var powers []controller.Power
		for _, row := range rows {
			var power controller.Power
			values := strings.Fields(row)
			power.StationID := "A"
			checkchargerID, err := strconv.Atoi(values[1])
			if err != nil {
				fmt.Printf("無法將字符轉換成整數1: %v\n", err)
			}
			power.ChargerID = checkchargerID
			checkpower, err := strconv.Atoi(values[2])
			if err != nil {
				fmt.Printf("無法將字符轉換成整數2: %v\n", err)
			}
			power.Power = checkpower
			checkstate, err := strconv.Atoi(values[3])
			if err != nil {
				fmt.Printf("無法將字符轉換成整數3: %v\n", err)
			}
			power.State = checkstate
			checktimestamp, err := strconv.Atoi(values[4])
			if err != nil {
				fmt.Printf("無法將字符轉換成整數4: %v\n", err)
			}
			power.TimeStamp = checktimestamp
			powers = append(powers, power)
		}

		// B場最佳化函式3
		cmd := exec.Command("python", "../FCS/call_FCS.py", strconv.Itoa(i), "2")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("執行失敗：", err)
			return
		}
		matrixStr := string(output)
		matrixStr = strings.TrimSpace(matrixStr)
		rows := strings.Split(matrixStr, "\n")
	
		for _, row := range rows {
			var power controller.Power
			values := strings.Fields(row)
			power.StationID := "B"
			checkchargerID, err := strconv.Atoi(values[1])
			if err != nil {
				fmt.Printf("無法將字符轉換成整數1: %v\n", err)
			}
			power.ChargerID = checkchargerID
			checkpower, err := strconv.Atoi(values[2])
			if err != nil {
				fmt.Printf("無法將字符轉換成整數2: %v\n", err)
			}
			power.Power = checkpower
			checkstate, err := strconv.Atoi(values[3])
			if err != nil {
				fmt.Printf("無法將字符轉換成整數3: %v\n", err)
			}
			power.State = checkstate
			checktimestamp, err := strconv.Atoi(values[4])
			if err != nil {
				fmt.Printf("無法將字符轉換成整數4: %v\n", err)
			}
			power.TimeStamp = checktimestamp
			powers = append(powers, power)
		}

		// C場最佳化函式3
		cmd := exec.Command("python", "../COMMERCIAL/call_commercial.py", strconv.Itoa(i), "2")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("執行失敗：", err)
			return
		}
		matrixStr := string(output)
		matrixStr = strings.TrimSpace(matrixStr)
		rows := strings.Split(matrixStr, "\n")
	
		for _, row := range rows {
			var power controller.Power
			values := strings.Fields(row)
			power.StationID := "C"
			checkchargerID, err := strconv.Atoi(values[1])
			if err != nil {
				fmt.Printf("無法將字符轉換成整數1: %v\n", err)
			}
			power.ChargerID = checkchargerID
			checkpower, err := strconv.Atoi(values[2])
			if err != nil {
				fmt.Printf("無法將字符轉換成整數2: %v\n", err)
			}
			power.Power = checkpower
			checkstate, err := strconv.Atoi(values[3])
			if err != nil {
				fmt.Printf("無法將字符轉換成整數3: %v\n", err)
			}
			power.State = checkstate
			checktimestamp, err := strconv.Atoi(values[4])
			if err != nil {
				fmt.Printf("無法將字符轉換成整數4: %v\n", err)
			}
			power.TimeStamp = checktimestamp
			powers = append(powers, power)
		}

		// 將各個充電樁的功率上鏈
		fmt.Printf("第%d區間上鍊開始\n",i)
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
