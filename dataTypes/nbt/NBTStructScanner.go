package nbt

import (
	"fmt"
	"reflect"
	"strings"
)

// NBTStructScan converts an NBT Compound into an object
func NBTStructScan(obj interface{}, compound *Compound) {
	v := reflect.ValueOf(obj).Elem()
	t := reflect.TypeOf(obj).Elem()

	fieldMap := make(map[string]string)
	for fieldIndex := 0; fieldIndex < t.NumField(); fieldIndex++ {
		field := t.Field(fieldIndex)
		tag, exists := field.Tag.Lookup("nbt")
		if !exists || tag == "-" {
			continue
		}
		nbtSettings := strings.Split(tag, ",")
		fieldMap[nbtSettings[0]] = field.Name
	}

	for name, tag := range compound.Data {
		fieldName := name
		_, exists := t.FieldByName(name)
		if !exists {
			newFieldName, ok := fieldMap[name]
			if !ok {
				continue
			}
			fieldName = newFieldName
		}

		field := v.FieldByName(fieldName)
		fieldType, _ := t.FieldByName(fieldName)

		// map[string]* can be used as an alternative for compounds
		if tag.GetType() == 10 && fieldType.Type.Kind() == reflect.Map {
			compound := tag.(*Compound)
			newMap := reflect.MakeMap(fieldType.Type)
			for key, val := range compound.Data {
				newMap.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val.GetValue()))
			}
			field.Set(newMap)
		} else if tag.GetType() == 9 {
			list := tag.(*List)
			newSlice := reflect.MakeSlice(fieldType.Type, len(list.Data), len(list.Data))
			for i, value := range list.Data {
				newSlice.Index(i).Set(ReflectValueFromNBT(value, fieldType.Type.Elem()))
			}
			field.Set(newSlice)
		} else {
			setValue := ReflectValueFromNBT(tag, fieldType.Type)
			if fieldType.Type != setValue.Type() {
				fmt.Printf("Type mismatch! Field %s is of type %s but NBT value is %s\n", fieldName, field.Type().Name(), setValue.Type())
			}
			field.Set(setValue)
		}
	}
}

// NBTStructToCompound converts an object into an NBT Compound
func NBTStructToCompound(obj interface{}) Compound {
	v := reflect.ValueOf(obj).Elem()
	t := reflect.TypeOf(obj).Elem()

	compoundMap := make(map[string]NBTValue)

	for fieldIndex := 0; fieldIndex < t.NumField(); fieldIndex++ {
		field := t.Field(fieldIndex)
		tag, exists := field.Tag.Lookup("nbt")
		if tag == "-" {
			continue
		}

		if !exists {
			tag = field.Name
		}

		nbtSettings := strings.Split(tag, ",")

		fieldValue := v.Field(fieldIndex)

		if field.Type.Kind() == reflect.Slice {
			outputValues := make([]NBTValue, fieldValue.Len())
			for i := range outputValues {
				outputValues[i] = NBTValueFromReflect(fieldValue.Index(i))
			}
			nbtList := NewList(outputValues)
			compoundMap[nbtSettings[0]] = &nbtList
		} else {
			compoundMap[nbtSettings[0]] = NBTValueFromReflect(fieldValue)
		}
	}

	return NewCompound(compoundMap)
}

// NBTValueFromReflect converts a reflected value into an NBT value
func NBTValueFromReflect(v reflect.Value) NBTValue {
	id := IDFromType(v.Type().Kind().String())

	nbtValue := NewValue(id)

	if id == 10 {
		newNbtValue := NBTStructToCompound(v.Addr().Interface())
		nbtValue = &newNbtValue
	} else {
		if v.Kind() == reflect.Slice {
			interfaceArray := v.Interface().([]interface{})
			nbtValueArray := make([]NBTValue, len(interfaceArray))
			for i, value := range interfaceArray {
				nbtValueArray[i] = value.(NBTValue)
			}
			nbtValue.SetValue(nbtValueArray)
		} else {
			nbtValue.SetValue(v.Interface())
		}

	}

	return nbtValue
}

// ReflectValueFromNBT converts an NBT value into a reflected value
func ReflectValueFromNBT(v NBTValue, fieldType reflect.Type) reflect.Value {
	if v.GetType() != 10 {
		return reflect.ValueOf(v.GetValue())
	}

	objectValue := reflect.New(fieldType)
	NBTStructScan(objectValue.Interface(), v.(*Compound))
	return objectValue.Elem()
}

// NBTMarshal converts an object into NBT bytes
func NBTMarshal(obj interface{}) []byte {
	compound := NBTStructToCompound(obj)
	return compound.Write()
}
