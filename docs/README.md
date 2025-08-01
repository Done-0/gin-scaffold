接口文档

## 统一响应格式：

- 正确响应：

```json
{
  "data": any,
  "requestId": string,
  "timeStamp": number
}
```

- 错误响应：

```json
{
  "code": number,
  "msg": string,
  "data": any,
  "requestId": string,
  "timeStamp": number
}
```

## test 测试模块

1. **testPing** 测试接口
   - 请求方式：GET
   - 请求路径：/api/v1/test/testPing
   - 请求参数：无
   - 响应示例：
   ```text
   Pong successfully!
   ```
2. **testHello** 测试接口
   - 请求方式：GET
   - 请求路径：/api/v1/test/testHello
   - 请求参数：无
   - 响应示例：
   ```text
   Hello, metaphysics 🎉!
   ```
3. **testLogger** 测试接口
   - 请求方式：GET
   - 请求路径：/api/v1/test/testLogger
   - 请求参数：无
   - 响应示例：
   ```text
   测试日志成功!
   ```
4. **testRedis** 测试接口
   - 请求方式：GET
   - 请求路径：/api/v1/test/testRedis
   - 请求参数：无
   - 响应示例：
   ```text
   测试缓存功能完成!
   ```
5. **testSuccessRes** 测试接口
   - 请求方式：GET
   - 请求路径：/api/v1/test/testSuccessRes
   - 请求参数：无
   - 响应示例：
   ```json
   {
     "data": "测试成功响应成功!",
     "requestId": "XtZvqFlDtpgzwEAesJpFMGgJQRbQDXyM",
     "timeStamp": 1740118491
   }
   ```
6. **testErrRes** 测试接口
   - 请求方式：GET
   - 请求路径：/api/v1/test/testErrRes
   - 请求参数：无
   - 响应示例：
   ```json
   {
     "code": 10000,
     "msg": "服务端异常",
     "data": {},
     "requestId": "BRnzCMxAoprBllAuBGPWqoDNofArbuOX",
     "timeStamp": 1740118534
   }
   ```
7. **testErrorMiddleware** 测试接口

   - 请求方式：GET
   - 请求路径：/api/v1/test/testErrorMiddleware
   - 请求参数：无
   - 响应示例：无

8. **testLongReq** 测试接口
   - 请求方式：GET
   - 请求路径：/api/v2/test/testLongReq
   - 请求参数：无
   - 响应示例：
   ```text
   模拟耗时请求处理完成!
   ```
