package repositories

import (
	"github.com/fundraising/rest-api/helpers"
	"github.com/fundraising/rest-api/models"
)

var response helpers.ResponseFormatter

func GetRole() (helpers.ResponseFormatter, error) {

	model := &models.Role{}

	builder := &helpers.QueryBuilder{}

	result, err := builder.Table("business_roles").Select("role_id", "role").ObjRepresentations(model).Get()

	if err != nil {
		return result, err
	}

	return result, nil
}

func RoleDetail(role_id string) (helpers.ResponseFormatter, error) {

	model := &models.Role{}

	builder := &helpers.QueryBuilder{}

	response, err := builder.Table("business_roles").Select("*").Where("role_id", "=", role_id).ObjRepresentations(model).First()

	if err != nil {
		return response, err
	}

	return response, nil
}

func StoreRole(params map[string]interface{}) (helpers.ResponseFormatter, error) {

	builder := &helpers.QueryBuilder{}

	result, err := builder.Table("business_roles").Insert(params)

	if err != nil {
		return result, err
	}

	return result, nil
}

func UpdateRole(params map[string]interface{}) (helpers.ResponseFormatter, error) {
	builder := &helpers.QueryBuilder{}

	response, err := builder.Table("business_roles").Where("role_id", "=", params["key"]).Update(params)

	if err != nil {
		return response, err
	}

	return response, nil
}

func DeleteRole(role_id string) (helpers.ResponseFormatter, error) {
	builder := &helpers.QueryBuilder{}

	response, err := builder.Table("business_roles").Where("role_id", "=", role_id).Delete()

	if err != nil {
		return response, err
	}

	return response, nil
}
