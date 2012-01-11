package main

import (
	clean "./clean"
	gm "./fuzzy_gamodel"
	save "./save_redis"
	tools "./fuzzy_tools"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

var name_exe = flag.String("name", "simple", "name .exe")

const (
	N_genom_init = 20
	N_population = 400
)

var (
	original []float64
	values   map[int][]float64
	cur_pop  []int
	N_rule int
	list_gen_rule  map[int][]int

)

func verify_result(original, test []float64) float64 {
	//
	if len(original) != len(test) {
		panic("len(original) != len(test)")
	}
	error := 0.0
	for i := 0; i < len(test); i += 1 {
		error += math.Abs(original[i] - test[i])
	}
	return error
}

func run_genom(nomer, pop int, ext chan bool) {
	//create instance
	name := *name_exe + "!" + "genom_" + strconv.Itoa(nomer) + "_pop_" + strconv.Itoa(pop)
	save.Save(save.Ngenom, name)

	fmdl := gm.FModel{N_input: gm.N_input,
		N_output:   gm.N_output,
		N_fp:       gm.N_fp,
		Name_genom: name, 
		N_rule:		N_rule}
	fmdl.List_gen_rules = list_gen_rule

	//create fp param	
	var genom_fp, genom_fp_store []float64
	var genom_rule, genom_rule_store map[int][]bool

	if pop == 0 {
		fmdl.Init_fp()
		fmdl.Init_Rules()
	} else {
		//name := "genom_"+strconv.Itoa(nomer)+"_pop_"+strconv.Itoa(pop-1)
		genom_fp = save.Get_genom_fp_sortlist(nomer)
		genom_rule = save.Get_genom_rule_sortlist(nomer)

		if rand.Float64() < 0.6 {
			genom_fp = fmdl.Mutach_fp(genom_fp)
		}
		if rand.Float64() < 0.9 {
			genom_rule = fmdl.Mutach_rule(genom_rule)
		}

		genom_fp_store, genom_rule_store = save.Get_genom_sortlist(nomer, pop, 1, 0)
		//fmt.Printf("Selection_fp:::%v     %v\n", genom_fp, genom_fp_store)
		fmdl.Selection_fp(genom_fp, genom_fp_store)
		fmdl.Selection_rule(genom_rule, genom_rule_store)
		//fmt.Printf("genom_rule_store:%v\n", genom_rule_store)

	}
	//calc model
	fmdl.Calc_model(values)

	//calc error
	fmdl.Error = verify_result(original, fmdl.Result)

	//save
	save.Save(name+":error", strconv.FormatFloat(fmdl.Error, 'f', 2, 64))

	save.SaveSortList(*name_exe+"!"+"pop_"+strconv.Itoa(pop),
		fmdl.Error,
		name)

	save.SaveSortList(*name_exe+"!"+"genom_"+strconv.Itoa(nomer),
		fmdl.Error,
		name)

	save.SaveSortList(*name_exe+"!"+"itog",
		fmdl.Error,
		name)

	//fmt.Printf("END %d - %d \n",nomer, pop)
	//fmt.Printf("End-%d pop-%d\n", nomer, pop)
	ext <- true
	
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(2)

	start := time.Now()
	end := time.Now()

	original = []float64{1, 2, 3, 1, 2, 3, 1, 2, 3, 1}

	save.InitRedis()
	save.Set(*name_exe + "!" +"N_population", strconv.Itoa(N_population))
	save.Set(*name_exe + "!" +"N_genom_init", strconv.Itoa(N_genom_init))
	save.Name_exe = *name_exe

	N_rule = int(math.Pow(float64(gm.N_fp), float64(gm.N_input+gm.N_output)))
	list_gen_rule = tools.Generate_rules(N_rule,
		gm.N_fp,
		gm.N_input,
		gm.N_output)

	_ = save.Load_data_redis(*name_exe, gm.N_input, gm.N_output, gm.N_fp, gm.N_step)

	fmt.Printf("download data min:%.2f\n\n",
		float64(time.Now().Sub(start))/float64(1e9*60))
	start = time.Now()

	var res_bool bool
	res_bool = false
	for _, val := range save.GetNameExe_redis() {
		if *name_exe == val {
			res_bool = true
		}
	}
	if res_bool == false {
		save.Save_alone(save.Nexe, *name_exe)
	}

	//make test data
	values = make(map[int][]float64, gm.N_step)
	values[0] = make([]float64, gm.N_input)
	for i := 1; i < gm.N_step; i++ {
		values[i] = make([]float64, gm.N_input)
		for k := 0; k < gm.N_input; k += 1 {
			values[i][k] = values[i-1][k] + 10.0/float64(gm.N_step)
		}
	}

	var mass_chan [N_genom_init]chan bool
	for i := 0; i < N_genom_init; i += 1 {
		mass_chan[i] = make(chan bool)
	}

	var pop int
	for j := 0; j < N_population; j += 1 {

		pop = save.Get_pop()

		for i := 0; i < N_genom_init; i += 1 {
			go run_genom(i, pop, mass_chan[i])
			//fmt.Printf("Start-%d pop -%d\n", i, j)
		}
		for i :=0; i < N_genom_init; i += 1 {
			<-mass_chan[i]
			//fmt.Printf("Check-%d pop -%d\n", i, j)
		}
		save.Set_pop(pop + 1)
		
		fmt.Printf("Populate " + strconv.Itoa(pop) + " is " +
			strconv.Itoa(pop+N_population-j) + "\n")
		
		if math.Mod(float64(j), 5) == 0 && j != 0 {
			clean.Demon_clean()
		}
		if math.Mod(float64(j), 50) == 0 && j != 0 {
			//save.Backup_Redis()
		}
		fmt.Printf("Error:" + save.Get_best_error() + "\n")
		end = time.Now()
		fmt.Printf("time work min:%.2f\n\n",
			float64(end.Sub(start))/float64(1e9*60))
	}
	//save.Backup_Redis()
	fmt.Printf("Upload data min:%.2f\n\n",
		float64(time.Now().Sub(end))/float64(1e9*60))

}
