package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func makeRequest(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "gosession")
	empId := session.Values["employee_id"]
	requestType := r.FormValue("requestType")
	requestDetails := r.FormValue("requestDetails")
	//fmt.Println(requestDetails, requestType, empId)
	query := "INSERT INTO requests(`request_type`, `request_body`, `emp_id`)" +
		"VALUES(?,?,?)"
	stmt, err := database.Prepare(query)
	checkErr(err)
	_, err = stmt.Exec(requestType, requestDetails, empId)
	checkErr(err)
	rp := "Request sent!"
	var templates = template.Must(template.ParseFiles("views/report.gohtml", "views/employeebase.gohtml"))
	varmap := map[string]interface{}{
		"Report": rp,
		"Auth":   IsAuth(r),
	}
	_ = templates.ExecuteTemplate(w, "employeebase", varmap)

}

func requestForm(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Po bon")
	var templates = template.Must(template.ParseFiles("views/employeeBase.gohtml", "views/employeeRequest.gohtml"))
	varmap := map[string]interface{}{
		"Auth": IsAuth(r),
	}

	err := templates.ExecuteTemplate(w, "employeebase", varmap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func takeTask(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("a")
	session, _ := store.Get(r, "gosession")
	empId := session.Values["employee_id"]
	if err := r.ParseForm(); err != nil {
		fmt.Println("Gabim ne parsimin e formes")
	}
	for key, _ := range r.PostForm {
		//fmt.Println(empId)
		taskId, _ := strconv.Atoi(key)
		//fmt.Println(key, empId, taskId)
		stmt, err := database.Prepare("UPDATE tasks SET emp_id=? WHERE task_id=?")
		checkErr(err)
		_, err = stmt.Exec(empId, taskId)
		checkErr(err)
		//fmt.Println("A")
		rp := "Task taken!"
		var templates = template.Must(template.ParseFiles("views/report.gohtml", "views/employeebase.gohtml"))
		varmap := map[string]interface{}{
			"Report": rp,
			"Auth":   IsAuth(r),
		}
		_ = templates.ExecuteTemplate(w, "employeebase", varmap)


	}
}

func seeTask(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("pss")
	if IsAuth(r) {
		task := Task{}
		Tasks := []Task{}
		session, _ := store.Get(r, "gosession")
		mgrId := session.Values["employee_manager"]
		//fmt.Println("Qitu", task.TaskId, task.Name, task.Details, task.DateCreated, task.DueDate, mgrId)

		query := "SELECT `task_id`, `taskt_title`, `task_body`, `date_created`, `due_date`" +
			", `bonus` FROM tasks LEFT JOIN employees ON tasks.emp_id = employees.employee_id WHERE tasks.manager_id=? AND emp_id =0"

		rows, err := database.Query(query, mgrId)
		if err != nil {
			fmt.Println("Gabim me databazen!")
		} else {
			for rows.Next() {

				rows.Scan(&task.TaskId, &task.Name, &task.Details, &task.DateCreated, &task.DueDate, &task.Bonus)
				//fmt.Println("Qitu", task.TaskId, task.Name, task.Details, task.DateCreated, task.DueDate, mgrId)

				Tasks = append(Tasks, task)

			}

			var templates = template.Must(template.ParseFiles("views/employeebase.gohtml", "views/employeeTask.gohtml"))
			varmap := map[string]interface{}{
				"Tasks": Tasks,
				"Auth":  IsAuth(r),
			}

			err = templates.ExecuteTemplate(w, "employeebase", varmap)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func profile(w http.ResponseWriter, r *http.Request) {
	if IsAuth(r) {
		session, _ := store.Get(r, "gosession")
		empId := session.Values["employee_id"]
		objEmployee := Employee{}
		Employees := []Employee{}

		rows, err := database.Query("SELECT employee_id, `name`, surname, salary, bonuses FROM employees WHERE employee_id=?", empId)
		fmt.Println(empId)
		if err != nil {
			log.Println("Gabim me databazen!")
		} else {
			for rows.Next() {
				rows.Scan(&objEmployee.EmployeeId, &objEmployee.Name, &objEmployee.Surname, &objEmployee.Salary, &objEmployee.Bonuses)

				Employees = append(Employees, objEmployee)
				//fmt.Println(Employees)
			}

			var templates = template.Must(template.ParseFiles("views/employeeBase.gohtml", "views/employeeProfile.gohtml"))
			varmap := map[string]interface{}{
				"Employees": Employees,
				"Auth":      IsAuth(r),
			}

			err = templates.ExecuteTemplate(w, "employeebase", varmap)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}
