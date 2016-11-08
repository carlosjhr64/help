package help

import (
  "os"
  "fmt"
  "regexp"
  "strconv"
  "strings"
)

// CONSTANTS and GLOBALS

var Testing = false

var V,LV,H,LH = "-v","--version","-h","--help"
var USAGE,NOTES,TYPES = "usage","notes","types"
var ERROR = "Error: "
var BAD_USAGE = "Did not match usage."
var COLOR,NO_COLOR = "\x1b[31;1m","\x1b[0m"

var Types = map[string]*regexp.Regexp{
  "Float": regexp.MustCompile("^[+-]?\\d+\\.\\d+$"),
  "Int":   regexp.MustCompile("^[+-]?\\d+$") }

// UTILITIES

func kv(a string) (string,string) {
  k,v := a,""
  i := strings.Index(a,"=")
  if i>-1 {k,v = a[0:i],a[i+1:]}
  return k,v
}

// TODO: see uplc.go for better way
func squeeze(l string) string {
  line,previous := "",byte(' ')
  for _,c := range []byte(l) {
    if c=='\t' {break}
    if c==' ' && previous==' ' {continue}
    line += string(c)
    previous = c
  }
  if previous==' ' {return line[0:len(line)-1]}
  return line
}

func a_count(a string) bool {
  // TODO best way?
  i,_ := strconv.Atoi(a)
  return i>0
}

// Chars struct for Parse and METHODS

type Chars struct {
  Bytes []byte
  Index int
  Byte byte
}

func NewChars(line string) *Chars {
  return &Chars{[]byte(line),0,0}
}

func (chars *Chars) Next() bool {
  if chars.Index <  len(chars.Bytes){
    chars.Byte = chars.Bytes[chars.Index]
    chars.Index++
    return true
  }
  chars.Byte = 0
  return false
}

func (chars *Chars) Reset() {
  chars.Index=0
}

func (chars *Chars) Parse() []interface{} {
  tokens,token := make([]interface{},0),""
  for chars.Next() {
    c := chars.Byte
    switch c {
    case ' ','[',']':
      if token != "" {
        tokens = append(tokens,token)
        token = ""
      }
      if c=='[' {tokens = append(tokens,chars.Parse())}
      if c==']' {return tokens}
    default:
      token += string(c)
    }
  }
  if token != "" {tokens = append(tokens,token)}
  return tokens
}

func Parse(line string) []interface{} {
  chars := NewChars(line)
  return chars.Parse()
}

// Options struct AND METHODS

type Options struct {
  Version string
  Help string
  Args []string
  Hash map[string]string
  Keys []string
  Dict map[string][][]string
  Usage [][]interface{}
  Cache map[string]interface{}
}

func (options *Options) Has(k string) bool {
  e := false
  if []byte(k)[0]=='-' {
    _,e = options.Hash[k]
  } else {
    _,e = options.Cache[k]
  }
  return e
}

func (options *Options) Get(k string) string {
  word := ""
  if []byte(k)[0]=='-' {
    word = options.Hash[k]
  } else {
    words := options.Cache[k]
    switch words.(type){
    case string:
      word = words.(string)
    default:
      panic(fmt.Sprintf("'%s' is []string, not string.",k))
    }
  }
  return word
}

func (options *Options) Float(k string) float64 {
  f,e := strconv.ParseFloat(options.Get(k), 64)
  if e!=nil {panic(e)}
  return f
}

func (options *Options) Int(k string) int {
  i,e := strconv.Atoi(options.Get(k))
  if e!=nil {panic(e)}
  return i
}

func (options *Options) Gets(k string) []string {
  w,e := options.Cache[k]
  if !e {panic(fmt.Sprintf("'%s' not in Cache",k))}
  switch w.(type){
  case string:
    panic(fmt.Sprintf("'%s' not []string",k))
  default:
    return w.([]string)
  }
}

func (options *Options) Floats(k string) []float64 {
  as := options.Gets(k)
  fs := make([]float64,len(as))
  for i,a := range as {
    f,e := strconv.ParseFloat(a,64)
    if e!=nil {panic(e)}
    fs[i]=f
  }
  return fs
}

func (options *Options) Ints(k string) []int {
  as := options.Gets(k)
  ns := make([]int,len(as))
  for i,a := range as {
    n,e := strconv.Atoi(a)
    if e!=nil {panic(e)}
    ns[i]=n
  }
  return ns
}

func (options *Options) Synonyms() {
  hash := options.Hash
  for name,words := range options.Dict {
    if name==TYPES || name==USAGE || name=="-" {continue}
    for _, synonyms := range words {
      if len(synonyms)!=2 {continue}
      s,l := synonyms[0],synonyms[1]
      if len(s) < 2 {continue}
      b := []byte(s)
      if b[0]!='-' {continue}
      first,es := hash[s]
      if b[1]=='-' { //long given a default
        if !es {hash[s]=l}
      } else { // first=long and vice-versa
        long,el := hash[l]
        if el && !es {
          hash[s] = long
        }else if !el && es{
          hash[l] = first
        }
      }
    }
  }
}

func match_strings(rgx *regexp.Regexp, words ...string) bool {
  for _,word := range words {
    if !rgx.MatchString(word) { return false }
  }
  return true
}

func (options *Options) Types() string {
  word := ""
  for _,keys := range options.Dict[TYPES] {
    rgs := keys[0]
    rgx,e := Types[rgs]
    if !e {rgx = regexp.MustCompile(rgs)}
    for _,key := range keys[1:] {
      if []byte(key)[0]=='-' {
        word,e = options.Hash[key]
        if !e {continue}
      } else {
        wrd,e := options.Cache[key]
        if !e {continue}
        switch wrd.(type) {
        case string:
          word = wrd.(string)
        default:
          words := wrd.([]string)
          if !match_strings(rgx,words...) {
            return fmt.Sprintf("%s !~ /%s/",key,rgs)
          }
          continue
        }
      }
      if !match_strings(rgx,word) {
        return fmt.Sprintf("%s=%s !~ /%s/",key,word,rgs)
      }
    }
  }
  return "OK"
}

func (options *Options) InDict(name, word string) bool {
  if word[0]!='-' {return false}
  list := options.Dict[name]
  for _,words := range list {
    for _,synonym := range words {
      if synonym==word {return true}
    }
  }
  return false
}

func (options *Options) Matches(pattern []interface{}, i int) int {
  name,keys,hash,cache := "",options.Keys,options.Hash,options.Cache
  for _,token := range pattern {
    switch token.(type) {
    case []interface{}:
      j := options.Matches(token.([]interface{}), i)
      if j>0 {i=j}
      continue
    case string:
      if i>=len(keys) {return 0}
      key := keys[i]
      s := token.(string)
      b,l := []byte(s),len(s)
      x,z := b[0],b[l-1]
      y := x; if l>1 {y=b[l-2]}
      if x==':' && (y!=':' && y!='+') {
        // Selection
        if z=='+' {name=s[1:l-1]} else {name=s[1:l]}
        if !options.InDict(name,key) {return 0}
        if z=='+' {
          for i+1<len(keys) && options.InDict(name,keys[i+1]) {i++}
        }
      } else if x=='<' && (z=='>' || (y=='>' && z=='+')) {
        // Variable
        if !a_count(key) {return 0} // a no match error
        if z=='+' {
          name=s[1:l-2]
          words := make([]string,1)
          words[0] = hash[key]
          for i+1<len(keys) && a_count(keys[i+1]) {
            i++
            words = append(words, hash[keys[i]])
          }
          cache[name] = words
        } else {
          name=s[1:l-1]
          cache[name] = hash[key]
        }
      } else {
        // Literal
        if x=='-' {
          if s!=key {return 0}
        } else {
          if s!=hash[key] {return 0}
        }
      }
    default:
      panic("expected either string or []interface{}")
    }
    i++
  }
  return i
}

func (options *Options) Valid() bool {
  options.Cache = make(map[string]interface{})
  for _, pattern := range options.Usage {
    if options.Matches(pattern,0)==len(options.Keys) {return true}
  }
  return false
}

func (options *Options) Dictionary() {
  dict := make(map[string][][]string)
  usage := make([][]interface{},0)
  name := "-"
  lines := strings.Split(options.Help, "\n")
  for _,line := range lines {
    line = squeeze(line)
    chars := []byte(line)
    if chars[0]=='#' {continue}
    last := len(chars)-1
    if chars[last]==':' {
      name = strings.ToLower(line[0:last])
      if name == USAGE {continue}
      if name == NOTES {break}
      dict[name] = make([][]string,0)
      continue
    }
    if name==USAGE {
      usage = append(usage, Parse(line))
      continue
    }
    dict[name] = append(dict[name], strings.Split(line," "))
  }
  options.Dict = dict
  options.Usage = usage
}

func (options *Options) Rehash() {
  hash := make(map[string]string)
  keys := make([]string,0)
  n := 0
  for i,a := range options.Args {
    b := []byte(a)
    if b[0]=='-' {
      if a=="-" {
        hash[a] = strings.Join(options.Args[(i+1):]," ")
        keys = append(keys, a)
        break // done!
      } else {
        if b[1]=='-' {
	  k,v := kv(a)
          hash[k] = v
          keys = append(keys, k)
	} else {
          for _,c := range b[1:] {
            k := "-"+string(c)
            hash[k] = ""
            keys = append(keys, k)
          }
	}
      }
    } else {
      k := strconv.Itoa(n)
      hash[k] = a
      keys = append(keys, k)
      n++
    }
  }
  options.Hash = hash
  options.Keys = keys
}

func (options *Options) Do() string {
  options.Rehash()
  if options.Has(V) || options.Has(LV) {
    if Testing {return options.Version}
    fmt.Println(options.Version)
    os.Exit(0)
  }
  if options.Has(H) || options.Has(LH) {
    if Testing {return options.Help}
    fmt.Println(options.Help)
    os.Exit(0)
  }
  options.Dictionary()
  if !options.Valid() {
    if Testing {return BAD_USAGE}
    fmt.Fprintln(os.Stderr, COLOR+ERROR+BAD_USAGE+NO_COLOR)
    fmt.Fprintln(os.Stderr, options.Help)
    os.Exit(64) // Usage Error
  }
  msg := options.Types()
  if msg != "OK" {
    if Testing {return msg}
    fmt.Fprintln(os.Stderr, COLOR+ERROR+msg+NO_COLOR)
    os.Exit(64) // Usage Error
  }
  options.Synonyms()
  return "OK"
}

func New(version, help string) *Options {
  options := new(Options)
  options.Version = version
  options.Help = help
  options.Args = os.Args
  options.Do()
  return options
}
