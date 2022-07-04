1. 使用 redis benchmark 工具, 测试 10 20 50 100 200 1k 5k 字节 value 大小，redis get set 性能

| value大小（单位：字节） | set的效率 | get的效率 |
| --- | --- | --- |
| 10 | 12500.00 requests per second | 11111.11 requests per second |
| 20 | 14285.71 requests per second | 12500.00 requests per second |
| 50 | 12500.00 requests per second | 14285.71 requests per second |
| 100 | 14285.71 requests per second | 14285.71 requests per second |
| 200 | 12500.00 requests per second | 12500.00 requests per second |
| 1000 | 10000.00 requests per second | 14285.71 requests per second |
| 5000 | 14285.71 requests per second | 14285.71 requests per second |
2.写入一定量的 kv 数据, 根据数据大小 1w-50w 自己评估, 结合写入前后的 info memory 信息 , 分析上述不同 value 大小下，平均每个 key 的占用内存空间。 代码工作原理，写入不同数量不同长度的value, 分析内存占用,  相同长度的value在写入数量越多情况下，平均每个value占用内存更多

| value大小 | count  | size  |
| --- | --- | --- |
| 长度为10的字符串 | 10000 | 867k |
| 长度为10 | 50000 | 3.77M |
| 10 | 500000 | 33.01M |
| 1000 | 10000 | 10.51M |
| 1000 | 50000 | 54.97M |
| 1000 | 500000 | 575.51M |

