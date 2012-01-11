package fuzzy_var

import (
	lib_fp "./fuzzy_fp"
	save "./save_redis"
	"strconv"
)

type Variable struct {
	Name_var, Name_genom string
	Fp                   map[string]lib_fp.FP
}

func (self *Variable) Save_fps(value1 float64) {
	for _, fp := range self.Fp {
		//_ = strconv.Ftoa64(fp.Triangle(value1), 'e', 2)
		save.Save(self.Name_genom+":"+self.Name_var+":"+fp.Name_fp,
			strconv.FormatFloat(fp.Triangle(value1), 'f', 2, 64))
	}
}

func (self *Variable) Add_fp(namefp string, a, k, b float64) {
	fp1 := lib_fp.FP{Name_fp: namefp}
	fp1.A, fp1.B, fp1.K = a, b, k
	save.Save(save.NFP+":"+self.Name_genom+":"+self.Name_var, namefp)
	self.Fp[namefp] = fp1
}

func Make_var(Name_var, Name_genom string) Variable {
	x1 := Variable{Name_var: Name_var,
		Name_genom: Name_genom}
	x1.Fp = make(map[string]lib_fp.FP)
	save.Save(save.NVAR+":"+Name_genom, Name_var)
	return x1
}
