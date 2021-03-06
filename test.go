package main

//this is a package for test
import (
	"bufio"
	"container/list"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// "test/appcli"
	"test/packageone"
	"time"
	// "github.com/jinzhu/now"
	"github.com/cihub/seelog"
	"gopkg.in/redis.v3"
	// "test/reflect"
)

func init() {
	runtime.GOMAXPROCS(2)
}

func testGetenv() {
	glog.Infoln("RESOLVE_PAYOUT====", os.Getenv("RESOLVE_PAYOUT"))
}

var makes int
var frees int

func makeBuffer() []byte {
	makes += 1
	return make([]byte, rand.Intn(5000000)+500000)
}

type queued struct {
	when  time.Time
	slice []byte
}

func makeRecycler() (get, give chan []byte) {
	get = make(chan []byte)
	give = make(chan []byte)
	go func() {
		q := new(list.List)
		for {
			if q.Len() == 0 {
				q.PushFront(queued{when: time.Now(), slice: makeBuffer()})
			}
			e := q.Front()
			timeout := time.NewTimer(time.Minute)
			select {
			case b := <-give:
				timeout.Stop()
				q.PushFront(queued{when: time.Now(), slice: b})

			case get <- e.Value.(queued).slice:
				timeout.Stop()
				q.Remove(e)

			case <-timeout.C:
				e := q.Front()
				for e != nil {
					n := e.Next()
					if time.Since(e.Value.(queued).when) > time.Minute {
						q.Remove(e)
						e.Value = nil
					}
					e = n
				}
			}
		}
	}()
	return
}

func testDecreseGC() {
	pool := make([][]byte, 20)
	get, give := makeRecycler()
	var m runtime.MemStats
	for {
		b := <-get
		i := rand.Intn(len(pool))
		if pool[i] != nil {
			give <- pool[i]
		}
		pool[i] = b
		time.Sleep(time.Second)
		bytes := 0
		for i := 0; i < len(pool); i++ {
			if pool[i] != nil {
				bytes += len(pool[i])
			}
		}
		runtime.ReadMemStats(&m)
		fmt.Printf("%d, %d, %d, %d, %d, %d, %d\n", m.HeapSys, bytes, m.HeapAlloc, m.HeapIdle, m.HeapReleased, makes, frees)
	}
}

func testList() {
	l := list.New()
	e4 := l.PushBack(4)
	e1 := l.PushFront(1)
	l.InsertBefore(3, e4)
	l.InsertAfter(2, e1)

	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

type person interface {
	say() int
}

type student struct {
	p person
}

func (s *student) say() {
	fmt.Println("this is a test for interface usage")
}

func (s *student) readbook() {
	s.say()
}

func testInterface() {
	var stu = student{}
	stu.readbook()
}

/**************Mongodb****************/
const URL = "localhost:27017"

var session *mgo.Session

func testMongodb() {
	var err error
	session, err = mgo.Dial(URL)
	if err != nil {
		glog.Infoln("cannot connect to database, please check")
		return
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("test")
	collection := db.C("c_one")
	//update and insert
	selector := bson.M{"name": "tang"}
	// updata := bson.M{"$set": bson.M{"email": []string{"11111@email.com", "22222@qq.com", "333333@163.com"}}}
	updata := bson.M{"$addToSet": bson.M{"email": "222222@163.com"}}
	updateInfo, err := collection.Upsert(selector, updata)
	if err != nil {
		glog.Infoln("upsert failed:", err)
	}
	fmt.Println("update info:", updateInfo)
}

type feeInfo struct {
	Vip  float64 `bson:"vip"`
	Norm float64 `bson:"normal"`
}

func testFindParamNil() {
	var err error
	session, err = mgo.Dial(URL)
	if err != nil {
		glog.Infoln("cannot connect to database, please check")
		return
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	fees := feeInfo{}

	db := session.DB("etc_pool")
	collection := db.C("fee_info")
	// updata := bson.M{"$set": bson.M{"email": []string{"11111@email.com", "22222@qq.com", "333333@163.com"}}}
	err = collection.Find(nil).One(&fees)
	if err != nil {
		glog.Infoln("upsert failed:", err)
	}
	fmt.Println("fee info:", fees)
}

func testmongoConnect() {
	var err error
	session, err = mgo.Dial(URL)
	defer session.Close()
	if err != nil {
		glog.Infoln("cannot connect to database, please check")
		return
	}
	session.SetMode(mgo.Monotonic, true)
	testMongodb()
}

/******************************/
//test array and slice
func testArrayAndSlice() {
	array := [5]int{}
	slice := []int{}
	fmt.Println("arry len:", len(array))
	fmt.Println("slice len:", len(slice))
	slice = append(slice, 344)
	fmt.Println("arry len:", len(array))
	fmt.Println("slice len:", len(slice))
}

//***************************/
//test reflect
type datastruct1 struct {
	name string
	age  int
}

func getdata(data interface{}) {
	typ := reflect.TypeOf(data)
	fmt.Println("type:", typ)
	fmt.Println("type:", reflect.TypeOf(typ))
}

func testReflect() {
	var tempdata = datastruct1{"gavin", 25}
	getdata(tempdata)
}

func testTrim() {
	var temp = "#sdfsdfsdf#"
	glog.Infoln(strings.Trim(temp, "#"))
}

type structnil struct {
	name string
	age  int
}

func testStructnil() {
	var temp = structnil{}
	if temp == (structnil{}) {
		glog.Infoln("struct is nil")
	}
}

func getString(value *string) {
	*value = "my name is Gavin"
}

func testPassStringAddr() {
	var str = ""

	getString(&str)
	glog.Infoln("get string value is:", str)
}

func writeUsertoMongo() {
	var err error
	session, err = mgo.Dial(URL)
	if err != nil {
		glog.Infoln("cannot connect to database, please check")
		return
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	var email = "gavin.tang@btcc.com"
	selector := bson.M{"email": email}

	db := session.DB("etc_pool")
	collection := db.C("user_info")

	collection.Upsert(selector, bson.M{"$addToSet": bson.M{"walletAddress": "0x11111111111111111111111"}})
}

func testwriteTimetoMongo() {
	var err error
	session, err = mgo.Dial(URL)
	if err != nil {
		glog.Infoln("cannot connect to database, please check")
		return
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	var email = "gavin.tang@btcc.com"
	selector := bson.M{"email": email}

	db := session.DB("etc_pool")
	collection := db.C("off_line")

	collection.Upsert(selector, bson.M{"$set": bson.M{"offlineTime": time.Now().Format("2006-01-02 15:04:05")}})
}

type offinfo struct {
	Email       string `bson:"email"`
	Offlinetime string `bson:"offlineTime"`
}

type userinfo struct {
	Email  string   `bson:"email"`
	Wallet []string `bson:"walletAddress"`
}

func testgettimestamp() {
	var err error
	session, err = mgo.Dial(URL)
	if err != nil {
		glog.Infoln("cannot connect to database, please check")
		return
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	info := offinfo{}
	var email = "gavin.tang@btcc.com"
	selector := bson.M{"email": email}

	db := session.DB("etc_pool")
	collection := db.C("off_line")
	collection.Find(selector).One(&info)
	glog.Infoln("time stamp is:", info)

}

func testgetuserinfo() {
	var err error
	session, err = mgo.Dial(URL)
	if err != nil {
		glog.Infoln("cannot connect to database, please check")
		return
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	info := userinfo{}
	var email = "gavin.tang@btcc.com"
	selector := bson.M{"email": email}

	db := session.DB("etc_pool")
	collection := db.C("user_info")
	collection.Find(selector).One(&info)
	glog.Infoln("info is:", info)

}

func emailcheck(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	glog.Infoln("request url:", req.URL.Path, req.RequestURI)
	glog.Infoln("request method:", req.Method)
	glog.Infoln("request post data:", req.PostFormValue("email"))
}

func testHttpServer() {
	http.HandleFunc("/etcpool/usercheck", emailcheck)
	err := http.ListenAndServe("localhost:8090", nil)
	if err != nil {
		glog.Infoln("http server error:", err)
	}
}

type jsonmsg struct {
	Result bool   `json:"result"`
	Msg    string `json:"msg"`
}

func testHttpClient() {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, time.Second*3)
				if err != nil {
					glog.Infoln("dail btcc server timeout", err)
					return nil, err
				}
				return c, nil
			},
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * 2,
		},
	}
	form := url.Values{}
	var ts = strconv.FormatInt(time.Now().Unix(), 10)
	var email = "tangzc64@163.com"
	form.Set("email", email)
	form.Set("ts", ts)

	str := "e10adc3949" + "#" + ts + "#" + email
	h := md5.New()
	h.Write([]byte(str))
	hs := h.Sum(nil)
	var hashstring = hex.EncodeToString(hs)
	form.Set("hash", hashstring)
	glog.Infoln("md5 hash:", hashstring)
	glog.Infoln("format int:", strconv.FormatInt(time.Now().Unix(), 10))
	b := strings.NewReader(form.Encode())
	request, _ := http.NewRequest("POST", "https://api.btcc.com/proapi/checkemail", b)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		glog.Infoln("client do request error:", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		glog.Infoln("read body error:", err)
	}
	glog.Infoln("get response data:", string(body))
	res := jsonmsg{}
	json.Unmarshal(body, &res)
	glog.Infoln("response result:", res.Result)
	if res.Result {
		glog.Infoln("hahahhahaha, email exist")
	}

}

func testFormatTime() {
	// t := time.Now()
	// glog.Infoln("time.now:", t)
	// time.Sleep(time.Second*1)
	// t1 := time.Now()
	// glog.Infoln("t1:", t1)
	// glog.Infoln("time duration:", t1.Sub(t))
	// glog.Infoln("************************")
	glog.Infoln("time Now():", time.Now(), "time Unix():", time.Now().Unix())                    //time Now(): 2016-11-02 18:03:00.79858886 +0800 CST time Unix(): 1478080980
	glog.Infoln("time.Now().format():", time.Now().Format("2006-01-02 15:04:05"))                //time.Now().format(): 2016-11-02 18:03:00
	glog.Infoln("time.Unix().format():", time.Unix(1478079819, 0).Format("2006-01-02 15:04:05")) // time.Unix().format(): 2016-11-02 17:43:39
	glog.Infoln("***********************")
	timestap, _ := time.Parse("2006-01-02 15:04:05", "2016-11-02 17:47:10")
	h, m, s := timestap.Clock()
	glog.Infoln("hour:", h, "min:", m, "second:", s)
	year, month, day := timestap.Date()
	glog.Infoln("year:", year, "month:", month, "day:", day)

	times, _ := time.Parse("2006-01-02 15:04:05", "2016-12-12 17:47:10")
	y1, m1, d1 := times.Date()
	if y1 == year {
		glog.Infoln("year equal")
	}
	if m1 != month {
		glog.Infoln("month not equal")
	}
	if d1 == day {
		glog.Infoln("day equal")
	}
	timestap2, _ := time.Parse("2006-01-02 15:04:05", "2016-11-02 16:47:10")

	glog.Infoln("time duration:", timestap.Sub(timestap2))

	glog.Infoln("************************************")

	glog.Infoln("time.Unix().format():", time.Unix(time.Now().Unix(), 0).Format("20060102"))
}

func testSeelog() {
	logger, err := seelog.LoggerFromConfigAsFile("logconfig.xml")

	if err != nil {
		seelog.Critical("err parsing config log file", err)
		return
	}
	seelog.ReplaceLogger(logger)

	for i := 0; ; i++ {
		seelog.Error("seelog", i)
		seelog.Info("seelog info", i)
		seelog.Debug("seelog debug", i)
		time.Sleep(time.Minute * 1)
	}

}

func testSeelogConfigButLogtoWriteLog() {
	logger, err := seelog.LoggerFromConfigAsFile("logconfig.xml")
	if err != nil {
		seelog.Critical("err parsing config log file", err)
		return
	}
	seelog.ReplaceLogger(logger)

	defer seelog.Flush()
	one.One()
	params := "btcc"

	for i := 0; i < 999999; i++ {
		seelog.Error("seelog", i)
		seelog.Info("seelog info", i)
		seelog.Debug("seelog debug", i)
		seelog.Critical("critical error test", i)
		seelog.Infof("this is a string %s", params)
		if i == 99999 {
			time.Sleep(time.Second * 1)
		}
	}

}

func testToLower() {
	str := "0x0eE4c03776EFe873465cF3d999f09552a124c841"
	str = strings.ToLower(str)
	fmt.Println("str:", str)
}

//testRedis 使用散列来存储用户数据, value的类型为自定义结构体，将其序列化为json数据后存储值散列中
type redisClient struct {
	cli    *redis.Client
	prefix string
}

type userShare struct {
	Sharetimes int32
	Value      int64
}

func join(args ...interface{}) string {
	s := make([]string, len(args))
	for i, v := range args {
		switch v.(type) {
		case string:
			s[i] = v.(string)
		case int64:
			s[i] = strconv.FormatInt(v.(int64), 10)
		case uint64:
			s[i] = strconv.FormatUint(v.(uint64), 10)
		case float64:
			s[i] = strconv.FormatFloat(v.(float64), 'f', 0, 64)
		case bool:
			if v.(bool) {
				s[i] = "1"
			} else {
				s[i] = "0"
			}
		default:
			panic("Invalid type specified for conversion")
		}
	}
	return strings.Join(s, ":")
}

func (redis *redisClient) formatKey(args ...interface{}) string {
	return join(redis.prefix, join(args...))
}

func (redis *redisClient) setRedis(usersMap map[string]userShare) {
	tx := redis.cli.Multi()
	defer tx.Close()

	tx.Exec(func() error {
		//value to json object
		for key, value := range usersMap {
			v, _ := json.Marshal(value)
			tx.HSet(redis.formatKey("usershares"), key, string(v))
		}
		return nil
	})
}

func (redis *redisClient) getRedis() {
	cmd := redis.cli.HGetAllMap(redis.formatKey("usershares"))
	if cmd.Err() != nil {
		fmt.Println(cmd.Err())
		return
	}
	userMap, _ := cmd.Result()
	for key, value := range userMap {
		fmt.Println("key:", key)
		temp := []byte(value)
		var data userShare
		json.Unmarshal(temp, &data)
		fmt.Println("shareTimes :", data.Sharetimes, "value :", data.Value)
	}
}

func testRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
		PoolSize: 10,
	})
	backend := redisClient{cli: client, prefix: "test"}
	usersMap := make(map[string]userShare)
	usersMap["111111"] = userShare{Sharetimes: int32(10), Value: int64(111111)}
	usersMap["222222"] = userShare{Sharetimes: int32(10), Value: int64(222222)}

	// backend.setRedis(usersMap)
	backend.getRedis()
}

//=================================================
func testRedis2() {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
		PoolSize: 10,
	})
	backend := redisClient{cli: client, prefix: "test"}

	//pop and push new block to the redis
	blockshare := make(map[string]int64)
	blockshare["7777"] = 7777
	blockshare["8888"] = 8888
	backend.setRedis2(blockshare)
	backend.getRedis2()
}

func (backend *redisClient) setRedis2(blockshare map[string]int64) {
	var shareN = 3
	v, _ := json.Marshal(blockshare)

	lencmd := backend.cli.LLen("test:nblocksshares")
	n, _ := lencmd.Result()
	if lencmd.Err() == redis.Nil || int(n) < shareN {
		backend.cli.RPush("test:nblocksshares", string(v))
		fmt.Println("rpush success")
		return
	} else if lencmd.Err() != nil {
		fmt.Println(lencmd.Err())
	}

	tx := backend.cli.Multi()
	defer tx.Close()
	_, err := tx.Exec(func() error {
		tx.LPop("test:nblocksshares")
		tx.RPush("test:nblocksshares", string(v))
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (backend *redisClient) getRedis2() {
	cmd := backend.cli.LRange("test:nblocksshares", 0, -1)
	if cmd.Err() == redis.Nil { //first time, this key dont exist, so return nil
		fmt.Println(cmd.Err())
	} else if cmd.Err() != nil {
		fmt.Println(cmd.Err())
		return
	}
	nblocks := make([]map[string]int64, 0)
	stringArray, _ := cmd.Result()
	for _, substring := range stringArray {
		var data map[string]int64
		json.Unmarshal([]byte(substring), &data)
		nblocks = append(nblocks, data)
	}
	for _, block := range nblocks {
		for key, value := range block {
			fmt.Println("user key :", key, "user share :", value)
		}
		fmt.Println("==========================")
	}
}

//=================================================

func testPass(arg ...string) {
	t := reflect.TypeOf(arg)
	fmt.Println("type:", t)

	fmt.Println(arg)
}

func testPassParamPPP() {
	strs := []string{"1", "2", "3"}
	testPass("nihao", "tac")
	testPass(strs...)
}

func testStructUnderLine() {
	temp := one.A{}
	temp.C = 16
}

type b struct {
	one.B
}

func (b *b) say(word string) {
	fmt.Println("say:", word)
}

func testInheir() {
	temp := b{}
	temp.say("youyou")
}

func addvalue(m map[string]int) {
	m["2"] = 2
}

type struct1 struct {
	A int
}

func testMap() {
	temp := make(map[string]int)
	temp["1"] = 1
	addvalue(temp)
	temp["3"] = 3
	for key, v := range temp {
		fmt.Println("key:", key, "value:", v)
	}
	fmt.Println("==============")
	tempmap := make(map[string]struct1)
	tempmap["1"] = struct1{A: 111}
	tempmap["2"] = struct1{A: 222}
	for key, v := range tempmap {
		fmt.Println("key:", key, "value:", v)
	}
}

//===================test map to json object and Unmarshal
type structMap struct {
	A int32
	B int64
}

func testMapToJSON() {
	tempMap := make(map[string]structMap)
	tempMap["111111"] = structMap{A: int32(123), B: int64(321)}
	tempMap["222222"] = structMap{A: int32(456), B: int64(654)}

	v, _ := json.Marshal(tempMap)

	fmt.Println("Marshal value :", v)

	var data map[string]structMap
	json.Unmarshal(v, &data)
	for key, value := range data {
		fmt.Println("key:", key)
		fmt.Println("struct A:", value.A, "struct B:", value.B)
	}
}

func testInt64() {
	var num int64 = 601
	num = num * (1 / 2)
	fmt.Println("value:", num)
}

func funcA(temp *[]map[int]int) {
	var tempmap = make(map[int]int)
	tempmap[1] = 1
	*temp = append(*temp, tempmap)
}

func funcB(temp *[]map[int]int) {
	tempmap := make(map[int]int)
	tempmap[2] = 2
	*temp = append(*temp, tempmap)
}

func funcC() []map[int]int {
	temp := make([]map[int]int, 0)
	for i := 0; i < 3; i++ {
		t := make(map[int]int)
		t[i] = i
		temp = append(temp, t)
	}
	return temp
}

func funcD(temp *[]map[int]int) {
	// *temp = make([]map[int]int, 0)
	for i := 0; i < 3; i++ {
		t := make(map[int]int)
		t[i] = i
		*temp = append(*temp, t)
	}
}

func testVarArea() {
	a := 1
	var temp []map[int]int
	if a == 1 {
		funcA(&temp)
	} else if a == 2 {
		funcB(&temp)
	} else {
		// temp = funcC()
		funcD(&temp)
	}
	fmt.Println("array:", len(temp))
	fmt.Println("value:", temp)
}

//测试整点时间以及localTime
func testTimeIntegerTime() {
	t := time.Now().Local()
	str := t.Format("2006-01-02 15:04:05")
	fmt.Println(str)
	for {
		t = time.Now()
		h, m, s := t.Clock()
		fmt.Println("hour:", h, "min:", m, "second:", s)
		if m == 0 && s <= 10 {
			fmt.Println("recording...., it's a integer clock")
		}
		time.Sleep(time.Second * 10)
	}

}

func testStringSplit() {
	tempstring := "123456:asdfgh:kkkk"
	subs := strings.Split(tempstring, ":")
	for _, str := range subs {
		fmt.Println(str)
	}
}

func testTimestampDecrease() {
	t1, _ := time.Parse("2006-01-02 15:04:05", "2016-11-30 23:59:59")

	t2, _ := time.Parse("2006-01-02 15:04:05", "2016-12-01 23:59:59")
	glog.Infoln("time decrese:", t2.Unix()-t1.Unix())
}

func testGetIntegerTime() {
	ti := time.Now().Unix()
	t := time.Unix(ti, 0)
	str := t.Format("2006-01-02 15:00:00")
	fmt.Println("time:", str)

	t2, err := time.Parse("2006-01-02 15:00:00", str)
	if err != nil {
		fmt.Println("error:", err)
	}
	temp := t2.Unix() - 3600
	t = time.Unix(temp, 0)
	str2 := t.Format("2006-01-02 15:00:00")
	fmt.Println("time:", str2)

	fmt.Println("==============")
	stamp := time.Now().Local().Unix()
	fmt.Println(time.Unix(stamp, 0).Format("2006-01-02 15:00:00"))
	stamp = stamp - 3600
	fmt.Println(time.Unix(stamp, 0).Format("2006-01-02 15:00:00"))
	stamp = stamp - 3600
	fmt.Println(time.Unix(stamp, 0).Format("2006-01-02 15:00:00"))

	fmt.Println("============")
	ttt := time.Now().Local()
	strr := ttt.Format("2006-01-02 15:00:00")
	fmt.Println("time:", strr)
}

type reflectS struct {
	a string
	b int
}

func testReflectType() {
	temp := &reflectS{}
	t := reflect.TypeOf(temp)
	fmt.Println("type:", t)
}

func testReverse() {
	var maparray []map[string]interface{}
	for i := 0; i < 7; i++ {
		temp := make(map[string]interface{})
		temp[strconv.Itoa(i)] = i
		maparray = append(maparray, temp)
	}
	fmt.Println("before reverse:", maparray)
	middle := make(map[string]interface{})
	length := len(maparray)
	for i := 0; i < length/2; i++ {
		middle = maparray[i]
		maparray[i] = maparray[length-i-1]
		maparray[length-i-1] = middle
	}
	fmt.Println("after sort:", maparray)
}

func testInterfaceMaptoString() {
	tempmap := make(map[string]interface{})
	tempmap["a"] = "this is a string"
	tempmap["b"] = 12
	tempmap["c"] = structMap{A: 66, B: 77}
	tempmap["d"] = int64(44)
	for key, v := range tempmap {
		switch v.(type) {
		case string:
			fmt.Println("key:", key, "this is type string:", v.(string))
		case int:
			fmt.Println("key:", key, "this is type int:", v.(int))
		case int64:
			fmt.Println("key:", key, "this is type int64:", v.(int64))
		default:
			fmt.Println("key:", key, "this is type default:", v)
		}
	}
}

func testMapIsNil() {
	tempmap := make(map[string]string)
	tempmap["a"] = "a"
	fmt.Println("map key b:", tempmap["b"])
	fmt.Println("key b type is :", reflect.TypeOf(tempmap["b"]))
	v, _ := strconv.ParseInt(tempmap["b"], 10, 32)
	fmt.Println("value:", v)
}

func testStringlen() {
	var str = ""
	fmt.Println("string length:", len(str))
}

type InfoA struct {
	A int
	B int
}

type InfoB struct {
	C InfoA
	D int
}

func testStructValuePass() {
	tempA := []InfoA{InfoA{1, 2}, InfoA{3, 4}}
	fmt.Println("temp A length:", len(tempA))

	tempB := make([]InfoB, len(tempA))
	for idx, value := range tempA {
		tempB[idx].C = value
	}
	for _, v := range tempB {
		fmt.Println("tempB InfoA value:", v.C)
	}
}

type reflectStructA struct {
	Name string  `json:"name"`
	Age  int32   `json:"age"`
	S    float64 `json:"s"`
}

func testStructReflect() {
	temp := reflectStructA{Name: "gavin"}
	fmt.Println("age:", temp.Age)
	fmt.Println("s:", temp.S)
	if temp.S == 0 {
		fmt.Println("hahahahahh")
	}
	val := reflect.ValueOf(&temp).Elem()

	fmt.Println("number filed :", val.NumField())
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		typefield := val.Type().Field(i)

		fmt.Println("type:", typefield.Type.String())
		fmt.Println("filed name :", typefield.Name, "field value :", field.Interface())
	}
}

func testAtoi() {
	v, err := strconv.Atoi("")
	fmt.Println("value:", v)
	fmt.Println("err:", err)
}

func testSlice() {
	slice := make([]int, 0)
	for i := 0; i < 10; i++ {
		slice = append(slice, i)
	}
	fmt.Println("slice data:", slice[0:0])
}

func testRedisPubSub() {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
		PoolSize: 10,
	})
	pubsub, err := client.Subscribe("pool.blocknotify1st")
	if err != nil {
		glog.Info("subscribe redis error:", err)
	}
	defer pubsub.Close()
	var dur time.Duration
	for {
		t := time.Now()
		glog.Info("waiting time...:", dur)
		msg, err := pubsub.ReceiveMessage()
		if err != nil {
			glog.Info("get message error:", err)
		}
		dur = time.Since(t)
		glog.Info("channel:", msg.Channel, "msg :", msg.Payload)
	}
}

func testInterfaceCommaOK() string {
	var temp interface{}
	temp = "this is a test string"
	return temp.(string)
}

var (
	ee string
	ff = testInterfaceCommaOK()
)

func testVar() {
	ee = "this is ee"
	fmt.Println(ee)
	fmt.Println(ff)
}

func testStringCompare() {
	str1 := "201707140017"
	datestring := time.Unix(time.Now().Unix(), 0).Format("20060102")
	if str1[:8] == datestring {
		glog.Infoln("date :", str1[:8])
	}
	var number string
	number = str1[8:]
	num, _ := strconv.Atoi(number)
	glog.Infoln("out int :", num)
	d := strconv.Itoa(num)
	glog.Infoln("num :", d, "len :", len(d))
	for len(d) < 4 {
		d = "0" + d
	}
	glog.Infoln("num :", d, "len :", len(d))
	str1b := []byte(str1)
	str1b = append(str1b, []byte(","+"201707140018")...)

	glog.Infoln("after append :", string(str1b))
}

func testStringArray() {
	var strs []string
	for _, val := range []string{"1", "2", "3"} {
		strs = append(strs, val)
	}
	glog.Infoln("strs :", strs)

	strs = []string{"4", "5", "6"}
	glog.Infoln("strs :", strs)
}

func testArrayCopy() {
	var arr1 [3]int
	var arr2 [3]int
	for i := 0; i < 3; i++ {
		arr1[i] = i
	}
	arr2 = arr1
	fmt.Printf("arr1 %v ,address arr1 %p :\n", arr1, &arr1)
	fmt.Printf("arr2 %v ,address arr2 %p :\n", arr2, &arr2)
	fmt.Println("=========================")

	var slice1 []int
	var slice2 []int
	for i := 0; i < 5; i++ {
		slice1 = append(slice1, i)
	}
	slice2 = slice1
	fmt.Println(slice1, &slice1)
	fmt.Println(slice2, &slice2)
}

func testCopySlice() {
	slice := make([]int, 10)
	for i := 0; i < 10; i++ {
		slice[i] = i
	}
	fmt.Println("origin :", slice)
	changeslice(slice)
	fmt.Println("now :", slice)
}

func changeslice(slice []int) {
	slice[0] = 88
}

func testInterfaceSwitch() {
	var temp interface{}
	temp = 123
	switch temp.(type) {
	case int:
		fmt.Println("this is int")
	case string:
		fmt.Println("this is string")
	case float32:
		fmt.Println("this is float")
	default:
		fmt.Println("default type")
	}
}

func testdefer() {
	defer fmt.Println("1")
	defer fmt.Println("2")
	defer fmt.Println("3")
	fmt.Println("let's begin ...")
}

func testSliceandArrayCopy() {
	slice := make([]int, 10)
	for i := 0; i < 10; i++ {
		slice[i] = i
	}

	slice2 := make([]int, 10)
	copy(slice2[:], slice[:])
	fmt.Println("slice2 :", slice2)
	////array

	var array [10]int
	for i := 0; i < 10; i++ {
		array[i] = i
	}

	var array2 [10]int
	copy(array2[:], array[:])

	fmt.Println("array2 :", array2)
}

func testRoutine() {
	for {
		go fmt.Print(1)
		fmt.Print(0)
	}

}

type Payload struct {
	Abi      string        `json:"abi"`
	Function string        `json:"function"`
	Params   []interface{} `json:"params"`
}

func transferObjecttoString() {
	load := Payload{
		Abi:      "this is abi string",
		Function: "this is function",
		Params:   []interface{}{123, 456},
	}
	b, err := json.Marshal(load)
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	fmt.Println("marshal string :", string(b))
}

func genEvmAccount() {
	//touch new file
	f, err := os.Create("./accounts.txt")
	if err != nil {
		fmt.Println("err when create file:", err)
		return
	}
	defer f.Close()

	// cmd.Stdout = f
	//gen evm account
	for i := 1; i < 201; i++ {
		w := bufio.NewWriter(f)
		num := fmt.Sprintf("%d", i)
		fmt.Fprintln(w, num)
		cmd := exec.Command("./chorustool", "account", "geneth")
		bs, err := cmd.Output()
		if err != nil {
			fmt.Println("write key err :", err)
			return
		}
		fmt.Println(string(bs))
		fmt.Fprintln(w, string(bs))
		w.Flush()
	}

	//write to file
}

func testLabel() {
	fmt.Println("this is test for label")
OUTER:
	for i := 0; i < 10; i++ {
		fmt.Println("outer for ", i)
		for j := 0; j < 5; j++ {
			if j > 3 {
				fmt.Println("inner for ", j)
				//break OUTER
				goto OUTER
			}

		}
	}
	fmt.Println("over...")
}

//测试内存模型
func testMemoryModel() {
	wg := sync.WaitGroup{}
	var count int
	//var ch = make(chan bool, 1)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			//ch <- true
			count++
			time.Sleep(time.Millisecond)
			count--
			//<-ch
			wg.Done()
		}()
	}
	wg.Wait()
}

func testCloseChannel() {
	//donnt judge close channel
	//c := make(chan bool, 5)
	//c <- true
	//c <- true
	//c <- true
	//close(c)
	//for {
	//	fmt.Println(<-c)
	//	time.Sleep(time.Millisecond * 500)
	//}

	//judge close channel
	c := make(chan int, 10)
	c <- 1
	c <- 2
	c <- 3
	close(c)

	for {
		i, isClose := <-c
		fmt.Println(i, " ", isClose)
		time.Sleep(time.Millisecond * 500)
	}
}

func testMakeMap() {
	m := make(map[string]int, 10)
	m["s"] = 10
	fmt.Println("s :", m["s"])
}

type Compare struct {
	A int
	B string
}

func testStructCompare() {
	a := Compare{A: 1, B: "a"}
	b := Compare{A: 1, B: "a"}
	fmt.Println(a == b)
}

func main() {
	// flag.Set("log_dir", "./logs")
	// flag.Set("alsologtostderr", "true")
	flag.Set("logtostderr", "true")
	flag.Parse()

	//testGetenv()
	//testDecreseGC()
	//testList()
	//testInterface()
	// testMongodb()
	//testArrayAndSlice()
	// testmongoConnect()
	// testReflect()
	//testTrim()
	//testStructnil()
	// testPassStringAddr()
	//testwriteTimetoMongo()
	// testFormatTime()
	// testgettimestamp()
	//testgetuserinfo()
	//  writeUsertoMongo()
	// testHttpServer()
	// testHttpClient()
	//testSeelog()
	// testSeelogConfigButLogtoWriteLog()
	// testToLower()
	// testRedis()
	// testRedis2()
	// testPassParamPPP()
	//testStructUnderLine()
	// testInheir()
	// testMap()
	// testMapToJSON()
	// testInt64()
	// testVarArea()
	// testTimeIntegerTime()
	// testStringSplit()
	// testTimestampDecrease()
	// testGetIntegerTime()
	// testReflectType()
	// testReverse()
	// testInterfaceMaptoString()
	// testMapIsNil()
	// testStringlen()
	// testStructValuePass()
	// testStructReflect()
	// testAtoi()
	// testSlice()
	// testFindParamNil()
	// testRedisPubSub()
	// fmt.Println(testInterfaceCommaOK())
	// testVar()
	// cliTest.TestCLi()
	// ref.TestReflect1()
	// ref.TestReflect2()
	// testStringCompare()
	// testStringArray()
	//testArrayCopy()
	//testCopySlice()

	//testInterfaceSwitch()
	//testdefer()
	//testSliceandArrayCopy()
	// testRoutine()
	//transferObjecttoString()
	//genEvmAccount()
	//testLabel()
	//testMemoryModel()
	//testCloseChannel()
	//testMakeMap()
	testStructCompare()
}
