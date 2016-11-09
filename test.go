package main
//this is a package for test
import (
  "github.com/golang/glog"
  "fmt"
  "runtime"
  "os"
  "flag"
  "net"
  "math/rand"
  "time"
  "container/list"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "reflect"
  "strings"
  "net/http"
  "net/url"
  "io/ioutil"
  "strconv"
  "crypto/md5"
  "encoding/hex"
  "encoding/json"
  // "github.com/jinzhu/now"
)

func testGetenv(){
  glog.Infoln("RESOLVE_PAYOUT====",os.Getenv("RESOLVE_PAYOUT"))
}

var makes int
var frees int

func makeBuffer()[]byte{
  makes += 1
  return make([]byte, rand.Intn(5000000) + 500000)
}

type queued struct {
  when time.Time
  slice []byte
}

func makeRecycler()(get, give chan []byte){
  get = make(chan []byte)
  give = make(chan []byte)
  go func(){
    q := new(list.List)
    for {
      if q.Len() == 0{
         q.PushFront(queued{when: time.Now(), slice: makeBuffer()})
      }
      e := q.Front()
      timeout := time.NewTimer(time.Minute)
      select {
      case b := <-give:
        timeout.Stop()
        q.PushFront(queued{when:time.Now(), slice: b})

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

func testDecreseGC(){
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
    for i := 0; i < len(pool); i++{
      if pool[i] != nil {
        bytes += len(pool[i])
      }
    }
    runtime.ReadMemStats(&m)
    fmt.Printf("%d, %d, %d, %d, %d, %d, %d\n", m.HeapSys, bytes, m.HeapAlloc, m.HeapIdle, m.HeapReleased, makes, frees)
  }
}

func testList(){
  l := list.New()
  e4 := l.PushBack(4)
  e1 := l.PushFront(1)
  l.InsertBefore(3, e4)
  l.InsertAfter(2, e1)

  for e := l.Front(); e != nil; e = e.Next(){
    fmt.Println(e.Value)
  }
}

type person interface {
    say() int
}

type student struct {
    p person
}

func (s *student)say(){
    fmt.Println("this is a test for interface usage")
}

func (s *student)readbook(){
    s.say()
}

func testInterface(){
    var stu = student{}
    stu.readbook()
}
/**************Mongodb****************/
const URL = "localhost:27017"
var session *mgo.Session
func testMongodb(){
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

func testmongoConnect(){
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
func testArrayAndSlice(){
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
  age int
}

func getdata(data interface{}){
  typ := reflect.TypeOf(data)
  fmt.Println("type:", typ)
  fmt.Println("type:", reflect.TypeOf(typ))
}

func testReflect(){
  var tempdata = datastruct1{"gavin", 25}
  getdata(tempdata)
}

func testTrim(){
  var temp = "#sdfsdfsdf#"
  glog.Infoln(strings.Trim(temp, "#"))
}

type structnil struct {
    name string
    age int
}

func testStructnil(){
  var temp = structnil{}
  if temp == (structnil{}) {
    glog.Infoln("struct is nil")
  }
}

func getString(value *string){
  *value = "my name is Gavin"
}

func testPassStringAddr(){
  var str = ""

  getString(&str)
  glog.Infoln("get string value is:", str)
}

func writeUsertoMongo(){
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

func testwriteTimetoMongo(){
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
  Email string `bson:"email"`
  Offlinetime string `bson:"offlineTime"`
}

type userinfo struct{
  Email   string  `bson:"email"`
  Wallet  []string `bson:"walletAddress"`
}

func testgettimestamp(){
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

func testgetuserinfo(){
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

func emailcheck(res http.ResponseWriter, req *http.Request){
  req.ParseForm()
  glog.Infoln("request url:", req.URL.Path, req.RequestURI)
  glog.Infoln("request method:", req.Method)
  glog.Infoln("request post data:", req.PostFormValue("email"))
}

func testHttpServer(){
  http.HandleFunc("/etcpool/usercheck", emailcheck)
  err := http.ListenAndServe("localhost:8090", nil)
  if err != nil{
    glog.Infoln("http server error:", err)
  }
}

type jsonmsg struct {
  result bool
  msg string
}

func testHttpClient(){
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
  var email = "007pig@gmail.com"
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
  request, _ := http.NewRequest("POST", "http://10.0.20.246/proapi/checkemail", b)
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
  glog.Infoln("response result:", res.result)

}

func testFormatTime(){
  // t := time.Now()
  // glog.Infoln("time.now:", t)
  // time.Sleep(time.Second*1)
  // t1 := time.Now()
  // glog.Infoln("t1:", t1)
  // glog.Infoln("time duration:", t1.Sub(t))
  // glog.Infoln("************************")
  glog.Infoln("time Now():", time.Now(),"time Unix():", time.Now().Unix())//time Now(): 2016-11-02 18:03:00.79858886 +0800 CST time Unix(): 1478080980
  glog.Infoln("time.Now().format():", time.Now().Format("2006-01-02 15:04:05"))//time.Now().format(): 2016-11-02 18:03:00
  glog.Infoln("time.Unix().format():", time.Unix(1478079819, 0).Format("2006-01-02 15:04:05"))// time.Unix().format(): 2016-11-02 17:43:39
  glog.Infoln("***********************")
  timestap, _ := time.Parse("2006-01-02 15:04:05", "2016-11-02 17:47:10")

  timestap2, _ := time.Parse("2006-01-02 15:04:05", "2016-11-02 16:47:10")


  glog.Infoln("time duration:", timestap.Sub(timestap2))
}

func main(){
    flag.Set("logtostderr", "true")
    flag.Parse()

    //testGetenv()
    //testDecreseGC()
    //testList()
    //testInterface()
    //testMongodb()
    //testArrayAndSlice()
    // testmongoConnect()
    //testReflect()
    //testTrim()
    //testStructnil()
    // testPassStringAddr()
    //testwriteTimetoMongo()
    //testFormatTime()
    // testgettimestamp()
    //testgetuserinfo()
     writeUsertoMongo()
    // testHttpServer()
    //testHttpClient()
}
