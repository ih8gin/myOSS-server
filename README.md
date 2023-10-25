## [English](./README_en.md)
# 简介
使用<b>go</b>实现的分布式对象存储服务，使用<b>es</b>做元数据存储、<b>rabbitmq</b>做不同节点之间的异步通信
参考[分布式对象存储：原理、架构及Go语言实现/胡世杰著.-北京:人民邮电出版社,2018.6]一书
# 开发环境
go 1.21.1

rabbitMQ 3.12.4

elasticsearch 8.9.0
# 详细介绍
整个项目包括三个模块：
* apiServer模块负责对外接受调用请求
* dataServer负责数据的存储与读取
* maintenanceUtils包含一些实现了数据清洗、复原等保持服务健康的脚本文件
### apiServer
>cd [工程路径]<BR>
RABBITMQ_SERVER=[rabbitMQ服务地址] ES_SERVER=[es服务地址] LISTEN_ADDRESS=[服务监听地址与端口] go run apiServer/apiServerApplication.go
>
### dataServer
>cd [工程路径]<BR>
RABBITMQ_SERVER=[rabbitMQ服务地址] ES_SERVER=[es服务地址] LISTEN_ADDRESS=[服务监听地址与端口] STORAGE_ROOT=[本地文件存储路径] go run dataServer/dataServerApplication.go
>
### maintenanceUtils
>cd [工程路径]<BR>
RABBITMQ_SERVER=[rabbitMQ服务地址] ES_SERVER=[es服务地址] go run ./maintenanceUtils/deleteOldMetadata/versionLimit.go <BR>
RABBITMQ_SERVER=[rabbitMQ服务地址] ES_SERVER=[es服务地址] STORAGE_ROOT=[本地文件存储路径] LISTEN_ADDRESS=127.0.0.1:55556 go run ./maintenanceUtils/deleteOrphanObject/deleteOrphanObject.go <BR>
RABBITMQ_SERVER=[rabbitMQ服务地址] ES_SERVER=[es服务地址] STORAGE_ROOT=[本地文件存储路径] go run ./maintenanceUtils/ObjectScanner/scan.go
>
# 使用方法
使用上述指令成功启动apiServer与dataServer后，可以按如下api调用服务

|                           路径                           | header                                                                | body  | reply         | 功能              |
|:------------------------------------------------------:|:----------------------------------------------------------------------|:------|:--------------|:----------------|
|          GET [apiServer]/objects/[objectName]          | range：[起始字节数]-[结束字节数]（可选）<br/>Accept-Encoding:gzip（可选）                |       | 对象            | 下载对象            |
| GET [apiServer]/objects/[objectName]?version=[version] | range：[起始字节数]-[结束字节数]（可选）<br/>Accept-Encoding:gzip（可选）                |       | 对象            | 下载指定版本的对象       |
|          PUT [apiServer]/objects/[objectName]          | Digest:SHA256=[上传文件的SHA256哈希值的base64编码字符串]（必须）                        | 待上传文件 |               | 上传对象            |
|                GET [apiServer]/versions                |                                                                       |       | 版本信息          | 获取全部已存储对象的有效版本号 |
|         GET [apiServer]/versions/[objectName]          |                                                                       |       | 版本信息          | 获取指定名称对象的有效版本号  |
|        DELETE [apiServer]/objects/[objectName]         |                                                                       |       |               | 删除对象            |
|         POST [apiServer]/versions/[objectName]         | Digest:SHA256=[上传文件的SHA256哈希值的base64编码字符串]（必须）<br/>Size:[待传输文件大小]（必须） |       | 临时对象标识tempKey | 请求断续上传对象        |
|             PUT [apiServer]/temp/[tempKey]             | Size:[待传输文件大小]（必须）                                                    | 待上传文件 |               | 断续上传文件          |
|            HEAD [apiServer]/temp/[tempKey]             |                                                                       |       | 已上传文件大小       | 检查断续上传文件已上传的大小  |