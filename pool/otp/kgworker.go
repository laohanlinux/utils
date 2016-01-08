package otp

type kgWorker struct {
	inChan     chan interface{}
	outputChan chan string
}

// args[0] :bool => true: can restart, false: not restart
// args[1] :
func (kgw kgWorker) start(args ...interface{}) error {
	go kgw.work()
	return nil
}

func (kgw kgWorker) terminate() {

}

func (kgw kgWorker) reply() {

}

func (kgw kgWorker) noReply() {

}

func (kgw kgWorker) work() {
	//do some thing
	// select {
	// case opChan := <-kgw.inChan:
	// 	workChan, ok := opChan.(string)
	// 	if !ok {
	//
	// 	}
	// }
}
