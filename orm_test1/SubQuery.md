# 第十周作业：设计子查询

**作业：**
支持子查询

## 场景分析

1. 作为from后面的衍生表出现，类似于如下这种
```azure
SELECT ... FROM (subquery) [AS] tbl_name ...
```

2. 在where后面的比较符出现(having也是同理)
```azure
SELECT * FROM t1
  WHERE column1 = (SELECT MAX(column2) FROM t2);
```
3. 在上面的场景下，还会添加any，some,all 而in exist这种前面没有操作符
```azure
## 这几种类型前面都会有比较符然后再接子查询，使用形式如下
SELECT s1 FROM t1 WHERE s1 > ANY (SELECT s1 FROM t2);
SELECT s1 FROM t1 WHERE s1 <> SOME (SELECT s1 FROM t2);
SELECT s1 FROM t1 WHERE s1 > ALL (SELECT s1 FROM t2);
## IN 和 EXIST没有操作符
SELECT column1 FROM t1 WHERE EXISTS (SELECT * FROM t2);
```
4. 嵌套查询
- 子查询和join查询
```azure
join查询左右连接这一个子查询
  SELECT `sub`.`item_id` FROM (`order` JOIN (SELECT * FROM `order_detail`) AS `sub` ON `id` = `sub`.`order_id`)

```
这里涉及到一个问题。join查询时涉及到的两个表相同字段，怎么指定是哪个的字段呢。让subquery也实现tablereference接口。







5. 作为标量使用（这种使用很少不支持）


6. 对select 后面的字段名进行校验，当select的字段不存在于后面的表里面时需要报错，这就涉及到一种情况。子查询可能会在后面出现了select detail_id 。。。这种情况。子查询外面的字段使用子查询表中其他但没有在子查询的返回结果中。







## 行业方案
GORM
```azure

gorm是将整个子查询作为参数传入where中,gorm在buildsql时会断言你的参数是不是db类型。

db.Where("amount > (?)", db.Table("orders").Select("AVG(amount)")).Find(&orders)
// SELECT * FROM "orders" WHERE amount > (SELECT AVG(amount) FROM "orders");

subQuery := db.Select("AVG(age)").Where("name LIKE ?", "name%").Table("users")
db.Select("AVG(age) as avgage").Group("name").Having("AVG(age) > (?)", subQuery).Find(&results)
// SELECT AVG(age) as avgage FROM `users` GROUP BY `name` HAVING AVG(age) > (SELECT AVG(age) FROM `users` WHERE name LIKE "name%")
```


