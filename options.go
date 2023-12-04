package goany

import (
	"github.com/pkg/errors"
	"reflect"
	"time"
)

var (
	ErrDecodeStop           = errors.New("decode stop")
	ErrBasic                = "unable to convert %#v(type %[1]T) to "
	ErrInToOut              = "unable to convert %#v(type %[1]T) to %s"
	ErrFieldNoFound         = "the specified field was not found key: %s"
	ErrUnSupportType        = "unsupported out type %v"
	ErrUnableConvertBasic   = ErrBasic + "basic type"
	ErrUnableConvertInt64   = ErrBasic + "int64"
	ErrUnableConvertUint64  = ErrBasic + "uint64"
	ErrUnableConvertFloat64 = ErrBasic + "float64"
	ErrUnableConvertString  = ErrBasic + "string"
	ErrUnableConvertBool    = ErrBasic + "bool"
	ErrUnableConvertTime    = ErrBasic + "time"
	ErrNotJson              = "the input %#v(type %[1]T) is not json, or not map or slice"
	ErrInNotPtr             = errors.New("if want to export a unexported field, the input must be of pointer type")
)

const (
	TagIgnore = "-"
)

const (
	DecodeContinue = iota // continue with decoding as normal
	DecodeSkip            // skip decoding of the field
	DecodeStop            // stop decoding
)

// some field can customize the parsing, such as time.Duration, net.IP, net.IPNet.
// return DecodeStop if the hook has handled the decoding of the field.
type HookFunc func(in interface{}, out reflect.Value) (int, error)

// Options is a struct for specifying configuration options for any client.
type Options struct {
	location   *time.Location //time zone default is "UTC"
	timeFormat string         //time format default is "2006-01-02 15:04:05"

	mapKeyField  string //map key field,default is index
	mapKeyToList bool   //map key to list,default is false

	tagName string //default is json

	exportedUnExported bool //exported lower field,default is false

	structToMapDetail bool //if out is interface, convert all nest to any,default is false

	assignKey map[string]string //assign key

	hooks []HookFunc //customize the parsing
}

// NewOptions creates a new options. The default options are:
func NewOptions() *Options {
	return &Options{
		location:   time.UTC,
		timeFormat: "2006-01-02 15:04:05",
		tagName:    "json",
	}
}

func (op *Options) SetLocation(v *time.Location) *Options {
	op.location = v
	return op
}

func (op *Options) SetTimeFormat(v string) *Options {
	op.timeFormat = v
	return op
}

func (op *Options) SetMapKeyField(v string) *Options {
	op.mapKeyField = v
	return op
}

func (op *Options) SetMapKeyToList(v bool) *Options {
	op.mapKeyToList = v
	return op
}

func (op *Options) SetTagName(v string) *Options {
	op.tagName = v
	return op
}

func (op *Options) SetExportedUnExported(v bool) *Options {
	op.exportedUnExported = v
	return op
}

func (op *Options) SetStructToMapDetail(v bool) *Options {
	op.structToMapDetail = v
	return op
}

func (op *Options) SetAssignKey(v map[string]string) *Options {
	op.assignKey = v
	return op
}

func (op *Options) AddHook(v HookFunc) *Options {
	op.hooks = append(op.hooks, v)
	return op
}

type anyClient struct {
	options *Options
}

// NewAnyClient creates a new any client.
func newAnyClient(options ...Options) *anyClient {
	cli := &anyClient{
		options: NewOptions(),
	}
	if len(options) > 0 {
		cli.options = &options[0]
	}
	return cli
}
