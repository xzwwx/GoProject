package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
)

func StartHttpServer() bool {

	http.HandleFunc("/login", login)

	return true
}

var pool *redis.Pool

// redis List 用的key名
var listkey = "ListKey"

// // 连接池
// func newPool(server, password string) *redis.Pool {
// 	return &redis.Pool{
// 		MaxIdle:     64,
// 		MaxActive:   1000,
// 		IdleTimeout: 240 * time.Second,
// 		Dial: func() (redis.Conn, error) {
// 			c, err := redis.Dial("tcp", server)
// 			if err != nil {
// 				fmt.Println("err:", err)
// 				return nil, err
// 			}
// 			if _, err := c.Do("AUTH", password); err!= nil{
// 				c.Close()
// 				fmt.Println("DB Auth Err:", err)
// 				return nil, err
// 			}
// 			return c, err
// 		},
// 		TestOnBorrow: func(c redis.Conn, t time.Time) error {
// 			if time.Since(t) < time.Minute {
// 				return nil
// 			}
// 			_, err := c.Do("PING")
// 			return err
// 		},
// 	}
// }

func RecordHandler(w http.ResponseWriter, r *http.Request) {
	// 获取text的内容
	values := r.FormValue("text")
	if values == "" {
		fmt.Fprintf(w, "try to use \"/add?text=xxx!\"")
		return
	}
	// 返回text内容
	fmt.Fprintf(w, values)
	// 记录IP:时间:text
	temp_values := fmt.Sprintf("%s:%s:%s", r.RemoteAddr, GetDateFormat(), values)
	// 连接池不为空
	if pool != nil {
		// 获取连接
		conn := pool.Get()
		defer conn.Close()
		// 将该条记录放入redis list结构中
		if _, err := conn.Do("rpush", listkey, temp_values); err != nil {
			fmt.Fprintln(w, " DB err:", err)
		}
	}
}

func EchoAllHandler(w http.ResponseWriter, r *http.Request) {
	if pool != nil {
		// 获取连接
		conn := pool.Get()
		defer conn.Close()

		// 使用LLEN 命令获取list长度
		listLen, err := redis.Int(conn.Do("LLEN", listkey))
		if err != nil {
			fmt.Println("Redis LLEN error!")
			return
		}

		// 获取list的所有数据
		values, err := redis.MultiBulk(conn.Do("LRANGE", listkey, 0, listLen))
		items, err := redis.Strings(values, nil)
		for _, s := range items {
			fmt.Fprintf(w, s+"\n")
		}

	}

}

func ClearAllHandler(w http.ResponseWriter, r *http.Request) {
	if pool != nil {
		conn := pool.Get()
		defer conn.Close()

		// 删除该key，即可删除所有记录
		_, err := conn.Do("DEL", listkey)
		if err != nil {
			fmt.Println("Redis DEL error!")
			return
		} else {
			fmt.Fprintf(w, "clear success!")
		}
	}
}

// 时间戳转年月日 时分秒
func GetDateFormat() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// func main1() {
// 	pool = newPool("127.0.0.1:6379", "123456")

// 	mux := http.NewServeMux()
// 	// 注册
// 	mux.HandleFunc("/add", RecordHandler)
// 	mux.HandleFunc("/all", EchoAllHandler)
// 	mux.HandleFunc("/clean", ClearAllHandler)

// 	server := &http.Server{
// 		Addr:    ":8081",
// 		Handler: mux,
// 	}

// 	// 创建系统信号接收器
// 	done := make(chan os.Signal, 1)
// 	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
// 	go func() {
// 		<-done

// 		if err := server.Shutdown(context.Background()); err != nil {
// 			log.Fatal("Shutdown server:", err)
// 		}
// 	}()

// 	err := server.ListenAndServe()
// 	if err != nil {
// 		if err == http.ErrServerClosed {
// 			log.Print("Server closed under request")
// 		} else {
// 			log.Fatal("Server closed unexpected")
// 		}
// 	}
// }
