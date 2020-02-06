package main

import (
    "encoding/json"
    "log"
    "net/http"
    "io/ioutil"
    "../constants"
)

type Message []struct {
	TraceID     string `json:"traceId"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	ParentID    string `json:"parentId,omitempty"`
	Annotations []struct {
		Timestamp int64  `json:"timestamp"`
		Value     string `json:"value"`
		Endpoint  struct {
			ServiceName string `json:"serviceName"`
			Ipv4        string `json:"ipv4"`
		} `json:"endpoint"`
	} `json:"annotations"`
	BinaryAnnotations []struct {
		Key      string `json:"key"`
		Value    string `json:"value"`
		Endpoint struct {
			ServiceName string `json:"serviceName"`
			Ipv4        string `json:"ipv4"`
		} `json:"endpoint"`
	} `json:"binaryAnnotations"`
}

func ingestSpan(rw http.ResponseWriter, req *http.Request) {
    
    /*get JSON body from handler*/
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        log.Printf("verbose error info: %#v", err)
    }

    /*Get Traces from Span by unmarshall JSON*/
    traces := &Message{}
    err = json.Unmarshal(body, &traces)
    if err != nil {
        log.Printf("error decoding span input: %v", err)
        if e, ok := err.(*json.SyntaxError); ok {
            log.Printf("syntax error at byte offset %d", e.Offset)
        }
        log.Printf("span response: %q", body)
    }

    extractDataFromSpan(traces)
}

/*Now that we have the JSON data in our custom type struct, we can extract out the needed data to generate BT Snapshot information showing the overall transaction
and related traces*/
func extractDataFromSpan(traces *Message){

    var resultSet = make(constants.ExtractedData, len(*traces))
    var startTime int64
    var endTime int64
    var counter int = 0
    var index int

    for _, trace := range *traces {

        /*Index manipulation to make sure Parent is first index (0) slot*/
        index = counter + 1
        counter += 1

        startTime = 0
        endTime = 0
        
        if(trace.ParentID == ""){
            /*Set the top of the []extractedData for the Parent Information*/
            index = 0
        }
        resultSet[index].TraceID = trace.TraceID
        resultSet[index].ParentID = trace.ParentID
        resultSet[index].SpanID = trace.ID

        for _, annotation := range trace.Annotations {
            if(annotation.Value =="ss"){
                startTime = annotation.Timestamp
                resultSet[index].StartTime = annotation.Timestamp
                resultSet[index].TierName = annotation.Endpoint.ServiceName
                resultSet[index].NodeName = annotation.Endpoint.Ipv4 
            }
            if(annotation.Value =="sr"){
                endTime = annotation.Timestamp
            }  
        }
        for _, binaryAnn := range trace.BinaryAnnotations {
            if(binaryAnn.Key =="http.url"){
                resultSet[index].BTName = binaryAnn.Value
            }
        }

        resultSet[index].Duration = startTime-endTime
        
    }

    /*For Debugging ResultSet from JSON import
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
    */

    /*Calling Span stitching for BT Snapshots function*/
}

func main() {
	http.HandleFunc("/span", ingestSpan)
    log.Fatal(http.ListenAndServe(":3030", nil))
}