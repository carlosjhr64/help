### Sinopsis ###

    import "github.com/carlosjhr64/help"
    const HELP = `### cmd ###
    Usage:
      cmd [:options+] <args>+
    Options:
      -q --quiet   `+"\t"+`Be vewy quiet.
      --n 2        `+"\t"+`Example: --n=3, defaults to 2.
      --price 100.0
    Types:
      Int --n
      Float --price
    Notes:
      Everything after "Notes:" is ignored.
      Everything  to the right of a tab(\t) is ignored.`
    const VERSION = "1.2.3"
    # Automatically exits on -v,--version,-h,--help.
    # Will exit(64) on usage errors.
    var opt = help.New(VERSION,HELP)

### Getter methods ###

    opt.Has("-q") // bool
    opt.Get("--n") // string
    opt.Gets("args") // []string
    opt.Float("--price") // float64
    opt.Int("--n") // int
    opt.Floats("args") // []float64 if args can be interpreted that way
    opt.Ints("args") // []int if args can be interpreted that way

### go doc dump ###

package help // import "github.com/carlosjhr64/help"

var BAD_USAGE = "Did not match usage."
var COLOR = "\x1b[31;1m"
var ERROR = "Error: "
var Testing = false
var Types = map[string]*regexp.Regexp{
	"Float": regexp.MustCompile("^[+-]?\\d+\\.\\d+$"),
	"Int":   regexp.MustCompile("^[+-]?\\d+$")}
var USAGE = "usage"
var V = "-v"
func Parse(line string) []interface{}
func NewChars(line string) *Chars
func New(version, help string) *Options
type Chars struct { ... }
type Options struct { ... }
type Chars struct {
	Bytes []byte
	Index int
	Byte  byte
}

func NewChars(line string) *Chars
func (chars *Chars) Next() bool
func (chars *Chars) Parse() []interface{}
func (chars *Chars) Reset()
type Options struct {
	Version string
	Help    string
	Args    []string
	Hash    map[string]string
	Keys    []string
	Dict    map[string][][]string
	Usage   [][]interface{}
	Cache   map[string]interface{}
}

func New(version, help string) *Options
func (options *Options) Dictionary()
func (options *Options) Do() string
func (options *Options) Float(k string) float64
func (options *Options) Floats(k string) []float64
func (options *Options) Get(k string) string
func (options *Options) Gets(k string) []string
func (options *Options) Has(k string) bool
func (options *Options) InDict(name, word string) bool
func (options *Options) Int(k string) int
func (options *Options) Ints(k string) []int
func (options *Options) Matches(pattern []interface{}, i int) int
func (options *Options) Rehash()
func (options *Options) Synonyms()
func (options *Options) Types() string
func (options *Options) Valid() bool
