package goany

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestToInt64E(t *testing.T) {
	t.Run("Test with byte", func(t *testing.T) {
		v, err := ToInt64E([]byte("14"))
		assert.NoError(t, err)
		assert.Equal(t, int64(14), v)
	})

	t.Run("Test with *int", func(t *testing.T) {
		var a *int
		v, err := ToInt64E(a)
		assert.NoError(t, err)

		b := 123
		a = &b
		v, err = ToInt64E(a)
		assert.NoError(t, err)
		assert.Equal(t, int64(123), v)
	})

	t.Run("Test with int", func(t *testing.T) {
		v, err := ToInt64E(123)
		assert.NoError(t, err)
		assert.Equal(t, int64(123), v)
	})

	t.Run("Test with int64", func(t *testing.T) {
		v, err := ToInt64E(int64(123))
		assert.NoError(t, err)
		assert.Equal(t, int64(123), v)
	})

	t.Run("Test with string", func(t *testing.T) {
		v, err := ToInt64E("123")
		assert.NoError(t, err)
		assert.Equal(t, int64(123), v)
	})

	t.Run("Test with bool", func(t *testing.T) {
		v, err := ToInt64E(true)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), v)

		v, err = ToInt64E(false)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), v)
	})

	t.Run("Test with float64", func(t *testing.T) {
		v, err := ToInt64E(123.45)
		assert.NoError(t, err)
		assert.Equal(t, int64(123), v) // Note: loss of precision
	})

	t.Run("Test with uint64", func(t *testing.T) {
		v, err := ToInt64E(uint64(123))
		assert.NoError(t, err)
		assert.Equal(t, int64(123), v)
	})

	t.Run("Test with nil", func(t *testing.T) {
		v, err := ToInt64E(nil)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), v)
	})

	t.Run("Test with unsupported type", func(t *testing.T) {
		_, err := ToInt64E([]int{1, 2, 3})
		assert.Error(t, err)
	})

	t.Run("Test with string, no decimal", func(t *testing.T) {
		_, err := ToInt64E("abc")
		assert.Error(t, err)
	})

	t.Run("Test with new []byte unsupported type", func(t *testing.T) {
		type NByte byte
		_, err := ToInt64E([]NByte("123"))
		assert.Error(t, err)
	})

	t.Run("Test with time", func(t *testing.T) {
		now := time.Now()
		v, err := ToInt64E(now)
		assert.NoError(t, err)
		assert.Equal(t, now.Unix(), v)
	})
}

func TestToUint64E(t *testing.T) {
	t.Run("Test with byte", func(t *testing.T) {
		v, err := ToUint64E([]byte("14"))
		assert.NoError(t, err)
		assert.Equal(t, uint64(14), v)
	})

	t.Run("Test with int", func(t *testing.T) {
		v, err := ToUint64E(123)
		assert.NoError(t, err)
		assert.Equal(t, uint64(123), v)
	})

	t.Run("Test with int64", func(t *testing.T) {
		v, err := ToUint64E(int64(123))
		assert.NoError(t, err)
		assert.Equal(t, uint64(123), v)
	})

	t.Run("Test with string", func(t *testing.T) {
		v, err := ToUint64E("123")
		assert.NoError(t, err)
		assert.Equal(t, uint64(123), v)
	})

	t.Run("Test with bool", func(t *testing.T) {
		v, err := ToUint64E(true)
		assert.NoError(t, err)
		assert.Equal(t, uint64(1), v)

		v, err = ToUint64E(false)
		assert.NoError(t, err)
		assert.Equal(t, uint64(0), v)
	})

	t.Run("Test with float64", func(t *testing.T) {
		v, err := ToUint64E(123.45)
		assert.NoError(t, err)
		assert.Equal(t, uint64(123), v) // Note: loss of precision
	})

	t.Run("Test with uint64", func(t *testing.T) {
		v, err := ToUint64E(uint64(123))
		assert.NoError(t, err)
		assert.Equal(t, uint64(123), v)
	})

	t.Run("Test with nil", func(t *testing.T) {
		v, err := ToUint64E(nil)
		assert.NoError(t, err)
		assert.Equal(t, uint64(0), v)
	})

	t.Run("Test with unsupported type", func(t *testing.T) {
		_, err := ToUint64E([]int{1, 2, 3})
		assert.Error(t, err)
	})

	t.Run("Test with string, no decimal", func(t *testing.T) {
		_, err := ToUint64E("abc")
		assert.Error(t, err)
	})

	t.Run("Test with new []byte unsupported type", func(t *testing.T) {
		type NByte byte
		_, err := ToUint64E([]NByte("123"))
		assert.Error(t, err)
	})
}

func TestToFloat64E(t *testing.T) {
	t.Run("Test with byte", func(t *testing.T) {
		v, err := ToFloat64E([]byte("14"))
		assert.NoError(t, err)
		assert.Equal(t, float64(14), v)
	})

	t.Run("Test with int", func(t *testing.T) {
		v, err := ToFloat64E(123)
		assert.NoError(t, err)
		assert.Equal(t, float64(123), v)
	})

	t.Run("Test with int64", func(t *testing.T) {
		v, err := ToFloat64E(int64(123))
		assert.NoError(t, err)
		assert.Equal(t, float64(123), v)
	})

	t.Run("Test with string", func(t *testing.T) {
		v, err := ToFloat64E("123.45")
		assert.NoError(t, err)
		assert.Equal(t, 123.45, v)
	})

	t.Run("Test with bool", func(t *testing.T) {
		v, err := ToFloat64E(true)
		assert.NoError(t, err)
		assert.Equal(t, float64(1), v)

		v, err = ToFloat64E(false)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), v)
	})

	t.Run("Test with float64", func(t *testing.T) {
		v, err := ToFloat64E(123.45)
		assert.NoError(t, err)
		assert.Equal(t, 123.45, v) // Note: loss of precision
	})

	t.Run("Test with uint64", func(t *testing.T) {
		v, err := ToFloat64E(uint64(123))
		assert.NoError(t, err)
		assert.Equal(t, float64(123), v)
	})

	t.Run("Test with nil", func(t *testing.T) {
		v, err := ToFloat64E(nil)
		assert.NoError(t, err)
		assert.Equal(t, float64(0), v)
	})

	t.Run("Test with unsupported type", func(t *testing.T) {
		_, err := ToFloat64E([]int{1, 2, 3})
		assert.Error(t, err)
	})

	t.Run("Test with string, no decimal", func(t *testing.T) {
		_, err := ToFloat64E("abc")
		assert.Error(t, err)
	})

	t.Run("Test with new []byte unsupported type", func(t *testing.T) {
		type NByte byte
		_, err := ToFloat64E([]NByte("123"))
		assert.Error(t, err)
	})
}

func TestToStringE(t *testing.T) {
	t.Run("Test with nil", func(t *testing.T) {
		v, err := ToStringE(nil)
		assert.NoError(t, err)
		assert.Equal(t, "", v)
	})

	t.Run("Test with byte", func(t *testing.T) {
		v, err := ToStringE([]byte("14"))
		assert.NoError(t, err)
		assert.Equal(t, "14", v)
	})

	t.Run("Test with int", func(t *testing.T) {
		v, err := ToStringE(123)
		assert.NoError(t, err)
		assert.Equal(t, "123", v)
	})

	t.Run("Test with string", func(t *testing.T) {
		v, err := ToStringE("test")
		assert.NoError(t, err)
		assert.Equal(t, "test", v)
	})

	t.Run("Test with float", func(t *testing.T) {
		v, err := ToStringE(123.45)
		assert.NoError(t, err)
		assert.Equal(t, "123.45", v)
	})

	t.Run("Test with bool", func(t *testing.T) {
		v, err := ToStringE(true)
		assert.NoError(t, err)
		assert.Equal(t, "true", v)
	})

	t.Run("Test with time", func(t *testing.T) {
		now := time.Now().UTC()
		v, err := ToStringE(now, Options{timeFormat: "2006-01-02"})
		assert.NoError(t, err)
		assert.Equal(t, now.Format("2006-01-02"), v)

		v, err = ToStringE(now)
		assert.NoError(t, err)
		assert.Equal(t, now.Format("2006-01-02 15:04:05"), v)
		fmt.Println(v)
	})

	t.Run("Test with unsupported type", func(t *testing.T) {
		v, err := ToStringE([]int{1, 2, 3})
		assert.NoError(t, err)
		assert.Equal(t, "[1,2,3]", v)
	})
}

func TestToBoolE(t *testing.T) {
	t.Run("Test with nil", func(t *testing.T) {
		v, err := ToBoolE(nil)
		assert.NoError(t, err)
		assert.Equal(t, false, v)
	})

	t.Run("Test with *int", func(t *testing.T) {
		var a *int
		v, err := ToBoolE(a)
		assert.NoError(t, err)
		assert.Equal(t, false, v)
	})

	t.Run("Test with string", func(t *testing.T) {
		v, err := ToBoolE("true")
		assert.NoError(t, err)
		assert.Equal(t, true, v)
	})

	t.Run("Test with int", func(t *testing.T) {
		v, err := ToBoolE(1)
		assert.NoError(t, err)
		assert.Equal(t, true, v)
	})

	t.Run("Test with uint", func(t *testing.T) {
		v, err := ToBoolE(uint(0))
		assert.NoError(t, err)
		assert.Equal(t, false, v)
	})

	t.Run("Test with float", func(t *testing.T) {
		v, err := ToBoolE(2.0)
		assert.NoError(t, err)
		assert.Equal(t, true, v)
	})

	t.Run("Test with bool", func(t *testing.T) {
		v, err := ToBoolE(true)
		assert.NoError(t, err)
		assert.Equal(t, true, v)
	})

	t.Run("Test with byte slice", func(t *testing.T) {
		v, err := ToBoolE([]byte("true"))
		assert.NoError(t, err)
		assert.Equal(t, true, v)
	})

	t.Run("Test with invalid input", func(t *testing.T) {
		_, err := ToBoolE([]int{1, 2, 3})
		assert.Error(t, err)

		_, err = ToBoolE(time.Now())
		assert.Error(t, err)
	})
}

func TestBasicOther(t *testing.T) {
	t.Run("Test ToInt64", func(t *testing.T) {
		v64 := ToInt64([]byte("14"))
		assert.Equal(t, int64(14), v64)
	})

	t.Run("Test ToInt32", func(t *testing.T) {
		v32 := ToInt32(nil)
		assert.Equal(t, int32(0), v32)
	})

	t.Run("Test ToInt16", func(t *testing.T) {
		v16 := ToInt16("16")
		assert.Equal(t, int16(16), v16)
	})

	t.Run("Test ToInt8", func(t *testing.T) {
		v8 := ToInt8("8")
		assert.Equal(t, int8(8), v8)
	})

	t.Run("Test ToInt", func(t *testing.T) {
		v := ToInt("1")
		assert.Equal(t, 1, v)
	})

	t.Run("Test ToUint64", func(t *testing.T) {
		vu64 := ToUint64("14")
		assert.Equal(t, uint64(14), vu64)
	})

	t.Run("Test ToUint32", func(t *testing.T) {
		vu32 := ToUint32("32")
		assert.Equal(t, uint32(32), vu32)
	})

	t.Run("Test ToUint16", func(t *testing.T) {
		vu16 := ToUint16("16")
		assert.Equal(t, uint16(16), vu16)
	})

	t.Run("Test ToUint8", func(t *testing.T) {
		vu8 := ToUint8("8")
		assert.Equal(t, uint8(8), vu8)
	})

	t.Run("Test ToUint", func(t *testing.T) {
		vu := ToUint("1")
		assert.Equal(t, uint(1), vu)
	})

	t.Run("Test ToFloat64", func(t *testing.T) {
		vf64 := ToFloat64("14.0")
		assert.Equal(t, float64(14), vf64)
	})

	t.Run("Test ToFloat32", func(t *testing.T) {
		vf32 := ToFloat32("32.0")
		assert.Equal(t, float32(32), vf32)
	})

	t.Run("Test ToBool", func(t *testing.T) {
		var a *int
		var b = 1
		a = &b
		vb := ToBool(a)
		assert.Equal(t, true, vb)
	})

	t.Run("Test ToString", func(t *testing.T) {
		vs := ToString(1)
		assert.Equal(t, "1", vs)
	})
}
