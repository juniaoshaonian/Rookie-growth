| value大小（单位：字节） | set的效率 | get的效率 |
| --- | --- | --- |
| 10 | 12500.00 requests per second | 11111.11 requests per second |
| 20 | 14285.71 requests per second | 12500.00 requests per second |
| 50 | 12500.00 requests per second | 14285.71 requests per second |
| 100 | 14285.71 requests per second | 14285.71 requests per second |
| 200 | 12500.00 requests per second | 12500.00 requests per second |
| 1000 | 10000.00 requests per second | 14285.71 requests per second |
| 5000 | 14285.71 requests per second | 14285.71 requests per second |


| value大小 | count  | size  |
| --- | --- | --- |
| 长度为10的字符串 | 10000 | 867k |
| 长度为10 | 50000 | 3.77M |
| 10 | 500000 | 33.01M |
| 1000 | 10000 | 10.51M |
| 1000 | 50000 | 54.97M |
| 1000 | 500000 | 575.51M |

