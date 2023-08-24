/*
*

	author: Jerry
*/
package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/eep/service"
)
type Option struct {
	StationID string           `json:"stationID"`
	ChargerID int              `json:"chargerID"`
    MaxSoC int                 `json:"maxSoC"`
    Price int                  `json:"price"`
}
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

func setCookie(w http.ResponseWriter, name, value string) {
    cookie := http.Cookie{
        Name:    name,
        Value:   value,
        Expires: time.Now().Add(3 * time.Hour),
    }
    http.SetCookie(w, &cookie)
}
func clearCookies(w http.ResponseWriter, cookieNames ...string) {
    for _, name := range cookieNames {
        cookie := http.Cookie{
            Name:    name,
            Value:   "",
            Expires: time.Now().Add(-time.Hour), // 設置過期時間為過去的時間
        }
        http.SetCookie(w, &cookie)
    }
}
func cookiesExist(r *http.Request, name string) bool {
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == name {
			return true
		}
	}
	return false
}
// 將時間轉為區間段
func toTimeInt(Time string) (string, error) {
    layout := "15:04"
    t, err := time.Parse(layout, Time)
    if err != nil {
        return "", err
    }
    TimeInSeconds := t.Hour()*3600 + t.Minute()*60 + t.Second()
    Time2 := strconv.Itoa((TimeInSeconds / 300) + 1)
    return Time2, nil
}

func (app *Application) CreateAccountView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "createAccount.html", nil)
}
func (app *Application) IndexView(w http.ResponseWriter, r *http.Request) {
	if cookiesExist(r, "now_userID") {
		fmt.Printf("[使用者登出] %s\n", now_userID)
	}
	clearCookies(w, "now_userID", "now_offerID")
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
		FlagA:false,
		MsgA1:"",
		MsgA2:"",
		MsgA3:"",
		FlagB:false,
		MsgB1:"",
		MsgB2:"",
		MsgB3:"",
		FlagC:false,
		MsgC1:"",
		MsgC2:"",
		MsgC3:"",
	}
	if cookiesExist(r, "optionA") {
		var optionA Option
		cookie, _ := r.Cookie("optionA")
		json.Unmarshal([]byte(cookie.Value), &optionA)
		data.FlagA = true
		data.MsgA1 = strconv.Itoa(optionA.MaxSoC)
		data.MsgA2 = strconv.Itoa(optionA.Price)
		data.MsgA3 = strconv.Itoa(optionA.Price)
	}
	if cookiesExist(r, "optionB") {
		var optionB Option
		cookie, _ := r.Cookie("optionB")
		json.Unmarshal([]byte(cookie.Value), &optionB)
		data.FlagB = true
		data.MsgB1 = strconv.Itoa(optionB.MaxSoC)
		data.MsgB2 = strconv.Itoa(optionB.Price)
		data.MsgB3 = strconv.Itoa(optionB.Price)
	}
	if cookiesExist(r, "optionC") {
		var optionC Option
		cookie, _ := r.Cookie("optionC")
		json.Unmarshal([]byte(cookie.Value), &optionC)
		data.FlagC = true
		data.MsgC1 = strconv.Itoa(optionC.MaxSoC)
		data.MsgC2 = strconv.Itoa(optionC.Price)
		data.MsgC3 = strconv.Itoa(optionC.Price)
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
	showView(w, r, "request4.html", nil)
}
func (app *Application) TrackNoView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "trackNo.html", nil)
}
func (app *Application) TrackYesView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "trackYes.html", nil)
}
func (app *Application) TrackView(w http.ResponseWriter, r *http.Request) {
	if cookiesExist(r, "now_offerID") {
		app.TrackYesView(w, r)
	} else {
		app.TrackNoView(w, r)
	}
}
func (app *Application) HistoryListView(w http.ResponseWriter, r *http.Request) {
	showView(w, r, "historyList.html", nil)
}

var index int = 1
func (app *Application) Schedule() {
	// 設定每5分鐘執行
	ticker := time.NewTicker(5 * time.Minute)
	time.Sleep(5 * time.Minute)
	for {
		select {
		case <-ticker.C:
			// 顯示現在時間
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

	// 查看車牌是否已註冊過
	_ , err := app.Fabric.GetUserIDbyCarID(carID)
	if err == nil {
		err = errors.New("車牌曾經註冊!請登入!")
	}else{
		_, err = app.Fabric.Register(carID, userName, capacity, password)
	}

	if err != nil {
		// 註冊失敗，要求重新輸入
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
		// 註冊成功則登入
		now_userID, _ := app.Fabric.Login(carID, password)
		setCookie(w, "now_userID", now_userID)
		
		fmt.Printf("[新使用者註冊並登入] %s\n", now_userID)
		showView(w, r, "mainPage.html", nil)
	}
}

func (app *Application) Login(w http.ResponseWriter, r *http.Request)  {
	carID := r.FormValue("carID")
	password := r.FormValue("password")

	now_userID, err := app.Fabric.Login(carID, password)

	if err != nil {
		// 登入失敗，要求重新輸入
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
		// 登入成功
		setCookie(w, "now_userID", now_userID)
		fmt.Printf("[使用者登入] %s\n", now_userID)
		showView(w, r, "mainPage.html", nil)
	}
}

func (app *Application) Offer(w http.ResponseWriter, r *http.Request)  {
	// 將進場時間轉換為區間段
	arrTime := r.FormValue("arrTime")
	arrTime2, err := toTimeInt(arrTime)
	if err != nil {
		http.Error(w, "Invalid arrtime format", http.StatusBadRequest)
	}

	// 將離場時間轉換為區間段
	depTime := r.FormValue("depTime")
	depTime2, err := toTimeInt(depTime)
	if err != nil {
		http.Error(w, "Invalid depTime format", http.StatusBadRequest)
	}

	arrSoC := r.FormValue("arrSoC")
	depSoC := r.FormValue("depSoC")
	acdc := r.FormValue("acdc")
	cookie, _ := r.Cookie("now_userID")
	now_userID := cookie.Value
	now_offerID, err := app.Fabric.Offer(arrTime2, depTime2, arrSoC, depSoC, acdc, now_userID)

	if err != nil {
		// 申請充電失敗
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
		// 申請充電成功
		setCookie(w, "now_offerID", now_offerID)
		fmt.Printf("[發送新請求] %s: %s\n", now_userID, now_offerID)
		
		// 呼叫最佳化函式1
		optionA := Option{StationID: "A", ChargerID: 1, MaxSoC: 100, Price: 100,}
		optionB := Option{StationID: "B", ChargerID: 1, MaxSoC: 100, Price: 50,}
		optionC := Option{StationID: "C", ChargerID: 1, MaxSoC: 100, Price: 30,}
		optionAAsBytes, _ := json.Marshal(optionA)
		optionBAsBytes, _ := json.Marshal(optionB)
		optionCAsBytes, _ := json.Marshal(optionC)
		setCookie(w, "optionA", optionAAsBytes)
		setCookie(w, "optionB", optionBAsBytes)
		setCookie(w, "optionC", optionCAsBytes)
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