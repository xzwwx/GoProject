package main
import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"context"
)
//var (
//	rdb *redis.Client
//)
//func initClient() (err error){
//	rdb = redis.NewClient(&redis.Options{
//		Addr:     "127.0.0.1:6379",
//		Password: "123456",
//		DB:       0,
//		PoolSize: 100,
//	})
//
//	ctx, cancel :=context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	_, err = rdb.Ping(ctx).Result()
//	return err
//}

//
//func createClient() *redis.Client {
//	client := redis.NewClient(&redis.Options{
//		Addr:     "127.0.0.1:6379",
//		Password: "123456",
//		DB:       0,
//		PoolSize: 100,
//	})
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
//	defer cancel()
//	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
//	_, err := client.Ping(ctx).Result()
//	if err != nil{
//		panic(err)
//	}
//
//	return client
//}

func generate_Incr(client *redis.Client, incrName string)(incr_Key string){

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	result, err := client.Incr(ctx,"incrName").Result()
	if err != nil{
		panic(err)
	}
	// incrName :="userID"
	incr_Key = incrName + ":" + strconv.FormatInt(result, 10)

	return incr_Key
}

func get_lately_Incr(client *redis.Client, incrName string)(incr_Key string){
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	incr_Key, err :=client.Get(ctx, incrName).Result()
	if err != nil{
		panic(err)
	}


	return incr_Key
}

//func password_md5(password, salt string)string{
//	h := md5.New()
//	h.Write([]byte(salt + password))
//	return fmt.Sprintf("%x", h.Sum(nil))
//}
//
//type User struct{
//	userName string
//	userEmail string
//	userPasswd string
//	userTel string
//	userKey string
//	isNew bool
//}

func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func register_user(client *redis.Client, user User)string{
	if user.userName == ""  || user.userPasswd == ""{
		return "用户名或者密码不能为空"
	}
	if VerifyEmailFormat(user.userEmail) == false {
		return "邮箱格式不对"
	}

	nameKey := "user_Name:"+strings.ToLower(strings.TrimSpace(user.userName))

	//ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	//defer cancel()

	ctx := context.Background()

	cmd, err := client.Exists(ctx, nameKey).Result()
	if err != nil{
		panic(err)
	}
	if 0 != cmd{
		return "该用户名已经被注册"
	}

	emailkey := "user_Email:"+user.userEmail
	cmd, err = client.Exists(ctx, emailkey).Result()
	if err != nil{
		panic(err)
	}
	if 0 != cmd{
		return "该邮箱已经被注册"
	}

	uid_Key := generate_Incr(client, "userID")
	sta, err := client.Set(ctx, nameKey,uid_Key, 0).Result()
	if err != nil{
		panic(err)
	}
	if sta != "OK"{
		panic(sta)
	}
	sta, err = client.Set(ctx, emailkey,uid_Key, 0).Result()
	if err != nil{
		panic(err)
	}
	if sta != "OK"{
		panic(sta)
	}

	mess := map[string]interface{}{
		"user_Name":user.userName,
		"user_Email" : user.userEmail,
		"user_Passwd": user.userPasswd,
		"create_Time" :time.Now().Format("2006-01-02 15:04:05"),
	}

	if user.userTel != ""{
		mess["user_Tel"] = user.userTel
	}
	if user.userKey != ""{
		mess["user_Key"] = user.userKey
	}

	_, err = client.HMSet(ctx, uid_Key, mess).Result()
	if err != nil{
		panic(err)
	}
	return  "register successfully"
}



//func register(client *redis.Client, newResister User){
//	//var newResister User
//	newResister.userName = "xzw"
//	newResister.userEmail = "11111@qq.com"
//	newResister.userPasswd = password_md5("123456", "yan")
//	newResister.userTel = "17826111082"
//
//	status := register_user(client, newResister)
//	if status != "register successfully"{
//		panic(status)
//	}
//
//}

//jieshou kehuduan de POST xinxi
func rigester(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	// 1. 请求类型是application/x-www-form-urlencoded时解析form数据
	//r.ParseForm()
	//fmt.Println(r.PostForm) // 打印form数据
	//fmt.Println(r.PostForm.Get("name"), r.PostForm.Get("age"))
	// 2. 请求类型是application/json时从r.Body读取数据
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read request.Body failed, err:%v\n", err)
		return
	}
	data:= make(map[string]string)
	err = json.Unmarshal([]byte(body), &data)
	if err !=nil {
		fmt.Println("body To Json failed, err:", err)
		panic("Error")
	}
	fmt.Println(data)

	//rigest
	client := createClient()
	fmt.Println(client)
	var newResister User

	newResister.userName = data["userName"]
	newResister.userEmail = data["userEmail"]
	newResister.userPasswd = password_md5(data["userPasswd"], data["userName"])
	newResister.userTel = data["userTel"]
	newResister.userKey = data["userKey"]

	status := register_user(client, newResister)
	if status != "register successfully"{
		panic(status)
	}
	fmt.Println(status)

	answer := `{"status": "ok"}`
	w.Write([]byte(answer))
}

//func main(){
//
//	err := http.ListenAndServe(":9091", nil)
//	log.Fatal(err)
//
//}
