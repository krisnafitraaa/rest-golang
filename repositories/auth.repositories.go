package repositories

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/fundraising/rest-api/database"
	"github.com/fundraising/rest-api/helpers"
	"github.com/fundraising/rest-api/models"
)

func RegisterNewUser(params map[string]interface{}) (helpers.ResponseFormatter, error) {
	var response helpers.ResponseFormatter
	var listData []interface{}
	var lastInsertedIdVolunteer interface{}
	// var objectRepresentation models.User

	con := database.CreateConnection()

	counter := 0

	index := 0

	length := len(params)

	// queryCheck := "SELECT email FROM users WHERE email = ?"

	query := "INSERT INTO volunteers (volunteer_id, volunteer_name) VALUES (?,?)"

	query2 := "INSERT INTO users ("

	addedQuery2 := ""

	// query3 := "INSERT INTO user_activations (user_id, activation_token, valid_until,ip_address,user_agent) VALUES (?,?,?,?,?)"

	// v := validator.New()

	// users := models.User{
	// 	UserID:      params["user_id"].(string),
	// 	RoleID:      params["role_id"].(string),
	// 	VolunteerID: params["volunteer_id"].(string),
	// 	Fullname:    params["fullname"].(string),
	// 	Email:       params["email"].(string),
	// 	Password:    params["password"].(string),
	// 	Provider:    params["provider"].(string),
	// }

	// err := v.Struct(users)

	// if err != nil {
	// 	return response, err
	// }

	for key, value := range params {
		if counter < length-1 {
			if value != nil || value != "" {
				query2 += key + ","
				if counter == 0 {
					addedQuery2 += " VALUES (?,"
				} else {
					addedQuery2 += "?,"
				}
			}
		} else {
			if value != nil || value != "" {
				query2 += key + ")"
				addedQuery2 += "?)"
			}
		}

		if value != nil || value != "" {
			listData = append(listData, value)
			index++
		}

		counter++
	}

	query2 = query2 + addedQuery2

	//check existing email in database
	// err = con.QueryRow(queryCheck, params["email"]).Scan(&objectRepresentation.Email)

	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// if objectRepresentation.Email != "" {
	// 	response.Status = http.StatusOK
	// 	response.Message = "email is already used"
	// 	response.Data = map[string]bool{
	// 		"status": false,
	// 	}

	// 	return response, nil
	// }

	println("First section")

	trans, err := con.Begin()

	if err != nil {
		return response, err
	}

	if params["role_id"] == "user_role" {

		stmt, err := trans.Prepare(query)

		println("Second section")

		if err != nil {
			trans.Rollback()
			return response, err
		}

		_, err = stmt.Exec(params["volunteer_id"], params["fullname"])

		println("Third section")

		if err != nil {
			trans.Rollback()
			return response, err
		}

		lastInsertedIdVolunteer = params["volunteer_id"]
	}

	// statement, err := con.Prepare(query2)

	// println("Fourth section")

	// if err != nil {
	// 	trans.Rollback()
	// 	return response, err
	// }

	// println("Fifth section")

	// println(query2)

	// _, e := statement.Exec(listData...)

	// println("Six section")

	// if e != nil {
	// 	trans.Rollback()
	// 	return response, err
	// }

	println("Middle section")

	lastInsertedId := params["user_id"]

	//insert user activation
	// token := helpers.MakeRandomString(32)
	// valid_until := time.Now().Add(time.Hour * 24).Unix()

	// statement2, err := con.Prepare(query3)

	// if err != nil {
	// 	trans.Rollback()
	// 	return response, err
	// }

	// _, err = statement2.Exec(params["user_id"], token, valid_until, "127.0.0.1", "Mozilla")

	// if err != nil {
	// 	trans.Rollback()
	// 	return response, err
	// }

	trans.Commit()

	println("Last section")

	response.Status = http.StatusOK
	response.Message = "success create new user"
	response.Data = map[string]interface{}{
		"last_inserted_id":        lastInsertedId,
		"last_inserted_volunteer": lastInsertedIdVolunteer,
		"activation_token":        "token",
	}

	return response, nil
}

func CheckLogin(params map[string]string) (helpers.ResponseFormatter, error) {
	var response helpers.ResponseFormatter
	var objectRepresentation models.User

	email := params["email"]
	password := params["password"]

	con := database.CreateConnection()

	statement := "SELECT u.fullname,u.email,u.password,u.role_id,u.volunteer_id FROM users u LEFT JOIN volunteers v ON u.volunteer_id=v.volunteer_id WHERE u.email = ?"

	err := con.QueryRow(statement, email).Scan(&objectRepresentation.Fullname, &objectRepresentation.Email, &objectRepresentation.Password, &objectRepresentation.RoleID, &objectRepresentation.VolunteerID, &objectRepresentation.IsActive)

	if err == sql.ErrNoRows {
		response.Status = http.StatusOK
		response.Message = "user not found"
		response.Data = map[string]bool{
			"status": false,
		}

		return response, err
	}

	if err != nil {
		return response, err
	}

	isPasswordValid := helpers.VerifyPassword(password, objectRepresentation.Password)

	if !isPasswordValid {
		//insert to login attempts

		response.Status = http.StatusOK
		response.Message = "wrong password"
		response.Data = map[string]bool{
			"status": false,
		}

		return response, errors.New("wrong password")
	}

	//user apakah aktif?
	if objectRepresentation.IsActive == 0 {
		response.Status = http.StatusOK
		response.Message = "user not active"
		response.Data = map[string]bool{
			"status": false,
		}

		return response, errors.New("user not active")

	} else if objectRepresentation.IsActive == 2 {
		response.Status = http.StatusOK
		response.Message = "user has been blocked by admin"
		response.Data = map[string]bool{
			"status": false,
		}

		return response, errors.New("user has been blocked by admin")
	}

	response.Status = http.StatusOK
	response.Message = "success login"
	response.Data = map[string]bool{
		"status": true,
	}

	return response, nil
}
