package appdcollector

import (
	"log"
	"fmt"
	"time"
	//"math/rand"
	appd "../appd"
	constants "../constants"
)


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


func CreateBusinessTransaction(resultSet constants.ExtractedData) {

	log.Printf("%v\n", len(resultSet))
        fmt.Println(len(resultSet))
	//create BT
	fmt.Println("Creating a Business Transaction")
	initializeAppD(resultSet[0].TierName, resultSet[0].NodeName)
	//appd_correlation := ""
	bt := appd.BT_begin(resultSet[0].BTName, "")

	//appd.BT_override_start_time_ms(bt, uint64(e.StartTime))
	//s1 := rand.NewSource(time.Now().UnixNano())
	//r1 := rand.New(s1)
	appd.BT_override_time_ms(bt, uint64(resultSet[0].Duration))

	if appd.BT_is_snapshotting(bt) != 0 {
		fmt.Println("adding data")
		appd.BT_set_url(bt, resultSet[0].BTName)
		appd.BT_add_user_data(bt, "trace_id", resultSet[0].TraceID)
		appd.BT_add_user_data(bt, "span_id", resultSet[0].SpanID)
		appd.BT_add_user_data(bt, "parent_id", resultSet[0].ParentID)
	}
	log.Printf("%v\n", resultSet[0].Duration)
                        log.Printf(resultSet[0].BTName)
                        log.Printf(resultSet[0].TierName)
                        log.Printf(resultSet[0].NodeName)
                        log.Printf(resultSet[0].TraceID)
                        log.Printf(resultSet[0].ParentID)
                        log.Printf(resultSet[0].SpanID)
                        log.Printf("END OF RESULT")
	for i := 1; i < len(resultSet); i++ {
    //fmt.Printf("%x ", resultSet[i])
    log.Printf("%v\n", resultSet[i].StartTime)
                        log.Printf("%v\n", resultSet[i].Duration)
                        log.Printf(resultSet[i].BTName)
                        log.Printf(resultSet[i].TierName)
                        log.Printf(resultSet[i].NodeName)
                        log.Printf(resultSet[i].TraceID)
                        log.Printf(resultSet[i].ParentID)
                        log.Printf(resultSet[i].SpanID)
                        log.Printf("END OF RESULT")
                        log.Printf("")
		// Attach Exit Calls

		// APPD Backend
		appd.Backend_declare(appd.BACKEND_HTTP, resultSet[i].BTName)
		props := appd.ID_properties_map {
			"Time": "1867 ms",
			"From": "A",
			"To": "B",
			"Details": "http://devops-data-services-mysql:8080/devops-legacy-data-services/data/v1/accountProfile/20329",

		}
		appd.Backend_set_identifying_properties(resultSet[i].BTName, props)
		//if(rc != 0) {
		//	log.Fatal("Backend ", rc)
		//}
		appd.Backend_add(resultSet[i].BTName)
		//if(rc != 0) {
		//	log.Fatal("Backend ", rc)
		//}
		fmt.Println("Backend added")

                sleep()
		exit := appd.Exitcall_begin(bt, resultSet[i].BTName)

		//appd.Exitcall_override_start_time_ms(exit, uint64(resultSet[i].StartTime))
		appd.Exitcall_override_time_ms(exit, uint64(resultSet[i].Duration))

		fmt.Println("Exit handle", exit)
                sleep()
		appd.Exitcall_set_details(exit, resultSet[i].TraceID)
		//appd.Exitcall_set_details(exit, e.ParentID)
		//appd.Exitcall_set_details(exit, e.SpanID)

		//appd.BT_add_error(bt, appd.ERROR_LEVEL_ERROR, "no match", 1)
		appd.Exitcall_end(exit)


	}

	//appd.BT_add_error(bt, appd.ERROR_LEVEL_ERROR, "no match", 1)
  appd.BT_end(bt)

}


func processTrace(resultSet constants.ExtractedData) {

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

	log.Printf("%v\n", len(resultSet))
	for _, result := range resultSet {
			log.Printf("%v\n", result.StartTime)
			log.Printf("%v\n", result.Duration)
			log.Printf(result.BTName)
			log.Printf(result.TierName)
			log.Printf(result.NodeName)
			log.Printf(result.TraceID)
			log.Printf(result.ParentID)
			log.Printf(result.SpanID)
			log.Printf("")
			log.Printf("")
			log.Printf("END OF RESULT")
			log.Printf("")
			log.Printf("")
	}


}

func sleep() {

	fmt.Print(".")
	milliseconds := 250
	time.Sleep(time.Duration(milliseconds) * time.Millisecond)

}

//func main() {
//	fmt.Println("Testing Creation of BT via ZipKin incoming Trace")
//	input1 := ExtractedData{StartTime: 1580944070, Duration: 1024, BTName: "Lester'sTrace", TierName: "GoLang", NodeName: "openTrace1"}
//  input2 := ExtractedData{StartTime: 1580944070, Duration: 1024, BTName: "Lester'sExitCall", TierName: "GoLang", NodeName: "openTrace1"}
//	btCount := 0
//	maxBtCount := 200
//	fmt.Print("Doing something")
//	for btCount < maxBtCount {
//	    		createBusinessTransaction(input1, input2)
//		}
//
//}

//func main() {
//    	fmt.Println("Testing Creation of BT via ZipKin incoming Trace")
//    	btCount := 0
//	maxBtCount := 200
//	fmt.Print("Doing something")
//	for btCount < maxBtCount {
//    		createBusinessTransaction()
//	}
//}
