/**
  author: jerry
 */
package web

import (
	"net/http"
	"fmt"
	"github.com/hyperledger/fabric/eep/web/controller"
)

func  WebStart(app *controller.Application)  {

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", app.MyHomeView)
	http.HandleFunc("/myHome.html", app.MyHomeView)
	http.HandleFunc("/myMatch.html", app.MyMatchView)

    http.HandleFunc("/register", app.Register)
    http.HandleFunc("/showAllUser", app.ShowAllUser)
    http.HandleFunc("/login", app.Login)
    http.HandleFunc("/showNowUser", app.ShowNowUser)
    http.HandleFunc("/offer", app.Offer)
	http.HandleFunc("/showAllOffer", app.ShowAllOffer)
    http.HandleFunc("/match", app.Match)
	http.HandleFunc("/showAllMatch", app.ShowAllMatch)
    http.HandleFunc("/power", app.Power)
	http.HandleFunc("/showAllPower", app.ShowAllPower)
	http.HandleFunc("/showPowerbyCharger", app.ShowPowerbyCharger)


	fmt.Println("Start test WEB, port: 9000")

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("Start WEB Error")
	}

}
