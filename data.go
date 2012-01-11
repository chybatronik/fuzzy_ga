package data

import (
	"fmt"
	"sort"
	"sync"
)

type NF struct {
	Name  string
	Value float64
}

type DB struct {
	lmp map[string][]string
	smp map[string]string
	zmp map[string][]NF
	mu  sync.Mutex
}

func (self *DB) Init() {
	//
	N := 1000000
	self.lmp = make(map[string][]string, N)
	self.smp = make(map[string]string, N)
	self.zmp = make(map[string][]NF, N)
}

func (self *DB) Get_zmp() map[string][]NF {
	//
	return self.zmp
}

func (self *DB) Get_lmp() map[string][]string {
	//
	return self.lmp
}

func (self *DB) Get_smp() map[string]string {
	//
	return self.smp
}

func (self *DB) Zrem(name, name_list string) bool {
	//
	self.mu.Lock()
	//self.Zsort(name)
	for i := 0; i < len(self.zmp[name]); i += 1 {
		if self.zmp[name][i].Name == name_list {
			self.zmp[name] = append(self.zmp[name][:i], self.zmp[name][i+1:]...)

			self.mu.Unlock()
			return true
		}
	}
	self.mu.Unlock()
	return false
}

func (self *DB) Zcard(name string) int {
	//
	self.mu.Lock()
	res := len(self.zmp[name])
	self.mu.Unlock()
	return res
}

func (self *DB) Zsort(name string) {
	//
	//self.mu.Lock()

	values := make([]float64, len(self.zmp[name]))
	tm := make(map[float64][]NF, len(self.zmp[name]))
	for i := 0; i < len(self.zmp[name]); i += 1 {
		value := self.zmp[name][i].Value
		values[i] = value
		tm[value] = append(tm[value], self.zmp[name][i])
	}
	a := sort.Float64Slice(values[0:])
	sort.Sort(a)

	var temp []NF
	var old float64
	for i := 0; i < len(a); i += 1 {
		if old != a[i] {
			temp = append(temp, tm[a[i]]...)
		}
		old = a[i]
	}
	self.zmp[name] = temp
	//self.mu.Unlock()
}

func (self *DB) Zrangebyscore(name string, min, max float64) []string {
	//
	self.mu.Lock()
	self.Zsort(name)
	var start, end int
	start, end = 0, 0
	for i := 0; i < len(self.zmp[name]); i += 1 {
		if self.zmp[name][i].Value >= min && start == 0 {
			start = i
		}
		if self.zmp[name][i].Value > max && end == 0 {
			end = i
		}
	}
	if start != 0 && end == 0 {
		end = len(self.zmp[name])
	}
	//self.mu.Unlock()
	//fmt.Printf("Zrangebyscore start: %d  end:%d\n", start, end)
	res := self.Zrange_not_lock(name, start, end)
	self.mu.Unlock()
	return res
}

func (self *DB) Sadd(name, value string) {
	//
	self.mu.Lock()
	var res bool
	res = false
	for _, val := range self.lmp[name] {
		if val == value {
			res = true
		}
	}
	if res == false {
		self.lmp[name] = append(self.lmp[name], value)
	}
	self.mu.Unlock()
}

func (self *DB) Zadd(name string, value float64, name_sortlist string) {
	//
	self.mu.Lock()
	var t NF
	t.Name = name_sortlist
	t.Value = value
	self.zmp[name] = append(self.zmp[name], t)

	//self.mu.Unlock()
	//self.Zsort(name)
	self.mu.Unlock()
}

func (self *DB) Zrange_not_lock(name string, start, end int) []string {
	//
	//self.mu.Lock()
	self.Zsort(name)
	if end == start {
		end = start + 1
	}
	if end > len(self.zmp[name]) {
		end = len(self.zmp[name])
	}
	//fmt.Printf("ZRANGE:start - %d : end - %d: len = %d\n",start,  end, len(self.zmp[name]))
	temp := self.zmp[name][start:end]
	lt := make([]string, len(temp))
	for i := 0; i < len(temp); i += 1 {
		lt[i] = temp[i].Name
	}
	//self.mu.Unlock()
	return lt
}

func (self *DB) Zrange(name string, start, end int) []string {
	//
	self.mu.Lock()
	self.Zsort(name)
	if end == start {
		end = start + 1
	}
	if end > len(self.zmp[name]) {
		end = len(self.zmp[name])
	}
	//fmt.Printf("ZRANGE:start - %d : end - %d: len = %d\n",start,  end, len(self.zmp[name]))
	temp := self.zmp[name][start:end]
	lt := make([]string, len(temp))
	for i := 0; i < len(temp); i += 1 {
		lt[i] = temp[i].Name
	}
	self.mu.Unlock()
	return lt
}

func (self *DB) LDelete_name(name string) {
	//
	self.mu.Lock()
	//var tm []string
	delete(self.lmp, name)
	self.mu.Unlock()
}

func (self *DB) ZDelete_name(name string) {
	//
	self.mu.Lock()
	//var tm []NF
	delete(self.zmp, name)
	self.mu.Unlock()
}

func (self *DB) Rpush(name, value string) {
	//
	self.mu.Lock()
	self.lmp[name] = append(self.lmp[name], value)
	self.mu.Unlock()
}

func (self *DB) Lrange(name string, start, end int) []string {
	//
	self.mu.Lock()
	if end == start {
		end = start + 1
	}
	if end > len(self.lmp[name]) {
		end = len(self.lmp[name])
	}
	res := self.lmp[name][start:end]
	self.mu.Unlock()
	return res
}

func (self *DB) Llen(name string) int {
	//
	self.mu.Lock()
	res := len(self.lmp[name])
	self.mu.Unlock()
	return res
}

func (self *DB) Lpop(name string) (string, bool) {
	//
	self.mu.Lock()
	if len(self.lmp[name]) == 0 {
		self.mu.Unlock()
		return "", false
	}
	x := self.lmp[name][len(self.lmp[name])-1]
	self.lmp[name] = self.lmp[name][:len(self.lmp[name])-1]
	self.mu.Unlock()
	return x, true
}

func (self *DB) Lrem(name string, count int, name_d string) {
	//
	self.mu.Lock()
	temp := self.lmp[name]
	for i := 0; i < len(temp); i += 1 {
		if temp[i] == name_d {
			temp = append(temp[:i], temp[i+1:]...)
		}
	}
	self.lmp[name] = temp
	self.mu.Unlock()
}

func (self *DB) Lindex(name string, index int) string {
	//
	//fmt.Printf("\n\nLINDEX:%d: %s: :%v\n\n", index, name, self.lmp[name])
	self.mu.Lock()
	res := self.lmp[name][index]
	self.mu.Unlock()
	return res
}

func (self *DB) Set(name, value string) {
	//
	self.mu.Lock()
	self.smp[name] = value
	self.mu.Unlock()
}

func (self *DB) Get(name string) string {
	//
	self.mu.Lock()
	res := self.smp[name]
	self.mu.Unlock()
	return res
}

func main() {
	tm := new(DB)
	tm.Init()
/*	
		tm.Rpush("asd", "asdd1")
		tm.Rpush("asd", "asdd1")
		tm.Rpush("asd", "asdd3")
		tm.Rpush("asd", "asdd1")

		fmt.Printf("Lrange : %s:%v:LEN:%d\n", "asd", tm.Lrange("asd", 0, 1000), tm.Llen("asd"))

		//tm.Lrem("asd", 0, "asdd1")
		t, _ := tm.Lpop("asd")
		fmt.Printf("Lpop1 : %s\n", t)
		fmt.Printf("Lrange : %s:%v:LEN:%d\n", "asd", tm.Lrange("asd", 0, 1000), tm.Llen("asd"))
		t, _ = tm.Lpop("asd")
		fmt.Printf("Lpop2 : %s\n", t)
		fmt.Printf("Lrange : %s:%vLEN:%d\n", "asd", tm.Lrange("asd", 0, 1000), tm.Llen("asd"))
		t, _ = tm.Lpop("asd")
		fmt.Printf("Lpop3 : %s\n", t)
		fmt.Printf("Lrange : %s:%v:LEN:%d\n", "asd", tm.Lrange("asd", 0, 1000), tm.Llen("asd"))
		t, _ = tm.Lpop("asd")
		fmt.Printf("Lpop4 : %s\n", t)
		fmt.Printf("Lrange : %s:%v:LEN:%d\n", "asd", tm.Lrange("asd", 0, 1000), tm.Llen("asd"))
		t, _ = tm.Lpop("asd")
		fmt.Printf("Lpop5 : %s\n", t)
		fmt.Printf("Lrange : %s:%v:LEN:%d\n", "asd", tm.Lrange("asd", 0, 1000), tm.Llen("asd"))

		tm.Rpush("asd1", "asdd1")
		tm.Rpush("asd1", "asdd2")
		tm.Rpush("asd1", "asdd3")
		tm.Rpush("asd1", "asdd4")

		tm.Zadd("zasd", 1.9, "name5")
		tm.Zadd("zasd", 1.8, "name4")
		tm.Zadd("zasd", 1.7, "name3")
		tm.Zadd("zasd", 1.6, "name2")
		tm.Zadd("zasd", 1.6, "name222")
		tm.Zadd("zasd", 1.6, "name22")
		tm.Zadd("zasd", 1.5, "name1")
		tm.Zadd("zasd", 1.5, "name11")

		fmt.Printf("Zadd :%v\n", tm.Zrange("zasd", 0, 1000))

		fmt.Printf("OK ZREM?:%v\n", tm.Zrem("zasd", "name4"))

		fmt.Printf("Zadd :%v\n", tm.Zrange("zasd", 0, 1000))

		res := tm.Lindex("asd1", 0)
		fmt.Printf("Lindex :%s\n", res)

		fmt.Printf("Zrangebyscore :%v\n", tm.Zrangebyscore("zasd", 1.7, 1.7))
		fmt.Printf("Lindex :%v\n", tm.Lindex("asd1", 0))
	
	tm.Sadd("Sadd", "asdd1")
	tm.Sadd("Sadd", "asdd1")
	tm.Sadd("Sadd", "asdd1")
	tm.Sadd("Sadd", "asdd1")
	tm.Sadd("Sadd", "asdd2")
	tm.Sadd("Sadd", "asdd1")
	tm.Sadd("Sadd", "asdd2")
	fmt.Printf("Zadd :%v\n", tm.Lrange("Sadd", 0, 1000))
	*/
		tm.Zadd("zasd", 1.9, "name5")
		tm.Zadd("zasd", 1.8, "name4")
		tm.Zadd("zasd", 1.8, "name4")
		tm.Zadd("zasd", 1.7, "name3")
		tm.Zadd("zasd", 1.6, "name2")
		tm.Zadd("zasd", 1.6, "name222")
		tm.Zadd("zasd", 1.6, "name22")
		tm.Zadd("zasd", 1.5, "name1")
		tm.Zadd("zasd", 1.5, "name11")
		fmt.Printf("Zrange :%v\n", tm.Zrange("zasd", 0, 1000))
		fmt.Printf("OK ZREM?:%v\n", tm.Zrem("zasd", "name4"))
		fmt.Printf("OK ZREM?:%v\n", tm.Zrem("zasd", "name4"))
		fmt.Printf("OK ZREM?:%v\n", tm.Zrem("zasd", "name22"))
		fmt.Printf("Zrange :%v\n", tm.Zrange("zasd", 0, 1000))
}
