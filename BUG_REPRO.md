# BUG_REPRO

## Bug 是什么
错误包装用 %v 丢掉哨兵；账单、案件、客户查询未命中返回新错误而不是 ErrNotFound，导致 errors.Is 失效。

## 如何触发
`cd backend && go test ./...`

## 错误信息
- service.TestErrorChainSentinels 失败：FindByID 未命中不再匹配 ErrNotFound，创建账单未报 not found。
