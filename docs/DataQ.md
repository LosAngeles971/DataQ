<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# DataQ

```go
import ".."
```

Main and only package of DataQ dataq\.go defines the Surfer object and its methods

Main and only package of DataQ helpers\.go includes utility functions\, especially the crucial getValueOf

## Index

- [Constants](<#constants>)
- [func Compare(f1 interface{}, f2 interface{}) (bool, error)](<#func-compare>)
- [type Surfer](<#type-surfer>)
  - [func NewSurfer(opts ...SurferOption) *Surfer](<#func-newsurfer>)
  - [func (s Surfer) Get(name string, source interface{}) (interface{}, error)](<#func-surfer-get>)
  - [func (s Surfer) GetBool(name string, source interface{}) (bool, error)](<#func-surfer-getbool>)
  - [func (s Surfer) GetFlatData(source interface{}) (map[string]interface{}, error)](<#func-surfer-getflatdata>)
  - [func (s Surfer) GetFloat64(name string, source interface{}) (float64, error)](<#func-surfer-getfloat64>)
  - [func (s Surfer) GetInt64(name string, source interface{}) (int64, error)](<#func-surfer-getint64>)
  - [func (s Surfer) GetString(name string, source interface{}) (string, error)](<#func-surfer-getstring>)
  - [func (s Surfer) GetVars(source interface{}) ([]string, error)](<#func-surfer-getvars>)
  - [func (s Surfer) SetBool(name string, value bool, source interface{}) error](<#func-surfer-setbool>)
  - [func (s Surfer) SetFloat64(name string, value float64, source interface{}) error](<#func-surfer-setfloat64>)
  - [func (s Surfer) SetInt64(name string, value int64, source interface{}) error](<#func-surfer-setint64>)
  - [func (s Surfer) SetString(name string, value string, source interface{}) error](<#func-surfer-setstring>)
- [type SurferOption](<#type-surferoption>)
  - [func WithSep(sep string) SurferOption](<#func-withsep>)


## Constants

```go
const (
    Default_sep = "."
)
```

## func [Compare](<https://github.com/LosAngeles971/DataQ/blob/main/helpers.go#L63>)

```go
func Compare(f1 interface{}, f2 interface{}) (bool, error)
```

Two fields comparison without first knowing their types

## type [Surfer](<https://github.com/LosAngeles971/DataQ/blob/main/dataq.go#L16-L18>)

```go
type Surfer struct {
    // contains filtered or unexported fields
}
```

### func [NewSurfer](<https://github.com/LosAngeles971/DataQ/blob/main/dataq.go#L229>)

```go
func NewSurfer(opts ...SurferOption) *Surfer
```

NewSurfer creates a pointer to a new Surfer object with default configuration

### func \(Surfer\) [Get](<https://github.com/LosAngeles971/DataQ/blob/main/dataq.go#L30>)

```go
func (s Surfer) Get(name string, source interface{}) (interface{}, error)
```

Get returns the value of the given field

### func \(Surfer\) [GetBool](<https://github.com/LosAngeles971/DataQ/blob/main/dataq.go#L82>)

```go
func (s Surfer) GetBool(name string, source interface{}) (bool, error)
```

GetBool returns the bool value of the given field

### func \(Surfer\) [GetFlatData](<https://github.com/LosAngeles971/DataQ/blob/main/dataq.go#L186>)

```go
func (s Surfer) GetFlatData(source interface{}) (map[string]interface{}, error)
```

GetFlatData returns a map of interface\{\} including all fields extracted from the source

### func \(Surfer\) [GetFloat64](<https://github.com/LosAngeles971/DataQ/blob/main/dataq.go#L39>)

```go
func (s Surfer) GetFloat64(name string, source interface{}) (float64, error)
```

GetBool returns the float64 value of the given field

### func \(Surfer\) [GetInt64](<https://github.com/LosAngeles971/DataQ/blob/main/dataq.go#L56>)

```go
func (s Surfer) GetInt64(name string, source interface{}) (int64, error)
```

GetInt64 returns the int64 value of the given field

### func \(Surfer\) [GetString](<https://github.com/LosAngeles971/DataQ/blob/main/dataq.go#L73>)

```go
func (s Surfer) GetString(name string, source interface{}) (string, error)
```

GetString returns the string value of the given field

### func \(Surfer\) [GetVars](<https://github.com/LosAngeles971/DataQ/blob/main/dataq.go#L143>)

```go
func (s Surfer) GetVars(source interface{}) ([]string, error)
```

GetVars extracts from the source a list of all exportable fields using their fully qualified names

### func \(Surfer\) [SetBool](<https://github.com/LosAngeles971/DataQ/blob/main/dataq.go#L133>)

```go
func (s Surfer) SetBool(name string, value bool, source interface{}) error
```

SetBool updates the given bool field

### func \(Surfer\) [SetFloat64](<https://github.com/LosAngeles971/DataQ/blob/main/dataq.go#L123>)

```go
func (s Surfer) SetFloat64(name string, value float64, source interface{}) error
```

SetFloat64 updates the given float64 field

### func \(Surfer\) [SetInt64](<https://github.com/LosAngeles971/DataQ/blob/main/dataq.go#L113>)

```go
func (s Surfer) SetInt64(name string, value int64, source interface{}) error
```

SetInt64 updates the given int64 field

### func \(Surfer\) [SetString](<https://github.com/LosAngeles971/DataQ/blob/main/dataq.go#L97>)

```go
func (s Surfer) SetString(name string, value string, source interface{}) error
```

SetBool updates the given string field

## type [SurferOption](<https://github.com/LosAngeles971/DataQ/blob/main/dataq.go#L20>)

```go
type SurferOption func(*Surfer)
```

### func [WithSep](<https://github.com/LosAngeles971/DataQ/blob/main/dataq.go#L23>)

```go
func WithSep(sep string) SurferOption
```

WithSep sets the separation string for the fully qualified name of the fields



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)