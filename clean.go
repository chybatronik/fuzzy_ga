package clean

import (
	//"time"
	//"fmt"
	"strconv"
	save "./save_redis"
	ga "./fuzzy_gamodel"
)

func Demon_clean(){
	//for {
	//	work()
	//	time.Sleep(1*1e9)
	//}
	work()
}

func work(){
	//удалить все популяции кроме 2 последних
	// и не удалять те геномы которые в 
	//списках itog первые 100
	//первые элементы genom
	
	pop := save.Get_pop()
	res, _ := strconv.Atoi(save.Get(save.Name_exe+"!"+"N_genom_init"))
	N_genom_init := res
	res, _ = strconv.Atoi(save.Get(save.Name_exe+"!"+"N_population"))
	N_population := res

	//fmt.Printf("Pop:%d  N_genom_init:%d    N_population:%d\n", pop, N_genom_init, N_population)
	
	var name_genom string
	good_list := get_spisok_good(N_population, N_genom_init, pop)
	//fmt.Printf("good_list:%v\n", good_list)
	for j:=pop - 6;j<pop;j+=1{
		for i:=0;i<N_genom_init;i+=1{
			name_genom = save.Name_exe+"!"+"genom_"+strconv.Itoa(i)+"_pop_"+strconv.Itoa(j)
			if !is_find(name_genom, good_list){				
				save.Delete_genom(i, j, ga.N_input, ga.N_output, ga.N_fp, ga.N_step)
				//fmt.Printf("DELETE ... \n")
			}			
		}
		//fmt.Printf("Populate "+strconv.Itoa(j)+" is "+ 
		//		strconv.Itoa(pop) + "\n")
	}
	//fmt.Printf("Wait ... \n")
}

func is_find(name string, mass []string) bool{
	//
	for _, val := range mass{
		if val == name{
			return true
		}
	}
	return false
}

func get_spisok_good(N_population, N_genom_init, pop int)[]string{
	//
	//список тех кто не удаляется
	result := make([]string, N_genom_init * 3 + save.N_itog)
	count := 0
	//2 последних популяции
	for i:=0;i<N_genom_init;i+=1{
		result[count] = save.Name_exe+"!"+"genom_"+strconv.Itoa(i)+"_pop_"+strconv.Itoa(pop)
		count += 1
		result[count] = save.Name_exe+"!"+"genom_"+strconv.Itoa(i)+"_pop_"+strconv.Itoa(pop-1)
		count += 1
	}
	//itog первые 100

	itog := save.Get_itog()
	//fmt.Printf("len(itog):%v\n", len(itog))
	for i:=0;i<len(itog);i+=1{
		result[count] = itog[i]
		count += 1
	}

	//первые элементы genom
	for i:=0;i<N_genom_init;i+=1{
	    result[count] = save.Get_1_genom(i)
	    count += 1
	}
	return result
}