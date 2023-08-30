/**
  author: Jerry
 */

package service

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (t *ServiceSetup) Offer(carNum, arrTime, depTime, arrSoC, depSoC, acdc, capacity, location_x, location_y string) (string, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "offer", Args: [][]byte{[]byte(carNum), []byte(arrTime), []byte(depTime), []byte(arrSoC), []byte(depSoC), []byte(acdc), []byte(capacity), []byte(location_x), []byte(location_y)}}
	respone, err := t.Client.Execute(req)
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

func (t *ServiceSetup) Match(stationID, chargerID, maxSoC, perPrice, now_offerID string) (string, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "match", Args: [][]byte{[]byte(stationID), []byte(chargerID), []byte(maxSoC), []byte(perPrice), []byte(parkPrice), []byte(tolPrice), []byte(now_offerID)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}
	return string(respone.TransactionID), nil
}
func (t *ServiceSetup) ShowAllMatch() (string, error){
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "showAllMatch", Args: [][]byte{}}
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