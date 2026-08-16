# BUG_REPRO

## Bug 是什么
文档上传把合法类型校验取反；按案件查文档排序方向写反；关键词查询改查 file_url 而不是标题；文档类型默认值写错。

## 如何触发
`cd backend && go test ./...`

## 错误信息
- service.TestDocumentUploadAndOrderChain 失败：合法类型被拒、文档列表顺序反了。
