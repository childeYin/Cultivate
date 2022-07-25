# 观察者模式


type Customer interface {
	updete()
}

type CustomerA struct {

}

func (*CustomerA) update() {
	fmt.Println("我是客户A，我收到报纸了")
}

type CustomerB struct {

}

func (*CustomerB) update() {
	fmt.Println("我是客户B，我收到报纸了")
}

type NewOffice struct {
	cutomers []Customer
}

func (n *NewOffice) newall() {
	n.notifyAllCusotmer()
}

func (n *NewOffice) addCustomer(cutomer Customer) {
	n.cutomers = append(n.cutomers, cutomer)
}

func (n *NewOffice)notifyAllCusotmer() {
	for _, customer := range n.cutomers {
		cutomer.update()
	}
}

func Main() {
	customerA := &CustomerA{}
	cusotomerB := &CustomerB{}

	office := &NewOffice{}

	office.addCustomer(customerA)
	office.addCustomer(cusotomerB)

	office.newall()
}