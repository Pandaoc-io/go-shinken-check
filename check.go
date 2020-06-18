package sknchk

import (
	"fmt"
	"math/big"
	"os"
	"strings"
)

//Standard Shinken/Nagios-like Return Code
const (
	RcOk Status = iota
	RcWarning
	RcCritical
	RcUnknwon
)

//Prefix used for the final output
const (
	PrefixCliOk        string = "[OK]"
	PrefixCliWarning   string = "[WARNING]"
	PrefixCliCritical  string = "[CRITICAL]"
	PrefixCliUnknown   string = "[UNKNOWN]"
	PrefixHTMLOk       string = `<span style="color:white; background-color: #28a745; display: inline-block; padding: .25em .4em; font-size: 75%; font-weight: 700; line-height: 1; text-align: center; white-space: nowrap; vertical-align: baseline; border-radius: .25rem;">OK</span>`
	PrefixHTMLWarning  string = `<span style="color:#212529; background-color: #ffc107; display: inline-block; padding: .25em .4em; font-size: 75%; font-weight: 700; line-height: 1; text-align: center; white-space: nowrap; vertical-align: baseline; border-radius: .25rem;">Warning</span>`
	PrefixHTMLCritical string = `<span style="color:white; background-color: #dc3545; display: inline-block; padding: .25em .4em; font-size: 75%; font-weight: 700; line-height: 1; text-align: center; white-space: nowrap; vertical-align: baseline; border-radius: .25rem;">Critical</span>`
	PrefixHTMLUnknown  string = `<span style="color:white; background-color: #6c757d; display: inline-block; padding: .25em .4em; font-size: 75%; font-weight: 700; line-height: 1; text-align: center; white-space: nowrap; vertical-align: baseline; border-radius: .25rem;">Unknown</span>`
)

//Status type used to define the status of the check
type Status int

//OutputMode to display the check result, it can be 'cli' or 'html'
type OutputMode struct {
	mode          string
	newLine       string
	bullet        string
	classOk       string
	classWarning  string
	classCritical string
}

//Output is a global variable that define the type of output, 'cli' by default
var Output *OutputMode = &OutputMode{
	mode:          "cli",
	newLine:       "",
	bullet:        " - ",
	classOk:       "",
	classWarning:  "",
	classCritical: "",
}

//Check struct
type Check struct {
	short    []string
	long     []string
	perfData []*PerfData
	rc       []Status
}

//Mode return the current output mode type
func (o *OutputMode) Mode() string {
	return o.mode
}

//SetHTML will set the output format to html format
func (o *OutputMode) SetHTML() {
	Output.mode = "html"
	Output.newLine = "<br />"
	Output.bullet = "&#8226;&#8194;"
	Output.classOk = "color: #28a745!important;"
	Output.classWarning = "color: #ffc107!important;"
	Output.classCritical = "color: #dc3545!important;"
}

//SetDebug will set the output format to debug format
func (o *OutputMode) SetDebug() {
	Output.mode = "debug"
	Output.newLine = "\n"
	Output.bullet = "- "
	Output.classOk = ""
	Output.classWarning = ""
	Output.classCritical = ""
}

func fmtOutput(str string, class string) string {
	if Output.mode == "html" {
		return fmt.Sprintf(`<span style="%v">%v</span>`, class, str)
	}
	return str
}

//FmtOk will format the OK output string depending on the output choosen mode
func FmtOk(str string) string {
	return fmtOutput(str, Output.classOk)
}

//FmtWarning will format the Warning output string depending on the output choosen mode
func FmtWarning(str string) string {
	return fmtOutput(str, Output.classWarning)
}

//FmtCritical will format the Critical output string depending on the output choosen mode
func FmtCritical(str string) string {
	return fmtOutput(str, Output.classCritical)
}

//AddShort add a new string to the short output
func (c *Check) AddShort(short string, bullet bool) {
	if bullet {
		c.short = append(c.short, Output.bullet+short)
	} else {
		c.short = append(c.short, short)
	}
}

//PrependShort add a new string at the begining of the short output
func (c *Check) PrependShort(short string, bullet bool) {
	if bullet {
		c.short = append([]string{Output.bullet + short}, c.short...)
	} else {
		c.short = append([]string{short}, c.short...)
	}
}

//AddLong add a new string to the long output
func (c *Check) AddLong(long string, bullet bool) {
	if bullet {
		c.long = append(c.long, Output.bullet+long)
	} else {
		c.long = append(c.long, long)
	}
}

//AddOk will add an Ok status to the Check structure
//to prepare the final output/RC
func (c *Check) AddOk() {
	c.rc = append(c.rc, RcOk)
}

//AddWarning will add an Warning status to the Check structure
//to prepare the final output/RC
func (c *Check) AddWarning() {
	c.rc = append(c.rc, RcWarning)
}

//AddCritical will add an Critical status with some short and long information to the Check structure
//to prepare the final output/RC
func (c *Check) AddCritical() {
	c.rc = append(c.rc, RcCritical)
}

//AddUnknown will add an Unknown status with some short and long information to the Check structure
//to prepare the final output/RC
func (c *Check) AddUnknown() {
	c.rc = append(c.rc, RcUnknwon)
}

//AddPerfData add a new perfdata to the check
func (c *Check) AddPerfData(name string, value *big.Float, unit string, warn *big.Float, crit *big.Float, min *big.Float, max *big.Float) {
	c.perfData = append(c.perfData, &PerfData{
		Name:  name,
		Value: value,
		Unit:  unit,
		Warn:  warn,
		Crit:  crit,
		Min:   min,
		Max:   max})
}

//Ok will exit the program with the OK status
func Ok(short string, long string) {
	var sLong []string
	if len(long) > 0 {
		sLong = append(sLong, long)
	}
	check := &Check{[]string{short}, sLong, nil, []Status{RcOk}}
	Exit(check)
}

//Warning will exit the program with the Warning status
func Warning(short string, long string) {
	var sLong []string
	if len(long) > 0 {
		sLong = append(sLong, long)
	}
	check := &Check{[]string{short}, sLong, nil, []Status{RcWarning}}
	Exit(check)
}

//Critical will exit the program with the Critical status
func Critical(short string, long string) {
	var sLong []string
	if len(long) > 0 {
		sLong = append(sLong, long)
	}
	check := &Check{[]string{short}, sLong, nil, []Status{RcCritical}}
	Exit(check)
}

//Unknown will exit the program with the Unknown status
func Unknown(short string, long string) {
	var sLong []string
	if len(long) > 0 {
		sLong = append(sLong, long)
	}
	check := &Check{[]string{short}, sLong, nil, []Status{RcUnknwon}}
	Exit(check)
}

//Rc return the Return Code of the check
func (c *Check) Rc() Status {
	var maxRc Status
	for _, value := range c.rc {
		if value > maxRc {
			maxRc = value
		}
	}
	return maxRc
}

//Exit quit the program displaying the short and long output with the
func Exit(c *Check) {
	rc := c.Rc()
	var prefix string
	if Output.mode == "html" {
		switch rc {
		case RcOk:
			prefix = PrefixHTMLOk
		case RcWarning:
			prefix = PrefixHTMLWarning
		case RcCritical:
			prefix = PrefixHTMLCritical
		case RcUnknwon:
			prefix = PrefixHTMLUnknown
		}
	} else {
		switch rc {
		case RcOk:
			prefix = PrefixCliOk
		case RcWarning:
			prefix = PrefixCliWarning
		case RcCritical:
			prefix = PrefixCliCritical
		case RcUnknwon:
			prefix = PrefixCliUnknown
		}
	}

	perfStr := ""
	if len(c.perfData) > 0 {
		perfStr = "|" + generatePerfOutput(c.perfData)
	}

	if len(c.long) > 0 {
		c.AddShort(Output.newLine+"For more details see long output.", false)
		fmt.Fprintf(os.Stdout, "%v %v\n%v%v", prefix, strings.Join(c.short, Output.newLine), strings.Join(c.long, Output.newLine), perfStr)
	} else {
		fmt.Fprintf(os.Stdout, "%v %v%v", prefix, strings.Join(c.short, Output.newLine), perfStr)
	}
	os.Exit(int(rc))
}
