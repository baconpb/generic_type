package main

import (
	"fmt"
)

type Int interface {
	int | int8 | int16 | int32 | int64
}
type Uint interface {
	uint | uint8 | uint16 | uint32 | uint64
}
type Float interface {
	float32 | float64
}
type SLice[T Int | Uint | Float] []T

//1 ~指定底层类型
// ~int: 所有以int为底层类型的类型都可以用于实例化
// ~后面的类型不能为接口
// ~后面的类型必须为基本类型

type IntLike interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}
type UintLike interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}
type FloatLike interface {
	~float32 | ~float64
}
type SLiceLike[T IntLike | UintLike | FloatLike] []T

//2 从方法集（Method set）到类型集（Type set）
/*
就如下面这个代码一样， ReadWriter 接口定义了一个接口(方法集)，这个集合中包含了 Read() 和 Write() 这两个方法。所有同时定义了这两种方法的类型被视为实现了这一接口。
可以这样理解：我们可以把 ReaderWriter 接口看成代表了一个 类型的集合，所有实现了 Read() Writer() 这两个方法的类型都在接口代表的类型集合当中
换个角度看，在我们眼中接口的定义就从方法集变为了类型集。
An interface type defines a type set (一个接口类型定义了一个类型集)
type Float interface {
    ~float32 | ~float64
}
type Slice[T Float] []T
用 类型集 的概念重新理解上面的代码的话就是：
接口类型 Float 代表了一个 类型集合， 所有以 float32 或 float64 为底层类型的类型，都在这一类型集之中
而 type Slice[T Float] []T 中， 类型约束 的真正意思是：
类型约束 指定了类型形参可接受的类型集合，只有属于这个集合中的类型才能替换形参用于实例化
*/

//2.1 接口实现（implement）定义的变化
/*
当满足以下条件时，我们可以说 类型 T 实现了接口 I ( type T implements interface I)：
 T 不是接口时：类型 T 是接口 I 代表的类型集中的一个成员 (T is an element of the type set of I)
 T 是接口时： T 接口代表的类型集是 I 代表的类型集的子集(Type set of T is a subset of the type set of I)
*/

//2.2 类型的并集
/*
 | 符号就是求类型的并集
type Uint interface {  // 类型集 Uint 是 ~uint 和 ~uint8 等类型的并集
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}
*/

//2.3 类型的交集
// 接口可以不止书写一行，如果一个接口有多行类型定义，那么取它们之间的 交集。
/*
type AllInt interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint32
}
type Uint interface {
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}
// ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
type A interface { // 接口A代表的类型集是 AllInt 和 Uint 的交集
    AllInt
    Uint
}
// ~int
type B interface { // 接口B代表的类型集是 AllInt 和 ~int 的交集
    AllInt
    ~int
}
*/

//2.4 空集
// 当多个类型的交集如下面 Bad 这样为空的时候， Bad 这个接口代表的类型集为一个空集.
/*
type Bad interface {
    int
    float32
} // 类型 int 和 float32 没有相交的类型，所以接口 Bad 代表的类型集为空
// 没有任何一种类型属于空集。虽然 Bad 这样的写法是可以编译的，但实际上并没有什么意义
*/

//2.5 空接口和any
/*
空接口代表了所有类型的集合.
Go1.18之后对于空接口的理解：
1 虽然空接口内没有写入任何的类型，但它代表的是所有类型的集合，而非一个 空集
2 类型约束中指定 空接口 的意思是指定了一个包含所有类型的类型集，并不是类型约束限定了只能使用 空接口 来做类型形参
// 空接口代表所有类型的集合。写入类型约束意味着所有类型都可拿来做类型实参
type Slice[T interface{}] []T
var s1 Slice[int]    // 正确
var s2 Slice[map[string]string]  // 正确
var s3 Slice[chan int]  // 正确
var s4 Slice[interface{}]  // 正确

Go1.18开始提供了一个和空接口 interface{} 等价的新关键词 any ，用来使代码更简单
type Slice[T any] []T // 代码等价于 type Slice[T interface{}] []T
*/

//2.6 comparable(可比较) 和 可排序(ordered)
/*
对于一些数据类型，我们需要在类型约束中限制只接受能 != 和 == 对比的类型，如map：
// 错误。因为 map 中键的类型必须是可进行 != 和 == 比较的类型
type MyMap[KEY any, VALUE any] map[KEY]VALUE

// Go直接内置了一个叫 comparable 的接口，它代表了所有可用 != 以及 == 对比的类型：
type MyMap[KEY comparable, VALUE any] map[KEY]VALUE // 正确
comparable 比较容易引起误解的一点是很多人容易把他与可排序搞混淆。
可比较指的是 可以执行 != == 操作的类型，并没确保这个类型可以执行大小比较（ >,<,<=,>= ）

可进行大小比较的类型被称为 Orderd 。目前Go语言并没有像 comparable 这样直接内置对应的关键词，所以想要的话需要自己来定义相关接口，
比如我们可以参考Go官方包golang.org/x/exp/constraints 如何定义：
// Ordered 代表所有可比大小排序的类型
type Ordered interface {
    Integer | Float | ~string
}
type Integer interface {
    Signed | Unsigned
}
type Signed interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64
}
type Unsigned interface {
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}
type Float interface {
    ~float32 | ~float64
}
*/

// 3 接口的两种类型: 基本接口和一般接口
/*
接口类型 ReadWriter 代表了一个类型集合，所有以 string 或 []rune 为底层类型，
并且实现了 Read() Write() 这两个方法的类型都在 ReadWriter 代表的类型集当中。
定义一个 ReadWriter 类型的接口变量，然后接口变量赋值的时候不光要考虑到方法的实现，还必须考虑到具体底层类型。
为了解决这个问题也为了保持Go语言的兼容性，Go1.18开始将接口分为了两种类型。
*/

type ReadWriter interface {
	~string | ~[]rune
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
}
type StringReadWriter string

// 类型 StringReadWriter 实现了接口 Readwriter
func (s StringReadWriter) Read(p []byte) (n int, err error) {
	return len(p), nil
}
func (s StringReadWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// 类型BytesReadWriter 没有实现接口 Readwriter
type BytesReadWriter []byte

func (s BytesReadWriter) Read(p []byte) (n int, err error) {
	return len(p), nil
}

func (s BytesReadWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// 3.1 基本接口
//接口定义中如果只有方法的话，那么这种接口被称为基本接口(Basic interface)。这种接口就是Go1.18之前的接口，用法也基本和Go1.18之前保持一致。
/*
最常用的，定义接口变量并赋值
type MyError interface { // 接口中只有方法，所以是基本接口
    Error() string
}
// 用法和 Go1.18之前保持一致
var err MyError = fmt.Errorf("hello world")

基本接口因为也代表了一个类型集，所以也可用在类型约束中
// io.Reader 和 io.Writer 都是基本接口，也可以用在类型约束中
type MySlice[T io.Reader | io.Writer]  []Slice
*/

type MyError interface {
	Error() string
}

type ReadWriter2 interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
}

type OneRead struct {
}

func (o *OneRead) Read(p []byte) (n int, err error) {
	return len(p), nil
}

func (o *OneRead) Write(p []byte) (n int, err error) {
	return len(p), nil
}

type Invoker interface {
	Call(int2 interface{}) string
}
type FnCall func(interface{}) string

func (f FnCall) Call(val interface{}) string {
	return f(val)
}

// 3.2 一般接口
// 如果接口内不光只有方法，还有类型的话，这种接口被称为 一般接口(General interface)
/*
type Uint interface { // 接口 Uint 中有类型，所以是一般接口
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}
type ReadWriter interface {  // ReadWriter 接口既有方法也有类型，所以是一般接口
    ~string | ~[]rune

    Read(p []byte) (n int, err error)
    Write(p []byte) (n int, err error)
}

// 一般接口类型不能用来定义变量，只能用于泛型的类型约束中。
type Uint interface {
    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}
var uintInf Uint // 错误。Uint是一般接口，只能用于类型约束，不得用于变量定义
// 这一限制保证了一般接口的使用被限定在了泛型之中，不会影响到Go1.18之前的代码，同时也极大减少了书写代码时的心智负担.
*/

// 4 泛型接口
type DataProcessor[T any] interface {
	Process(oriData T) (newData T)
	Save(data T) error
}

/*
DataProcessor[string]

// 实例化之后的接口定义相当于如下所示：
type DataProcessor[string] interface {
	Process(oriData string) (newData string)
	Save(data string) error
}

DataProcessor[string] 因为只有方法，所以它实际上就是个 基本接口(Basic interface)，这个接口包含两个能处理string类型的方法。
*/

type DataProcessor2[T any] interface {
	int | ~struct{ Data interface{} }
	Process(oriData T) (newData T)
	Save(data T) error
}

/*
DataProcessor2[string]

// 实例化后的接口定义可视为
type DataProcessor2[T string] interface {
    int | ~struct{ Data interface{} }
    Process(data string) (newData string)
    Save(data string) error
}

DataProcessor2[string] 因为带有类型并集所以它是 一般接口.所以实例化之后的这个接口代表的意思是：
1 只有实现了 Process(string) string 和 Save(string) error 这两个方法，并且以 int 或 struct{ Data interface{} } 为底层类型的类型才算实现了这个接口
2 一般接口(General interface) 不能用于变量定义只能用于类型约束，所以接口 DataProcessor2[string] 只是定义了一个用于类型约束的类型集

// XMLProcessor 虽然实现了接口 DataProcessor2[string] 的两个方法，但是因为它的底层类型是 []byte，所以依旧是未实现 DataProcessor2[string]
type XMLProcessor []byte

// JsonProcessor 实现了接口 DataProcessor2[string] 的两个方法，同时底层类型是 struct{ Data interface{} }。所以实现了接口 DataProcessor2[string]
type JsonProcessor struct {
    Data interface{}
}

// 错误，带方法的一般接口不能作为类型并集的成员.参考接口定义的种种限制规则
type StringProcessor interface {
    DataProcessor2[string] | DataProcessor2[[]byte]

    PrintString()
}
*/

//3 接口定义的种种限制规则

//3.1 用|连接多个类型时，类型之间不能有相交的部分
/*
type MyInt int
// 错误，MyInt的底层类型是int,和 ~int 有相交的部分
type _ interface {
    ~int | MyInt
}
但是相交的类型中是接口的话，则不受这一限制：
type MyInt int
type _ interface {
    ~int | interface{ MyInt }  // 正确
}
type _ interface {
    interface{ ~int } | MyInt // 也正确
}
type _ interface {
    interface{ ~int } | interface{ MyInt }  // 也正确
}
*/

//3.2 类型的并集中不能有类型形参
/*
type MyInf[T ~int | ~string] interface {
    ~float32 | T  // 错误。T是类型形参
}
type MyInf2[T ~int | ~string] interface {
    T  // 错误
}
*/

//3.3 接口不能直接或间接的并入自己
/*
type Bad interface {
    Bad // 错误，接口不能直接并入自己
}
type Bad2 interface {
    Bad1
}
type Bad1 interface {
    Bad2 // 错误，接口Bad1通过Bad2间接并入了自己
}
type Bad3 interface {
    ~int | ~string | Bad3 // 错误，通过类型的并集并入了自己
}
*/

//3.4 接口的并集成员个数大于1的时候不能直接或间接并入comparable接口
/*
type OK interface {
    comparable // 正确。只有一个类型的时候可以使用 comparable
}
type Bad1 interface {
    []int | comparable // 错误，类型并集不能直接并入 comparable 接口
}
type CmpInf interface {
    comparable
}
type Bad2 interface {
    chan int | CmpInf  // 错误，类型并集通过 CmpInf 间接并入了comparable
}
type Bad3 interface {
    chan int | interface{comparable}  // 理所当然，这样也是不行的
}
*/

//3.5 带方法的接口（无论是基本接口还是一般接口），都不能写入接口的并集中
/*
type _ interface {
    ~int | ~string | error // 错误，error是带方法的接口(一般接口) 不能写入并集中
}
type DataProcessor[T any] interface {
    ~string | ~[]byte

    Process(data T) (newData T)
    Save(data T) error
}
// 错误，实例化之后的 DataProcessor[string] 是带方法的一般接口，不能写入类型并集
type _ interface {
    ~int | ~string | DataProcessor[string]
}
type Bad[T any] interface {
    ~int | ~string | DataProcessor[T]  // 也不行
}
*/

func main() {
	var s SLice[int] = []int{1, 2, 3, 4, 5}
	fmt.Println(s)

	type MySlice int
	var s2 SLiceLike[MySlice] = []MySlice{1, 2, 3}
	fmt.Println(s2)

	var read ReadWriter2
	read = &OneRead{}
	leng, err := read.Read([]byte("abc"))
	if err != nil {
		panic(err)
	}
	fmt.Println(leng)

	var invoke Invoker
	invoke = FnCall(func(v interface{}) string {
		return fmt.Sprintln(v)
	})
	callVal := invoke.Call("hello world")
	fmt.Println("callVal: ", callVal)

}
