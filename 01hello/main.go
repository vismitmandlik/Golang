package main
import ("fmt"
"io" 
)

func main(){
	a := 8.8
	fmt.Println(a)
	m := []int{11,22,33}
	for _,i := range m {
        fmt.Println("range", i)
    }
}