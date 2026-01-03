package main

import (
	"fmt"
	"project/Model"
	"regexp"
	"strings"
	"unicode"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

// ===================== FK INFO =====================

// Foreign key info for PostgreSQL
type foreignKeyInfo struct {
	TableName        string `gorm:"column:tablename"`
	ColumnName       string `gorm:"column:columnname"`
	ReferencedTable  string `gorm:"column:referencedtable"`
	ReferencedColumn string `gorm:"column:referencedcolumn"`
}

// Get all FK relations (PostgreSQL compatible)
func getForeignKeys(db *gorm.DB) ([]foreignKeyInfo, error) {
	var results []foreignKeyInfo

	query := `
		SELECT 
			fk_tab.relname AS tablename,
			fk_col.attname AS columnname,
			pk_tab.relname AS referencedtable,
			pk_col.attname AS referencedcolumn
		FROM pg_constraint AS con
		JOIN pg_class AS fk_tab ON fk_tab.oid = con.conrelid
		JOIN pg_attribute AS fk_col 
			ON fk_col.attnum = ANY(con.conkey) 
			AND fk_col.attrelid = con.conrelid
		JOIN pg_class AS pk_tab ON pk_tab.oid = con.confrelid
		JOIN pg_attribute AS pk_col 
			ON pk_col.attnum = ANY(con.confkey) 
			AND pk_col.attrelid = con.confrelid
		WHERE con.contype = 'f';
	`

	err := db.Raw(query).Scan(&results).Error
	return results, err
}

// ===================== PK INFO (NEW) =====================

// NEW: Primary key info
type primaryKeyInfo struct {
	TableName  string `gorm:"column:table_name"`
	ColumnName string `gorm:"column:column_name"`
}

// NEW: ดึง PK ของทุก table
func getPrimaryKeys(db *gorm.DB) ([]primaryKeyInfo, error) {
	var keys []primaryKeyInfo

	query := `
		SELECT 
			kcu.table_name,
			kcu.column_name
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu 
			ON tc.constraint_name = kcu.constraint_name
			AND tc.table_schema = kcu.table_schema
		WHERE tc.constraint_type = 'PRIMARY KEY'
		  AND tc.table_schema = 'public';
	`

	err := db.Raw(query).Scan(&keys).Error
	return keys, err
}

// ===================== COLUMN INFO + LENGTH =====================

// column info สำหรับ gen validate tag ตาม schema
type columnInfo struct {
	TableName  string `gorm:"column:table_name"`
	ColumnName string `gorm:"column:column_name"`
	DataType   string `gorm:"column:data_type"`
	MaxLength  *int32 `gorm:"column:character_maximum_length"`
	IsNullable string `gorm:"column:is_nullable"`
}

// ดึงข้อมูลคอลัมน์จาก information_schema.columns
func getColumnInfos(db *gorm.DB) ([]columnInfo, error) {
	var cols []columnInfo

	query := `
		SELECT 
			table_name,
			column_name,
			data_type,
			character_maximum_length,
			is_nullable
		FROM information_schema.columns
		WHERE table_schema = 'public';
	`

	err := db.Raw(query).Scan(&cols).Error
	return cols, err
}

// แปลง []columnInfo -> map[table][column]columnInfo
func buildColumnMap(cols []columnInfo) map[string]map[string]columnInfo {
	result := make(map[string]map[string]columnInfo)

	for _, c := range cols {
		t := c.TableName
		if _, ok := result[t]; !ok {
			result[t] = make(map[string]columnInfo)
		}
		result[t][c.ColumnName] = c
	}

	return result
}

// ===================== FIX AUTOINCREMENT (NEW) =====================

// Fix the malformed nextval default tag for primary key columns
// GORM gen copies PostgreSQL's nextval('sequence') which breaks the struct tag
func buildAutoIncrementFixOpts(table string, pkMap map[string]map[string]bool) []gen.ModelOpt {
	var opts []gen.ModelOpt

	pks, ok := pkMap[table]
	if !ok {
		return opts
	}

	for colName := range pks {
		// For each primary key column, override the GORM tag to use autoIncrement:true
		// instead of the broken default:nextval(...) that GORM gen creates
		opts = append(opts, gen.FieldGORMTag(colName, func(tag field.GormTag) field.GormTag {
			// Remove the broken "default" tag if it contains nextval
			tag.Remove("default")
			// Ensure autoIncrement is set
			tag.Set("autoIncrement", "true")
			return tag
		}))
	}

	return opts
}

// ===================== VALIDATE OPTS =====================

// สร้าง gen.ModelOpt สำหรับ validate tag จาก columnInfo + pkMap
// PK ทุกตัวจะได้ validate:"" (ไม่ required)
func buildValidateOpts(table string, colMap map[string]map[string]columnInfo, pkMap map[string]map[string]bool) []gen.ModelOpt {
	var opts []gen.ModelOpt

	cols, ok := colMap[table]
	if !ok {
		return opts
	}

	for colName, info := range cols {

		// ถ้าเป็น Primary Key → ให้เป็น validate:"" (ไม่บังคับ)
		if pkMap[table][colName] {
			tag := field.Tag{
				"validate": "",
			}
			opts = append(opts, gen.FieldNewTag(colName, tag))
			continue
		}

		var tags []string

		// NOT NULL -> required
		if strings.EqualFold(info.IsNullable, "NO") {
			tags = append(tags, "required")
		}

		// varchar/char ที่มี length -> max=N
		if info.MaxLength != nil &&
			(strings.Contains(info.DataType, "character varying") ||
				strings.Contains(info.DataType, "character") ||
				strings.Contains(info.DataType, "varchar") ||
				strings.Contains(info.DataType, "char")) {

			tags = append(tags, fmt.Sprintf("max=%d", *info.MaxLength))
		}

		if len(tags) == 0 {
			continue
		}

		tag := field.Tag{
			"validate": strings.Join(tags, ","), // เช่น required,max=10
		}

		// ตรงนี้ใช้ชื่อ column (snake_case) ให้ gen ไปแมปกับ field ให้
		opts = append(opts, gen.FieldNewTag(colName, tag))
	}

	return opts
}

// ===================== MAIN =====================

func main() {
	var adapter *Model.Adapter
	db := adapter.GetGormIntance()

	// Configure generator
	g := gen.NewGenerator(gen.Config{
		OutPath:       "Technical_Service/Entity",
		ModelPkgPath:  "Entity/EntityStruct",
		Mode:          gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	g.WithJSONTagNameStrategy(toCamelJSON)
	g.UseDB(db)

	// Step 1: Load tables
	allTables, err := db.Migrator().GetTables()
	if err != nil {
		panic(fmt.Errorf("❌ Failed to list tables: %w", err))
	}

	// Step 1.5: Load column infos (สำหรับ validate tag)
	colInfos, err := getColumnInfos(db)
	if err != nil {
		fmt.Println("⚠️ Warning: failed to load column infos:", err)
	}
	colMap := buildColumnMap(colInfos)

	// NEW: Load primary keys
	pkInfos, err := getPrimaryKeys(db)
	if err != nil {
		fmt.Println("⚠️ Warning: failed to load primary keys:", err)
	}
	// สร้าง pkMap[table][column] = true
	pkMap := make(map[string]map[string]bool)
	for _, pk := range pkInfos {
		if pkMap[pk.TableName] == nil {
			pkMap[pk.TableName] = make(map[string]bool)
		}
		pkMap[pk.TableName][pk.ColumnName] = true
	}

	// Step 3: Get foreign keys
	fks, err := getForeignKeys(db)
	if err != nil {
		fmt.Println("⚠️ Warning: failed to load foreign keys:", err)
	}
	if len(fks) == 0 {
		fmt.Println("⚠️ No foreign keys detected. Check your schema or permissions.")
	}

	// ---------- นับจำนวน FK ต่อ (childTable, parentTable) ----------
	type fkPair struct {
		Child  string
		Parent string
	}
	fkCount := make(map[fkPair]int)
	for _, fk := range fks {
		key := fkPair{Child: fk.TableName, Parent: fk.ReferencedTable}
		fkCount[key]++
	}

	// Step 4: Group relation fields
	relations := map[string][]gen.ModelOpt{}

	// กันไม่ให้สร้าง field ซ้ำชื่อกันใน struct เดียวกัน
	addedChildField := make(map[string]map[string]bool)  // childTable -> fieldName
	addedParentField := make(map[string]map[string]bool) // parentTable -> fieldName

	for _, fk := range fks {
		parent := g.GenerateModelAs(fk.ReferencedTable, toPascal(fk.ReferencedTable))
		child := g.GenerateModelAs(fk.TableName, toPascal(fk.TableName))

		pair := fkPair{Child: fk.TableName, Parent: fk.ReferencedTable}
		count := fkCount[pair]

		var childFieldName string
		var parentFieldName string

		// ----- ตั้งชื่อ relation -----
		if fk.TableName == fk.ReferencedTable {
			// Self reference → rename both ends
			fmt.Printf("🔁 Self reference detected: %s(%s → %s)\n", fk.TableName, fk.ColumnName, fk.ReferencedColumn)
			childFieldName = "Parent" + toPascal(fk.ReferencedTable)
			parentFieldName = "Sub" + toPascal(fk.ReferencedTable) + "s"
		} else {
			// child (BelongsTo)
			if count == 1 {
				// เคสทั่วไป: 1 FK จาก child -> parent → ใช้ชื่อแบบเดิม (Module, Employee...)
				childFieldName = toPascal(fk.ReferencedTable)
			} else {
				// ถ้ามีหลาย FK ไป parent เดียวกัน → ใช้ role name จากชื่อ column
				// ex: assigned_checker -> AssignedChecker, current_owner -> CurrentOwner
				roleName := buildRoleName(fk.ColumnName)   // AssignedChecker
				parentName := toPascal(fk.ReferencedTable) // Employee
				childFieldName = roleName + parentName
			}

			// parent (HasMany) ใช้ชื่อ table ลูกแบบเดิม เช่น Approverequestforms
			parentFieldName = toPascal(fk.TableName)
		}

		// ---------- BelongsTo (child -> parent) ----------
		if addedChildField[fk.TableName] == nil {
			addedChildField[fk.TableName] = make(map[string]bool)
		}
		if !addedChildField[fk.TableName][childFieldName] {
			addedChildField[fk.TableName][childFieldName] = true

			childRel := gen.FieldRelate(
				field.BelongsTo,
				childFieldName,
				parent,
				&field.RelateConfig{
					RelatePointer: true,
					JSONTag:       toCamelJSON(childFieldName),
					GORMTag: field.GormTag{}.
						Set("foreignKey", fk.ColumnName).
						Set("references", fk.ReferencedColumn),
				},
			)

			relations[fk.TableName] = append(relations[fk.TableName], childRel)
		}

		// ---------- HasMany (parent -> children) ----------
		if addedParentField[fk.ReferencedTable] == nil {
			addedParentField[fk.ReferencedTable] = make(map[string]bool)
		}
		if !addedParentField[fk.ReferencedTable][parentFieldName] {
			addedParentField[fk.ReferencedTable][parentFieldName] = true

			parentRel := gen.FieldRelate(
				field.HasMany,
				parentFieldName,
				child,
				&field.RelateConfig{
					RelateSlicePointer: true,
					JSONTag:            toCamelJSON(parentFieldName),
					GORMTag: field.GormTag{}.
						Set("foreignKey", fk.ColumnName).
						Set("references", fk.ReferencedColumn),
				},
			)

			relations[fk.ReferencedTable] = append(relations[fk.ReferencedTable], parentRel)
		}
	}

	// Step 5: Rebuild models with relations + validate tags + autoIncrement fix
	var all []any
	for _, table := range allTables {

		structName := toPascal(table)
		if table == "AssetItemInventoryGoods" {
			structName = "AssetItemInventoryGoods"
		}

		rels := relations[table]
		validateOpts := buildValidateOpts(table, colMap, pkMap)
		autoIncrementOpts := buildAutoIncrementFixOpts(table, pkMap) // NEW: Fix nextval issue

		opts := make([]gen.ModelOpt, 0)

		if len(rels) > 0 {
			opts = append(opts, rels...)
		}
		if len(validateOpts) > 0 {
			opts = append(opts, validateOpts...)
		}
		if len(autoIncrementOpts) > 0 {
			opts = append(opts, autoIncrementOpts...)
		}

		model := g.GenerateModelAs(table, structName, opts...)
		all = append(all, model)
	}

	// Step 6: Execute
	g.ApplyBasic(all...)
	g.Execute()

	fmt.Println("✅ Scaffold generation complete with PostgreSQL relations + validate tags (PK not required)!")
}

// ===================== JSON CAMEL CASE =====================
func toPascal(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	// ถ้ามี underscore → ถือว่าเป็น snake_case: menu_category → MenuCategory
	if strings.Contains(s, "_") {
		parts := strings.Split(s, "_")
		for i := range parts {
			if parts[i] == "" {
				continue
			}
			runes := []rune(parts[i])
			runes[0] = unicode.ToUpper(runes[0])
			for j := 1; j < len(runes); j++ {
				runes[j] = unicode.ToLower(runes[j])
			}
			parts[i] = string(runes)
		}
		return strings.Join(parts, "")
	}

	// ถ้าไม่มี underscore → ถือว่าตัวเดิมอาจจะเป็น camelCase/PascalCase อยู่แล้ว
	// แค่บังคับตัวแรกให้เป็นพิมพ์ใหญ่ แล้ว "ไม่ไปยุ่งตัวที่เหลือ"
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// เอาชื่อ column มาเป็นชื่อ role
// ex: assigned_checker_id -> AssignedChecker
//
//	current_owner       -> CurrentOwner
func buildRoleName(col string) string {
	c := strings.ToLower(col)
	if strings.HasSuffix(c, "_id") {
		col = col[:len(col)-3] // ตัด "_id"
	}
	return toPascal(col)
}

func toCamelJSON(s string) string {
	if s == "" {
		return s
	}

	// 1. Replace snake_case with parts
	s = strings.ReplaceAll(s, "_", " ")

	// 2. Split words by uppercase boundaries (Go-safe regex)
	re := regexp.MustCompile(`[A-Z][a-z]*|[a-z]+|[A-Z]+`)
	parts := re.FindAllString(s, -1)

	if len(parts) == 0 {
		return s
	}

	// 3. Build camelCase
	for i := range parts {
		parts[i] = strings.ToLower(parts[i])

		if i > 0 {
			parts[i] = strings.Title(parts[i])
		}
	}

	return strings.Join(parts, "")
}
