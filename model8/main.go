

package main

import (
	"context"
	"fmt"
	redis "github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"sync"
	"time"
)

type Redis_conn struct {
	Host     string `yaml: "host"`
	Port     uint16 `yaml: "port"`
	Auth     string `yaml: "auth"`
	Database int    `yaml: "database"`
	Cur      int    `yaml: "cur"`
}
type RedisSingleObj struct {
	xx *Redis_conn
	Db *redis.Client
}

func (r *RedisSingleObj) InitSingleRedis(ctx context.Context) (err error) {
	redisAddr := fmt.Sprintf("%s:%d", r.xx.Host, r.xx.Port)
	fmt.Println(redisAddr)
	r.Db = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: r.xx.Auth,
		DB:       r.xx.Database,
	})
	res, err := r.Db.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Connect Failed! Err: %v\n", err)
		return err
	} else {
		fmt.Printf("Connect Successful! Ping => %v\n", res)
		return nil
	}
}
func setdemo(rdb *redis.Client, ctx context.Context, cur int, length int) {
	wg := sync.WaitGroup{}
	wg.Add(cur)

	for i := 0; i < cur; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 10000; i++ {
				n := rand.Intn(length)
				cmd := rdb.Set(ctx, randStr(n), randStr(n),time.Minute*1000)
				if err := cmd.Err(); err != nil {
					fmt.Println(err)
				}
			}

		}()
	}

	wg.Wait()

}

func main() {
	// 实例化RedisSingleObj结构体
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := ioutil.ReadFile("./host.yaml")
	if err != nil {
		fmt.Println(err)
	}
	conn := &Redis_conn{}
	err = yaml.Unmarshal(config, conn)
	fmt.Println(string(config), conn.Host)
	if err != nil {
		fmt.Println(err)
	}
	r := &RedisSingleObj{
		xx: conn,
	}
	// 初始化连接 Single Redis 服务端
	err = r.InitSingleRedis(ctx)
	if err != nil {
		panic(err)
	}
	setdemo(r.Db, ctx, 1, 10)
	setdemo(r.Db, ctx, 5, 10)
	setdemo(r.Db, ctx, 50, 10)
	setdemo(r.Db, ctx, 1, 1000)
	setdemo(r.Db, ctx, 5, 1000)
	setdemo(r.Db, ctx, 50, 1000)
	setdemo(r.Db, ctx, 1, 5000)
	setdemo(r.Db, ctx, 5, 5000)
	setdemo(r.Db, ctx, 50, 5000)
	// 程序执行完毕释放资源
	defer r.Db.Close()
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
