package main

import (
    "encoding/json"
    "log"
    "net/http"
    "io/ioutil"
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

type extractedData struct {
	StartTime int64    `json:"StartTime"`
	Duration  int64    `json:"Duration"`
	BTName    string `json:"BTName"`
	TierName  string `json:"TierName"`
	NodeName  string `json:"NodeName"`
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

/*Now that we have the JSON data in our custom type struct, we can extract out the needed data to generate BT naming, Response time and
 Tier and Node names*/
func extractDataFromSpan(traces *Message){
    for _, trace := range *traces {
        if(trace.ParentID == ""){
            /*Found the Parent Trace*/
            var startTime int64
            var endTime int64
            var BTName string
            var TierName string
            var NodeName string

            for _, annotation := range trace.Annotations {
                if(annotation.Value =="ss"){
                    startTime = annotation.Timestamp
                    TierName = annotation.Endpoint.ServiceName
                    NodeName = annotation.Endpoint.Ipv4 
                }
                if(annotation.Value =="sr"){
                    endTime = annotation.Timestamp
                }  
            }
            for _, binaryAnn := range trace.BinaryAnnotations {
                if(binaryAnn.Key =="http.url"){
                    BTName = binaryAnn.Value
                }
            }
            /*Calling BT compliers funcytion*/

            output := extractedData{StartTime: startTime, Duration: startTime-endTime, BTName: BTName, TierName: TierName, NodeName: NodeName}
            log.Printf(output.BTName)
        }
    }
}

func main() {
	http.HandleFunc("/span", ingestSpan)
    log.Fatal(http.ListenAndServe(":3030", nil))
}