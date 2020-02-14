# 抖音签名服务

此项目主要学习下Golang的交叉编译

主要原理通过 Golang 的Mobile交叉编译到越狱后的iOS设备上 加载抖音签名动态库调用对应函数 抖音动态签名库(AwemeDylib)自行砸壳

编译成功后将砸壳后的签名库放至越狱设备的`/var/lib/ibreaker/libs`目录下

将编译好的程序放至`/usr/libexec`下

## 准备
#### 一台越狱的iOS设备
#### 一台Mac或者黑苹果

## 调用

```json
POST http://<device-ip>:25583/sign/aweme/sign
Content-Type: application/json

{
    "URL": "需要Gorgon参数的完整URL地址",
    "Cookie": "可选参数如果没有不填",
    "Stub": "如果请求的原始地址为Post则需要对Body进行md5并传入此参数"
}
```

```json
{
    "X-Gorgon": "xxxxxxxxxxxxxxxxxxxxxxxx",
    "X-Khronos": "返回的timestamp"
}
```
