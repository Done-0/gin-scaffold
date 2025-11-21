# API Documentation

## Unified Response Format

- Successful Response:

```json
{
  "data": any,
  "requestId": string,
  "timeStamp": number
}
```

- Error Response:

```json
{
  "code": number,
  "msg": string,
  "data": any,
  "requestId": string,
  "timeStamp": number
}
```

## test Module

1. **testPing** Test Endpoint
   - HTTP Method: GET
   - Request Path: /api/v1/test/testPing
   - Request Parameters: None
   - Response Example:
   ```json
   {
     "data": {
       "time": "2025-09-26T01:46:57+08:00",
       "message": "Pong successfully!"
     },
     "requestId": "01d01617-cb23-46ec-85f1-777eeba3377c",
     "timeStamp": 1758822417
   }
   ```
2. **testHello** Test Endpoint
   - HTTP Method: GET
   - Request Path: /api/v1/test/testHello
   - Request Parameters: None
   - Response Example:
   ```json
   {
     "data": {
       "version": "1.0.0",
       "message": "Hello, gin-scaffold! 🎉!"
     },
     "requestId": "b42eb8af-b48d-48cd-8c15-f3cd52860d11",
     "timeStamp": 1758822421
   }
   ```
3. **testLogger** Test Endpoint
   - HTTP Method: GET
   - Request Path: /api/v1/test/testLogger
   - Request Parameters: None
   - Response Example:
   ```json
   {
     "data": {
       "level": "info",
       "message": "Log test succeeded!"
     },
     "requestId": "a74cfa1d-c313-45c4-bc1d-0a0c998d3e60",
     "timeStamp": 1758822424
   }
   ```
4. **testRedis** Test Endpoint
   - HTTP Method: POST
   - Request Path: /api/v1/test/testRedis
   - Request Parameters:
   ```json
   {
     "key": "test",
     "value": "hello",
     "ttl": 60
   }
   ```
   - Response Example:
   ```json
   {
     "data": {
       "key": "test",
       "value": "hello",
       "ttl": 60,
       "message": "Cache functionality test completed!"
     },
     "requestId": "XtZvqFlDtpgzwEAesJpFMGgJQRbQDXyM",
     "timeStamp": 1740118491
   }
   ```
5. **testSuccessRes** Test Endpoint
   - HTTP Method: GET
   - Request Path: /api/v1/test/testSuccessRes
   - Request Parameters: None
   - Response Example:
   ```json
   {
     "data": {
       "status": "success",
       "message": "Successful response validation passed!"
     },
     "requestId": "7f114931-51bc-47d5-922f-208ca9d86445",
     "timeStamp": 1758822431
   }
   ```
6. **testErrRes** Test Endpoint
   - HTTP Method: GET
   - Request Path: /api/v1/test/testErrRes
   - Request Parameters: None
   - Response Example:
   ```json
   {
     "data": {
       "code": 10001,
       "message": "Server exception"
     },
     "requestId": "79768196-75cc-4b9e-8286-998a4bd4218b",
     "timeStamp": 1758822435
   }
   ```
7. **testErrorMiddleware** Test Endpoint

   - HTTP Method: GET
   - Request Path: /api/v1/test/testErrorMiddleware
   - Request Parameters: None
   - Response Example: Recovery middleware handles panic and returns empty response

8. **testLongReq** Test Endpoint
   - HTTP Method: POST
   - Request Path: /api/v2/test/testLongReq
   - Request Parameters:
   ```json
   {
     "duration": 3
   }
   ```
   - Response Example:
   ```json
   {
     "data": {
       "duration": 3,
       "message": "Simulated long-running request completed!"
     },
     "requestId": "caecc92a-0e04-4b4a-ac9e-cdbba2cc34ad",
     "timeStamp": 1758822445
   }
   ```
