# BUG_REPRO

## Bug 是什么
账单开票状态检查写错成 pending；作废账单的重复作废校验被移除；账单汇总把已开票金额漏掉；账单状态默认值写成 paid。

## 如何触发
`cd backend && go test ./...`

## 错误信息
- service.TestBillingStatusAndSummaryChain 失败：已支付账单开票报冲突、重复作废未报错、汇总已收金额少了已开票部分。
