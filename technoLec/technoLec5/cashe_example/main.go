package main

import (
	"github.com/garyburd/redigo/redis"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"encoding/json"
	"reflect"
	"time"
)

var (
	c redis.Conn
	db *sql.DB
)

type CacheItem struct {
	Data interface{}
	Tags map[string]int
}

type Articles []Article

type Article struct {
	Name string
	From string
}

func GetCacheRecord(key string) string {
	item, err := redis.String(c.Do("GET", key))
	if err == redis.ErrNil {
		fmt.Println("Record not found")
		return ""
	} else if err != nil {
		PanicOnErr(err)
	}
	return item
}

func GetFioByID(id uint64) (string, error)  {
	fmt.Println("call sql")
	var fio string
	row := db.QueryRow("SELECT `fio` FROM `students` WHERE `id` = ?", id)
	err := row.Scan(&fio)
	return fio, err
}

func getCachedFio(ket string) (string, error)  {
	fmt.Println("regis get ", ket)
	data, err := c.Do("GET", ket)
	item, err := redis.String(data, err)
	if err == redis.ErrNil {
		fmt.Println("Record not found")
		return "", err
	} else if err != nil {
		return "", err
	}
	return item, nil
}

func main()  {
	var err error

	c, err = redis.DialURL("redis://user:@localhost:6379/0")
	PanicOnErr(err)
	defer c.Close()

	db, err = sql.Open("mysql", "root@tcp(localhost:3306)/golang_test?charset=utf8&interpolateParams=true")
	PanicOnErr(err)

	fmt.Println("---- taggedCache")
	taggedCache()
	fmt.Println("--- lockCacheRebuild")
	lockCacheRebuild()
}

func taggedCache()  {
	top := Articles{
		Article{"Джава и Докер - это должен знать каждый", "Хабр"},
		Article{"Как взрываются базовые станиции", "Гиктаймс"},
	}

	item := CacheItem{
		Data: top,
		// Tags это метки валидности записи в кеше
		// данная запись если эти метки в кеше имеют такое же значение
		Tags: map[string]int {
			"Habr": 1,
			"GT": 1,
		},
	}

	jsonData, _ := json.Marshal(item)
	fmt.Println("json to store: ", string(jsonData))

	mkey := "top_news_mobile"
	result, err := redis.String(c.Do("SET", mkey, jsonData))
	fmt.Println("result" ,result)
	if result != "OK" {
		panic("result not ok: " + result)
	}

	c.Do("SET", "Habr", 1)
	c.Do("SET", "GT", 1)

	// если раскомментировать эту строчку, то наш кеш перестанет быть валидным
	// c.Do("INCR", "Habr")

	topCache, err := getCachedFio(mkey)
	fmt.Println("top Cache", topCache)

	cItems := CacheItem{}
	_ = json.Unmarshal([]byte(topCache), &cItems)
	fmt.Printf("top Cache unpacked %+v\n", cItems)

	keys := make([]interface{}, 0)
	toCompare := make([]int, 0)
	for k, v := range cItems.Tags {
		keys = append(keys, k)
		toCompare = append(toCompare, v)
	}
	reply, err := redis.Ints(c.Do("MGET", keys...))
	PanicOnErr(err)

	fmt.Println("compare cached values", toCompare, "with current values", reply)
	fmt.Println("cache record is valid:", reflect.DeepEqual(toCompare, reply))
}

func lockCacheRebuild()  {
	userId := 7
	mkey := "top_user_" + string(userId)

	var fio string
	var err error
	// если кеш есть, то мы выходим сразу, если нет - у нас 4 попытки получить значение
	for i := 0; i < 4; i++ {
		fio, err = getCachedFio(mkey)
		if err == redis.ErrNil {
			// пытаемся сказать "я строю этот кеш, другие - ждите"
			lockStatus, _ := redis.String(c.Do("SET", mkey+"_lock", fio, "EX", 3, "NX"))
			if lockStatus != "OK" {
				// кто-то другой держит лок, подождём и попробуем получить запись еще раз
				fmt.Println("sleep", i)
				time.Sleep(time.Millisecond * 10)
			} else {
				// успешло залочились, можем строить кеш
				break
			}
		} else if err != nil {
			PanicOnErr(err)
		} else {
			//запись нашлас
			break
		}
	}
	// если записи нету, то надо её построить и положить туда
	if err == redis.ErrNil {
		fmt.Println("Create cache data")
		// основная работа по доставанию данных для кеша
		// потенциально тяжелая операция
		fio, err = GetFioByID(1)
		if err != nil {
			PanicOnErr(err)
		}

		ttl := 50
		// добавляет запись, https://redis.io/commands/set
		result, err := redis.String(c.Do("SET", mkey, fio, "EX", ttl))
		PanicOnErr(err)
		if result != "OK" {
			panic("result not ok: " + result)
		}
		// удаляем лок на построение
		n, err := redis.Int(c.Do("DEL", mkey+"_lock"))
		PanicOnErr(err)
		fmt.Println("loc deleted: ", n)
	}
	fmt.Println(fio)
}

func PanicOnErr(err error)  {
	if err != nil {
		panic(err)
	}
}