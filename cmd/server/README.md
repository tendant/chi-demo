# Start server

    $ air
    
# Test API

## Get Request

    curl --location --request GET 'localhost:3000/get?id=testid&name=test name&notExported=test should not show this value'
    
## Post Request

    curl --location --request POST 'localhost:3000/form' \
    --form 'name="test form name"' \
    --form 'description="test form description"'
    
## Post with JSON Body

    curl --location --request POST 'localhost:3000/body' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "name": "test name",
        "description": "test description"
    }'