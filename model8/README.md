
1. 使用 redis benchmark 工具, 测试 10 20 50 100 200 1k 5k 字节 value 大小，redis get set 性能

value大小（单位：字节）	         set的效率	          get的效率
10	                  12500.00 requests per second	11111.11 requests per second
20	                  14285.71 requests per second	12500.00 requests per second
50	                  12500.00 requests per second	14285.71 requests per second
100                   14285.71 requests per second	14285.71 requests per second
200	                  12500.00 requests per second	12500.00 requests per second
1000	                10000.00 requests per second	14285.71 requests per second
5000	                 14285.71 requests per second	14285.71 requests per second
bhj

2.
