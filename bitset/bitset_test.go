package bitset

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		length uint
	}
	tests := []struct {
		name     string
		args     args
		wantBset *BitSet
	}{
		// TODO: Put test cases.

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotBset := New(tt.args.length); !reflect.DeepEqual(gotBset, tt.wantBset) {
				t.Errorf("New() = %v, want %v", gotBset, tt.wantBset)
			}
		})
	}
}

func TestBitSetAndGet(t *testing.T) {
	v := New(1000)
	v.Set(100).Set(11)
	if !v.Contain(11) {
		t.Errorf("Bit %d is clear, and it shouldn't be.", 100)
	}

	if v.Intersection(New(100).Set(10)).Count() > 1 {
		fmt.Println("Intersection works.")
	}

	fmt.Println(Cap())
}

func TestStringer(t *testing.T) {
	v := New(0)
	for i := uint(0); i < 10; i++ {
		v.Set(i)
	}
	if v.String() != "{0,1,2,3,4,5,6,7,8,9}" {
		t.Error("bad string output")
	}
}
