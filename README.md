# mloginsvr
mloginsvr+Gin搭建的登入服务器架构。它意在帮助游戏服务端开发快速的搭建“登入服务器”。


## 环境
- **Golang**  1.15.6
- **Gin**  框架
- **mysql**
- **Redis**


## 服务器架构
![](https://github.com/stringMao/mloginsvr/raw/main/readme/image/url-1.jpg)
- 1：client账号密码登入
- 2：登入服务器进行数据库验证，并且生成token
- 3：将token和大厅服务器地址返回给client
- 4：client用token登入大厅
- 5：大厅向登入服务器验证token，并且获得账号信息



## 功能
1. 登入
   -  账号密码登入
   -  负载均衡
   -  渠道分发
   -  token验证
2. 注册
   - 手机号注册
     - 短信验证码
   - 修改默认昵称
     - 敏感词检测
   - 密码
     - 忘记密码
3. 第三方登入
   - 微信登入
4. 实名认证
   





