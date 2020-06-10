package sknchk

import (
	"math/big"
)

//PerfData is the nagios style perf element with the following output pattern : 'label'=value[UOM];[warn];[crit];[min];[max]
type PerfData struct {
	Name  string
	Value *big.Float
	Unit  string
	Warn  *big.Float
	Crit  *big.Float
	Min   *big.Float
	Max   *big.Float
}
