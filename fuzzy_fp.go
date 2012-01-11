package fuzzy_fp

import (
	"math"
)

type FP struct {
	Fn                             string
	A, K, B                        float64
	Name_genoma, Name_var, Name_fp string
}

func (self FP) Triangle(x float64) float64 {
	var res float64
	//res := math.Fmax(v.A - math.Fabs(v.K*x - v.B), 0)
	if self.A == 0 {
		//начало
		if x <= self.K {
			res = 1.0
		} else {
			res = math.Max((1/(self.K-self.B))*x-self.B/(self.K-self.B), 0)
		}
	} else if self.B == 10 {
		//конец
		if x >= self.K {
			res = 1.0
		} else {
			res = math.Max((1/(self.K-self.A))*x-self.A/(self.K-self.A), 0)
		}
	} else if x <= self.K {
		//левая палка
		res = math.Max((1/(self.K-self.A))*x-self.A/(self.K-self.A), 0)
	} else if x > self.K {
		//правая палка
		res = math.Max((1/(self.K-self.B))*x-self.B/(self.K-self.B), 0)
	}

	return res
}
