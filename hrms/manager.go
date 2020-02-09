package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func dashboard(w http.ResponseWriter, r *http.Request) {
	if IsAuth(r) {
		session, _ := store.Get(r, "gosession")
		mgrId := session.Values["manager_id"]
		objEmployee := Employee{}
		Employees := []Employee{}

		rows, err := database.Query("SELECT employee_id, `name`, surname, salary, bonuses FROM employees WHERE manager_id=?", mgrId)
		//fmt.Println(mgrId)
		if err != nil {
			log.Println("Gabim me databazen!")
		} else {
			for rows.Next() {
				rows.Scan(&objEmployee.EmployeeId, &objEmployee.Name, &objEmployee.Surname, &objEmployee.Salary, &objEmployee.Bonuses)

				Employees = append(Employees, objEmployee)
				//fmt.Println(Employees)
			}

			var templates= template.Must(template.ParseFiles("views/managerBase.gohtml", "views/managerDashboard.gohtml"))
			varmap := map[string]interface{}{
				"Employees": Employees,
				"Auth":      IsAuth(r),
			}

			err = templates.ExecuteTemplate(w, "managerBase", varmap)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func modifySalary(w http.ResponseWriter, r *http.Request){
	if IsAuth(r){
		salary := r.FormValue("salary")
		if err := r.ParseForm(); err != nil {
			fmt.Println("Gabim ne parsimin e formes")
		} else {
			for key, _ := range r.PostForm {
				if strings.HasPrefix(key, "modify") {
					empId, _ := strconv.Atoi(key[7 : len(key)-1])
					//fmt.Println(empId, key)
					stmt, err := database.Prepare("UPDATE employees SET salary=? WHERE employee_id=?")
					checkErr(err)
					_, err = stmt.Exec(salary, empId)
					checkErr(err)
					//fmt.Println("A")
					rp := "Salary changed!"
					var templates = template.Must(template.ParseFiles("views/report.gohtml", "views/managerBase.gohtml"))
					varmap := map[string]interface{}{
						"Report": rp,
						"Auth":   IsAuth(r),
					}
					_ = templates.ExecuteTemplate(w, "managerBase", varmap)
				}
			}
		}
	}else {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}

}

func fireEmployee(w http.ResponseWriter, r *http.Request){
	if IsAuth(r) {
		if err := r.ParseForm(); err != nil {
			fmt.Println("Gabim ne parsimin e formes")
		} else {
			for key, _ := range r.PostForm {
				if strings.HasPrefix(key, "fire") {
					empId, _ := strconv.Atoi(key[5 : len(key)-1])
					//fmt.Println(empId, key)
					stmt, err := database.Prepare("DELETE FROM employees WHERE employee_id=?")
					checkErr(err)
					_, err = stmt.Exec(empId)
					checkErr(err)
					//fmt.Println("A")
					rp := "Employee Fired"
					var templates= template.Must(template.ParseFiles("views/report.gohtml", "views/managerBase.gohtml"))
					varmap := map[string]interface{}{
						"Report": rp,
						"Auth":   IsAuth(r),
					}
					_ = templates.ExecuteTemplate(w, "managerBase", varmap)
				}
			}
		}
	}else {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func tasks(w http.ResponseWriter, r *http.Request) {
	if IsAuth(r) {
		session, _ := store.Get(r, "gosession")
		mgrId := session.Values["manager_id"]
		task := Task{}
		Tasks := []Task{}
		query := "SELECT `task_id`, `taskt_title`, `task_body`, `date_created`, `due_date`" +
			", `bonus`, `emp_id`, tasks.manager_id, employees.`name` FROM tasks LEFT JOIN employees " +
			"ON tasks.emp_id = employees.employee_id WHERE tasks.manager_id=?"
		rows, err := database.Query(query, mgrId)
		if err != nil {
			log.Println("Gabim me databazen!")
		} else {
			for rows.Next() {
				err := rows.Scan(&task.TaskId, &task.Name, &task.Details, &task.DateCreated, &task.DueDate, &task.Bonus, &task.EmpId, &task.ManagerId, &task.EmpName)
				if err!=nil {
					fmt.Println("Gabim ne shtimin e te dhenave ne struct \n", task)
				}
				if task.EmpId == 0{
					task.EmpName = ""
				}

				//fmt.Println(task.EmpId, task.EmpName)
				Tasks = append(Tasks, task)

			}

			var templates= template.Must(template.ParseFiles("views/managerBase.gohtml", "views/managerTasks.gohtml"))
			varmap := map[string]interface{}{
				"Tasks": Tasks,
				"Auth":  IsAuth(r),
			}

			err = templates.ExecuteTemplate(w, "managerBase", varmap)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}else {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func addTask(w http.ResponseWriter, r *http.Request){
	session, _ := store.Get(r, "gosession")
	mgrId := session.Values["manager_id"]
	taskTitle := r.FormValue("taskTitle")
	taskDetails := r.FormValue("taskDetails")
	dueDate := r.FormValue("dueDate")
	bonus := r.FormValue("bonus")

	currentTime := time.Now()
	dateCreated := currentTime.Format("2006-01-02")
	//fmt.Println(dateCreated)
	// TO BE 
	//fmt.Println(mgrId)
	query := "INSERT INTO tasks(`taskt_title`, `task_body`, `date_created`, `due_date`, `bonus`, `manager_id`)" +
		"VALUES(?,?,?,?,?,?)"
	stmt, err := database.Prepare(query)
	checkErr(err)
	_, err = stmt.Exec(taskTitle, taskDetails, dateCreated, dueDate, bonus, mgrId)
	checkErr(err)
	rp := "Task Added Successfully"
	var templates= template.Must(template.ParseFiles("views/report.gohtml", "views/managerBase.gohtml"))
	varmap := map[string]interface{}{
		"Report": rp,
		"Auth":   IsAuth(r),
	}
	_ = templates.ExecuteTemplate(w, "managerBase", varmap)

}

func seeRequests(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "gosession")
	mgrId := session.Values["manager_id"]
	request := Request{}
	Requests := []Request{}
	query := "SELECT request_id, request_type, request_body, `approval`, employees.name, emp_id FROM requests " +
		"left join employees ON emp_id = employee_id WHERE employees.manager_id =?"

	rows, err := database.Query(query, mgrId)
	if err != nil {
		log.Println("Gabim me databazen!")
	} else {
		for rows.Next() {
			rows.Scan(&request.RequestId, &request.RequestType, &request.RequestBody, &request.Approval, &request.EmpName, &request.EmpId)

			//fmt.Println("a", request.Approval)
			Requests = append(Requests, request)

		}
		var templates= template.Must(template.ParseFiles("views/managerBase.gohtml", "views/managerRequests.gohtml"))
		varmap := map[string]interface{}{
			"Requests": Requests,
			"Auth":  IsAuth(r),
		}

		err = templates.ExecuteTemplate(w, "managerBase", varmap)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}


func approveRequests(w http.ResponseWriter, r *http.Request) {
	var rp string
	if err := r.ParseForm(); err != nil {
		fmt.Println("Gabim ne parsimin e formes")
	} else {
		for key, _ := range r.PostForm {
			if strings.HasPrefix(key, "app") {
				requestId, _ := strconv.Atoi(key[4 : len(key)-1])
				//fmt.Println(requestId, key)
				stmt, err := database.Prepare("UPDATE requests SET approval=1 WHERE request_id=?")
				checkErr(err)
				_, err = stmt.Exec(requestId)
				checkErr(err)
				//fmt.Println("A")
				//fmt.Println("a")
				rp = "Approved"
			}
			var templates = template.Must(template.ParseFiles("views/report.gohtml", "views/managerBase.gohtml"))
			varmap := map[string]interface{}{
				"Report": rp,
				"Auth":   IsAuth(r),
			}
			_ = templates.ExecuteTemplate(w, "managerBase", varmap)
		}
	}
}

