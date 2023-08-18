/*
*

	author: Jerry
*/
package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/eep/service"
)

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
	showView(w, r, "request2.html", nil)
}
func (app *Application) Request3View(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "request3.html", nil)
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
var now_userID string = ""
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
		data.Msg = err.Error()
		showView(w, r, "createAccount.html", data)
	}else{
		now_userID, err = app.Fabric.Login(carID, password)
		showView(w, r, "mainPage.html", nil)
	}
}

func (app *Application) Login(w http.ResponseWriter, r *http.Request)  {
	carID := r.FormValue("carID")
	password := r.FormValue("password")

	var err error
	now_userID, err = app.Fabric.Login(carID, password)

	if err != nil {
		data := &struct {
			Flag bool
			Msg string
		}{
			Flag:true,
			Msg:"",
		}
		data.Msg = err.Error()
		showView(w, r, "index.html", data)
	}else{
		showView(w, r, "mainPage.html", nil)
	}
}

var now_offerID string = ""
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
	
	fmt.Println("now_userID: "+ now_userID)
	fmt.Println("arrTime2: "+ arrTime2)
	fmt.Println("depTime2: "+ depTime2)
	fmt.Println("arrSoC: "+ arrSoC)
	fmt.Println("depSoC: "+ depSoC)
	fmt.Println("acdc: "+ acdc)

	now_offerID, err = app.Fabric.Offer(arrTime2, depTime2, arrSoC, depSoC, acdc, now_userID)

	if err != nil {
		data := &struct {
			Flag bool
			Msg string
		}{
			Flag:true,
			Msg:"",
		}
		data.Msg = err.Error()
		showView(w, r, "request1.html", data)
	}else{
		showView(w, r, "request2.html", nil)
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
// 		data.Msg1 = err.Error()
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
// 		data.Msg1 = err.Error()
// 	}else{
// 		data.Msg1 = "ShowAllMatch success: " + msg
// 	}
// 	showView(w, r, "myMatch.html", data)
// }
// func (app *Application) Power(w http.ResponseWriter, r *http.Request)  {
// 	stationID := r.FormValue("stationID")
// 	chargerID := r.FormValue("chargerID")
// 	power := r.FormValue("power")
// 	state := r.FormValue("state")
// 	timestamp := r.FormValue("timestamp")

// 	transactionID, err := app.Fabric.Power(stationID, chargerID, power, state, timestamp)

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
// 		data.Msg2 = err.Error()
// 	}else{
// 		data.Msg2 = "Power success, Transaction ID: " + transactionID
// 	}
// 	showView(w, r, "myMatch.html", data)
// }
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
// 		data.Msg2 = err.Error()
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
// 		data.Msg2 = err.Error()
// 	}else{
// 		data.Msg2 = "ShowPowerbyCharger success: " + msg
// 	}
// 	showView(w, r, "myMatch.html", data)
// }