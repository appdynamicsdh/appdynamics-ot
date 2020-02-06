#!/bin/sh  
#usage:
#Following command will send test.json to AppD OT collector 30 times and sleep for 60 seconds between calls
#./sendJSONSpansLoop.sh test.json 30
#
for ((i=1;i<=$2;i++))
do   
curl -X POST -H "Content-Type: application/json" -d @$1 "http://localhost:3030/span"
sleep 60
done