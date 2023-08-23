/*
*

	author: Jerry
*/
package controller

import (
	"errors"
	"fmt"
	"log"
	"time"
	"strings"
	"net/http"
	"strconv"
	"encoding/json"

	"github.com/hyperledger/fabric/eep/service"
)
type Power struct {
	StationID string           `json:"stationID"`
	ChargerID int              `json:"chargerID"`
    Power int                  `json:"power"`
    State int                  `json:"state"`
    TimeStamp int              `json:"timeStamp"`
}

type Application struct {
	Fabric *service.ServiceSetup
}
func (app *Application) CreateAccountView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "createAccount.html", nil)
}
func (app *Application) HistoryListView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "historyList.html", nil)
}
func (app *Application) IndexView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "index.html", nil)
}
func (app *Application) MainPageView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "mainPage.html", nil)
}
func (app *Application) Request1View(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "request1.html", nil)
}
func (app *Application) Request2View(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		FlagA bool
		MsgA1 string
		MsgA2 string
		MsgA3 string
		FlagB bool
		MsgB1 string
		MsgB2 string
		MsgB3 string
		FlagC bool
		MsgC1 string
		MsgC2 string
		MsgC3 string
	}{
		FlagA:true,
		MsgA1:"100",
		MsgA2:"100",
		MsgA3:"100",
		FlagB:true,
		MsgB1:"90",
		MsgB2:"50",
		MsgB3:"70",
		FlagC:false,
		MsgC1:"-",
		MsgC2:"-",
		MsgC3:"-",
	}
	showView(w, r, "request2.html", data)
}
func (app *Application) Request3View(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		Msg1 string
		Msg2 string
		Msg3 string
		Msg4 string
		Msg5 string
		Msg6 string
	}{
		Msg1:"甲",
		Msg2:"20",
		Msg3:"快充",
		Msg4:"30",
		Msg5:"90",
		Msg6:"100",
	}
	showView(w, r, "request3.html", data)
}
func (app *Application) Request4View(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "request4.html", nil)
}
func (app *Application) TrackNoView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "trackNo.html", nil)
}
func (app *Application) TrackYesView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "trackYes.html", nil)
}
var index int = 1
func (app *Application) Schedule() {
	ticker := time.NewTicker(5 * time.Minute)
	time.Sleep(5 * time.Minute)
	for {
		select {
		case <-ticker.C:
			fmt.Println("Current time:", time.Now())
			// 將各個充電裝的功率上鏈
			fmt.Printf("第%d區間上鍊開始\n",index)
			var powers []Power
			for j := 1; j <= 12; j++{
				powers = append(powers, Power{StationID: "A", ChargerID: j, Power: 0, State: 0, TimeStamp: index})
			}
			for j := 1; j <= 6; j++{
				powers = append(powers, Power{StationID: "B", ChargerID: j, Power: 30, State: 1, TimeStamp: index})
			}
			for j := 1; j <= 20; j++{
				powers = append(powers, Power{StationID: "C", ChargerID: j, Power: 40, State: 1, TimeStamp: index})
			}
			app.Power(powers)
			fmt.Printf("第%d區間上鍊結束\n",index)
			index++
		}
	}
}

func (app *Application) Register(w http.ResponseWriter, r *http.Request)  {
	carID := r.FormValue("carID")
	userName := r.FormValue("userName")
	capacity := r.FormValue("capacity")
	password := r.FormValue("password")

	_ , err := app.Fabric.GetUserIDbyCarID(carID)	
	
	if err == nil {
		err = errors.New("Register ERROR! CarID already exist!")
	}else{
		_, err = app.Fabric.Register(carID, userName, capacity, password)
	}

	if err != nil {
		data := &struct {
			Flag bool
			Msg string
		}{
			Flag:true,
			Msg:"",
		}
		errMessage := strings.Split(err.Error(), "Description: ")
		data.Msg = errMessage[1]
		showView(w, r, "createAccount.html", data)
	}else{
		now_userID, _ := app.Fabric.Login(carID, password)

		cookie := http.Cookie{
			Name:  "now_userID",
			Value: now_userID,
			Expires: time.Now().Add(3 * time.Hour),
		}
		http.SetCookie(w, &cookie)
		
		fmt.Printf("[新使用者註冊並登入] %s\n", now_userID)
		showView(w, r, "mainPage.html", nil)
	}
}

func (app *Application) Login(w http.ResponseWriter, r *http.Request)  {
	carID := r.FormValue("carID")
	password := r.FormValue("password")

	now_userID, err := app.Fabric.Login(carID, password)

	if err != nil {
		data := &struct {
			Flag bool
			Msg string
		}{
			Flag:true,
			Msg:"",
		}
		errMessage := strings.Split(err.Error(), "Description: ")
		data.Msg = errMessage[1]
		showView(w, r, "index.html", data)
	}else{
		cookie := http.Cookie{
			Name:  "now_userID",
			Value: now_userID,
			Expires: time.Now().Add(3 * time.Hour),
		}
		http.SetCookie(w, &cookie)

		fmt.Printf("[使用者登入] %s\n", now_userID)
		showView(w, r, "mainPage.html", nil)
	}
}

func (app *Application) Offer(w http.ResponseWriter, r *http.Request)  {
	arrTime := r.FormValue("arrTime")
	layout := "15:04"
	t, err := time.Parse(layout, arrTime)
	if err != nil {
		http.Error(w, "Invalid arrtime format", http.StatusBadRequest)
	}
	arrTimeInSeconds := t.Hour()*3600 + t.Minute()*60 + t.Second()
	arrTime2 := strconv.Itoa((arrTimeInSeconds / 300) + 1)

	depTime := r.FormValue("depTime")
	t, err = time.Parse(layout, depTime)
	if err != nil {
		http.Error(w, "Invalid deptime format", http.StatusBadRequest)
	}
	depTimeInSeconds := t.Hour()*3600 + t.Minute()*60 + t.Second()
	depTime2 := strconv.Itoa((depTimeInSeconds / 300) + 1)

	arrSoC := r.FormValue("arrSoC")
	depSoC := r.FormValue("depSoC")
	acdc := r.FormValue("acdc")
	cookie, _ := r.Cookie("now_userID")
	now_userID := cookie.Value
	now_offerID, err := app.Fabric.Offer(arrTime2, depTime2, arrSoC, depSoC, acdc, now_userID)

	if err != nil {
		data := &struct {
			Flag bool
			Msg string
		}{
			Flag:true,
			Msg:"",
		}
		errMessage := strings.Split(err.Error(), "Description: ")
		data.Msg = errMessage[1]
		showView(w, r, "request1.html", data)
	}else{
		cookie := http.Cookie{
			Name:  "now_offerID",
			Value: now_offerID,
			Expires: time.Now().Add(3 * time.Hour),
		}
		http.SetCookie(w, &cookie)
		fmt.Printf("[發送新請求] %s: %s\n", now_userID, now_offerID)
		app.Request2View(w,r)
	}
}

// func (app *Application) Match(w http.ResponseWriter, r *http.Request)  {
// 	stationID := r.FormValue("stationID")
// 	maxSoC := r.FormValue("maxSoC")
// 	price := r.FormValue("price")

// 	transactionID, err := app.Fabric.Match(stationID, maxSoC, price, now_offerID)

// 	data := &struct {
// 		Flag1 bool
// 		Msg1 string
// 		Flag2 bool
// 		Msg2 string
// 	}{
// 		Flag1:true,
// 		Msg1:"",
// 		Flag2:false,
// 		Msg2:"",
// 	}
// 	if err != nil {
	// errMessage := strings.Split(err.Error(), "Description: ")
	// data.Msg1 = errMessage[1]
// 	}else{
// 		data.Msg1 = "Match success, Transaction ID: " + transactionID
// 	}
// 	showView(w, r, "myMatch.html", data)
// }

// func (app *Application) ShowAllMatch(w http.ResponseWriter, r *http.Request)  {
// 	msg, err := app.Fabric.ShowAllMatch()

// 	data := &struct {
// 		Flag1 bool
// 		Msg1 string
// 		Flag2 bool
// 		Msg2 string
// 	}{
// 		Flag1:true,
// 		Msg1:"",
// 		Flag2:false,
// 		Msg2:"",
// 	}
// 	if err != nil {
	// errMessage := strings.Split(err.Error(), "Description: ")
	// data.Msg1 = errMessage[1]
// 	}else{
// 		data.Msg1 = "ShowAllMatch success: " + msg
// 	}
// 	showView(w, r, "myMatch.html", data)
// }

func (app *Application) Power(powers []Power)  {
    PowersAsBytes, err := json.Marshal(powers)
    if err != nil {
        log.Fatalln(err)
    }
	_, err = app.Fabric.Power(PowersAsBytes)
	if err != nil {
		log.Fatalln(err)
	}
}
// func (app *Application) ShowAllPower(w http.ResponseWriter, r *http.Request)  {
// 	msg, err := app.Fabric.ShowAllPower()

// 	data := &struct {
// 		Flag1 bool
// 		Msg1 string
// 		Flag2 bool
// 		Msg2 string
// 	}{
// 		Flag1:false,
// 		Msg1:"",
// 		Flag2:true,
// 		Msg2:"",
// 	}
// 	if err != nil {
	// errMessage := strings.Split(err.Error(), "Description: ")
	// data.Msg2 = errMessage[1]
// 	}else{
// 		data.Msg2 = "ShowAllPower success: " + msg
// 	}
// 	showView(w, r, "myMatch.html", data)
// }
// func (app *Application) ShowPowerbyCharger(w http.ResponseWriter, r *http.Request)  {
// 	stationID := r.FormValue("stationID_search")
// 	chargerID := r.FormValue("chargerID_search")

// 	msg, err := app.Fabric.ShowPowerbyCharger(stationID, chargerID)

// 	data := &struct {
// 		Flag1 bool
// 		Msg1 string
// 		Flag2 bool
// 		Msg2 string
// 	}{
// 		Flag1:false,
// 		Msg1:"",
// 		Flag2:true,
// 		Msg2:"",
// 	}
// 	if err != nil {
	// errMessage := strings.Split(err.Error(), "Description: ")
	// data.Msg2 = errMessage[1]
// 	}else{
// 		data.Msg2 = "ShowPowerbyCharger success: " + msg
// 	}
// 	showView(w, r, "myMatch.html", data)
// }