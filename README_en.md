## [中文](./README.md)
# brief introduction
this is an object storage system implemented with <b>golang</b>. Metadata stored on <b>elasticsearch</b>. Nodes asynchronously communicate with others with <b>rabbitmq</b>.
reference:[分布式对象存储：原理、架构及Go语言实现/胡世杰著.-北京:人民邮电出版社,2018.6]
# dev env
go 1.21.1

rabbitMQ 3.12.4

elasticsearch 8.9.0
# detail of project
the whole project consists of three modules：
* apiServer receive http request and provide service
* dataServer process local stored data
* maintenanceUtils consists of several scripts to clean orphan data, recover damaged data and etc 
### apiServer
>cd [project dir]<BR>
RABBITMQ_SERVER=[rabbitMQ addr] ES_SERVER=[es addr] LISTEN_ADDRESS=[listen addr&port] go run apiServer/apiServerApplication.go
>
### dataServer
>cd [project dir]<BR>
RABBITMQ_SERVER=[rabbitMQ addr] ES_SERVER=[es addr] LISTEN_ADDRESS=[listen addr&port] STORAGE_ROOT=[storage root dir] go run dataServer/dataServerApplication.go
>
### maintenanceUtils
>cd [project dir]<BR>
RABBITMQ_SERVER=[rabbitMQ addr] ES_SERVER=[es addr] go run ./maintenanceUtils/deleteOldMetadata/versionLimit.go <BR>
RABBITMQ_SERVER=[rabbitMQ addr] ES_SERVER=[es addr] STORAGE_ROOT=[storage root dir] LISTEN_ADDRESS=127.0.0.1:55556 go run ./maintenanceUtils/deleteOrphanObject/deleteOrphanObject.go <BR>
RABBITMQ_SERVER=[rabbitMQ addr] ES_SERVER=[es addr] STORAGE_ROOT=[storage root dir] go run ./maintenanceUtils/ObjectScanner/scan.go
>
# how to use
deploy apiServer and dataServer, use api as below

|                           路径                           | header                                                                                         | body           | reply                         | 功能                                      |
|:------------------------------------------------------:|:-----------------------------------------------------------------------------------------------|:---------------|:------------------------------|:----------------------------------------|
|          GET [apiServer]/objects/[objectName]          | range：[startByte]-[endByte]（optional）<br/>Accept-Encoding:gzip（optional）                       |                | object                        | download object                         ||
| GET [apiServer]/objects/[objectName]?version=[version] | range：[startByte]-[endByte]（optional）<br/>Accept-Encoding:gzip（optional）                       |                | object                        | download object of certain version      |
|          PUT [apiServer]/objects/[objectName]          | Digest:SHA256=[base64-encoded SHA256-hash]（mandatory）                                          | uploading file |                               | upload object                           |
|                GET [apiServer]/versions                |                                                                                                |                | version info                  | get versions info of all objects        |
|         GET [apiServer]/versions/[objectName]          |                                                                                                |                | version info                  | get versions info of certain object     |
|        DELETE [apiServer]/objects/[objectName]         |                                                                                                |                |                               | delete object                           |
|         POST [apiServer]/versions/[objectName]         | Digest:SHA256=[base64-encoded SHA256-hash]（mandatory）<br/>Size:[size of file]（mandatory）       |                | tempKey                 | request for partial upload              |
|             PUT [apiServer]/temp/[tempKey]             | Size:[size of uploading file]（mandatory）                                                       | uploading file |                               | partial upload                          |
|            HEAD [apiServer]/temp/[tempKey]             |                                                                                                |                | size of already uploaded part | check the size of already uploaded part |