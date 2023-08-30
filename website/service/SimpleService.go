/**
  author: Jerry
*/

package service

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (t *ServiceSetup) Register(carID, userName, capacity, password string) (string, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "register", Args: [][]byte{[]byte(carID), []byte(userName), []byte(capacity), []byte(password)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}
	return string(respone.TransactionID), nil
}
// func (t *ServiceSetup) ShowAllUser() (string, error){
// 	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "showAllUser", Args: [][]byte{}}
// 	respone, err := t.Client.Query(req)
// 	if err != nil {
// 			return "", err
// 	}
// 	return string(respone.Payload), nil
// }

func (t *ServiceSetup) Login(carID, password string) (string, error){
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "login", Args: [][]byte{[]byte(carID), []byte(password)}}
	respone, err := t.Client.Query(req)
	if err != nil {
			return "", err
	}
	return string(respone.Payload), nil
}
func (t *ServiceSetup) GetUserIDbyCarID(carID string) (string, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "getUserIDbyCarID", Args: [][]byte{[]byte(carID)}}
	respone, err := t.Client.Query(req)
	if err != nil {
			return "", err
	}
	return string(respone.Payload), nil
}
func (t *ServiceSetup) ShowCarbyID(userID string) (string, error){
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "ShowCarbyID", Args: [][]byte{[]byte(userID)}}
	respone, err := t.Client.Query(req)
	if err != nil {
			return "", err
	}
	return string(respone.Payload), nil
}

func (t *ServiceSetup) Offer(arrTime, depTime, arrSoC, depSoC, acdc, userID string) (string, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "offer", Args: [][]byte{[]byte(arrTime), []byte(depTime), []byte(arrSoC), []byte(depSoC), []byte(acdc), []byte(userID)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}
	return string(respone.Payload), nil
}
func (t *ServiceSetup) ShowOfferbyID(offerID string) (string, error){
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "showOfferbyID", Args: [][]byte{[]byte(offerID)}}
	respone, err := t.Client.Query(req)
	if err != nil {
			return "", err
	}
	return string(respone.Payload), nil
}
func (t *ServiceSetup) ShowAllOffer() (string, error){
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "showAllOffer", Args: [][]byte{}}
	respone, err := t.Client.Query(req)
	if err != nil {
			return "", err
	}
	return string(respone.Payload), nil
}

func (t *ServiceSetup) Match(stationID, chargerID, date, arrTime, depTime, arrSoC, maxSoC, price, offerID string) (string, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "match", Args: [][]byte{[]byte(stationID), []byte(chargerID), []byte(date), []byte(arrTime), []byte(depTime), []byte(arrSoC), []byte(maxSoC), []byte(price), []byte(offerID)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}
	return string(respone.Payload), nil
}
func (t *ServiceSetup) ShowMatchbyUser(userID string) (string, error){
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "showMatchbyUser", Args: [][]byte{[]byte(userID)}}
	respone, err := t.Client.Query(req)
	if err != nil {
			return "", err
	}
	return string(respone.Payload), nil
}

func (t *ServiceSetup) Power(PowersAsBytes []byte) (string, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "power", Args: [][]byte{PowersAsBytes}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}
	return string(respone.TransactionID), nil
}
func (t *ServiceSetup) ShowAllPower() (string, error){
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "showAllPower", Args: [][]byte{}}
	respone, err := t.Client.Query(req)
	if err != nil {
			return "", err
	}
	return string(respone.Payload), nil
}
func (t *ServiceSetup) ShowPowerbyCharger(stationID, chargerID string) (string, error){
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "showPowerbyCharger", Args: [][]byte{[]byte(stationID), []byte(chargerID)}}
	respone, err := t.Client.Query(req)
	if err != nil {
			return "", err
	}
	return string(respone.Payload), nil
}
func (t *ServiceSetup) ShowPowerbyMatch(matchID string) (string, error){
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "showPowerbyMatch", Args: [][]byte{[]byte(matchID)}}
	respone, err := t.Client.Query(req)
	if err != nil {
			return "", err
	}
	return string(respone.Payload), nil
}