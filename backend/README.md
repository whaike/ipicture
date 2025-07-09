# ipicture-backend

## 启动方式

```bash
cd backend
# 启动服务
 go run main.go
```

## 主要接口

- POST /api/register  用户注册
  - 参数: {"username": "用户名", "password": "密码"}
- POST /api/login     用户登录
  - 参数: {"username": "用户名", "password": "密码"}
  - 返回: {"token": "JWT Token"}

后续会补充文件上传、图片/视频查询等接口。 