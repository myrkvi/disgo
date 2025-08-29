package optparse

import "reflect"

// derefType returns the non-pointer type and how many pointers were removed.
func derefType(t reflect.Type) (base reflect.Type, ptrDepth int) {
	base = t
	for base.Kind() == reflect.Ptr {
		base = base.Elem()
		ptrDepth++
	}
	return base, ptrDepth
}

// setValue assigns v to field, rebuilding pointers if ptrDepth > 0 and handling convertibility.
func setValue(field reflect.Value, v any, ptrDepth int) {
	value := reflect.ValueOf(v)

	// Rebuild pointer layers if field expects pointers.
	for i := 0; i < ptrDepth; i++ {
		p := reflect.New(value.Type())
		p.Elem().Set(value)
		value = p
	}

	if value.Type().AssignableTo(field.Type()) {
		field.Set(value)
		return
	}
	if value.Type().ConvertibleTo(field.Type()) {
		field.Set(value.Convert(field.Type()))
		return
	}
}

// Integer helpers with overflow checks and pointer support.
func setInt(field reflect.Value, vi int64, ptrDepth int) {
	target := field
	if ptrDepth > 0 {
		// Build pointer chain to the base kind
		base := field.Type()
		for i := 0; i < ptrDepth; i++ {
			base = base.Elem()
		}
		p := reflect.New(base)
		target = p.Elem()
		defer field.Set(p)
	}
	if target.OverflowInt(vi) {
		return
	}
	target.SetInt(vi)
}

func setUint(field reflect.Value, vu uint64, ptrDepth int) {
	target := field
	if ptrDepth > 0 {
		base := field.Type()
		for i := 0; i < ptrDepth; i++ {
			base = base.Elem()
		}
		p := reflect.New(base)
		target = p.Elem()
		defer field.Set(p)
	}
	if target.OverflowUint(vu) {
		return
	}
	target.SetUint(vu)
}

func setFloat(field reflect.Value, vf float64, ptrDepth int) {
	target := field
	if ptrDepth > 0 {
		base := field.Type()
		for i := 0; i < ptrDepth; i++ {
			base = base.Elem()
		}
		p := reflect.New(base)
		target = p.Elem()
		defer field.Set(p)
	}

	switch target.Kind() {
	case reflect.Float32:
		target.SetFloat(float64(float32(vf)))
	case reflect.Float64:
		target.SetFloat(vf)
	}
}
