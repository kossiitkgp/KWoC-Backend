rows, _ := db.Query("SELECT * FROM Student;")

    err := sqltocsv.WriteFile("student.csv", rows)
    if err != nil {
        panic(err)
    }

    columns, _ := rows.Columns()
    count := len(columns)
    values := make([]interface{}, count)
    valuePtrs := make([]interface{}, count)

    for rows.Next() {
        for i := range columns {
            valuePtrs[i] = &values[i]
        }

        rows.Scan(valuePtrs...)

        for i, col := range columns {
            val := values[i]

            b, ok := val.([]byte)
            var v interface{}
            if ok {
                v = string(b)
            } else {
                v = val
            }

            fmt.Println(col, v)
        }
    }
}
