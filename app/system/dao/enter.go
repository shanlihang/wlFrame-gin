package dao

import (
	"WlFrame-gin/app/system/model"
	"WlFrame-gin/utils/authentication"
	"WlFrame-gin/utils/global"
	"fmt"
	"gorm.io/gorm"
)

// 新增用户
func InsertUser(user model.SysUser, roleIDs []int64) *gorm.DB {
	user.Birthday = user.Birthday[:10]
	result := global.DB.Model(model.SysUser{}).Create(&user)
	// 创建用户角色关联
	for _, roleID := range roleIDs {
		userRole := model.RelateUserRole{
			SysUserID: user.ID,
			SysRoleID: uint(roleID),
		}
		global.DB.Model(model.RelateUserRole{}).Create(&userRole)
	}

	return result
}

// 查询用户列表
func SelectUserList(name string, phone string, email string, roles []string) ([]model.SysUser, error) {
	var users []model.SysUser
	tx := global.DB.Model(model.SysUser{}).Omit("Password")
	if name != "" {
		tx.Where("name Like ?", "%"+name+"%")
	}
	if phone != "" {
		tx.Where("phone Like ?", "%"+phone+"%")
	}
	if email != "" {
		tx.Where("email Like ?", "%"+email+"%")
	}
	if len(roles) != 0 {
		var sysUserIds []int64
		global.DB.Table("relate_user_role").
			Where("sys_role_id in (?)", roles).
			Select("sys_user_id").Scan(&sysUserIds)
		tx.Where("id in (?)", sysUserIds)
	}
	err := tx.Preload("Roles").Find(&users).Error
	return users, err
}

// 根据id查询用户
func SelectUserById(id int64) (model.SysUser, *gorm.DB) {
	var user model.SysUser
	result := global.DB.Where("ID = ?", id).Omit("Password").First(&user)
	return user, result
}

// 修改用户
func UpdateUser(user model.SysUser, roles []int64) *gorm.DB {
	result := global.DB.Save(&user)
	//删除之前用户的所有角色（后面覆盖）
	global.DB.Unscoped().Where("sys_user_id = ?", user.ID).Delete(&model.RelateUserRole{})
	//添加新的角色（覆盖）
	for _, roleID := range roles {
		userRole := model.RelateUserRole{
			SysUserID: user.ID,
			SysRoleID: uint(roleID),
		}
		global.DB.Model(model.RelateUserRole{}).Create(&userRole)
	}
	return result
}

// 删除用户
func DeleteUser(id int64) *gorm.DB {
	result := global.DB.Delete(&model.SysUser{}, id)

	//删除用户-角色关联
	global.DB.Model(model.RelateUserRole{}).
		Where("sys_user_id = ?", id).
		Delete(&model.RelateUserRole{})

	return result
}

// 根据账号查询密码
func SelectUserAndPass(username string) (model.SysUser, int64) {
	login := model.SysUser{}
	result := global.DB.Model(model.SysUser{}).Where("username =? ", username).First(&login)
	return login, result.RowsAffected
}

// 根据账号查询用户
func SelectUserByUserName(username string) int64 {
	user := model.SysUser{}
	result := global.DB.Model(model.SysUser{}).Where("username =? ", username).First(&user)
	return result.RowsAffected
}

// 新增角色
func InsertRole(role model.SysRole, jurisdictions []model.Jurisdiction) *gorm.DB {
	result := global.DB.Model(model.SysRole{}).Create(&role)
	for _, jurisdiction := range jurisdictions {
		rolePermission := model.RelateRolePermission{
			SysRoleID:       role.ID,
			SysPermissionID: jurisdiction.Value,
		}
		global.DB.Model(model.RelateRolePermission{}).Create(&rolePermission)
	}
	return result
}

// 删除角色
func DeleteRole(id int64) *gorm.DB {
	//查询角色名称
	sysRole, _ := SelectRoleById(id)
	//获取包含这个角色名称的权限 Policy
	var casbinRuleS = [][]string{}
	list, _ := authentication.Enforcer.GetPolicy()
	for _, vlist := range list {
		if vlist[0] == sysRole.Name {
			casbinRuleS = append(casbinRuleS, vlist)
		}
	}

	result := global.DB.Delete(&model.SysRole{}, id)

	//删除角色-权限关联
	global.DB.Model(model.RelateRolePermission{}).
		Where("sys_role_id = ?", id).
		Delete(&model.RelateRolePermission{})

	//删除权限 Policy
	for _, casbinRule := range casbinRuleS {
		ok, _ := authentication.Enforcer.RemovePolicy(casbinRule[0], casbinRule[1], casbinRule[2])
		if !ok {
			fmt.Println("Policy不存在")
		}
	}
	return result
}

// 修改角色
func UpdateRole(role model.SysRole, jurisdictions []model.Jurisdiction) *gorm.DB {
	result := global.DB.Save(&role)

	//删除已有的角色-权限关联
	global.DB.Model(model.RelateRolePermission{}).
		Where("sys_role_id = ?", role.ID).
		Delete(&model.RelateRolePermission{})

	//重新加入 角色-权限关联
	for _, jurisdiction := range jurisdictions {
		rolePermission := model.RelateRolePermission{
			SysRoleID:       role.ID,
			SysPermissionID: jurisdiction.Value,
		}
		global.DB.Model(model.RelateRolePermission{}).Create(&rolePermission)
	}
	return result
}

// 查询角色列表
func SelectRoleList() ([]model.SysRole, *gorm.DB) {
	var roles []model.SysRole
	result := global.DB.Model(model.SysRole{}).Find(&roles)
	return roles, result
}

// 根据id查询角色
func SelectRoleById(id int64) (model.SysRole, *gorm.DB) {
	var role model.SysRole
	result := global.DB.Where("ID = ?", id).First(&role)
	return role, result
}

// 根据用户id查询角色
func SelectRelateUserRoleById(id int64) ([]string, *gorm.DB) {
	var relateUserRoles []model.RelateUserRole
	result := global.DB.Model(model.RelateUserRole{}).Where("sys_user_id = ?", id).Find(&relateUserRoles)
	var roleNames []string
	for _, relateUserRole := range relateUserRoles {
		//根据角色id查询角色信息
		sysRole, _ := SelectRoleById(int64(relateUserRole.SysRoleID))
		roleNames = append(roleNames, sysRole.Name)
	}
	return roleNames, result
}

// 根据用户id查询权限ID
func SelectRelateUserRoleIDById(id int64) ([]uint, *gorm.DB) {
	var relateUserRoles []model.RelateUserRole
	result := global.DB.Model(model.RelateUserRole{}).Where("sys_user_id = ?", id).Find(&relateUserRoles)

	var relateUserRoleSNew []model.RelateRolePermission
	for _, relateUserRole := range relateUserRoles {
		//根据角色id查询权限id
		var relateUserRoleS []model.RelateRolePermission
		global.DB.Model(model.RelateRolePermission{}).Where("sys_role_id = ?", relateUserRole.SysRoleID).Find(&relateUserRoleS)
		relateUserRoleSNew = append(relateUserRoleSNew, relateUserRoleS...)
	}
	//去重
	relateUserRoleSNew = removeDuplicatesByRoleID(relateUserRoleSNew)

	var SysPermissionIDs []uint
	for _, permission := range relateUserRoleSNew {
		SysPermissionIDs = append(SysPermissionIDs, permission.SysPermissionID)
	}

	return SysPermissionIDs, result
}

// 定义一个函数来去除切片中的重复RoleID
func removeDuplicatesByRoleID(arr []model.RelateRolePermission) []model.RelateRolePermission {
	encountered := map[int]bool{}
	result := []model.RelateRolePermission{}

	for _, v := range arr {
		if !encountered[int(v.SysPermissionID)] {
			encountered[int(v.SysPermissionID)] = true
			result = append(result, v)
		}
	}
	return result
}

// 新增权限
func InsertPermission(permission model.SysPermission) *gorm.DB {
	result := global.DB.Model(model.SysPermission{}).Create(&permission)
	return result
}

// 删除权限
func DeletePermission(id int64) *gorm.DB {
	result := global.DB.Delete(&model.SysPermission{}, id)
	return result
}

// 修改权限
func UpdatePermission(permission model.SysPermission) *gorm.DB {
	result := global.DB.Updates(&permission)
	return result
}

// 查询顶级菜单列表
func SelectTopPermission() ([]*model.SysPermission, *gorm.DB) {
	var permissions []*model.SysPermission
	result := global.DB.Model(model.SysPermission{}).Where("type = ?", 0).Find(&permissions)
	return permissions, result
}

// 查询子级菜单列表
func SelectSubPermission(id uint) ([]*model.SysPermission, *gorm.DB) {
	var permissions []*model.SysPermission
	result := global.DB.Model(model.SysPermission{}).Where("parentId = ? ", id).Find(&permissions)
	return permissions, result
}

// 查询目录列表
func SelectDirectory() ([]model.SysPermission, *gorm.DB) {

	var directory []model.SysPermission
	result := global.DB.Model(model.SysPermission{}).Where("type <= ?", 1).Find(&directory)
	return directory, result
}

// 查询目录下的功能
func SelectFeatures(id uint) ([]model.SysPermission, *gorm.DB) {
	var features []model.SysPermission
	result := global.DB.Model(model.SysPermission{}).Where("parentId = ?", id).Find(&features)
	return features, result
}

// 根据id查询权限
func SelectPermissionById(id int64) (model.SysPermission, *gorm.DB) {
	var permission model.SysPermission
	result := global.DB.Where("ID = ?", id).First(&permission)
	return permission, result
}

// 新增验证码
func InsertCaptcha(captcha model.SysCaptcha) *gorm.DB {
	result := global.DB.Model(model.SysCaptcha{}).Create(&captcha)
	return result
}

// 删除验证码
func DeleteCaptcha(key string) *gorm.DB {
	result := global.DB.Where("verifyKey = ?", key).Delete(&model.SysCaptcha{})
	return result
}

// 根据权限id查询角色信息
func SelectPermissionList(permissionsID int64) ([]model.SysRole, *gorm.DB) {
	var relateRolePermission []model.RelateRolePermission
	result := global.DB.Model(model.RelateRolePermission{}).
		Where("sys_permission_id = ?", permissionsID).
		Find(&relateRolePermission)
	var roles []model.SysRole
	for _, rolePermission := range relateRolePermission {
		sysRole, _ := SelectRoleById(int64(rolePermission.SysRoleID))
		roles = append(roles, sysRole)
	}
	return roles, result
}

// 根据角色id获取权限
func selectPermissionByRoleID(roleIDsStr []string) {
	//根据角色名称获取id
	var roleIDs []int64
	for _, roleIDStr := range roleIDsStr {
		var role model.SysRole
		global.DB.Model(model.SysRole{}).Where("name = ?", roleIDStr).Find(&role)
		roleIDs = append(roleIDs, int64(role.ID))
	}
	//var rolePermission []model.RelateRolePermission
	//for _, roleID := range roleIDs {
	//	var rolePermission []model.RelateRolePermission
	//	global.DB.Model(model.RelateRolePermission{}).Where("sys_role_id = ?", roleID).Find(&rolePermission)
	//}
}
