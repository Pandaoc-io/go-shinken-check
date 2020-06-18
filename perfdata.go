package sknchk

import (
	"fmt"
	"math/big"
	"strings"
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

func generatePerfOutput(perf []*PerfData) string {
	var perfsSlice []string
	for _, p := range perf {
		if p.Value != nil {
			perfStr := fmt.Sprintf("%v=%.2f%v;%v;%v;%v;%v", p.Name, p.Value, p.Unit, p.Warn, p.Crit, p.Min, p.Max)
			perfsSlice = append(perfsSlice, perfStr)
		}
	}
	return strings.Join(perfsSlice, " ")
}
