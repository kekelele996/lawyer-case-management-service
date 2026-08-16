# BUG_REPRO

## Bug 是什么
案件状态机去掉了向前一步的回退；状态流转把管理员也纳入 canFlow 限制；按主办律师查询把等值过滤写成不等值；案件状态默认值写成 closed。

## 如何触发
`cd backend && go test ./...`

## 错误信息
- service.TestCaseStatusAndListChain 失败：回退流转被拒、管理员跨状态被拒、按律师查不到案件。
