package fuzzy_gamodel

import (
	"math"
	"math/rand"
	"strconv"
	//"fmt"
	model "./fuzzy_model"
	tools "./fuzzy_tools"
	fvar "./fuzzy_var"
	save "./save_redis"
)

const (
	max      = 50
	N_step   = 10
	N_fp     = 4
	N_input  = 2
	N_output = 1
)

func Zade_And(x, y float64) float64 {
	return math.Min(x, y)
}

func Zade_Or(x, y float64) float64 {
	return math.Max(x, y)
}

type FModel struct {
	N_input, N_output, N_fp, N_rule int
	List_gen_rules                  map[int][]int
	Result                          []float64
	Genom_fp                        []float64
	Genom_rules                     map[int][]bool
	Spisok_var                      []fvar.Variable
	Error                           float64
	Name_genom                      string
	Max                             int
}

func (self *FModel) Save_fp_genom(genom []float64) {
	//
	str_genom := tools.Mass_float64_to_str(genom)
	for _, str_g := range str_genom {
		save.Save(self.Name_genom+":genom_fp", str_g)
	}
}

func (self *FModel) Save_rule_genom(rule_genom map[int][]bool) {
	//
	save.Save(self.Name_genom+":n_rule", strconv.Itoa(self.N_rule))
	//fmt.Printf("Save_rule_genom_---N_rule:%d\n", self.N_rule)
	for i := 0; i < self.N_rule; i++ {
		str_genom := tools.Mass_bool_to_str(rule_genom[i])
		for _, str_g := range str_genom {
			save.Save(self.Name_genom+":genom_rule:"+strconv.Itoa(i), str_g)
		}
	}
}

func (self *FModel) Set_fp(genom []float64) {
	//save redis genom
	self.Save_fp_genom(genom)

	count := 0
	self.Spisok_var = make([]fvar.Variable, self.N_input+self.N_output)

	//определяем функции принадлежности
	for i := 0; i < self.N_input+self.N_output; i += 1 {
		self.Spisok_var[i] = fvar.Make_var("Variable_"+strconv.Itoa(i+1), self.Name_genom)
		for j := 1; j <= self.N_fp; j += 1 {
			self.Spisok_var[i].Add_fp("N"+strconv.Itoa(j), genom[count], genom[count+1], float64(genom[count+2]))
			count += 3
		}
	}
	var value1 float64
	value1 = 0.0
	for i := 0; i < max; i += 1 {
		for _, value := range self.Spisok_var {
			value.Save_fps(value1)
		}
		value1 += 10.0 / float64(max)
	}
}

func (self *FModel) Init_fp() {
	//
	N_genom := (self.N_input + self.N_output) * self.N_fp * 3
	self.Genom_fp = make([]float64, N_genom)
	count := 0
	for i := 0; i < self.N_input+self.N_output; i += 1 {
		for j := 0; j < self.N_fp; j += 1 {
			self.Genom_fp[count] = 2 * float64(j)     //a			
			self.Genom_fp[count+1] = 2 * float64(j+1) //0.8 + tools.Get_rand(-0.2, 0.2)//k
			self.Genom_fp[count+2] = 2 * float64(j+2) //float64(2.5*float64(j+1)) + tools.Get_rand(-0.2, 0.2)//b
			count += 3
		}
	}
	self.Set_fp(self.Genom_fp)
}

func (self *FModel) Init_Rules() {
	//
	//self.N_rule = int(math.Pow(float64(self.N_fp), float64(self.N_input+self.N_output)))
	//fmt.Printf("self.N_rule:%d\n", self.N_rule)
	//self.List_gen_rules = tools.Generate_rules(self.N_rule,
	//	self.N_fp,
	//		self.N_input,
	//	self.N_output)
	self.Genom_rules = make(map[int][]bool, self.N_rule)
	for i := 0; i < self.N_rule; i++ {
		N_genov := 2*self.N_input - 1
		self.Genom_rules[i] = make([]bool, N_genov)
		for j := 0; j < N_genov; j++ {
			self.Genom_rules[i][j] = tools.Get_rand_bool()
		}
	}
	self.Save_rule_genom(self.Genom_rules)
}

func (self *FModel) Selection_fp(genom1, genom2 []float64) {
	//
	/*count := 0
	for i:=0; i< (self.N_input + self.N_output)/2; i+=1{
		for j:=0; j < self.N_fp; j +=1{
			genom1[count] = genom2[count] //a
			genom1[count+1] = genom2[count+1]   //k
			genom1[count+2] = genom2[count+2] //b
			count += 3
		}
	}*/
	self.Genom_fp = self.Verify_fp(genom1)
	self.Set_fp(self.Genom_fp)
}

func (self *FModel) Selection_rule(genom1, genom2 map[int][]bool) {
	//
	/*self.N_rule = int(math.Pow(float64(self.N_fp), float64(self.N_input+self.N_output)))
	self.List_gen_rules = tools.Generate_rules(self.N_rule,
		self.N_fp,
		self.N_input,
		self.N_output)
	/*
		for i:=0;i<self.N_rule/2;i+=1{
			for j:=0;j<self.N_input*2-1;j+=1{
				genom1[i][j] = genom2[i][j]
			}		
		}*/

	N_selection := self.N_rule / 2
	for i := 0; i < N_selection; i += 1 {
		x := rand.Intn(self.N_rule)
		y := rand.Intn(self.N_input*2 - 1)
		genom1[x][y] = genom2[x][y]
	}

	self.Genom_rules = genom1
	self.Save_rule_genom(self.Genom_rules)
}

func (self *FModel) Mutach_rule(genom map[int][]bool) map[int][]bool {
	//
	//N_rule := int(math.Pow(float64(self.N_fp), 
	//				float64(self.N_input + self.N_output)))
	N_i := rand.Intn(int(self.N_rule))
	for i := 0; i < N_i; i += 1 {
		if rand.Float32() < 0.5 {
			genom[rand.Intn(self.N_rule)][rand.Intn(self.N_input*2-1)] = true
		} else {
			genom[rand.Intn(self.N_rule)][rand.Intn(self.N_input*2-1)] = false
		}
	}
	return genom
}

func (self *FModel) Verify_fp(genom []float64) []float64 {
	//
	var count_2, count int
	var old_gen float64
	for i := 0; i < self.N_input+self.N_output; i += 1 {
		old_gen = 0.0
		count = 2 + i*self.N_fp*3
		count_2 = 1
		for j := 0; j < self.N_fp*3; j += 1 {
			if old_gen >= genom[count-1] && j != (self.N_fp*3-1) && j != 0 {
				//fmt.Printf("ERROR:%v count==%d  j == %d i == %d\n",genom,  count, j, i)
				if genom[count-1] == old_gen {
					genom[count-1] = old_gen + tools.Get_rand(0, 0.1)
				} else if genom[count-1] < 10.0 {
					genom[count-1] = old_gen + tools.Get_rand(0, old_gen-genom[count-1])
				} else if genom[count-1] >= 10.0 {
					genom[count-1] = old_gen + tools.Get_rand(0, old_gen-10)
				}
				return self.Verify_fp(genom)
			}
			old_gen = genom[count-1]
			if count_2 == 2 {
				count_2 = 0
				count += -1
			} else {
				count_2 += 1
				count += 2
				if count > len(genom) {
					count += -1
				}
			}
		}
		//hack ?!
		genom[i*self.N_fp*3] = 0
		genom[(i+1)*self.N_fp*3-1] = 10
	}
	return genom
}

func (self *FModel) Mutach_fp(genom []float64) []float64 {
	//
	N_genom := (self.N_input + self.N_output) * self.N_fp * 3
	new_genom_fp := make([]float64, N_genom)
	count := 0
	for i := 0; i < self.N_input+self.N_output; i += 1 {
		for j := 0; j < self.N_fp; j += 1 {
			if j == 0 {
				new_genom_fp[count] = genom[count] //a
			} else {
				new_genom_fp[count] = genom[count] + tools.Get_rand(-0.5, 0.5)
			}

			new_genom_fp[count+1] = genom[count+1] + tools.Get_rand(-0.5, 0.5) //k
			if j == self.N_fp-1 {
				new_genom_fp[count+2] = genom[count+2]
			} else {
				new_genom_fp[count+2] = genom[count+2] + tools.Get_rand(-0.5, 0.5) //b
			}
			count += 3
		}
	}
	new_genom_fp = self.Verify_fp(new_genom_fp)
	//fmt.Printf("ERROR:%v \n",new_genom_fp)
	return new_genom_fp
}

func (self *FModel) Calc_model(mass_values map[int][]float64) {

	//N_rule := int(math.Pow(float64(self.N_fp), float64(self.N_input + self.N_output)))
	//list_gen_rules := tools.Generate_rules(self.N_rule,self.N_fp, self.N_input, self.N_output)

	values := make([]float64, self.N_input)
	self.Result = make([]float64, len(mass_values))

	var ysl, r2 float64
	var count, or_and, count_gen, count_rule int
	mdl := model.Make_model(self.N_rule)
	mdl.Max = max
	mdl.Name_genom = self.Name_genom
	for i := 0; i < len(mass_values); i++ {
		values = mass_values[i]
		count_rule = 0
		for _, gen := range self.List_gen_rules {
			count, count_gen = 0, 0
			ysl = self.Spisok_var[gen[0]-1].Fp["N"+strconv.Itoa(gen[1])].Triangle(values[gen[0]-1])
			ysl = verify_genom(ysl,
				self.Genom_rules[count_rule][count_gen],
				self.Genom_rules[count_rule][count_gen+1])
			count += 2
			for i := 0; i < self.N_input-1; i += 1 {
				//for gen
				count_gen += 2
				//
				or_and = gen[count]
				count += 1
				r2 = self.Spisok_var[gen[count]-1].Fp["N"+strconv.Itoa(gen[count+1])].Triangle(values[gen[count]-1])
				r2 = verify_genom(r2,
					self.Genom_rules[count_rule][count_gen],
					self.Genom_rules[count_rule][count_gen-1])
				switch or_and {
				case 0:
					ysl = Zade_Or(ysl, r2)
				case 1:
					ysl = Zade_And(ysl, r2)
				}
				count += 2
			}
			mdl.Add_rule(ysl, self.Spisok_var[gen[count]-1].Fp["N"+strconv.Itoa(gen[count+1])])
			count_rule += 1
		}
		y_val := mdl.Calc_rules(mdl.Max, "result_"+strconv.Itoa(i))
		self.Result[i] = y_val
	}
}

func verify_genom(res float64, gen, gen_oper bool) float64 {
	if gen {
		return res
	}
	if gen_oper {
		//AND -> min
		return 1.0
	}
	//OR -> max
	return 0.0
}
