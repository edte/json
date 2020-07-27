## 前言

这篇文章总结了我编写 json 库的思路，实现了序列化，对 value 的 set，get 操作，和检查 json 是否有效的功能。基本思路都是参考标准库和已有的轮子。

## json

### json 是什么

json 即 JavaScript Object Notation，是一种数据交换格式。json 在可以转换为对应语言中相应的数据结构，所以可以使用 json 在不同语言间传输数据。

### json 数据格式

在说明 json 的结构之前，我们需要先知道任何数据基本上都可以分为三种类型，即基本的 data，data 的集合，和 data 的映射。

回忆一下任何一门高级语言，基本上自带的数据类型都是这么几种情况，有了这个前置知识后，再来看 json 的数据类型。

在 json 中 data 被称为 value，同样 json 有 value ，value 的集合，和 value 的映射。

**bool**

value 为 false 和 true 的被称为 bool 类型。

```
false
true
```



**null**

value 为 null 的被称为 null 类型。

```
null
```



**string**

value 为 char 码点集合，并且包含在 `""` 中的序列被称为 string

![](https://upload-images.jianshu.io/upload_images/4329360-2488111f8b9d9d76.png?imageMogr2/auto-orient/strip|imageView2/2/w/611/format/webp)

```
"sdfisdf"
```



**number**

value 为数字的被称为 number

number 必须为十进制

![](https://upload-images.jianshu.io/upload_images/4329360-befffe8e46910c3f.png?imageMogr2/auto-orient/strip|imageView2/2/w/622/format/webp)

```
123
123.1
-123
```



**array**

value 的集合被称为 array，value 可以不同。

![](https://upload-images.jianshu.io/upload_images/4329360-e9024223202a8bcf.png?imageMogr2/auto-orient/strip|imageView2/2/w/627/format/webp)

value 间使用 `,` 连接，最后一个 value 后不加 `,`,所有 value 放在 `[]` 中

如

```
[
    1,
    2,
    "one",
    "two",
    null,
    true,
    false
]
```

**object**

value 的映射被称为 object，映射的自变量被称为 key，key 必须是 string 类型

![](https://upload-images.jianshu.io/upload_images/4329360-fcc53eef9154d1b3.png?imageMogr2/auto-orient/strip|imageView2/2/w/621/format/webp)

key 和 value 间使用 `:` 连接，不同的 key-value 间使用 `,` 连接，最后一个 key-value 不加 `,`, 所有 key-value 放在 `{}` 中

如

```
{
    "name":"Lucy",
    "age":25
}
```

**语法**

json 的基本类型就是上面六种，然后这几种类型进行任意组合形成 json，可以知道 json 是一个树形结构。

安装 sjon 的官方标准，json 的 root 节点需要是 array 或者 object，不过一般的 parser 都支持其他 type，所以我们在编写时也会支持，这样也会更方便一些。

注意，json 是没有注释语法的，当然你可以自己搞个 parser 来添加注释语法，不过标准是没有的。这也导致 json 写配置有一些问题，你不知道对应的格式为什么要这么写，这其实是历史原因，json 的起源并不是为了写配置，不过现在我们一般用来写配置比较多。

在 json 中，字符 \t, \n, \r, space 会被忽略，就像高级语言一样。

json 使用 utf-8 编码 。

json 的 tokens 有这几种格式

 `string` , `number` , `false` , `true` , `null` ,  `{` , `}` , `[` , `]`,`,`,`:` 

从编写语言到编译成功，一般有这么几个过程，词法分析，语法分析，语义分析，生成中间代码，链接。前三个只在字符间操作，所以被称为编译器前端，后面的一些操作一般被称为编译器后端。

一门语言由一些词组成，如 for，while，=，`,`，`[`，`false` 等等，这些词在编译中就被称为 token，词法分析器（tokenizer）就是把源码分成一些 tokens。然后语法分析器会把这些 token 放到抽象语法树中，再进行语法分析，然后又到语义分析器中分析语义，不过 json 的语法和 token 都比较简单，可以把这些放在一起实现。

## 序列化

> In computing, **serialization** (US spelling) or **serialisation** (UK spelling) is the process of translating a [data structure](https://en.wikipedia.org/wiki/Data_structure) or [object](https://en.wikipedia.org/wiki/Object_(computer_science)) state into a format that can be stored (for example, in a [file](https://en.wikipedia.org/wiki/Computer_file) or memory [data buffer](https://en.wikipedia.org/wiki/Data_buffer)) or transmitted (for example, across a [computer network](https://en.wikipedia.org/wiki/Computer_network)) and reconstructed later (possibly in a different computer environment).

序列化（serialization/marshalling）按照 wiki 的解释，是说把 data 转换成可以传输或能存储的格式的过程，在这里的意思就是把 go 中的数据类型转换成 byte 流。

要把 go 的数据转换成 byte 流，标准库是使用的 reflect 来实现的，我看了其他的一些库，对这个的实现比较少，一般都是反序列化和对 value 的操作比较多，因为反序列化使用 reflect 后性能比较低，一些库则实现了一个 parser 来实现反序列化，而序列化则性能提升不大，而同时标准库也没有提供对 value 操作的支持，所以 go 的 json 库比较多。

这里为了增强 reflect 包的使用，所以也使用 reflect 来实现。

为了实现这个功能，我们需要先来看一下 go 的数据类型。

以 json 的类型做参考

json           <-------->    go

string        <-------->    tring

number    <-------->    int,int8等,uint,uint8等，float32,float64

null           <-------->    无

array         <-------->    array,slice

object       <-------->    map, struct

有人可以会对 struct 有疑问，但是 struct 的 field 其实也是 key-value 映射类型，只是 map 的 key 和 value 是的类型固定的，而 struct 是不固定的

除此外，go 中还有 value 可以在运行中变化的 interface，还有任何类型地址的 ptr

有了上面的知识后，我们继续来看实现，在 go 中，我们最常见的操作就是把一个 struct 变成 byte 流了，所以我们就以 struct 来展开。

现在我们有一个 struct，要开始对它进行操作，我们最终的结果是得到一个 byte slice，所以需要用一个 encoder 来存这个数据。

```
// encodeState 是总的 encoder
type encodeState struct {
	// 这是 bytes buffer，encode 得到的 byte 就存这里面
	bytes.Buffer
}
```

然后我们需要知道，struct 使用 field 组成的，每个 field 有 key 和 value，所以对 struct 进行序列化的结果，最终是得到一个 object `{}` 的。

既然 struct 是由 field 组成，所以我们需要分别对这些 field 进行 encode，key 直接写为 string，但是 value 可以是任何 go 一种类型，这怎么办呢。

我们在上面列出了 go 有哪些数据类型，而每种数据类型都可能进行序列化，所以我们需要对每种类型都建立一个 encoder 来操作。（这里展开就支持了所有类型的 encode）

当我们对每种类型都建立好 encoder 后，再对每个 field 的 value 调用对应的 encoder 进行操作即可。

这里先展开基本 value 的 encode

如对 bool 类型 encode

```
// boolEncoder encode bool
func boolEncoder(e *encodeState, v reflect.Value) {
	if v.Bool() {
		e.WriteString("true")
	} else {
		e.WriteString("false")
	}
}
```

判断 bool 的类型，然后写如对应的值

如对 string 的操作

```
// stringEncoder encode string
func stringEncoder(e *encodeState, v reflect.Value) {
	e.WriteByte('"')
	e.WriteString(v.String())
	e.WriteByte('"')
}
```

可以直接写入，当然这里 string 可以有转义字符，如 `\n`，如果直接写入的话，就会有问题，不过这里只提供了简单的功能。

如对 number 型的操作，number 有 int 类，uint 类，float 类

```
// intEncoder encode int
func intEncoder(e *encodeState, v reflect.Value) {
	appendInt := strconv.AppendInt([]byte(nil), v.Int(), 10)
	e.Write(appendInt)
}

// uintEncoder encode uint
func uintEncoder(e *encodeState, v reflect.Value) {
	appendInt := strconv.AppendUint([]byte(nil), v.Uint(), 10)
	e.Write(appendInt)
}

// float32Encoder encode float32
func float32Encoder(e *encodeState, v reflect.Value) {
	b := strconv.AppendFloat([]byte(nil), v.Float(), 'f', -1, 32)
	e.Write(b)
}
```

等等，从上面的代码都可以发现，只要把这些类型 encode 为 byte 流再写入 buffer 中即可。

现在我们已经弄好了 string, number, bool 这几个基本的 type，我们可以写一个简单的 encoder 了。

现在再回到开始的时候，我们现在有了一个 基本的 struct

我们提供接口

```
func Marshal(v interface{}) ([]byte, error) {
```

传入 struct，然后声明一个总的 encoder，再开始操作传入的 struct。

然后我们需要一个中转站，用来判断传入的类型，得到对应的 encoder，再调用这个 encoder

```
// reflectValue 比较重要，利用了 closure 的性质，把 encoder 作为返回值
// 这里起了一个中间站的作用，先反射拿到 encoder，再使用这个 encoder
func (e *encodeState) reflectValue(v reflect.Value) {
	// get encoder
	encoder := valueEncoder(v)
	// call encoder
	encoder(e, v)
}
```

然后我们需要一个根据类型得到对应 encoder 的函数

```
// newTypeEncoder 返回对应 type 的 encoder
func newTypeEncoder(t reflect.Type) encoderFunc {
```

现在我们假设传入的是 struct，自然在这个函数中就要拿到 struct 的 encoder

所以我们需要创建一个 struct 的 encoder

```
// newStructEncoder encode struct
func newStructEncoder(t reflect.Type) encoderFunc {
	se := structEncoder{
		// 遍历存入 field
		field: typeFields(t),
	}
	return se.encode
}
```

根据我们上面说的思路，我们可以先遍历把 field 存好，然后得到对应 field 的 encoder，再遍历得到的 field，再分别调用 field 的 encoder，在上面的函数中就是这样实现的。

当然我们需要东西来表示 field 和 struct 的 field 组

```
// filed 表示 一个 filed
type field struct {
	name      string      // field name
	nameBytes []byte      // field name bytes
	encoder   encoderFunc // field type encoder
	value     reflect.Value
}

// structEncoder 表示一个 struct
type structEncoder struct {
	field []field
}
```

再根据我们上面说的，创建遍历得到 field 的函数

```
// typeFields 遍历存入 field
func typeFields(t reflect.Type) []field {
	fields := []field{}

	for j := 0; j < t.NumField(); j++ {
		fields = append(fields, field{
			name:      t.Field(j).Name,
			nameBytes: []byte(t.Field(j).Name),
			encoder:   typeEncoder(t.Field(j).Type),
			value:     reflect.ValueOf(t.Field(j)),
		})
	}

	return fields
}
```

然后创建遍历 field 调用对应 encoder 的函数

```
// encode encode struct
func (se structEncoder) encode(e *encodeState, v reflect.Value) {
```

现在我们再来分析一下流程

![](https://img2020.cnblogs.com/blog/1823594/202007/1823594-20200728020654996-1945328124.png)

现在我们可以编写简单的 struct 来看一下了，

如

```
type person struct{
	Age int
	Name string
}
p := person {
	Age: 18
	Name: "lily"
}
```

// 这里暂时不写

现在我们的框架已经搭好了，如果 struct 的 value 是上面不支持的类型，我们就还需要继续编写其他的 encoder

现在我们再简单的说一下和 struct 类型的 map 的思路，map 的 string 最终都是 string 格式，为了方便，这里就只支持 key 为 string 的格式了。

我们需要遍历 map，写入 key，然后得到 value，再调用 value 对应的 encoder，当然我们也可以像上面那样对 struct 那样，先遍历把 value 存好，再遍历 value 调用 encoder，不过对 struct 这样操作主要是为了增加 tag 功能，不过为了简单我们这里没有增加 struct 的 tag 功能，这里 map 就可以直接遍历的同时调用了。

```
// mapEncoder encode map
func mapEncoder(e *encodeState, v reflect.Value) {
```

注意，这里只是讲解了思路，具体怎么利用 reflect 得到数据，和对语法的一些细节操作，如什么是否写入 `,` 什么的这里就没有展开了。

有了上面的基础后，我们再来增加 slice 和 array 的 encoder，类似的，可以遍历一般 slice，然后调用对应 value 的 encoder，只是比 map 少了写入 key 的步骤，以及一些细微的语法不同。

```
// arrayOrSliceEncoder 用于 slice 和 array 的 encode
func arrayOrSliceEncoder(e *encodeState, v reflect.Value) {
```

现在我们只剩下 interface 和 ptr 的 encoder 没有写了，现在思考一个问题，我们在传入数据时，接受数据的函数接受的是什么类型的数据，因为传入的数据类型可能不同，其实接受的都是动态类型 interface，都会自动类型转换为 interface，不过我们可以通过 reflect 来知道具体某时刻是什么类型。

如果如果传入是 interface 类型，直接调用总的 encoder，再根据具体的运行类型调用 encoder 即可。

```
// interfaceEncoder encode interface
func interfaceEncoder(e *encodeState, v reflect.Value) {
	if v.IsNil() {
		e.WriteString("null")
		return
	}
	// 获得 interface 的动态类型
	e.reflectValue(v.Elem())
}
```

当传入的是 pointer 类型时又该怎么办呢，这个其实和 interface 类似，我们可以通过 reflect 得到具体指向的类型，然后调用总的 encoder 即可。

```
// ptrEncoder 用于指针的 encode
func ptrEncoder(e *encodeState, v reflect.Value) {
	if v.IsNil() {
		e.WriteString("null")
		return
	}
	// 使用 reflect.Elem() 方法可以获得 ptr 指向的 type，然后调用所有的 encoder，当做 interface 即可
	e.reflectValue(v.Elem())
}
```



现在我们再来看，基本上所有类型的 encoder 我们都支持了，除了对应类型的特殊情况和一些功能我们没有增加外，普通的类型我们都可以支持。要完善对应的功能的话，只需要在上面的基础上增加响应的代码即可。

## 检查是否有效

好困

## 对 value 的操作



## 总结



## 参考