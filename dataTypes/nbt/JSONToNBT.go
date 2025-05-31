package nbt

import (
	"encoding/json"
	"errors"
	"reflect"
)

func JSONToNBT(jsonString []byte) (Compound, error) {

	var jsonData map[string]interface{}
	err := json.Unmarshal(jsonString, &jsonData)
	if err != nil {
		return Compound{}, errors.Join(errors.New("failed to unmarshal json for nbt"), err)
	}

	return mapStringInterfaceToNBT(jsonData), nil
}

func mapStringInterfaceToNBT(input map[string]interface{}) Compound {
	outputCompound := Compound{
		Data: make(map[string]NBTValue),
	}
	for key, val := range input {

		reflectValue := reflect.ValueOf(val)
		reflectKind := reflectValue.Type().Kind()

		if reflectKind == reflect.Map {
			compoundValue := mapStringInterfaceToNBT(reflectValue.Interface().(map[string]interface{}))
			outputCompound.Data[key] = &compoundValue
		} else if reflectKind == reflect.Slice {
			outputCompound.Data[key] = arrayInterfaceToNBT(val.([]interface{}))
		} else {
			outputCompound.Data[key] = NBTValueFromReflect(reflectValue)
		}
	}

	return outputCompound
}

func arrayInterfaceToNBT(interfaceArray []interface{}) *List {
	list := List{
		Data: make([]NBTValue, len(interfaceArray)),
	}

	if len(interfaceArray) == 0 {
		return &list
	}

	isMixedTypeList := false
	for _, val := range interfaceArray {
		elemKind := reflect.TypeOf(val).Kind()
		if list.Type == 0 {
			list.Type = IDFromType(elemKind.String())
			continue
		}

		typeId := IDFromType(elemKind.String())
		if typeId != list.Type {
			isMixedTypeList = true
			list.Type = 10
			break
		}

		list.Type = IDFromType(elemKind.String())
	}

	for i, val := range interfaceArray {
		elemKind := reflect.TypeOf(interfaceArray[i]).Kind()
		if elemKind == reflect.Map {
			mapValue := mapStringInterfaceToNBT(val.(map[string]interface{}))
			list.Data[i] = &mapValue
			continue
		} else if elemKind == reflect.Slice || elemKind == reflect.Array {
			list.Data[i] = arrayInterfaceToNBT(val.([]interface{}))
		} else {
			list.Data[i] = NBTValueFromReflect(reflect.ValueOf(val))
		}

		if isMixedTypeList {
			compoundWrapper := NewCompound(map[string]NBTValue{
				"": list.Data[i],
			})
			list.Data[i] = &compoundWrapper
		}
	}

	return &list
}
