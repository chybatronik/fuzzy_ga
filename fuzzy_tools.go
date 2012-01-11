package fuzzy_tools

import (
	"math"
	"math/rand"
	"strconv"
)

func s_plus(mass []int, max, n int) []int {
	if n >= len(mass) {
		return mass
	}
	if n == len(mass)-1 && mass[n] == max {
		for i := 0; i < len(mass); i += 1 {
			mass[i] = 0
		}
		return mass
	}
	if mass[n] < max {
		mass[n] += 1
		return mass
	}
	mass[n] = 0
	return s_plus(mass, max, n+1)
}

func Plus(mass []int, max int) []int {
	return s_plus(mass, max, 0)
}

func all_combinate(mass []int, count int) map[int][]int {

	N := int(math.Pow(float64(len(mass)), float64(count)))
	tmp := make([]int, count)
	result := make(map[int][]int, N)

	for i := 0; i < count; i += 1 {
		tmp[i] = 0
	}
	for i := 0; i < N; i += 1 {
		result[i] = make([]int, len(mass))
		tmp = Plus(tmp, len(mass)-1)
		for k := 0; k < len(tmp); k += 1 {
			result[i][k] = mass[tmp[k]]
		}
	}
	return result
}

func Generate_rules(N_rule, N_fp, N_input, N_output int) map[int][]int {
	//
	list_gen_rules := make(map[int][]int, N_rule)
	N_gen := N_input*2 + N_input - 1 + N_output*2

	mass_mark := make([]int, N_fp) //write N_input + N_output
	for i := 0; i < N_fp; i++ {    //write N_input + N_output
		mass_mark[i] = i + 1
	}
	fp := all_combinate(mass_mark, N_input+N_output)

	for i := 0; i < N_rule; i += 1 {
		list_gen_rules[i] = make([]int, N_gen)
		list_gen_rules[i][0] = 1
		list_gen_rules[i][1] = fp[i][0]
		list_gen_rules[i][2] = 1 //AND
		for input := 1; input < N_input; input += 1 {
			list_gen_rules[i][input+2] = input + 1
			list_gen_rules[i][input+3] = fp[i][input]
			//AND
			if input != N_input-1 {
				list_gen_rules[i][input+4] = 1 //AND
			}
		}
		N_input_pos := N_input*2 + N_input - 1
		for output := 0; output < N_output; output += 1 {
			list_gen_rules[i][N_input_pos+output] = output + N_input + 1
			list_gen_rules[i][N_input_pos+output+1] = fp[i][N_input+output]
		}
	}
	return list_gen_rules
}

func Get_rand(start, end float64) float64 {
	temp := start + (end-start)*rand.Float64()
	return temp
}

func Mass_float64_to_str(mass []float64) []string {
	N := len(mass)
	result := make([]string, N)
	for i := 0; i < N; i += 1 {
		result[i] = strconv.FormatFloat(mass[i], 'f', 1, 64)
	}
	return result
}

func Mass_bool_to_str(mass []bool) []string {
	N := len(mass)
	result := make([]string, N)
	for i := 0; i < N; i += 1 {
		result[i] = strconv.FormatBool(mass[i])
	}
	return result
}

func Mass_str_to_float64(mass []string) []float64 {
	N := len(mass)
	result := make([]float64, N)
	for i := 0; i < N; i += 1 {
		result[i], _ = strconv.ParseFloat(mass[i], 64)
	}
	return result
}

func Mass_str_to_bool(mass []string) []bool {
	N := len(mass)
	result := make([]bool, N)
	for i := 0; i < N; i += 1 {
		result[i], _ = strconv.ParseBool(mass[i])
	}
	return result
}

func Get_rand_bool() bool {
	if rand.Float64() <= 0.5 {
		return true
	}
	return false

}
