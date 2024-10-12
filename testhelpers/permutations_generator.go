package testhelpers

import (
	"reflect"
)

// GenerateSlicePermutationsForTests
//
//	Generates simple shallow copies of a struct to be used with gomock.AnyOf().
//
//	It permutates all the slice fields.
//	Ignores slice inside a slice case.
func GenerateSlicePermutationsForTests[T any](in T) []any {
	generated := doGenerate(in)

	res := make([]any, 0, len(generated))
	for _, g := range generated {
		res = append(res, g)
	}
	return res
}

func doGenerate[T any](in T) []T {
	v := reflect.ValueOf(in)
	switch v.Kind() {
	case reflect.Struct:
	// will generate
	case reflect.Ptr:
		if v.Elem().Kind() != reflect.Struct {
			return []T{in}
		}
		// otherwise will generate
		v = v.Elem()
	default:
		return []T{in}
	}

	res := []T{in}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		switch f.Kind() {
		case reflect.Slice:
			if f.Len() == 0 {
				continue
			}
			it := f.Index(0)
			itemType := it.Type()

			items := make([]any, f.Len())
			for i := 0; i < f.Len(); i++ {
				item := f.Index(i)
				items[i] = item.Interface()
			}

			newFieldElems := permutateSlice(items, false)
			newFieldValues := make([]reflect.Value, 0, len(newFieldElems))
			for _, permItems := range newFieldElems {
				elemsValue := reflect.MakeSlice(reflect.SliceOf(itemType), 0, len(permItems))
				for _, perm := range permItems {
					value := reflect.ValueOf(perm).Convert(itemType)
					elemsValue = reflect.Append(elemsValue, value)
				}
				newFieldValues = append(newFieldValues, elemsValue)
			}

			newElems := generateNewElemsWithField(res, newFieldValues, i)
			res = append(res, newElems...)

		case reflect.Struct, reflect.Ptr:
			if f.Kind() == reflect.Ptr && f.Elem().Kind() != reflect.Struct {
				continue
			}
			newFieldElems := doGenerate(f.Interface())
			newFieldElems = newFieldElems[1:] // we skip the first result as it is always the original one
			newFieldValues := make([]reflect.Value, 0, len(newFieldElems))
			for _, elem := range newFieldElems {
				newFieldValues = append(newFieldValues, reflect.ValueOf(elem))
			}

			newElems := generateNewElemsWithField(res, newFieldValues, i)
			res = append(res, newElems...)
		}
	}

	return res
}

func generateNewElemsWithField[T any](
	oldElems []T,
	newFields []reflect.Value,
	indexToSet int,
) []T {
	newElems := make([]T, 0, len(newFields))

	for _, newF := range newFields {
		for _, oldEl := range oldElems {
			newValue := generateNewValue(oldEl, indexToSet)
			isPointer := false
			if newValue.Kind() == reflect.Ptr {
				newValue = newValue.Elem()
				isPointer = true
			}
			newValue.Field(indexToSet).Set(newF)

			if isPointer {
				newElems = append(newElems, newValue.Addr().Interface().(T))
			} else {
				newElems = append(newElems, newValue.Interface().(T))
			}
		}
	}

	return newElems
}

func generateNewValue[V any](oldEl V, doNotCopyIndex int) reflect.Value {
	oldV := reflect.ValueOf(oldEl)
	isPointer := false
	if oldV.Kind() == reflect.Ptr {
		oldV = oldV.Elem()
		isPointer = true
	}

	newEl := reflect.New(oldV.Type())
	newEl = reflect.Indirect(newEl)
	for m := 0; m < oldV.NumField(); m++ {
		if m != doNotCopyIndex && newEl.Field(m).CanSet() {
			oldF := oldV.Field(m)
			newEl.Field(m).Set(oldF)
		}
	}

	if isPointer {
		return newEl.Addr()
	}

	return newEl
}

// copied from https://go.dev/play/p/wKrAlP62_s_s
func permutateSlice[V any](data []V, includeOriginal bool) [][]V {
	if len(data) == 0 {
		return nil
	}

	permutation := make([]V, len(data))
	indexInUse := make([]bool, len(data))

	var res [][]V
	var f func(idx int)

	f = func(idx int) {
		if idx >= len(data) {
			arr := make([]V, len(data))
			copy(arr, permutation)
			res = append(res, arr)
			return
		}
		for i := 0; i < len(data); i++ {
			if !indexInUse[i] {
				indexInUse[i] = true
				permutation[idx] = data[i]
				f(idx + 1)
				indexInUse[i] = false
			}
		}
	}

	f(0)

	finalRes := make([][]V, 0, len(data))
	for _, r := range res {
		if includeOriginal || !reflect.DeepEqual(r, data) {
			finalRes = append(finalRes, r)
		}
	}

	return finalRes
}
