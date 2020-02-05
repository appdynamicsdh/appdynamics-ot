package appdcollector

import (
	"log"
	"fmt"
	"time"
	//"math/rand"
	appd "../appd"
)

type ExtractedData struct {
        StartTime int64    `json:"StartTime"`
        Duration  int64    `json:"Duration"`
        BTName    string `json:"BTName"`
        TierName  string `json:"TierName"`
        NodeName  string `json:"NodeName"`
}

func initializeAppD(tier string, node string) {

	// see if you can parameterize the initialization parameters
	appd.Init("OpenTrace", "a2c289ac-561f-433e-abee-83a47321fe21")
	appd.SetTierName(tier)
	appd.SetNodeName(node)
	appd.SetControllerHost("4516controllernoss-lesterhackx-abjhbb1k.srv.ravcloud.com")
	appd.SetControllerPort(8090)
	appd.SetControllerAccount("customer1")
	appd.SetControllerUseSSL(0)
	appd.SetInitTimeout(5000)

	rc := appd.Sdk_init()

	if rc != 0 {
		log.Fatal("Failed to initialize C++ SDK", rc)
	}
}


func CreateBusinessTransaction(e ExtractedData) {

	fmt.Println("Creating a Business Transaction")
	initializeAppD(e.TierName, e.NodeName)
	//appd_correlation := ""
	bt := appd.BT_begin(e.BTName, "")

	fmt.Print(".")
	milliseconds := 250
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)

	appd.BT_override_start_time_ms(bt, uint64(e.StartTime))
	//s1 := rand.NewSource(time.Now().UnixNano())
        //r1 := rand.New(s1)
        appd.BT_override_time_ms(bt, uint64(e.Duration))

	if appd.BT_is_snapshotting(bt) != 0 {
		fmt.Println("adding data")
		appd.BT_set_url(bt, r.URL.String())
		//appd.BT_add_user_data(bt, "postcode", code)
	}

	//appd.BT_add_error(bt, appd.ERROR_LEVEL_ERROR, "no match", 1)
        appd.BT_end(bt)


}

func attachExitCall() {

	//exit := appd.Exitcall_begin(bt, backendName)
	//fmt.Println("Exit handle", exit)

	//appd.Exitcall_set_details(exit, query)

	exit := appd.Exitcall_begin(bt, backendName)
	fmt.Println("Exit handle", exit)

	appd.Exitcall_set_details(exit, query)
	fmt.Println("Looking for ", code)
	// report error
	appd.BT_add_error(bt, appd.ERROR_LEVEL_ERROR, "no match", 1)

	appd.Exitcall_end(exit)


}

func processTrace() {

	// add logic to read the incoming trace and decide what to do
	// a create a BT
	// create an exit call

	// 1. Validate Trace Map
	// 2. Create a BT

	// retrieve BT start time
	// retrieve BT response time
	// retrive BT Name (http.url)
	// retrieve Service Name (tierName)
	// retrieve Node Name = ipv4



}

//func main() {
//    	fmt.Println("Testing Creation of BT via ZipKin incoming Trace")
//    	btCount := 0
//	maxBtCount := 200
//	fmt.Print("Doing something")
//	for btCount < maxBtCount {
//    		createBusinessTransaction()
//	}
//}
