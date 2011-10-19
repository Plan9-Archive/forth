// forth package
package forth

import (
	"strconv"
	"os"
	"runtime"
	"strings"
)

type forthop func (f Forth)

type forthstack struct {
	stack []string
}

var opmap = map[string] forthop {
	"+": plus, 
	"-": sub, 
	"*": times, 
	"/": div, 
	"%": mod, 
	"swap": swap,
	"ifelse": ifelse,
	"hostname": hostname,
	"hostbase": hostbase,
	"strcat": strcat,
	"roundup": roundup,
	"dup": dup,
}

type Forth  interface{
	Push(string)
	Pop() (string)
	Length() int
	Empty() bool
	Newop(string, forthop)
	Reset()
}

func New() Forth {
	f := new(forthstack)
	return f
}

func (f *forthstack) Newop(n string, op forthop) {
	opmap[n] = op
}

func (f *forthstack) Reset() {
	f.stack = f.stack[0:0]
}

func (f *forthstack) Push(s string) {
	f.stack = append(f.stack, s)
	//fmt.Printf("push: %v: stack: %v\n", s, f.stack)
}

func (f *forthstack) Pop() (ret string){

	if len(f.stack) < 1 {
		panic(os.NewError("Empty stack"))
	}
	ret = f.stack[len(f.stack)-1]
	f.stack = f.stack[0:len(f.stack)-1]
	//fmt.Printf("Pop: %v stack %v\n", ret, f.stack)
	return ret
}

func (f *forthstack) Length() int {
	return len(f.stack)
}

func (f *forthstack) Empty() bool {
	return len(f.stack) == 0
}
	
func errRecover(errp *os.Error){
	e := recover()
	if e != nil {
		if _, ok := e.(runtime.Error); ok {
			panic(e)
		}
		*errp = e.(os.Error)
	}
}

/* eval and return TOS and an error if the stack is not empty */
func Eval(f Forth, s string) (ret string, err os.Error) {
	defer errRecover(&err)
	for _, val := range(strings.Split(s, " ")) {
		/* two or more spaces can have odd results */
		if len(val) == 0 || val == " " {
			continue
		}
		fun := opmap[val]
		if fun != nil {
			fun(f)
		} else {
			f.Push(val)
		}
	}
	ret = f.Pop()
	return
	
}

func plus(f Forth) {
	x, _ := strconv.Atoi(f.Pop())
	y, _ := strconv.Atoi(f.Pop())
	z := x + y
	f.Push(strconv.Itoa(z))
}

func times(f Forth) {
	x, _ := strconv.Atoi(f.Pop())
	y, _ := strconv.Atoi(f.Pop())
	z := x * y
	f.Push(strconv.Itoa(z))
}

func sub(f Forth) {
	x, _ := strconv.Atoi(f.Pop())
	y, _ := strconv.Atoi(f.Pop())
	z := y-x
	f.Push(strconv.Itoa(z))
}

func div(f Forth) {
	x, _ := strconv.Atoi(f.Pop())
	y, _ := strconv.Atoi(f.Pop())
	z := y/x
	f.Push(strconv.Itoa(z))
}

func mod(f Forth) {
	x, _ := strconv.Atoi(f.Pop())
	y, _ := strconv.Atoi(f.Pop())
	z := y%x
	f.Push(strconv.Itoa(z))
}

func roundup(f Forth) {
	rnd, _ := strconv.Atoi(f.Pop())
	v, _ := strconv.Atoi(f.Pop())
	v = ((v + rnd-1)/rnd)*rnd
	f.Push(strconv.Itoa(v))
}

func swap(f Forth) {
	x := f.Pop()
	y := f.Pop()
	f.Push(x)
	f.Push(y)
}

func strcat(f Forth) {
	x := f.Pop()
	y := f.Pop()
	f.Push(y+x)
}

func dup(f Forth) {
	x := f.Pop()
	f.Push(x)
	f.Push(x)
}

func ifelse(f Forth) {
	x,_ := strconv.Atoi(f.Pop())
	y := f.Pop()
	z := f.Pop()
	if x != 0 {
		f.Push(y)
	} else {
		f.Push(z)
	}
}

func hostname(f Forth){
	h, err := os.Hostname()
	if err != nil {
		panic("No hostname")
	}
	f.Push(h)
}

func hostbase(f Forth) {
	host := f.Pop()
	f.Push(strings.TrimLeft(host, "abcdefghijklmnopqrstuvwxyz -"))
}
