package micro

import "github.com/fabriqs/go-micro/messaging"

type Feature struct {
	Name string
	//@deprecated
	Init func(env *Env) error
}

type App struct {
	Name     string
	Features []Feature
	// Router   router.Router
	Env *Env
}

type Env struct {
	Ctx
	Conf interface{}
	//
	Router Router
	//Policy    policy.Manager
	Scheduler  Scheduler
	DataSource DataSource
	Mailer     messaging.Mailer
}

type AppCfg struct {
	Name     string
	Features []Feature
	Router   Router
	DB       DataSource
}

func (e *Env) DB() DataSource {
	return e.DataSource
}

func (e *Env) Config() interface{} {
	return e.Conf
}

func (e *Env) Tx(callback func(tx DataSource) (interface{}, error)) (interface{}, error) {
	var result interface{}
	err := e.DB().Transaction(func(tx DataSource) error {
		result0, err0 := callback(tx)
		result = result0
		return err0
	})
	return result, err
}
