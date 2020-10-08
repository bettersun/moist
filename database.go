package moist

import (
	"database/sql"
	"fmt"
)

// SQL查询结果放到Map中
// 参考（引用）：
//   https://www.cnblogs.com/mafeng/p/6207702.html
//   https://stackoverflow.com/questions/19991541/dumping-mysql-tables-to-json-with-golang
func GetDataMap(db *sql.DB, sSql string) ([]map[string]interface{}, error) {

	// 准备查询语句
	stmt, err := db.Prepare(sSql)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// 查询
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// 结果列（包含表以外的列）
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	count := len(columns)

	// 返回值 Map的List
	tableData := make([]map[string]interface{}, 0)
	// 值
	values := make([]interface{}, count)
	// 列名
	valuePtrs := make([]interface{}, count)
	for rows.Next() {

		// Rows的所有项和值
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		// 表中列的对应项和值
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}

	return tableData, nil
}

