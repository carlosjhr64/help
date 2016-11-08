package help

import (
  "testing"
)

func TestUnit(test *testing.T) {
  a,b,bad := "","",test.Error
  // kv
  a,b = kv("abc=xyz")
  if a!="abc" || b!="xyz" {bad("kv1")}
  a,b = kv("abc")
  if a!="abc" || b!="" {bad("kv2")}
  a,b = kv("abc=xyz=123")
  if a!="abc" || b!="xyz=123" {bad("kv3")}
  a,b = kv("")
  if a!="" || b!="" {bad("kv4")}
  // squeeze
  a = squeeze("abc=xyz")
  if a!="abc=xyz" || b!="" {bad("squeeze1")}
  a = squeeze(" abc   xyz ")
  if a!="abc xyz" || b!="" {bad("squeeze2")}
  a = squeeze("   abc   xyz\tNotes!")
  if a!="abc xyz" || b!="" {bad("squeeze3")}
  // a_count
  if !a_count("123") {bad("a_count1")}
  if a_count("abc") {bad("a_count2")}
  if a_count("a123") {bad("a_count3")}
  if a_count("123a") {bad("a_count4")}
  if a_count("-123") {bad("a_count5")} // counts are positive numbers > 0
  if !a_count("+123") {bad("a_count6")}
  if a_count("0") {bad("a_count7")} // counts are positive numbers > 0
}

func TestChars(test *testing.T) {
  bad := test.Error
  chars := NewChars("A1xb2Yc3z")
  if string(chars.Bytes)!="A1xb2Yc3z" {bad("Chars1")}
  if chars.Byte!=0 {bad("Chars2")}
  if chars.Index!=0 {bad("Chars3")}
  if !chars.Next() {bad("Chars4")}
  if chars.Byte!='A' {bad("Chars5")}
  if chars.Index!=1 {bad("Chars6")}
  chars.Next(); chars.Next(); chars.Next() // 1,x,b
  if chars.Byte!='b' {bad("Chars7")}
  if chars.Index!=4 {bad("Chars8")}
  chars.Next(); chars.Next(); chars.Next(); chars.Next() // 2,Y,c,3
  if chars.Byte!='3' {bad("Chars7")}
  if chars.Index!=8 {bad("Chars8")}
  if !chars.Next() {bad("Chars9")} // z
  if chars.Next() {bad("Chars10")} // end of chars
  if chars.Byte!=0 {bad("Chars11")}
  if chars.Next() {bad("Chars12")} // end of chars
  if chars.Byte!=0 {bad("Chars13")}
  if chars.Index!=9 {bad("Chars14")} // does not get past string's length
  chars.Reset()
  if chars.Index!=0 {bad("Chars15")} // we get to restart
}

func TestParse(test *testing.T) {
  bad := test.Error
  tokens := Parse("a b c")
  if len(tokens)!=3 {bad("Parse1")}
  if tokens[0].(string)!="a" || tokens[1].(string)!="b" || tokens[2].(string)!="c" {
    bad("Parse2")
  }
  tokens = Parse("abc [xyz 123]")
  if len(tokens)!=2 {bad("Parse3")}
  if tokens[0].(string)!="abc" {bad("Parse4")}
  if tokens[1].([]interface{})[0].(string)!="xyz" {bad("Parse5")}
  if tokens[1].([]interface{})[1].(string)!="123" {bad("Parse6")}
  tokens = Parse("-a [-b -c] --d")
  if len(tokens)!=3 {bad("Parse7")}
  if tokens[2].(string)!="--d" {bad("Parse8")}
}

var T = "\t"
var TEST_HELP = `### The Test Help ###
Usage:
  cmd run <program>
  cmd [:options+] <name>
  cmd <name>+
  cmd2 <prices>+
  cmd3 <numbers>+
Options:
  -a --abc
  --xyz 5.0
  --n 3
  --i
  -j
Types:
  Float --xyz prices
  Int --n numbers
  ^[A-Z][a-z]+$ name
Defaults:
  name Carlos
Notes:
  Blah, blah...
  BLAH!`

func TestGetters(test *testing.T) {
  bad := test.Error
  Testing = true; defer func() {Testing = false}()

  options := new(Options)
  options.Help = TEST_HELP
  options.Version = "1.2.3"
  options.Args = []string{"cmd","-a","Hernandez"}
  a := options.Do()
  if a!="OK" {bad("#1")}
  // Has
  if !options.Has("-a") {bad("#2")}
  if options.Has("--i") {bad("#3")}
  if !options.Has("name") {bad("#4")}
  if !options.Has("--xyz") {bad("#5")} // note that --xyz is defaulted
  // Get
  b := options.Get("name")
  if b!="Hernandez" {bad("#6")}
  b = options.Get("--xyz")
  if b!="5.0" {bad("#7")}
  // Float
  c := options.Float("--xyz")
  if c!=5.0 {bad("#8")}

  options.Args = []string{"cmd","--n=314","Carlos"}
  a = options.Do()
  if a!="OK" {bad("#9")}
  d := options.Int("--n")
  if d!=314 {bad("#10")}
  // Int

  options.Args = []string{"cmd","Carlos","Hernandez"}
  a = options.Do()
  if a!="OK" {bad("#11")}
  e := options.Gets("name")
  if len(e)!=2 {bad("#12")}
  if e[0]!="Carlos" || e[1]!="Hernandez" {bad("#13")}

  options.Args = []string{"cmd2","45.55","63.12", "9.99"}
  a = options.Do()
  if a!="OK" {bad("#14")}
  f := options.Floats("prices")
  if len(f)!=3 {bad("#15")}
  if f[0]!=45.55 || f[1]!=63.12 || f[2]!=9.99 {bad("#16")}

  options.Args = []string{"cmd3","5","7","13","17"}
  a = options.Do()
  if a!="OK" {bad("#17")}
  g := options.Ints("numbers")
  if len(g)!=4 {bad("#18")}
  if g[0]!=5 || g[1]!=7 || g[2]!=13 || g[3]!=17 {bad("#19")}
}


func TestSynonyms(test *testing.T) {
  bad := test.Error
  Testing = true; defer func() {Testing = false}()

  options := new(Options)
  options.Help = TEST_HELP
  options.Version = "1.2.3"

  options.Args = []string{"cmd","Jose"}
  a := options.Do()
  if a!="OK" {bad("A.")}
  if options.Has("-a") {bad("A.1")}
  if options.Has("--abc") {bad("A.2")}

  options.Args = []string{"cmd","-a","Jose"}
  a = options.Do()
  if a!="OK" {bad("B.")}
  if !options.Has("-a") {bad("B.1")}
  if !options.Has("--abc") {bad("B.2")}

  options.Args = []string{"cmd","--abc","Jose"}
  a = options.Do()
  if a!="OK" {bad("C.")}
  if !options.Has("-a") {bad("C.1")}
  if !options.Has("--abc") {bad("C.2")}
}

func TestTypes(test *testing.T) {
  bad := test.Error
  Testing = true; defer func() {Testing = false}()

  options := new(Options)
  options.Help = TEST_HELP
  options.Version = "1.2.3"

  options.Args = []string{"cmd","carlos"}
  a := options.Do()
  if a!="name=carlos !~ /^[A-Z][a-z]+$/" {bad("Types1")}

  options.Args = []string{"cmd","Carlos"}
  a = options.Do()
  if a!="OK" {bad("Types2")}

  options.Args = []string{"cmd2","1.23","345","324.5"}
  a = options.Do()
  if a!="prices !~ /Float/" {bad("Types3")}

  options.Args = []string{"cmd2","1.23","3.45","324.5"}
  a = options.Do()
  if a!="OK" {bad("Types4")}

  options.Args = []string{"cmd","--xyz=123","Okdokie"}
  a = options.Do()
  if a!="--xyz=123 !~ /Float/" {bad("Types5")}

  options.Args = []string{"cmd","--xyz=12.3","Okdokie"}
  a = options.Do()
  if a!="OK" {bad("Types6")}
}

func TestIdiotProofing(test *testing.T) {
  bad := test.Error
  Testing = true; defer func() {Testing = false}()

  options := new(Options)
  options.Help = TEST_HELP
  options.Version = "1.2.3"

  options.Args = []string{"cmd","--yzx=1.0","Carlos"}
  a := options.Do()
  if a!="Did not match usage." {bad("Idiot1")}

  options.Args = []string{"cmd","--xyz=1.0","Carlos"}
  a = options.Do()
  if a!="OK" {bad("Idiot2")}
}

func TestLiteral(test *testing.T) {
  bad := test.Error
  Testing = true; defer func() {Testing = false}()

  options := new(Options)
  options.Help = TEST_HELP
  options.Version = "1.2.3"

  options.Args = []string{"cmd","run","awesome"} // notice that first match executes
  a := options.Do()
  if a!="OK" {bad("Lit1")}
  b := options.Get("program")
  if b!="awesome" {bad("Lit2")}
}

func TestBasic(test *testing.T) {
  bad := test.Error
  Testing = true; defer func() {Testing = false}()

  options := new(Options)
  options.Help = TEST_HELP
  options.Version = "1.2.3"

  options.Args = []string{"cmd","-v"}
  a := options.Do()
  if a!="1.2.3" {bad("Version")}

  options.Args = []string{"cmd","--version"}
  a = options.Do()
  if a!="1.2.3" {bad("Long Version")}

  options.Args = []string{"cmd","-h"}
  a = options.Do()
  if a[0:21]!="### The Test Help ###" {bad("Help")}

  options.Args = []string{"cmd","--help"}
  a = options.Do()
  if a[0:21]!="### The Test Help ###" {bad("Long Help")}
}

func TestNew(test *testing.T) {
  bad := test.Error
  Testing = true; defer func() {Testing = false}()
  // Note that os.Args is empty so only works b/c of Testing set to true.
  options := New("a.b.c",TEST_HELP)
  if options.Version != "a.b.c" {bad("New1")}
  if options.Help != TEST_HELP {bad("New2")}
}

var SNEAKY_HELP = `### The Sneaky One ###
Usage:
  delta [:options] <stock> [<factor>]
Options:
  --decay 1008
  --days
  --div   1
Defaults:
  factor  1.0
Types:
  ^[a-z]+$   stock
  Float      factor
  Int        --decay --days --div`

func TestSneaky(test *testing.T) {
  bad := test.Error
  Testing = true; defer func() {Testing = false}()

  options := new(Options)
  options.Help = SNEAKY_HELP
  options.Version = "0.0.1611070"
  options.Args = []string{"delta","spy"}
  a := options.Do()
  if a!="OK" {bad("Sneaky got me!")}
}
