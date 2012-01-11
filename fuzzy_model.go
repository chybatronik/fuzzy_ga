package fuzzy_model

import (
	lib_fp "./fuzzy_fp"
	save "./save_redis"
	"math"
	"strconv"
)

type Model struct {
	Output_fp  []lib_fp.FP
	Input      []float64
	Count      int
	M          chan bool
	Name_genom string
	Max        int
}

func Make_model(count_rule int) Model {
	mdl := Model{Count: 0}
	mdl.Input = make([]float64, count_rule)
	mdl.Output_fp = make([]lib_fp.FP, count_rule)
	mdl.M = make(chan bool)
	return mdl
}

func (self *Model) Add_rule(input float64, output lib_fp.FP) {
	self.Output_fp[self.Count] = output
	self.Input[self.Count] = input
	self.Count += 1
}

func (self *Model) Calc_rules(cnt int, name_rule string) float64 {
	u_tmp, u_sum, y_sum, y_val, u_main := 0.0, 0.0, 0.0, 0.0, 0.0
	//save godis
	save.Save(save.NVAR+":"+self.Name_genom, name_rule)
	save.Save(save.NFP+":"+self.Name_genom+":"+name_rule, "main")
	for y := 0; y < cnt; y += 1 {
		u_tmp = 0
		u_main = 0
		for i := 0; i < self.Count; i += 1 {
			val := math.Min(self.Input[i], self.Output_fp[i].Triangle(y_val))
			u_tmp += val
			u_main = math.Max(u_main, val)
		}
		save.Save(self.Name_genom+":"+name_rule+":main", strconv.FormatFloat(u_main, 'f', 2, 64))
		y_sum += y_val * u_tmp
		u_sum += u_tmp
		y_val += 10.0 / float64(self.Max)
	}
	self.Count = 0
	if u_sum == 0 {
		save.Save(self.Name_genom+":"+name_rule+":result",
			strconv.FormatFloat(0.0, 'f', 2, 64))
		return 0.0
	}
	save.Save(self.Name_genom+":"+name_rule+":result",
		strconv.FormatFloat(y_sum/u_sum, 'f', 2, 64))
	return y_sum / u_sum
}
