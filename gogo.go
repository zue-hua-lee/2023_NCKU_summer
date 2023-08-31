package main

import (
	"fmt"
	"strings"
	"strconv"
	"os/exec"
	"io"
	"os"
	"log"
	"sort"
	"math/rand"
	"time"
	"encoding/csv"
)
type Offer struct {
    CarNum int               `json:"carNum"`
    ArrTime int              `json:"arrTime"`
    DepTime int              `json:"depTime"`
    ArrSoC int               `json:"arrSoC"`
    DepSoC int               `json:"depSoC"`
    Acdc int                 `json:"acdc"`
    Capacity int             `json:"capacity"`
	Location_x float64       `json:"location_x"`
	Location_y float64       `json:"location_y"`
}
type Option struct {
	StationID string           `json:"stationID"`
	ChargerID int              `json:"chargerID"`
    MaxSoC int                 `json:"maxSoC"`
    TolPrice int               `json:"perPrice"`
    PerPrice int               `json:"perPrice"`
    ParkPrice int              `json:"parkPrice"`
}
func loadAllOffer() []Offer {
	// 讀取csv檔案
	FilePath := "./ev_schedule_1.csv"
	file, err := os.OpenFile(FilePath, os.O_RDONLY, 0777)
	if err != nil {
		log.Fatalln("找不到CSV檔案路徑:", FilePath, err)
	}
	defer file.Close()

	// 讀取第一行文字並忽略
	r := csv.NewReader(file)
	r.Read()

	// 讀取請求資訊，逐行存到陣列中
	var offers []Offer
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		var offer Offer
		checkCarNum, err := strconv.Atoi(record[0])
		if err != nil {
			log.Fatalln(err)
		}
		offer.CarNum = checkCarNum
		checkArrTime, err := strconv.Atoi(record[1])
		if err != nil {
			log.Fatalln(err)
		}
		offer.ArrTime = checkArrTime
		checkDepTime, err := strconv.Atoi(record[2])
		if err != nil {
			log.Fatalln(err)
		}
		offer.DepTime = checkDepTime
		checkArrSoC, err := strconv.Atoi(record[3])
		if err != nil {
			log.Fatalln(err)
		}
		offer.ArrSoC = checkArrSoC
		checkDepSoC, err := strconv.Atoi(record[4])
		if err != nil {
			log.Fatalln(err)
		}
		offer.DepSoC = checkDepSoC
		checkCapacity, err := strconv.Atoi(record[5])
		if err != nil {
			log.Fatalln(err)
		}
		offer.Capacity = checkCapacity
		checkAcdc, err := strconv.Atoi(record[6])
		if err != nil {
			log.Fatalln(err)
		}
		offer.Acdc = checkAcdc
		checkLocation_x, err := strconv.ParseFloat(record[7], 64)
		if err != nil {
			log.Fatalln(err)
		}
		offer.Location_x = checkLocation_x
		checkLocation_y, err := strconv.ParseFloat(record[8], 64)
		if err != nil {
			log.Fatalln(err)
		}
		offer.Location_y = checkLocation_y
		offers = append(offers, offer)
	}
	// 將所有請求依照抵達先後順序排列
	sort.Slice(offers, func(i, j int) bool {
		return offers[i].ArrTime < offers[j].ArrTime
	})
	return offers
}
func chooseOption(options []Option) Option {
	if len(options) == 0 {
		return Option{}
	}
	// 選擇1: 最高的承諾SoC
	maxSoC := options[0].MaxSoC
	for _, opt := range options {
		if opt.MaxSoC > maxSoC {
			maxSoC = opt.MaxSoC
		}
	}
	var sameMaxSoC []Option
	for _, opt := range options {
		if opt.MaxSoC == maxSoC {
			sameMaxSoC = append(sameMaxSoC, opt)
		}
	}
	if len(sameMaxSoC) == 1 {
		return sameMaxSoC[0]
	}
	
	// 選擇2: 價格最低者
	rand.Seed(time.Now().UnixNano())
	totalInverse := 0.0
	for _, opt := range sameMaxSoC {
		totalInverse += 1.0 / float64(opt.TolPrice)
	}
	randomValue := rand.Float64() * totalInverse

	currentValue := 0.0
	for _, opt := range sameMaxSoC {
		currentValue += 1.0 / float64(opt.TolPrice)
		if randomValue <= currentValue {
			return opt
		}
	}
	return Option{}
}
func main() {
	fmt.Println("start")
	// 讀取新充電請求(csv)按時間排序存到陣列中
	offers := loadAllOffer()
	var num int = 0
	for i := 1; i <= 288; i++ {
		fmt.Printf("time: %d\n", i)
		for ; num < len(offers) && offers[num].ArrTime == i; num++ {
			// 新充電申請上鏈
			fmt.Printf("time: %d, num: %d, EV: %d, ArrSoC:%d\n", i, num, offers[num].CarNum, offers[num].ArrSoC)

			// 呼叫最佳化函式1 (先慢1後快2) 
			// 輸入參數: 新車資訊、時段
			// 回傳參數: 每個廠回傳充電樁編號、承諾SoC、單價、佔位價
			var optionA Option
			var optionB Option
			var optionC Option

			checknowtime := strconv.Itoa(i)
			checkcarnum := strconv.Itoa(offers[num].CarNum)
			checkarrtime := strconv.Itoa(offers[num].ArrTime)
			checkdeptime := strconv.Itoa(offers[num].DepTime)
			checkarrsoc := strconv.Itoa(offers[num].ArrSoC)
			checkdepsoc := strconv.Itoa(offers[num].DepSoC)
			checkcapacity := strconv.Itoa(offers[num].Capacity)
			checkacdc := strconv.Itoa(offers[num].Acdc)
			checklocation_x := strconv.FormatFloat(offers[num].Location_x, 'f', -1, 64)
			checklocation_y := strconv.FormatFloat(offers[num].Location_y, 'f', -1, 64)

			// C場最佳化函式1
			cmd := exec.Command("C:/Users/ACER/anaconda3/python.exe", "./COMMERCIAL/call_commercial_newEV.py", checknowtime, checkcarnum, checkarrtime, checkdeptime, checkarrsoc, checkdepsoc, checkcapacity, checkacdc, checklocation_x, checklocation_y)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("執行失敗3：", err)
				return
			}
			matrixStr := string(output)
			matrixStr = strings.TrimSpace(matrixStr)
			rows := strings.Split(matrixStr, "\n")
			if len(rows)!=3{
				fmt.Println(matrixStr)
			}
			values := strings.Fields(rows[3])
			optionC.StationID = "C"
			for j, value := range values {
				fmt.Printf("C%d: %s\n", j, value)
				checkvalue, err := strconv.Atoi(value)
				if err != nil {
					fmt.Printf("無法將字符轉換成整數: %v\n", err)
					return
				}
				switch j {
				case 0:
					optionC.ChargerID = checkvalue
				case 1:
					optionC.MaxSoC = checkvalue
				case 2:
					optionC.PerPrice = checkvalue
				case 3:
					optionC.ParkPrice = checkvalue
				case 4:
					optionC.TolPrice = checkvalue
				}
			}

			// A場最佳化函式1
			if offers[num].Acdc == 1 {
				cmd = exec.Command("C:/Users/ACER/anaconda3/python.exe", "./ROAD/call_road_newEV.py", checknowtime, checkcarnum, checkarrtime, checkdeptime, checkarrsoc, checkdepsoc, checkcapacity, checkacdc, checklocation_x, checklocation_y)
				output, err = cmd.CombinedOutput()
				if err != nil {
					fmt.Println("執行失敗1：", err)
					return
				}
				matrixStr = string(output)
				matrixStr = strings.TrimSpace(matrixStr)
				rows = strings.Split(matrixStr, "\n")
				if len(rows)!=3{
					fmt.Println(matrixStr)
				}
				values = strings.Fields(rows[3])
				optionA.StationID = "A"
				for j, value := range values {
					fmt.Printf("A%d: %s\n", j, value)
					checkvalue, err := strconv.Atoi(value)
					if err != nil {
						fmt.Printf("無法將字符轉換成整數: %v\n", err)
						return
					}
					switch j {
					case 0:
						optionA.ChargerID = checkvalue
					case 1:
						optionA.MaxSoC = checkvalue
					case 2:
						optionA.PerPrice = checkvalue
					case 3:
						optionA.ParkPrice = checkvalue
					case 4:
						optionA.TolPrice = checkvalue
					}
				}
			} else {
				optionA.MaxSoC = 0
			}
			
			// B場最佳化函式1
			if offers[num].Acdc == 2 {
				cmd = exec.Command("C:/Users/ACER/anaconda3/python.exe", "./FCS/call_FCS_newEV.py", checknowtime, checkcarnum, checkarrtime, checkdeptime, checkarrsoc, checkdepsoc, checkcapacity, checkacdc, checklocation_x, checklocation_y)
				output, err = cmd.CombinedOutput()
				if err != nil {
					fmt.Println("執行失敗2：", err)
					return
				}
				matrixStr = string(output)
				matrixStr = strings.TrimSpace(matrixStr)
				rows = strings.Split(matrixStr, "\n")

if len(rows)!=3{
	fmt.Println(matrixStr)
}
				values = strings.Fields(rows[3])
				optionB.StationID = "B"
				for j, value := range values {
					fmt.Printf("B%d: %s\n", j, value)
					checkvalue, err := strconv.Atoi(value)
					if err != nil {
						fmt.Printf("無法將字符轉換成整數: %v\n", err)
						continue
					}
					switch j {
					case 0:
						optionB.ChargerID = checkvalue
					case 1:
						optionB.MaxSoC = checkvalue
					case 2:
						optionB.PerPrice = checkvalue
					case 3:
						optionB.ParkPrice = checkvalue
					case 4:
						optionB.TolPrice = checkvalue
					}
				}
			} else {
				optionB.MaxSoC = 0
			}
			options := []Option{optionA, optionB, optionC}

			// 媒合階段
			// 呼叫函式2
			// 輸入參數: 時段、(1/-1)
			// 回傳參數:	
			option := chooseOption(options)
			emptyOption := Option{}
			if option == emptyOption {
				fmt.Println("充電站都沒有位置了!")
				continue
			}

			fmt.Printf("station: %s, tolPrice: %d\n", option.StationID, option.TolPrice)

			var choice [3]string
			switch option.StationID {
			case "A":
				choice[0] = "1"
				choice[1] = "-1"
				choice[2] = "-1"
			case "B":
				choice[0] = "-1"
				choice[1] = "1"
				choice[2] = "-1"
			case "C":
				choice[0] = "-1"
				choice[1] = "-1"
				choice[2] = "1"
			}

			// A場最佳化函式2
			cmd = exec.Command("C:/Users/ACER/anaconda3/python.exe", "./ROAD/call_road.py", strconv.Itoa(i), choice[0])
			_, err = cmd.CombinedOutput()
			if err != nil {
				fmt.Println("執行失敗4：", err)
				return
			}
			// B場最佳化函式2
			cmd = exec.Command("C:/Users/ACER/anaconda3/python.exe", "./FCS/call_FCS.py", strconv.Itoa(i), choice[1])
			_, err = cmd.CombinedOutput()
			if err != nil {
				fmt.Println("執行失敗5：", err)
				return
			}
			// C場最佳化函式3
			cmd = exec.Command("C:/Users/ACER/anaconda3/python.exe", "./COMMERCIAL/call_commercial.py", strconv.Itoa(i), choice[2])
			_, err = cmd.CombinedOutput()
			if err != nil {
				fmt.Println("執行失敗6：", err)
				return
			}
		}

		// 5分鐘最佳化排程
		// 呼叫最佳化函式3
		// 輸入參數: 時段、(1/-1之外)
		// 回傳參數: 每個樁的功率、目前場內汽車的SoC
		// A場最佳化函式3
		cmd := exec.Command("C:/Users/ACER/anaconda3/python.exe", "./ROAD/call_road.py", strconv.Itoa(i), "2")
		_, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("執行失敗7：", err)
			return
		}

// B場最佳化函式3
		cmd = exec.Command("C:/Users/ACER/anaconda3/python.exe", "./FCS/call_FCS.py", strconv.Itoa(i), "2")
		_, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Println("執行失敗8：", err)
			return
		}

// C場最佳化函式3
		cmd = exec.Command("C:/Users/ACER/anaconda3/python.exe", "./COMMERCIAL/call_commercial.py", strconv.Itoa(i), "2")
		_, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Println("執行失敗9：", err)
			return
		}
	}
}
