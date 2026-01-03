package Technical_Service

import (
	"reflect"
	"time"

	// "project/Technical_Service/Entity/EntityStruct"

	config "project/Config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var globalAdapterInstance *gorm.DB

type Adapter struct {
}

func (adapter *Adapter) newAdapter() {}
func ParseIntNil(number int32) *int32 {
	return &number
}
func ParseInt(number *int32) int32 {
	if number == nil {
		return 0
	}
	return *number
}
func ParseStringNil(text string) *string {
	return &text
}
func ParseString(text *string) string {
	if text == nil {
		return ""
	}
	return *text
}
func ParseBoolNil(bo bool) *bool {
	return &bo
}
func ParseBool(bo *bool) bool {
	if bo == nil {
		return false
	}
	return *bo
}
func ParseTimeNil(t time.Time) *time.Time {
	return &t
}
func IntsToIntPtrs(src []int) []*int {
	out := make([]*int, len(src))
	for i, v := range src {
		vv := v
		out[i] = &vv
	}
	return out
}
func IsZeroOrNull(number *int32) bool {
	if number == nil {
		return true
	} else if *number == 0 {
		return true
	}
	return false
}
func IsStringEmptyOrNull(number *string) bool {
	if number == nil {
		return true
	} else if *number == "" {
		return true
	}
	return false
}
func (adapter *Adapter) GetAdapterIntance() *Adapter {
	if globalAdapterInstance == nil {
		cfg := config.LoadConfig()

		db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
		if err != nil {
			panic("failed to connect database: " + err.Error())
		}

		globalAdapterInstance = db
		adapter = new(Adapter)
		adapter.newAdapter()
	}
	return adapter
}

// func (adapter *Adapter) GetDictionnaryByName(serviceName string, systemName string) ([]*EntityStruct.DataMappingDictionary, error) {
// 	var err error
// 	var dmc []*EntityStruct.DataMappingDictionary
// 	var tx *gorm.DB

// 	tx = globalAdapterInstance
// 	dmc = make([]*EntityStruct.DataMappingDictionary, 0)
// 	err = nil

// 	if err = tx.Where("service_name = ? AND system_name = ?", serviceName, systemName).Find(&dmc).Error; err != nil {
// 		return nil, err
// 	}

// 	//

// 	return dmc, nil
// }

func CheckDuplicateDynamic[T any](list []T, fields []string) (bool, []map[string]*string) {
	seenMap := make(map[string]map[any]int)
	results := make([]map[string]*string, len(list))
	anyDup := false

	for _, f := range fields {
		seenMap[f] = make(map[any]int)
	}

	for i, item := range list {
		row := make(map[string]*string)
		v := reflect.ValueOf(item)

		if v.Kind() == reflect.Pointer {
			v = v.Elem()
		}

		for _, field := range fields {

			fv := v.FieldByName(field)
			if !fv.IsValid() {
				row[field] = ParseStringNil("invalid_field")
				continue
			}

			var value interface{}

			// ⭐⭐ FIX สำคัญ — ดึงค่าจาก pointer ถ้ามี ⭐⭐
			if fv.Kind() == reflect.Pointer && !fv.IsNil() {
				value = fv.Elem().Interface()
			} else {
				value = fv.Interface()
			}

			seenMap[field][value]++

			if seenMap[field][value] > 1 {
				row[field] = ParseStringNil("ข้อมูลซ้ำ กรุณาตรวจสอบ")
				anyDup = true
			} else {
				row[field] = nil
			}
		}

		results[i] = row
	}

	return anyDup, results
}
func (adapter *Adapter) IsAccessable(roleID int32, controllername string) (bool, error) {
	// Declartion
	var exists bool

	// Definition
	exists = false

	err := globalAdapterInstance.
		Table(`"Controller" as c`).
		Joins(`JOIN "Module" m ON m.Module_ID = c.Module_ID`).
		Joins(`JOIN "ModuleAccessible" ma ON ma.Module_ID = m.Module_ID`).
		Where("ma.Role_ID = ?", roleID).
		Where("ma.Accessable = ?", true).
		Where("c.Controller_Name = ?", controllername).
		Select("1").
		Scan(&exists).Error

	if err != nil {
		return false, err
	}

	return exists, nil
}

// ข้อมูลตัวอย่าง
// var mockDictionary = []*EntityStruct.DataMappingDictionary{
// 	{
// 		ID:            1,
// 		InCommingName: "Coltyp",
// 		InternalName:  "AssetTypeId",
// 		ServiceName:   "LnColl",
// 		SystemName:    "DBD",
// 	},
// 	{
// 		ID:            2,
// 		InCommingName: "Coltypdesc",
// 		InternalName:  "AssetTypeName",
// 		ServiceName:   "LnColl",
// 		SystemName:    "DBD",
// 	},
// }

// // ฟังก์ชัน mock dictionary
// func (adapter *Adapter) GetDictionnaryByNameMock(serviceName string, systemName string) ([]*EntityStruct.DataMappingDictionary, error) {
// 	res := make([]*EntityStruct.DataMappingDictionary, 0)
// 	for _, d := range mockDictionary {
// 		if d.ServiceName == serviceName && d.SystemName == systemName {
// 			res = append(res, d)
// 		}
// 	}
// 	return res, nil
// }
