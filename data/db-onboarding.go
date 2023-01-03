package data

import (
	"bufio"
	"context"
	"log"
	"os"
	"strconv"
	"strings"
)

type TenantInfo struct {
	SubaccountId string
	Subdomain    string
}

type Tenant struct {
	TenantId     string
	Schema       string
	Password     string
	Subdomain    string
	SubaccountId string
}

/**
 * Onboard subaccount. Fails if this tenant id was already saved or schema/user name already used.
 * Will create new user and schema, then create all needed tables and finally create CONFIGURATION entry with default data.
 * @param tenantId
 * @param subdomain will be used as schema name
 */
func Onboard(ctx context.Context, tenantId string, tenantInfo TenantInfo) string {

	// DB schema should be in  upper case
	schema := strings.ToUpper(tenantInfo.Subdomain)
	msg := "onboarding of tenant " + tenantId + " schema " + schema + " and subdomain " + tenantInfo.Subdomain + " DB admin " + os.Getenv("HANA_SECRET_ADMIN")
	log.Println(msg)

	//Validate uniqueness of tenant
	//Validate uniqueness of subdomain
	//Validate uniqueness of schema

	tenant := Tenant{tenantId, schema, os.Getenv("HANA_SECRET_PASSWORD"), tenantInfo.Subdomain, tenantInfo.SubaccountId}
	dba, err := db.getDb()
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Adfter db...")
	tx, err := dba.BeginTx(ctx, nil)
	log.Println("Adfter Tx...")
	q := "CREATE USER " + schema + " PASSWORD " + "\"" + tenant.Password + "\"" + " NO FORCE_FIRST_PASSWORD_CHANGE;"
	result, err := tx.Exec(q)
	if err != nil {
		log.Println("error1")
		log.Println(err.Error())
	}
	res, _ := result.RowsAffected()
	log.Println("result1: ", strconv.FormatInt(res, 10))
	q = "ALTER USER " + schema + " DISABLE PASSWORD LIFETIME;"
	result, err = tx.Exec(q)
	if err != nil {
		log.Println("error2")
		log.Println(err.Error())
	}
	res, _ = result.RowsAffected()
	log.Println("result2: ", strconv.FormatInt(res, 10))
	admin := os.Getenv("HANA_SECRET_ADMIN")
	q = "INSERT INTO " + admin + ".TENANTS " + "VALUES " + "('" + tenant.TenantId + "','" + tenant.Schema + "','" + tenant.Password + "','" + tenant.Subdomain + "','" + tenant.SubaccountId + "')"
	result, err = tx.Exec(q)
	if err != nil {
		log.Println("error3")
		log.Println(err.Error())
	}
	res, _ = result.RowsAffected()
	log.Println("result3: ", strconv.FormatInt(res, 10))
	tx.Commit()
	dbt, err := GetTenantDb(schema)
	if err != nil {
		log.Println(err.Error())
	}
	q = "SELECT CURRENT_SCHEMA FROM DUMMY;"
	var current_schema string
	rres := dbt.QueryRow(q)
	err = rres.Scan(&current_schema)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Current Schema: ", current_schema)
	q = "SELECT CURRENT_USER FROM DUMMY;"
	var current_user string
	rres = dbt.QueryRow(q)
	err = rres.Scan(&current_user)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Current User: ", current_user)
	file, err := os.Open("./data/data/create.sql")
	if err != nil {
		log.Println(err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var sql_total string
	for scanner.Scan() {
		sql_total += scanner.Text()
	}
	//log.Println(sql_total)
	sql_total = strings.ReplaceAll(sql_total, "__SCHEMANAME__", schema)
	sqls := strings.Split(sql_total, ";")
	log.Println("==============================")
	tx, err = dbt.BeginTx(ctx, nil)
	for _, sql := range sqls {
		q = sql
		_, err := tx.Exec(q)
		if err != nil {
			log.Println(err.Error())
		}
	}
	tx.Commit()
	log.Println("==============================")
	return "Success"
}

func Offboard(ctx context.Context, tenantId string) string {

	admin := os.Getenv("HANA_SECRET_ADMIN")
	dba, err := db.getDb()
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("After dba...")

	q := "SELECT TOP 1 * FROM " + admin + "." + "TENANTS WHERE TENANTID='" + tenantId + "';"
	var tenant Tenant
	rres := dba.QueryRow(q)
	err = rres.Scan(&tenant.TenantId, &tenant.Schema, &tenant.Password, &tenant.Subdomain, &tenant.SubaccountId)
	if err != nil {
		log.Println("error1")
		log.Println(err.Error())
	}

	msg := "offboarding of tenant: " + tenantId + " Schema: " + tenant.Schema + " and subdomain " + tenant.Subdomain + " DB admin " + admin
	log.Println(msg)

	//Validate uniqueness of tenant
	//Validate uniqueness of subdomain
	//Validate uniqueness of schema

	tx, err := dba.BeginTx(ctx, nil)
	log.Println("Adfter Tx...")
	q = "DROP USER " + tenant.Schema + " CASCADE;"
	result, err := tx.Exec(q)
	if err != nil {
		log.Println("error1")
		log.Println(err.Error())
	}
	res, _ := result.RowsAffected()
	log.Println("result1: ", strconv.FormatInt(res, 10))
	q = "DELETE FROM " + admin + ".TENANTS" + " WHERE " + "TENANTID=" + "'" + tenant.TenantId + "';"
	result, err = tx.Exec(q)
	if err != nil {
		log.Println("error2")
		log.Println(err.Error())
	}
	res, _ = result.RowsAffected()
	log.Println("result2: ", strconv.FormatInt(res, 10))
	tx.Commit()
	return "Success"
}
