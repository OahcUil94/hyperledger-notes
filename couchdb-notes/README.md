# CouchDB

web界面操作窗口: http://ip:port/_utils

## docker运行

docker run -e COUCHDB_USER=admin -e COUCHDB_PASSWORD=admin123 -p 5984:5984 --name mycouchdb -d couchdb:3.1

## 基本操作

创建数据库:
curl -X PUT http://127.0.0.1:5984/baseball

## 常见问题

### 更新文档

```json
{
  "error": "conflict",
  "reason": "Document update conflict."
}
```

在更新文档的时候, 需要提交上一次创建或更新生成的版本号, 也就是附带_rev字段
