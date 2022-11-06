* Start server with air

    $ air
    
* Test event with curl

    curl --location --request POST 'localhost:3000' \
    --header 'ce-type: com.example.repro.create' \
    --header 'Content-Type: application/cloudevents+json' \
    --data-raw '{
        "specversion" : "1.0",
        "type" : "com.github.pull_request.opened",
        "source" : "https://github.com/cloudevents/spec/pull",
        "subject" : "123",
        "id" : "A234-1234-1234",
        "time" : "2018-04-05T17:31:00Z",
        "comexampleextension1" : "value",
        "comexampleothervalue" : 5,
        "datacontenttype" : "text/xml",
        "data" : "<much wow=\"xml\"/>"
    }'