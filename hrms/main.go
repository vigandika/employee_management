package  main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	//"time"
)


const (
	DBHost     = "127.0.0.1"
	DBPort     = ":3306"
	DBUser     = "xxx"
	DBPass     = "yyy"
	DBDbase    = "management"
)

type Employee struct {
	EmployeeId int     `json:"employee_id"`
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	ManagerId  int     `json:"password"`
	Salary     float32 `json:"salary"`
	Bonuses    float32 `json:"bonuses"`
}

type Manager struct {
	ManagerId int    `json:"manager_id"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Task struct {
	TaskId 		  int     `json:"task_id"`
	Name 		  string  `json:"task_title"`
	Details       string  `json:"task_body"`
	DateCreated   string  `json:"date_created"`
	DueDate       string  `json:"due_date"`
	Bonus  		  float32 `json:"bonus"`
	EmpId  		  int 	  `json:"EmpId"`
	EmpName 	  string  `json:"EmpName"`
	ManagerId  	  int     `json:"manager_id"`
}

type Request struct {
	RequestId	  int     `json:"request_id"`
	RequestType	  string  `json:"request_type"`
	RequestBody   string  `json:"request_body"`
	RequestDate   string  `json:"request_date"`
	Approval	  bool 	  `json:"approval"`
	EmpId  		  int 	  `json:"emp_id"`
	EmpName 	  string  `json:"EmpName"`
}

var (
	database *sql.DB
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("ahfy7634*&^%23DFERko923456df0&")
	store = sessions.NewCookieStore(key)
	auth  bool
)




func main() {
	dbConn := fmt.Sprintf("%s:%s@tcp(%s)/%s", DBUser, DBPass, DBHost, DBDbase)
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		panic("Couldn't connect!")
	}
	database = db

	routes := mux.NewRouter()
	routes.HandleFunc("/", Home).Methods("GET")
	routes.HandleFunc("/logincheck", loginCheck).Methods("POST")
	routes.HandleFunc("/register", registrationForm).Methods("GET")
	routes.HandleFunc("/registrationCheck", registrationCheck).Methods("POST")
	routes.HandleFunc("/logout", logout).Methods("GET")

	routes.HandleFunc("/manager/dashboard", dashboard).Methods("GET")
	routes.HandleFunc("/manager/salaries", modifySalary).Methods("POST")
	routes.HandleFunc("/manager/fire", fireEmployee).Methods("POST")
	routes.HandleFunc("/manager/tasks", tasks).Methods("GET")
	routes.HandleFunc("/manager/addTask", addTask).Methods("POST")
	routes.HandleFunc("/manager/requests", seeRequests).Methods("GET")
	routes.HandleFunc("/manager/approve", approveRequests).Methods("POST")


	routes.HandleFunc("/employee/profile", profile).Methods("GET")
	routes.HandleFunc("/employee/tasks", seeTask).Methods("GET")
	routes.HandleFunc("/employee/takeTask", takeTask).Methods("POST")
	routes.HandleFunc("/employee/request", requestForm).Methods("GET")
	routes.HandleFunc("/employee/makeRequest", makeRequest).Methods("POST")


	log.Println(http.ListenAndServe(":8085", routes))
}






func Home(w http.ResponseWriter, r *http.Request) {
	var templates = template.Must(template.ParseFiles("views/landing.gohtml", "views/login.gohtml"))

	err := templates.ExecuteTemplate(w, "landing", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}


func loginCheck(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "gosession")
	email := r.FormValue("email")
	password := r.FormValue("password")
	role := r.FormValue("role")

	//U := interface{}
	if role == "manager"{
		// okej fmt.Println("A")
		manager := Manager{}
		err := database.QueryRow("SELECT manager_id, `name`, surname, email, password FROM managers WHERE email=?", email).
			Scan(&manager.ManagerId, &manager.Name, &manager.Surname, &manager.Email, &manager.Password)
		if err != nil {
			var templates = template.Must(template.ParseFiles("views/report.gohtml", "views/landing.gohtml"))
			varmap := map[string]interface{}{
				"Report": "Authentication failed",
			}
			err := templates.ExecuteTemplate(w, "landing", varmap)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			if CheckPasswordHash(password, manager.Password) {
				session.Values["authenticated"] = true
				session.Values["manager_id"] = manager.ManagerId
				session.Values["manager_name"] = manager.Name
				session.Values["manager_email"] = manager.Email

				session.Save(r, w)
				http.Redirect(w, r, "/manager/dashboard", http.StatusSeeOther)
			} else {
				//fmt.Println(password, "\n" , manager.Password)
				varmap := map[string]interface{}{
					"Report": "Authentication failed",
				}
				var templates = template.Must(template.ParseFiles("views/report.gohtml", "views/landing.gohtml"))
				err := templates.ExecuteTemplate(w, "landing", varmap)
				if err != nil{
					fmt.Println(err)
				}
			}
		}
	}else if role == "employee"{
		//fmt.Println("A")
		employee := Employee{}
		err := database.QueryRow("SELECT employee_id, `name`, surname, email, password, manager_id, salary, bonuses FROM employees WHERE email=?", email).
			Scan(&employee.EmployeeId, &employee.Name, &employee.Surname, &employee.Email, &employee.Password, &employee.ManagerId, &employee.Salary, &employee.Bonuses)

		if err != nil {
			fmt.Println(err)
			var templates = template.Must(template.ParseFiles("views/report.gohtml", "views/landing.gohtml"))
			varmap := map[string]interface{}{
				"Report": "Authentication failed",
			}
			_ = templates.ExecuteTemplate(w, "landing", varmap)
		} else {
			//fmt.Println("Qitu")
			if CheckPasswordHash(password, employee.Password) {
				session.Values["authenticated"] = true
				session.Values["employee_id"] = employee.EmployeeId
				session.Values["employee_name"] = employee.Name
				session.Values["employee_email"] = employee.Email
				session.Values["employee_manager"] = employee.ManagerId

				session.Save(r, w)

				http.Redirect(w, r, "/employee/profile", http.StatusSeeOther)
			} else {
				varmap := map[string]interface{}{
					"Report": "Authentication failed",
				}
				var templates = template.Must(template.ParseFiles("views/report.gohtml", "views/landing.gohtml"))
				_ = templates.ExecuteTemplate(w, "landing", varmap)
			}
		}
	}



}


func registrationForm(w http.ResponseWriter, r *http.Request) {
	var templates = template.Must(template.ParseFiles("views/registration.gohtml", "views/landing.gohtml"))

	_ = templates.ExecuteTemplate(w, "landing", nil)
}

func registrationCheck(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	surname := r.FormValue("surname")
	email := r.FormValue("email")
	password1 := r.FormValue("password1")
	password2 := r.FormValue("password2")
	managerEmail := r.FormValue("managerEmail")
	//fmt.Println(managerEmail)
	// Gjej menaxherin ne databaze
	var managerId int

	err := database.QueryRow("SELECT manager_id FROM managers WHERE email=?", managerEmail).Scan(&managerId)
	//fmt.Println(managerId)
	if err != nil {
		fmt.Println(err)
	}

	if password1 == password2 {
		sqlquery := "INSERT INTO employees(`name`, surname, email, password, manager_id, salary, bonuses) VALUES (?,?,?,?,?,?,?)"

		stmt, err := database.Prepare(sqlquery)
		if err != nil {
			fmt.Println(err)
		}

		password, _ := HashPassword(password1)
		//fmt.Println(name, surname, email, password, managerId)
		_, err = stmt.Exec(name, surname, email, password, managerId, 0,0)
		if err != nil {
			fmt.Println(err)
		}

		http.Redirect(w,r, "/", http.StatusSeeOther)
	} else {
		rp := "Fjalekalimet nuk perputhen"
		var templates = template.Must(template.ParseFiles("views/report.gohtml", "views/landing.gohtml"))
		varmap := map[string]interface{}{
			"Report": rp,
		}
		_ = templates.ExecuteTemplate(w, "landing", varmap)

	}

}



func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "gosession")

	session.Values["authenticated"] = false
	session.Values["user_name"] = ""
	session.Values["user_email"] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}


func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func IsAuth(r *http.Request) bool {
	session, _ := store.Get(r, "gosession")
	var ok bool
	auth, ok = session.Values["authenticated"].(bool)

	if ok {
		return auth
	}
	return false
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	//fmt.Println(string(bytes))
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

