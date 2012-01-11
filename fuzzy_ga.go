package main

import (
	"./pplus"
	"fmt"
	"godis"
	"math"
	"runtime"
	"strconv"
	"time"
)

const (
	NFP      = "name_"
	NVAR     = "name_fp"
	NxAxis   = "axis"
	Ngenom   = "genom"
	max      = 10
	N_step   = 10
	N_input  = 2
	N_output = 1
)

var ch chan Message
var original []float64
var values map[int][]float64

func save(name, value string) {
	m := Message{name: name, value: value}
	ch <- m
}

type Message struct {
	name, value string
}

func daemon(m <-chan Message, ex chan<- bool) {
	c := godis.New("tcp:127.0.0.1:6379", 0, "")
	c.Flushdb()
	for {
		select {
		case mm, ok := <-m:
			if !ok {
				fmt.Printf("Read from channel failed")
			}
			c.Rpush(mm.name, mm.value)
		case <-time.After(0.1e9):
			ex <- true
		}
	}
	fmt.Printf("END DAEMON \n")
}

func verify_result(original, test []float64) float64 {
	if len(original) != len(test) {
		panic("len(original) != len(test)")
	}
	error := 0.0
	for i := 0; i < len(test); i += 1 {
		error += math.Abs(original[i] - test[i])
	}
	return error
}

func run_genom(name string) {
	//create instance
	save(Ngenom, name[:len(name)-1])
	fmdl := FModel{N_input: N_input,
		N_output: N_output,
		N_fp:     3,
		name:     name}

	//create fp param
	fmdl.init_fp()

	//calc model
	fmdl.calc_model(values)
	//fmt.Printf("result:%v\n", fmdl.result)
	//calc error
	fmdl.error = verify_result(original, fmdl.result)
	//fmt.Printf("error:%f\n", fmdl.error)
}

func main() {

	runtime.GOMAXPROCS(2)
	start := time.Now()

	original = []float64{0, 1, 2, 3, 0, 0, 0, 0, 0, 0}

	//make test data
	values = make(map[int][]float64, N_step)
	values[0] = make([]float64, N_input)
	for i := 1; i < N_step; i++ {
		values[i] = make([]float64, N_input)
		for k := 0; k < N_input; k += 1 {
			values[i][k] = values[i-1][k] + 10.0/float64(N_step)
		}
	}

	//for redis
	ex := make(chan bool)
	ch = make(chan Message, 100000)
	go daemon(ch, ex)
	//

	for i := 0; i < 2; i += 1 {
		go run_genom("genom_" + strconv.Itoa(i) + ":")
	}

	//finish transfer to redis
	<-ex

	end := time.Now()
	fmt.Printf("time work ms:%d", (end.Sub(start))/1e6)

}
