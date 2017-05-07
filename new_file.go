package main 

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {
}

type booking struct {
	ObjectType	string `json:"docType"` //docType is used to distinguish the various types of objects in state database  not sure why its needed
	AgentID		string `json:"agentId"`   //unique agend id
	BookingID	int	   `json:"bookingId"` //booking ID
	TotalBill	float32	`json:"totalBill"`  //no of days client will be staying		
	IsSettled	bool   `json:"isSettled"`  //is booking is pending or completed
	Payables	float32`json:"payables"`   //final payables to agent
}


func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil{
		fmt.Printf("Error starting SimpleChaincode")
		}
}

//Init initializes chaincode
//==============================

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Printf("Hello World")
	return shim.Success(nil)
}

//Invoke - entry points for invokations
//================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function,args:=stub.GetFunctionAndParameters()
	fmt.Println("invoke is running" +function)
	
	caller:= args[0]
	
	//divert different functions
	if function == "initBooking" && caller == "TA" {  //create a new booking
		return t.initBooking(stub,args)
	} else if function == "updateSettled" && caller == "HOTEL" {
		return t.updateSettled(stub,args)
	} else if function == "updatePayables" && caller == "CTAC" {
		return t.updatePayables(stub,args)
	}
	
}

func(t *SimpleChaincode) updatePayables(stub shim.ChaincodeStubInterface, args string[]) pb.Response {
	// arg0 will be caller
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting Booking ID")
	}
	
	bookingID := args[1]
	bookingAsBytes, err := stub.getState(bookingID)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + bookingID + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp := "{\"Error\":\"Booking ID does not exist: " + bookingID + "\"}"
		return shim.Error(jsonResp)
	} 
	
	bookingTemp := booking{}
	
	err:= json.Unmarshal(bookingAsBytes, &bookingTemp)
	
	if err != nil {
		return shim.Error(err.Error())
	}
	
	
	if bookingTemp.IsSettled  != true {
		
	}
	
	
}

func (t *SimpleChaincode) updateSettled(stub shim.ChaincodeStubInterface, args string[]) pb.Response {
	// arg0 will be caller arg1 booking id  arg2 totalbill
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting Booking ID")
	}
	
	bookingID := args[1]
	
	bookingAsBytes, err := stub.getState(bookingID)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + bookingID + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp := "{\"Error\":\"Booking ID does not exist: " + bookingID + "\"}"
		return shim.Error(jsonResp)
	} 
	
	
	bookingTemp := booking{}
	err:= json.Unmarshal(bookingAsBytes, &bookingTemp)
	
	if err != nil {
		return shim.Error(err.Error())
	}
	
	bookingTemp.IsSettled=true
	bookingTemp.Totalbill=args[2]
	
	bookingAsJsonBytes, err := json.Marshal(bookingTemp)
	err = stub.putState(bookingID, bookingAsJsonBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("settled status and total bill are successfully updated")
	
}


func (t *SimpleChaincode) initBooking(stub shim.ChaincodeStubInterface, args string[]) pb.Response {
	
	// when calling init booking we need to pass all 5 arguments 
	//arg0 - caller type
	var err error
	
	agentID := strings.ToLower(args[1])
	bookingID := strconv.Atoi(args[2])
	totalBill := args[3])
	isSettled := args[4]
	payables := args[5]
	
	
	objectType := "booking"
	booking := &booking(objectType, agentID, bookingID, totalBill, isSettled, payables)
	bookingAsJsonBytes,err := json.Marshal(booking)
	if err != nil {
		return shim.Error(err.Error())
	}
	//put the marshal into chaincode
	err = stub.putState(bookingID, bookingAsJsonBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	// make index later
	
}

