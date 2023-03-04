package repositories

import (
	"net/http"

	"github.com/fundraising/rest-api/database"
	"github.com/fundraising/rest-api/helpers"
	"github.com/fundraising/rest-api/models"
)

var menu models.Menu
var submenu models.Submenu
var menuLists []models.Menu
var res helpers.ResponseFormatter
var submenuLists []models.Submenu
var err error

func FetchAllMenus() (helpers.ResponseFormatter, error) {
	menuId := 0
	counter := 0
	index := 0

	menuLists = nil
	submenuLists = nil

	con := database.CreateConnection()
	rows, err := con.Query("SELECT IFNULL(am.menu_id,0),am.menu,am.menu_for,am.icon,am.position,am.menu_url,am.is_active,am.created_at,IFNULL(asm.submenu_id,0),IFNULL(asm.submenu, ''),IFNULL(asm.submenu_url, ''),IFNULL(asm.icon,''), IFNULL(asm.position, 0), IFNULL(asm.is_active,0) FROM app_menus am LEFT JOIN app_submenus asm ON am.menu_id=asm.menu_id")
	defer rows.Close()

	if err != nil {
		return res, err
	}

	for rows.Next() {
		err = rows.Scan(&menu.MenuID, &menu.Menu, &menu.MenuFor, &menu.Icon, &menu.Position, &menu.MenuURL,
			&menu.IsActive, &menu.CreatedAt, &submenu.SubmenuID, &submenu.Submenu, &submenu.SubmenuURL, &submenu.Icon, &submenu.Position, &submenu.IsActive)

		if err != nil {
			return res, err
		}

		createNewSubmenu(menu, submenu)

		if menuId != menu.MenuID {
			if counter > 0 {
				submenuLists = nil
				createNewSubmenu(menu, submenu)
				index++
			}
			menuLists = append(menuLists, menu)
		} else {
			menuLists[index].Submenus = submenuLists
		}

		menuId = menu.MenuID
		counter++
	}

	res.Status = http.StatusOK
	res.Message = "success"
	res.Data = menuLists

	return res, nil
}

func createNewSubmenu(m models.Menu, sm models.Submenu) {
	submenuLists = append(submenuLists, sm)
	m.Submenus = submenuLists
}

func StoreMenu(params map[string]interface{}) (helpers.ResponseFormatter, error) {
	builder := &helpers.QueryBuilder{}

	res, err = builder.Table("app_menus").Insert(params)

	if err != nil {
		return res, err
	}

	return res, nil
}

func GetMenu() (helpers.ResponseFormatter, error) {
	var object models.Menu

	res, err = helpers.FetchData("app_menus", &object, &object.MenuID, &object.Menu, &object.Icon, &object.Position, &object.MenuFor, &object.MenuURL,
		&object.IsActive, &object.CreatedAt, &object.UpdatedAt)

	if err != nil {
		return res, err
	}

	return res, nil
}
