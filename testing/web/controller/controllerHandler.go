/**
  author: Jerry
 */
package controller

import (
	"io"
	"os"
	"fmt"
	"log"
	"sort"
	"strconv"
	"math/rand"
	"time"
	"encoding/csv"
    "encoding/json"
	"github.com/hyperledger/fabric/eep/service"
)

type Application struct {
	Fabric *service.ServiceSetup
}
//////////////////////
func (app *Application) Testing() {
	fmt.Println("123456789")
}
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
type Power struct {
	StationID string           `json:"stationID"`
	ChargerID int              `json:"chargerID"`
    Power int                  `json:"power"`
    State int                  `json:"state"`
    TimeStamp int              `json:"timeStamp"`
}

func (app *Application) LoadAllOffer() []Offer {
	// 讀取csv檔案
	FilePath := "./web/static/csv/ev_schedule_1.csv"
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

func (app *Application) ChooseOption(options []Option) Option {
	// 選擇1: 最高的承諾SoC
	if options[0].MaxSoC > options[1].MaxSoC && options[0].MaxSoC > options[2].MaxSoC {
		return options[0]
	}
	if options[1].MaxSoC > options[0].MaxSoC && options[1].MaxSoC > options[2].MaxSoC {
		return options[1]
	}
	if options[2].MaxSoC > options[0].MaxSoC && options[2].MaxSoC > options[1].MaxSoC {
		return options[2]
	}

	// 選擇2: 價格最低者
	rand.Seed(time.Now().UnixNano())
	totalInverse := 0.0
	for _, opt := range options {
		totalInverse += 1.0 / float64(opt.TolPrice)
	}
	randomValue := rand.Float64() * totalInverse

	currentValue := 0.0
	for _, opt := range options {
		currentValue += 1.0 / float64(opt.TolPrice)
		if randomValue <= currentValue {
			return opt
		}
	}
	return options[len(options)-1]
}

var now_offerID string = ""
func (app *Application) SetOffer(offer Offer) {
	var err error
	now_offerID, err = app.Fabric.Offer(strconv.Itoa(offer.CarNum), strconv.Itoa(offer.ArrTime), strconv.Itoa(offer.DepTime),
										strconv.Itoa(offer.ArrSoC), strconv.Itoa(offer.DepSoC),
					 					strconv.Itoa(offer.Acdc), strconv.Itoa(offer.Capacity),
					 					strconv.FormatFloat(offer.Location_x, 'f', -1, 64),
					 					strconv.FormatFloat(offer.Location_y, 'f', -1, 64))
	if err != nil {
		log.Fatalln(err)
	}
}
func (app *Application) ShowAllOffer()  {
	msg, err := app.Fabric.ShowAllOffer()
	if err != nil {
		log.Fatalln(err)
	}else{
		fmt.Println(msg)
	}
}

func (app *Application) SetMatch(option Option)  {
	_, err := app.Fabric.Match(option.StationID, strconv.Itoa(option.ChargerID), strconv.Itoa(option.MaxSoC), strconv.Itoa(option.PerPrice), strconv.Itoa(option.ParkPrice), strconv.Itoa(option.TolPrice), now_offerID)
	if err != nil {
		log.Fatalln(err)
	}
}
func (app *Application) ShowAllMatch()  {
	msg, err := app.Fabric.ShowAllMatch()
	if err != nil {
		log.Fatalln(err)
	}else{
		fmt.Println(msg)
	}
}

func (app *Application) SetPower(powers []Power)  {
    PowersAsBytes, err := json.Marshal(powers)
    if err != nil {
        log.Fatalln(err)
    }
	_, err = app.Fabric.Power(PowersAsBytes)
	if err != nil {
		log.Fatalln(err)
	}
}
func (app *Application) ShowAllPower()  {
	msg, err := app.Fabric.ShowAllPower()
	if err != nil {
		log.Fatalln(err)
	}else{
		fmt.Println(msg)
	}
}