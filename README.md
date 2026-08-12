# EventBatch

采集事件 → 按大小分批 → 导出到 sink（文件/内存）。用于演示批处理管道。

## 测试

```bash
set GOTOOLCHAIN=local
go test ./... -count=1
```

## 演示

```bash
go run ./cmd/batcher -input ./testdata/sample.jsonl -batch 3
```
