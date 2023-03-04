package helpers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/fundraising/rest-api/database"
)

type Database struct {
	Query          string
	Table          string
	Select         string
	Limit          int
	StartFrom      int
	OrderBy        string
	ObjectManifest interface{}
	Fields         []interface{}
	Where          string
	WhereParams    []interface{}
	Params         map[string]interface{}
}

type QueryBuilder struct {
	Database
}

var response ResponseFormatter
var objArr []interface{}

func FetchData(table string, model interface{}, args ...interface{}) (ResponseFormatter, error) {

	var objectRepresentationalData []interface{}

	con := database.CreateConnection()

	query := "SELECT * FROM " + table

	rows, err := con.Query(query)

	defer rows.Close()

	if err != nil {
		return response, err
	}

	for rows.Next() {
		err := rows.Scan(args...)

		if err != nil {
			return response, err
		}

		objectRepresentationalData = append(objectRepresentationalData, model)
	}

	response.Status = http.StatusOK
	response.Message = "success"
	response.Data = objectRepresentationalData

	return response, nil
}

func (qb *QueryBuilder) Query(statement string) *QueryBuilder {
	qb.Database.Query = statement
	return qb
}

func (qb *QueryBuilder) Where(field string, condition string, params interface{}) *QueryBuilder {

	if qb.Database.Where == "" {
		qb.Database.Where = " WHERE " + field + " " + condition + " ? "
	} else {
		qb.Database.Where += "AND " + field + " " + condition + " ? "
	}

	qb.Database.WhereParams = append(qb.Database.WhereParams, params)

	return qb
}

func (qb *QueryBuilder) WhereOR(col string, cond string, params interface{}) *QueryBuilder {
	if qb.Database.Where == "" {
		qb.Database.Where = " WHERE " + col + " " + cond + " ? "
	} else {
		qb.Database.Where += "OR " + col + " " + cond + " ? "
	}

	qb.Database.WhereParams = append(qb.Database.WhereParams, params)

	return qb
}

func (qb *QueryBuilder) Insert(params map[string]interface{}) (ResponseFormatter, error) {
	qb.Database.Params = params

	var query string
	var parameters []interface{}
	var counter = 0
	var index = 0
	var length = len(qb.Database.Params)
	var addedQuery = ""

	if qb.Database.Query != "" {
		query = qb.Database.Query
	} else {
		if qb.Database.Table != "" {
			if qb.Database.Params != nil {
				query = "INSERT INTO " + qb.Database.Table + "("
				for key, value := range qb.Database.Params {
					if counter < length-1 {
						if value != "" || value != nil {
							query += key + ","
							if counter == 0 {
								addedQuery += " VALUES (?,"
							} else {
								addedQuery += "?,"
							}
						}
					} else {
						if value != nil || value != "" {
							query += key + ")"
							addedQuery += "?)"
						}
					}

					if value != "" || value != nil {
						parameters = append(parameters, value)
						index++
					}
					counter++
				}

				query += addedQuery

			} else {
				fmt.Println("Params cannot set to be nil")
				return response, errors.New("Params cannot set to be nil")
			}
		} else {
			fmt.Println("Table cannot empty")
			return response, errors.New("Table cannot empty!")
		}
	}

	con := database.CreateConnection()
	statement, err := con.Prepare(query)

	if err != nil {
		fmt.Println("statement err")
		return response, err
	}

	result, err := statement.Exec(parameters...)

	fmt.Println("Inserting data ...")

	if err != nil {
		return response, err
	}

	lastInsertedID, _ := result.LastInsertId()

	fmt.Println("Done")

	response.Status = http.StatusOK
	response.Message = "success"
	response.Data = map[string]interface{}{
		"status":           true,
		"last_inserted_id": lastInsertedID,
		"created_at":       time.Now(),
	}

	return response, nil
}

func (qb *QueryBuilder) Get() (ResponseFormatter, error) {
	var query = ""
	var err error
	var rows *sql.Rows
	var counter = 0

	//set obj arr
	objArr = nil

	if qb.Database.Query == "" {
		if qb.Database.Table != "" {
			if qb.Database.Select != "" {
				query += "SELECT " + qb.Database.Select + " FROM " + qb.Database.Table

				if qb.Database.Where != "" {
					query += qb.Database.Where
				}

				if qb.Database.OrderBy != "" {
					query += qb.Database.OrderBy
				}

				if qb.Database.Limit != 0 {
					if qb.Database.StartFrom != 0 {
						query += " LIMIT " + strconv.Itoa(qb.Database.StartFrom) + "," + strconv.Itoa(qb.Database.Limit)
					} else {
						query += " LIMIT " + strconv.Itoa(qb.Database.Limit)
					}
				}
			} else {
				return response, errors.New("Select cannot set to be nil")
			}
		} else {
			return response, errors.New("Table cannot set to be nil")
		}
	} else {
		query += qb.Database.Query
	}

	con := database.CreateConnection()
	fmt.Println("Prepare connection pool")

	if qb.Database.Where != "" {
		rows, err = con.Query(query, qb.Database.WhereParams...)
	} else {
		rows, err = con.Query(query)
	}

	if err != nil {
		fmt.Println("err : ", err.Error())
		return response, err
	}

	for rows.Next() {
		errs := rows.Scan(qb.Database.Fields...)

		if errs != nil {
			fmt.Println("err : ", errs.Error())
			return response, errs
		}

		appendData(qb.Database.ObjectManifest)

		//logging
		fmt.Println("")
		fmt.Println("---------- Fetching data row :", strconv.Itoa(counter), " ----------")
		fmt.Println(qb.Database.ObjectManifest)
		fmt.Println("")

		counter++
	}

	response.Status = http.StatusOK
	response.Message = "success"
	response.Data = objArr

	return response, nil
}

func appendData(a interface{}) {
	rf := reflect.Indirect(reflect.ValueOf(a))
	if rf.Interface() != nil {
		objArr = append(objArr, rf.Interface())
	}
}

func (qb *QueryBuilder) Update(params map[string]interface{}) (ResponseFormatter, error) {
	qb.Database.Params = params
	var counter = 0
	var length = len(params)
	var query = ""
	var allParams []interface{}

	if qb.Database.Query == "" {
		if qb.Database.Table != "" {
			if qb.Database.Where != "" {
				query = "UPDATE " + qb.Database.Table + " SET "
				for key, value := range params {
					if key != "key" {
						if counter < length-1 {
							query += key + " = ?,"
						} else {
							query += key + " = ?"
						}
						allParams = append(allParams, value)
					}
					counter++
				}

				for _, v := range qb.Database.WhereParams {
					allParams = append(allParams, v)
				}

				query += qb.Database.Where
			} else {
				return response, errors.New("where statement cannot set to be nil")
			}
		} else {
			return response, errors.New("where statement cannot set to be nil")
		}
	} else {
		query = qb.Database.Query
	}

	con := database.CreateConnection()

	statement, err := con.Prepare(query)

	if err != nil {
		fmt.Println(err.Error())
		return response, err
	}

	result, err := statement.Exec(allParams...)

	fmt.Println("Update data ...")

	if err != nil {
		fmt.Println(err.Error())
		return response, err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		fmt.Println(err.Error())
		return response, err
	}

	fmt.Println("Done")

	response.Status = http.StatusOK
	response.Message = "success"
	response.Data = map[string]interface{}{
		"status":        true,
		"rows_affected": rowsAffected,
	}

	return response, nil
}

func (qb *QueryBuilder) Delete() (ResponseFormatter, error) {
	var query = ""

	if qb.Database.Query == "" {
		if qb.Database.Table != "" {
			if qb.Database.Where != "" {
				query = "DELETE FROM " + qb.Database.Table + qb.Database.Where
			} else {
				return response, errors.New("Where condition cannot set to be empty")
			}
		} else {
			return response, errors.New("Table cannot set to be empty")
		}
	} else {
		query = qb.Database.Query
	}

	con := database.CreateConnection()

	statement, err := con.Prepare(query)

	if err != nil {
		return response, err
	}

	result, err := statement.Exec(qb.Database.WhereParams...)

	if err != nil {
		return response, err
	}

	rowsAffected, _ := result.RowsAffected()

	response.Status = http.StatusOK
	response.Message = "success delete data"
	response.Data = map[string]interface{}{
		"status":        true,
		"rows_affected": rowsAffected,
		"deleted_at":    time.Now(),
	}
	return response, nil
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.Database.Limit = limit
	return qb
}

func (qb *QueryBuilder) StartFrom(start int) *QueryBuilder {
	qb.Database.StartFrom = start
	return qb
}

func (qb *QueryBuilder) Table(table string) *QueryBuilder {
	qb.Database.Table = table
	return qb
}

func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	var length = len(columns)

	for index, value := range columns {
		if length > 1 {
			if index < length-1 {
				qb.Database.Select += value + ","
			} else {
				qb.Database.Select += value
			}
		} else {
			qb.Database.Select += value
		}
	}
	return qb
}

func (qb *QueryBuilder) ObjRepresentations(obj interface{}) *QueryBuilder {
	qb.Database.ObjectManifest = obj
	qb.Database.Fields = scanObject(qb.Database.Select, qb.Database.ObjectManifest)
	return qb
}

func scanObject(query string, obj interface{}) []interface{} {
	var v []interface{}

	val := reflect.ValueOf(obj).Elem()
	tag := reflect.Indirect(reflect.ValueOf(obj)).Type()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		tagField := tag.Field(i).Tag.Get("json")
		isAnyField := false

		if query == "*" {
			isAnyField = true
		} else {
			if strings.Contains(query, tagField) {
				isAnyField = true
			}
		}

		if isAnyField {
			v = append(v, valueField.Addr().Interface())
		}

	}

	return v
}

func (qb *QueryBuilder) First() (ResponseFormatter, error) {
	var query = ""
	objArr = nil

	if qb.Database.Query == "" {
		if qb.Database.Table != "" {
			if qb.Database.Select != "" {
				if qb.Database.Where != "" {
					query += "SELECT " + qb.Database.Select + " FROM " + qb.Database.Table
					query += qb.Database.Where
				} else {
					return response, errors.New("Must set where condition!")
				}
			} else {
				return response, errors.New("Select cannot set to be nil!")
			}
		} else {
			return response, errors.New("Table cannot set to be nil!")
		}
	} else {
		query += qb.Database.Query
	}

	con := database.CreateConnection()

	err := con.QueryRow(query, qb.Database.WhereParams...).Scan(qb.Database.Fields...)

	if err == sql.ErrNoRows {
		response.Status = http.StatusOK
		response.Message = "success"
		response.Data = make([]interface{}, 0)

		return response, nil
	}

	if err != nil && err != sql.ErrNoRows {
		return response, err
	}

	response.Status = http.StatusOK
	response.Message = "success"
	response.Data = qb.Database.ObjectManifest

	return response, nil
}

func (qb *QueryBuilder) OrderBy(column string, ordertype string) *QueryBuilder {
	qb.Database.OrderBy += " ORDER BY " + column + " " + ordertype
	return qb
}
