package save_redis

import (
	data "./data"
	"godis"
	"strconv"
	"strings"
	tools "./fuzzy_tools"
	"fmt"
)

const (
	NFP    = "name_"
	NVAR   = "name_fp"
	NxAxis = "axis"
	Ngenom = "genom"
	Nexe   = "exe"
	N_itog = 100
)

var (
	cs       *data.DB
	redis_cs *godis.Client
	Name_exe string
)

func InitRedis() {
	//
	cs = new(data.DB)
	cs.Init()

	redis_cs = godis.New("tcp:localhost:6379", 0, "")
	//cs.Flushdb()
}

func Get_zmp() map[string][]data.NF {
	//
	return cs.Get_zmp()
}

func Get_lmp() map[string][]string {
	//
	return cs.Get_lmp()
}

func Get_smp() map[string]string {
	//
	return cs.Get_smp()
}

func Backup_Redis() {
	//
	//redis_cs.Flushdb()

	zmp := Get_zmp()
	lmp := Get_lmp()
	smp := Get_smp()

	for name, mass_val := range zmp {
		for _, value := range mass_val {
			redis_cs.Zadd(name, value.Value, value.Name)
		}
	}

	for name, mass_val := range lmp {
		for _, value := range mass_val {
			redis_cs.Rpush(name, value)
		}
	}

	for name, value := range smp {
		redis_cs.Set(name, value)
	}
}

func Save(name, value string) {
	//
	cs.Rpush(name, value)
}

func SaveSortList(name string, value float64, name_sortlist string) {
	//
	cs.Zadd(name, value, name_sortlist)
}

func Save_alone(name, value string) {
	//
	cs.Sadd(name, value)
}

func Load_data_redis(NameExe string, N_input, N_output, N_fp, N_step int) bool {
	//2 последних популяции
	//список pop_
	//списках itog первые 100
	//
	nexe_list, _ := redis_cs.Lrange(Nexe, 0, 10000)
	nexe_list_string := nexe_list.StringArray()
	var check bool
	check = false
	for _, val := range nexe_list_string{
		if NameExe == val{
			check = true
		}
	}
	//fmt.Printf("NameExe:%s\n", NameExe)
	if !check{
		Save_alone(Nexe, NameExe)
	}
	//Save_alone(Nexe, NameExe)

	res, _ := redis_cs.Get(NameExe + ":Pop")
	pop, _ := strconv.Atoi(res.String())

	if pop == 0 {
		return false
	}
	elem, _ := redis_cs.Get(NameExe+"!"+"N_genom_init")
	N_genom_init := int(elem.Int64())

	res_list, _ := redis_cs.Zrange(NameExe+"!"+"pop_"+strconv.Itoa(pop), 0, 1000000)
	res_list_string := res_list.StringArray()

	itog_list, _ := redis_cs.Zrange(Name_exe+"!"+"itog", 0, 100)
	itog_list_string := itog_list.StringArray()

	N := N_genom_init*2 + len(res_list_string) + len(itog_list_string)
	name_genom_s := make(map[string]float64, N)

	//load genom pop
	for i := pop - 2; i < pop; i += 1 {
		for j := 0; j < N_genom_init; j += 1 {
			name := NameExe + "!" + "genom_" + strconv.Itoa(j) + "_pop_" + strconv.Itoa(i)
			name_genom_s[name] = 0
		}
	}
	//load pop
	for _, val := range res_list_string {
		res_float, _ := redis_cs.Zscore(NameExe+"!"+"pop_"+strconv.Itoa(pop), val)
		SaveSortList(NameExe+"!"+"pop_"+strconv.Itoa(pop), res_float, val)
		name_genom_s[val] = 0
	}
	//load itog
	for _, val := range itog_list_string {
		res_float, _ := redis_cs.Zscore(Name_exe+"!"+"itog", val)
		SaveSortList(Name_exe+"!"+"itog", res_float, val)
		name_genom_s[val] = 0
	}
	//Set_pop(pop)
	cs.Set(NameExe+":Pop", strconv.Itoa(pop))
	var nomerget, popget int
	for name, _ := range name_genom_s {
		nomerget, popget = get_nomer_and_pop_name_genom(name)
		save_redis_to_memory(nomerget, popget, N_input, N_output, N_fp, N_step)
	}
	//fmt.Printf("LOAD:%v\n", name_genom_s)
	return true
}

func get_nomer_and_pop_name_genom(name string) (nomer, pop int) {
	//
	str_mass := strings.Split(name, "_")
	nomer, _ = strconv.Atoi(str_mass[1])
	pop, _ = strconv.Atoi(str_mass[3])
	return
}

func save_redis_to_memory(nomer, pop, N_input, N_output, N_fp, N_step int) {
	//
	name_genom := Name_exe + "!" + "genom_" + strconv.Itoa(nomer) + "_pop_" + strconv.Itoa(pop)
	var res godis.Elem

	name_db := NVAR + ":" + name_genom
	N, _ := redis_cs.Llen(name_db)
	for k := 0; k < int(N); k += 1 {
		res, _ = redis_cs.Lindex(name_db, k)
		Save(name_db, res.String())
	}
	for i := 0; i < N_input+N_output; i += 1 {
		name_var := "Variable_" + strconv.Itoa(i+1)
		for j := 1; j <= N_fp; j += 1 {
			name_fp := "N" + strconv.Itoa(j)

			name_db = name_genom + ":" + name_var + ":" + name_fp
			//len
			N, _ = redis_cs.Llen(name_db)
			for k := 0; k < int(N); k += 1 {
				res, _ = redis_cs.Lindex(name_db, k)
				Save(name_db, res.String())
			}
		}
		name_db = NFP + ":" + name_genom + ":" + name_var
		N, _ = redis_cs.Llen(name_db)
		for k := 0; k < int(N); k += 1 {
			res, _ = redis_cs.Lindex(name_db, k)
			Save(name_db, res.String())
		}
	}
	//из fuzzy_model
	name_db = NVAR + ":" + name_genom
	N, _ = redis_cs.Llen(name_db)
	for k := 0; k < int(N); k += 1 {
		res, _ = redis_cs.Lindex(name_db, k)
		Save(name_db, res.String())
	}

	for i := 0; i < N_step; i += 1 {
		name_rule := "result_" + strconv.Itoa(i)
		name_db = NFP + ":" + name_genom + ":" + name_rule
		N, _ = redis_cs.Llen(name_db)
		for k := 0; k < int(N); k += 1 {
			res, _ = redis_cs.Lindex(name_db, k)
			Save(name_db, res.String())
		}

		name_db = name_genom + ":" + name_rule + ":main"
		N, _ = redis_cs.Llen(name_db)
		for k := 0; k < int(N); k += 1 {
			res, _ = redis_cs.Lindex(name_db, k)
			Save(name_db, res.String())
		}

		name_db = name_genom + ":" + name_rule + ":result"
		N, _ = redis_cs.Llen(name_db)
		for k := 0; k < int(N); k += 1 {
			res, _ = redis_cs.Lindex(name_db, k)
			Save(name_db, res.String())
		}
	}

	name_db = name_genom + ":genom_fp"
	N, _ = redis_cs.Llen(name_db)
	for k := 0; k < int(N); k += 1 {
		res, _ = redis_cs.Lindex(name_db, k)
		Save(name_db, res.String())
	}

	name_db = name_genom + ":n_rule"
	N, _ = redis_cs.Llen(name_db)

	var N_rule int64
	for k := 0; k < int(N); k += 1 {
		N_rule_str, _ := redis_cs.Lindex(name_db, k)
		N_rule = N_rule_str.Int64()
	}
	//fmt.Printf("n_rule!!!!!!!!!%d\n", int(N_rule))
	Save(name_db, strconv.Itoa(int(N_rule)))
	for i := 0; i < int(N_rule); i += 1 {
		name_db = name_genom + ":genom_rule:" + strconv.Itoa(i)
		N, _ = redis_cs.Llen(name_db)
		for k := 0; k < int(N); k += 1 {
			res, _ = redis_cs.Lindex(name_db, k)
			Save(name_db, res.String())
		}
	}

	name_db = Ngenom
	//_, _ = cs.Lrem(name_db, 0, name_genom)
	Save(Ngenom, name_genom)

	name_db = name_genom + ":error"
	N, _ = redis_cs.Llen(name_db)
	var error float64
	for k := 0; k < int(N); k += 1 {
		res, _ = redis_cs.Lindex(name_db, k)
		error, _ = strconv.ParseFloat(res.String(), 64)
		Save(name_db, res.String())
	}

	name_db = Name_exe + "!" + "pop_" + strconv.Itoa(pop)

	//_, _ = cs.Zrem(name_db, name_genom)
	SaveSortList(name_db, error, name_genom)

	name_db = Name_exe + "!" + "genom_" + strconv.Itoa(nomer)

	//_, _ = cs.Zrem(name_db, name_genom)
	SaveSortList(name_db, error, name_genom)

	name_db = Name_exe + "!" + "itog"

	//_, _ = cs.Zrem(name_db, name_genom)
	SaveSortList(name_db, error, name_genom)
}

func Get_genom_fp_sortlist(nomer_genom int) []float64 {
	//
	res := cs.Zrange(Name_exe+"!"+"genom_"+strconv.Itoa(nomer_genom), 0, 0)
	name_genoms := res
	if len(res) == 0 {
		fmt.Printf("ERROR_Zrange len(res)=%d", len(res))
	}
	return Get_genom_fp(name_genoms[0])
}

func Get_genom_fp_sortlist_redis(nomer_genom int) []float64 {
	//
	res, _ := redis_cs.Zrange(Name_exe+"!"+"genom_"+strconv.Itoa(nomer_genom), 0, 0)
	name_genoms := res.StringArray()
	return Get_genom_fp_redis(name_genoms[0])
}

func Get_genom_fp_redis(name string) []float64 {
	//
	res, _ := redis_cs.Lrange(name+":genom_fp", 0, 10000)
	genom_fp := tools.Mass_str_to_float64(res.StringArray())
	return genom_fp
}

func Get_genom_rule_sortlist_redis(nomer_genom int) map[int][]bool {
	//
	res, _ := redis_cs.Zrange(Name_exe+"!"+"genom_"+strconv.Itoa(nomer_genom), 0, 0)
	name_genoms := res.StringArray()
	return Get_genom_rule_redis(name_genoms[0])
}

func Get_genom_rule_redis(name_genom string) map[int][]bool {
	//
	res, _ := redis_cs.Lindex(name_genom+":n_rule", 0)
	N_rule := int(res.Int64())
	result := make(map[int][]bool, N_rule)
	for i := 0; i < N_rule; i += 1 {
		res, _ := redis_cs.Lrange(name_genom+":genom_rule:"+strconv.Itoa(i), 0, 10000)
		result[i] = tools.Mass_str_to_bool(res.StringArray())
	}
	return result
}

func Get_genom_rule_sortlist(nomer_genom int) map[int][]bool {
	//
	res := cs.Zrange(Name_exe+"!"+"genom_"+strconv.Itoa(nomer_genom), 0, 0)
	name_genoms := res
	if len(res) == 0 {
		fmt.Printf("ERROR_Zrange len(res)=%d", len(res))
	}
	return Get_genom_rule(name_genoms[0])
}

func Get_best_error() string {
	//
	res := cs.Zrange(Name_exe+"!"+"itog", 0, 0)
	//fmt.Printf("Get_best_error:::%s\n",res[0] )
	elem := cs.Lindex(res[0]+":error", 0)
	return elem
}

func Get_best_error_redis() string {
	//
	res, _ := redis_cs.Zrange(Name_exe+"!"+"itog", 0, 0)
	res_str := res.StringArray()
	//fmt.Printf("Get_best_error:::%s\n",res[0] )
	elem, _ := redis_cs.Lindex(res_str[0]+":error", 0)
	return elem.String()
}

func Get_genom_fp(name string) []float64 {
	//
	res := cs.Lrange(name+":genom_fp", 0, 10000)
	genom_fp := tools.Mass_str_to_float64(res)
	return genom_fp
}

func Get_genom_rule(name_genom string) map[int][]bool {
	//
	res := cs.Lindex(name_genom+":n_rule", 0)
	//fmt.Printf("Get_genom_rule:::::%s\n", res)
	N_rule, _ := strconv.Atoi(res)
	result := make(map[int][]bool, N_rule)
	for i := 0; i < N_rule; i += 1 {
		res := cs.Lrange(name_genom+":genom_rule:"+strconv.Itoa(i), 0, 10000)
		result[i] = tools.Mass_str_to_bool(res)
	}
	return result
}

func Delete_genom_fp_sortlist(pop, nomer int, name_genom_fp string) bool {
	//
	res1 := cs.Zrem(Name_exe+"!"+"pop_"+strconv.Itoa(pop), name_genom_fp)
	if !res1 {
		return false
	}
	return true
}

func Lenght_sortlist(nomer_genom int) int {
	//
	res := cs.Zcard(Name_exe + "!" + "genom_" + strconv.Itoa(nomer_genom))
	return int(res)
}

func Set_pop(nomer int) {
	//
	cs.Set(Name_exe+":Pop", strconv.Itoa(nomer))
}

func Get_pop() int {
	//
	res := cs.Get(Name_exe + ":Pop")
	res1, _ := strconv.Atoi(res)
	return int(res1)
}

func Set(name, value string) {
	//
	cs.Set(name, value)
}

func Get(name string) string {
	//
	res := cs.Get(name)
	return res
}

func Get_itog() []string {
	//
	N := N_itog
	res := cs.Zrange(Name_exe+"!"+"itog", 0, N)
	name_genoms := res
	return name_genoms
}

func Get_1_genom(nomer int) string {
	//
	res := cs.Zrange(Name_exe+"!"+"genom_"+strconv.Itoa(nomer), 0, 0)
	name_genoms := res
	return name_genoms[0]
}

func Get_genom_sortlist(nomer, pop, delta_min, delta_max int) ([]float64, map[int][]bool) {
	//

	name_genom := Name_exe + "!" + "genom_" + strconv.Itoa(nomer) + "_pop_" + strconv.Itoa(pop-1)
	elem := cs.Lindex(name_genom+":error", 0)

	error, _ := strconv.ParseFloat(elem, 64)

	res := cs.Zrangebyscore(Name_exe+"!"+"pop_"+strconv.Itoa(pop-1),
		error-float64(delta_min),
		error+float64(delta_max))
	temp_list_name_genom := res

	//fmt.Printf("temp_list_name_genom:%v\n", temp_list_name_genom)
	//fmt.Printf("name_genom:%s\n", name_genom)
	if len(temp_list_name_genom) == 0 {
		//fmt.Printf("len(temp_list_name_genom) == 0\n")
		return Get_genom_sortlist(nomer, pop, delta_min+1, delta_max+1)
	}
	if len(temp_list_name_genom) == 1 && temp_list_name_genom[0] == name_genom {
		if delta_min >= 10 {
			fmt.Printf("delta_min >= 10::::%s\n", name_genom)
			return Get_genom_fp(name_genom), Get_genom_rule(name_genom)
		}
		return Get_genom_sortlist(nomer, pop, delta_min+1, delta_max+1)
	}
	for _, nm_genom := range temp_list_name_genom {
		if nm_genom != name_genom {
			genom_fp := Get_genom_fp(nm_genom)
			genom_rule := Get_genom_rule(nm_genom)
			rep := Delete_genom_fp_sortlist(pop-1, nomer, nm_genom)
			if !rep {
				//fmt.Printf("!Ok\n")
				return Get_genom_sortlist(nomer, pop, delta_min, delta_max)
			}

			//fmt.Printf("!!!!!%v, ------%v \n", genom_fp, genom_rule)
			return genom_fp, genom_rule
		}
	}

	fmt.Printf("return \n")
	return make([]float64, 1), make(map[int][]bool, 1)
}

func Get_pop_redis() int {
	//
	res, _ := redis_cs.Get(Name_exe + ":Pop")
	return int(res.Int64())
}

func Delete_genom_fp_sortlist_redis(pop, nomer int, name_genom_fp string) bool {
	//
	res1, _ := redis_cs.Zrem(Name_exe+"!"+"pop_"+strconv.Itoa(pop), name_genom_fp)
	if !res1 {
		return false
	}
	return true
}

func Get_genom_sortlist_redis(nomer, pop, delta_min, delta_max int) ([]float64, map[int][]bool) {
	//
	name_genom := Name_exe + "!" + "genom_" + strconv.Itoa(nomer) + "_pop_" + strconv.Itoa(pop)
	elem, _ := redis_cs.Lindex(name_genom+":error", 0)
	error := elem.Float64()

	res, _ := redis_cs.Zrangebyscore(Name_exe+"!"+"pop_"+strconv.Itoa(pop-1),
		strconv.FormatFloat(error-float64(delta_min), 'f', 2, 64),
		strconv.FormatFloat(error+float64(delta_max), 'f', 2, 64))
	temp_list_name_genom := res.StringArray()
	if len(temp_list_name_genom) == 0 {
		return Get_genom_sortlist_redis(nomer, pop, delta_min+1, delta_max+1)
	}
	if len(temp_list_name_genom) == 1 && temp_list_name_genom[0] == name_genom {
		if delta_min >= 10 {
			return Get_genom_fp_redis(name_genom), Get_genom_rule_redis(name_genom)
		}
		return Get_genom_sortlist_redis(nomer, pop, delta_min+1, delta_max+1)
	}
	for _, nm_genom := range temp_list_name_genom {
		if nm_genom != name_genom {
			genom_fp := Get_genom_fp_redis(nm_genom)
			genom_rule := Get_genom_rule_redis(nm_genom)
			rep := Delete_genom_fp_sortlist_redis(pop-1, nomer, nm_genom)
			if !rep {
				return Get_genom_sortlist_redis(nomer, pop, delta_min, delta_max)
			}
			return genom_fp, genom_rule
		}
	}
	return make([]float64, 1), make(map[int][]bool, 1)
}

func GetNameExe_redis() []string {
	//
	nexe_list, _ := redis_cs.Lrange(Nexe, 0, 10000)
	nexe_list_string := nexe_list.StringArray()
	return nexe_list_string
}

func GetNameExe() []string {
	//
	return cs.Lrange(Nexe, 0, 10000)
}

func Delete_genom(nomer, pop, N_input, N_output, N_fp, N_step int) {
	//
	//из fuzzy_var
	name_genom := Name_exe + "!" + "genom_" + strconv.Itoa(nomer) + "_pop_" + strconv.Itoa(pop)

	name_db := NVAR + ":" + name_genom
	N := cs.Llen(name_db)
	for k := 0; k < int(N); k += 1 {
		_, _ = cs.Lpop(name_db)
	}
	cs.LDelete_name(name_db)

	for i := 0; i < N_input+N_output; i += 1 {
		name_var := "Variable_" + strconv.Itoa(i+1)
		for j := 1; j <= N_fp; j += 1 {
			name_fp := "N" + strconv.Itoa(j)

			name_db = name_genom + ":" + name_var + ":" + name_fp
			//len
			N = cs.Llen(name_db)
			for k := 0; k < int(N); k += 1 {
				_, _ = cs.Lpop(name_db)
			}
			cs.LDelete_name(name_db)
		}
		name_db = NFP + ":" + name_genom + ":" + name_var
		N = cs.Llen(name_db)
		for k := 0; k < int(N); k += 1 {
			_, _ = cs.Lpop(name_db)
		}
		cs.LDelete_name(name_db)
	}
	//из fuzzy_model
	name_db = NVAR + ":" + name_genom
	N = cs.Llen(name_db)
	for k := 0; k < int(N); k += 1 {
		_, _ = cs.Lpop(name_db)
	}
	cs.LDelete_name(name_db)

	for i := 0; i < N_step; i += 1 {
		name_rule := "result_" + strconv.Itoa(i)
		name_db = NFP + ":" + name_genom + ":" + name_rule
		N = cs.Llen(name_db)
		for k := 0; k < int(N); k += 1 {
			_, _ = cs.Lpop(name_db)
		}
		cs.LDelete_name(name_db)

		name_db = name_genom + ":" + name_rule + ":main"
		N = cs.Llen(name_db)
		for k := 0; k < int(N); k += 1 {
			_, _ = cs.Lpop(name_db)
		}
		cs.LDelete_name(name_db)

		name_db = name_genom + ":" + name_rule + ":result"
		N = cs.Llen(name_db)
		for k := 0; k < int(N); k += 1 {
			_, _ = cs.Lpop(name_db)
		}
		cs.LDelete_name(name_db)

	}

	name_db = name_genom + ":genom_fp"
	N = cs.Llen(name_db)
	for k := 0; k < int(N); k += 1 {
		_, _ = cs.Lpop(name_db)
	}
	cs.LDelete_name(name_db)

	name_db = name_genom + ":n_rule"
	N = cs.Llen(name_db)

	var N_rule int
	for k := 0; k < int(N); k += 1 {
		N_rule_str, _ := cs.Lpop(name_db)
		res1, _ := strconv.Atoi(N_rule_str)
		N_rule = res1
	}
	cs.LDelete_name(name_db)

	//fmt.Printf("N_rule:%d\n", N_rule)
	for i := 0; i < int(N_rule); i += 1 {
		name_db = name_genom + ":genom_rule:" + strconv.Itoa(i)
		N = cs.Llen(name_db)
		for k := 0; k < int(N); k += 1 {
			_, _ = cs.Lpop(name_db)
		}
		cs.LDelete_name(name_db)
	}

	name_db = Ngenom
	cs.Lrem(name_db, 0, name_genom)

	name_db = name_genom + ":error"
	N = cs.Llen(name_db)
	for k := 0; k < int(N); k += 1 {
		_, _ = cs.Lpop(name_db)
	}
	cs.LDelete_name(name_db)

	name_db = Name_exe + "!" + "pop_" + strconv.Itoa(pop)
	_ = cs.Zrem(name_db, name_genom)

	name_db = Name_exe + "!" + "genom_" + strconv.Itoa(nomer)
	_ = cs.Zrem(name_db, name_genom)

	name_db = Name_exe + "!" + "itog"
	_ = cs.Zrem(name_db, name_genom)
}
