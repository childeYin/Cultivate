# Array

##### 1.数组是由相同类型元素的集合组成的数据结构，初始化之后，大小无法改变。

* cmd/compile/internal/types/type.go
```
// Array contains Type fields specific to array types.
type Array struct {
	Elem  *Type // element type
	Bound int64 // number of elements; <0 if unknown yet
}


// NewArray returns a new fixed-length array Type.
func NewArray(elem *Type, bound int64) *Type {
	if bound < 0 {
		Fatalf("NewArray: invalid bound %v", bound)
	}
	t := New(TARRAY)
	t.Extra = &Array{Elem: elem, Bound: bound}
	t.SetNotInHeap(elem.NotInHeap())
	return t
}

```

##### 2.初始化

```
arr1 := [3]int{1, 2, 3}
arr2 := [...]int{1, 2, 3}

```

* cmd/compile/internal/gc/typecheck.go
```
// The result of typecheckcomplit MUST be assigned back to n, e.g.
// 	n.Left = typecheckcomplit(n.Left)
func typecheckcomplit(n *Node) (res *Node) {
	...

	// Need to handle [...]T arrays specially.
	if n.Right.Op == OTARRAY && n.Right.Left != nil && n.Right.Left.Op == ODDD {
		n.Right.Right = typecheck(n.Right.Right, ctxType)
		if n.Right.Right.Type == nil {
			n.Type = nil
			return n
		}
		elemType := n.Right.Right.Type

		length := typecheckarraylit(elemType, -1, n.List.Slice(), "array literal")

		n.Op = OARRAYLIT
		n.Type = types.NewArray(elemType, length)
		n.Right = nil
		return n
	}

	.....

	switch t.Etype {
	default:
		yyerror("invalid composite literal type %v", t)
		n.Type = nil

	case TARRAY:
		typecheckarraylit(t.Elem(), t.NumElem(), n.List.Slice(), "array literal")
		n.Op = OARRAYLIT
		n.Right = nil

	case TSLICE:
		...
	case TMAP:
		...
	case TSTRUCT:
		...

	return n
}
```
##### 3.数组越界是非常严重的问题

* 无论是在栈上还是静态存储区，数组在内存中都是一连串的内存空间
* 数组和字符串的一些简单越界错误都会在编译期间发现，使用变量访问，编译器无法提前知道
* 访问数组的索引是非整数时，报错 “non-integer array index %v”
* 访问数组的索引是负数时，报错 “invalid array index %v (index must be non-negative)"
* 访问数组的索引越界时，报错 “invalid array index %v (out of bounds for %d-element array)"

```
// The result of typecheck1 MUST be assigned back to n, e.g.
// 	n.Left = typecheck1(n.Left, top)
func typecheck1(n *Node, top int) (res *Node) {

	...
	switch n.Op {

		case OINDEX:
			ok |= ctxExpr
			n.Left = typecheck(n.Left, ctxExpr)
			n.Left = defaultlit(n.Left, nil)
			n.Left = implicitstar(n.Left)
			l := n.Left
			n.Right = typecheck(n.Right, ctxExpr)
			r := n.Right
			t := l.Type
			if t == nil || r.Type == nil {
				n.Type = nil
				return n
			}
			switch t.Etype {
			default:
				yyerror("invalid operation: %v (type %v does not support indexing)", n, t)
				n.Type = nil
				return n

			case TSTRING, TARRAY, TSLICE:
				n.Right = indexlit(n.Right)
				if t.IsString() {
					n.Type = types.Bytetype
				} else {
					n.Type = t.Elem()
				}
				why := "string"
				if t.IsArray() {
					why = "array"
				} else if t.IsSlice() {
					why = "slice"
				}

				if n.Right.Type != nil && !n.Right.Type.IsInteger() {
					yyerror("non-integer %s index %v", why, n.Right)
					break
				}

				if !n.Bounded() && Isconst(n.Right, CTINT) {
					x := n.Right.Int64()
					if x < 0 {
						yyerror("invalid %s index %v (index must be non-negative)", why, n.Right)
					} else if t.IsArray() && x >= t.NumElem() {
						yyerror("invalid array index %v (out of bounds for %d-element array)", n.Right, t.NumElem())
					} else if Isconst(n.Left, CTSTR) && x >= int64(len(strlit(n.Left))) {
						yyerror("invalid string index %v (out of bounds for %d-byte string)", n.Right, len(strlit(n.Left)))
					} else if n.Right.Val().U.(*Mpint).Cmp(maxintval[TINT]) > 0 {
						yyerror("invalid %s index %v (index too large)", why, n.Right)
					}
				}

			case TMAP:
				n.Right = assignconv(n.Right, t.Key(), "map index")
				n.Type = t.Elem()
				n.Op = OINDEXMAP
				n.ResetAux()
			}
	...
	}
}
```

##### 4. transfer

```
	func ArrayToSlice() {
		var a = [5]int{0, 1, 2, 3, 4}
		var b = a[:]
		b[1] += 10
		fmt.Printf("%v\n", b) // [0 11 2 3 4]
	}


```
##### 5. 参考文章

1.[数组](https://draveness.me/golang/docs/part2-foundation/ch03-datastructure/golang-array/)
2.[Array vs Slice: accessing speed](https://stackoverflow.com/questions/30525184/array-vs-slice-accessing-speed)








