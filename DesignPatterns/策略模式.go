# 策略模式

/*
@Time : 2022/5/16 9:45 AM
@Author : zhangjun
@File : celve.go
@Description:
@Run:
*/
package main

import (
	"fmt"
	"strings"
)

type CookerInterface interface {
	fire()
	cooke()
	outfire()
}

type Oper struct {
	cooker CookerInterface
}

type Cooker struct {

}

func (*Cooker) fire() {
	fmt.Println("fire")
fmt.Println(	strings.Count("firefff","f"))
	a := strings.Fields("vvvv fire aaa nnn")
	fmt.Println(a[0])
}

func (*Cooker) cooke() {
	fmt.Println("cooke")
}

func (*Cooker) outfire() {
	fmt.Println("outfire")
}

type Cooker2 struct {

}

func (*Cooker2) fire() {
	fmt.Println("fire 2 ")
}

func (*Cooker2) cooke() {
	fmt.Println("cooke 2 ")
}

func (*Cooker2) outfire() {
	fmt.Println("outfire 2")
}


func (oper *Oper) SetCooker(cooker CookerInterface) {
	oper.cooker = cooker
}

func (oper *Oper) Cooking() {
	oper.cooker.fire()
	oper.cooker.cooke()
	oper.cooker.outfire()
}


func main(){
	oper := &Oper{}
	oper.SetCooker(&Cooker2{})
	oper.Cooking()

	oper.SetCooker(&Cooker{})
	oper.Cooking()


		s := []int{1,2}
		s = append(s, 4,5,6)
		fmt.Printf("len=%d,cap=%d\n",len(s),cap(s))

}