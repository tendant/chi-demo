# HTTP Status Code

```
1. 200: Ok, normal response
   201: Created. Created record, usually for Post request.
2. 301: Moved permanently.
   302: Found.
   303: See Other.
4. 400: Bad Request. Wrong parameters
   401: Unauthorized. No token or invalid token.
   403: Forbidden. Authenticated user with valid token has no permission to access the resource.
   404: Not Found. No resource found for requested URI.
   429: Too Many Requests. Exceeded rate limit.
5. 500: Internal Server Error. Server error, e.g. database connection issue, server side code/data issue.
```

# Error response body structure


```
{
    "message": "Human readable error message",
    "code": "error code, defined by server",
    "errors": [
        {
            "message": "Human readable error message 1",
            "code": "error code, defined by server"
        },
        {
            "message": "Human readable error message 2",
            "code": "error code, defined by server"
        }
    ]
}
```


