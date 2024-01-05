package pkg

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"net/http"
	"time"
)

func SetRedis() {
	// 创建 Redis 客户端实例
	rdb := redis.NewClient(&redis.Options{
		Addr:     "113.31.115.135:6379", // Redis 服务器地址
		Password: "123456",              // Redis 访问密码，如果没有密码则留空
		DB:       0,                     // Redis 数据库索引
	})

	// 使用 Ping 命令测试与 Redis 的连接
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("无法连接到 Redis:", err)
		return
	}
	fmt.Println("连接成功:", pong)

	// 设置键值对
	err = rdb.Set(context.Background(), "mykey", "myvalue", 0).Err()
	if err != nil {
		fmt.Println("设置键值对失败:", err)
		return
	}
	fmt.Println("键值对设置成功")

	// 获取键的值
	val, err := rdb.Get(context.Background(), "mykey").Result()
	if err != nil {
		fmt.Println("获取键值失败:", err)
		return
	}
	fmt.Println("键值为:", val)

	// 设置带过期时间的键值对
	err = rdb.Set(context.Background(), "mykey2", "myvalue2", 60*time.Second).Err()
	if err != nil {
		fmt.Println("设置带过期时间的键值对失败:", err)
		return
	}
	fmt.Println("带过期时间的键值对设置成功")

	// 等待一段时间，使带过期时间的键过期
	time.Sleep(15 * time.Second)

	// 获取过期的键
	val, err = rdb.Get(context.Background(), "mykey2").Result()
	if err == redis.Nil {
		fmt.Println("键已过期")
	} else if err != nil {
		fmt.Println("获取键值失败:", err)
		return
	} else {
		fmt.Println("键值为:", val)
	}

	// 关闭与 Redis 的连接
	err = rdb.Close()
	if err != nil {
		fmt.Println("关闭连接失败:", err)
		return
	}
	fmt.Println("连接已关闭")
}

func HandleRedis(w http.ResponseWriter, r *http.Request) {

}
