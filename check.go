package sknchk

import (
	"fmt"
	"math/big"
	"os"
	"strings"
)

//Standard Shinken/Nagios-like Return Code
const (
	RC_OK Status = iota
	RC_WARNING
	RC_CRITICAL
	RC_UNKNOWN
)

//Prefix used for the final output
const (
	PREFIX_CLI_OK        string = "[OK]"
	PREFIX_CLI_WARNING   string = "[WARNING]"
	PREFIX_CLI_CRITICAL  string = "[CRITICAL]"
	PREFIX_CLI_UNKNOWN   string = "[UNKNOWN]"
	PREFIX_HTML_OK       string = `<span style="color:white; background-color: #28a745; display: inline-block; padding: .25em .4em; font-size: 75%; font-weight: 700; line-height: 1; text-align: center; white-space: nowrap; vertical-align: baseline; border-radius: .25rem;">OK</span>`
	PREFIX_HTML_WARNING  string = `<span style="color:#212529; background-color: #ffc107; display: inline-block; padding: .25em .4em; font-size: 75%; font-weight: 700; line-height: 1; text-align: center; white-space: nowrap; vertical-align: baseline; border-radius: .25rem;">Warning</span>`
	PREFIX_HTML_CRITICAL string = `<span style="color:white; background-color: #dc3545; display: inline-block; padding: .25em .4em; font-size: 75%; font-weight: 700; line-height: 1; text-align: center; white-space: nowrap; vertical-align: baseline; border-radius: .25rem;">Critical</span>`
	PREFIX_HTML_UNKNOWN  string = `<span style="color:white; background-color: #6c757d; display: inline-block; padding: .25em .4em; font-size: 75%; font-weight: 700; line-height: 1; text-align: center; white-space: nowrap; vertical-align: baseline; border-radius: .25rem;">Unknown</span>`
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

func formatOutput(str string, class string) string {
	if Output.mode == "html" {
		return fmt.Sprintf(`<span style="%v">%v</span>`, class, str)
	}
	return str
}

//FmtOk will format the OK output string depending on the output choosen mode
func FmtOk(str string) string {
	return formatOutput(str, Output.classOk)
}

//FmtWarning will format the Warning output string depending on the output choosen mode
func FmtWarning(str string) string {
	return formatOutput(str, Output.classWarning)
}

//FmtCritical will format the Critical output string depending on the output choosen mode
func FmtCritical(str string) string {
	return formatOutput(str, Output.classCritical)
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
	c.rc = append(c.rc, RC_OK)
}

//AddWarning will add an Warning status to the Check structure
//to prepare the final output/RC
func (c *Check) AddWarning() {
	c.rc = append(c.rc, RC_WARNING)
}

//AddCritical will add an Critical status with some short and long information to the Check structure
//to prepare the final output/RC
func (c *Check) AddCritical() {
	c.rc = append(c.rc, RC_CRITICAL)
}

//AddUnkown will add an Unknown status with some short and long information to the Check structure
//to prepare the final output/RC
func (c *Check) AddUnkown() {
	c.rc = append(c.rc, RC_UNKNOWN)
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
	check := &Check{[]string{short}, sLong, nil, []Status{RC_OK}}
	Exit(check)
}

//Warning will exit the program with the Warning status
func Warning(short string, long string) {
	var sLong []string
	if len(long) > 0 {
		sLong = append(sLong, long)
	}
	check := &Check{[]string{short}, sLong, nil, []Status{RC_WARNING}}
	Exit(check)
}

//Critical will exit the program with the Critical status
func Critical(short string, long string) {
	var sLong []string
	if len(long) > 0 {
		sLong = append(sLong, long)
	}
	check := &Check{[]string{short}, sLong, nil, []Status{RC_CRITICAL}}
	Exit(check)
}

//Unknwown will exit the program with the Unknown status
func Unknwown(short string, long string) {
	var sLong []string
	if len(long) > 0 {
		sLong = append(sLong, long)
	}
	check := &Check{[]string{short}, sLong, nil, []Status{RC_UNKNOWN}}
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
		case RC_OK:
			prefix = PREFIX_HTML_OK
		case RC_WARNING:
			prefix = PREFIX_HTML_WARNING
		case RC_CRITICAL:
			prefix = PREFIX_HTML_CRITICAL
		case RC_UNKNOWN:
			prefix = PREFIX_HTML_UNKNOWN
		}
	} else {
		switch rc {
		case RC_OK:
			prefix = PREFIX_CLI_OK
		case RC_WARNING:
			prefix = PREFIX_CLI_WARNING
		case RC_CRITICAL:
			prefix = PREFIX_CLI_CRITICAL
		case RC_UNKNOWN:
			prefix = PREFIX_CLI_UNKNOWN
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

func generatePerfOutput(perf []*PerfData) string {
	var perfsSlice []string
	for _, p := range perf {
		perfStr := fmt.Sprintf("%v=%.2f;%v;%v;%v;%v;%v", p.Name, p.Value, p.Unit, p.Warn, p.Crit, p.Min, p.Max)
		perfsSlice = append(perfsSlice, perfStr)
	}
	return strings.Join(perfsSlice, " ")
}
