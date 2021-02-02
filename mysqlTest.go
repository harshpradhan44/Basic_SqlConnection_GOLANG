package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "log"
	"net/http"
	"strconv"
)

type sqlConnection struct{
	Db *sql.DB
}

const (
	username = "root"
	host = "127.0.0.1"
	password = "Harsh@123"
	database = "harshDB"
)

type Employee struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Age int `json:"age"`
	Gender string `json:"gender"`
	NewRole int `json:"role"`
}

func connectionString() string{
	return fmt.Sprintf("%v:%v@(%v)/%v",username,password,host,database)
}

/*
func dbms(){
	fmt.Println("DBMS Testing!")

	db, err := sql.Open("mysql",connectionString())

	if err != nil {
		fmt.Println("Here")
		panic(err.Error())
	}
	query := `CREATE TABLE employee(
				ID int AUTO_INCREMENT PRIMARY KEY ,
				NAME VARCHAR (20) NOT NULL ,
				AGE INT,
				GENDER VARCHAR (1),
               ROLE int ,
               FOREIGN key(ROLE) references role(ID) ON update cascade
				);  `
	query = `INSERT INTO employee(name ,age,gender,role ) values("Harsh",23,"M",3)`
	//query = `ALTER table Employee ALTER column set role = "SDE" where id==1"`
	//query = `select Employee.n
	//ame,role.field from Employee INNER JOIN role ON Employee.id=role.id`
	//query = `drop table employee`
	query = `select * from employee`
	//query = `CREATE TABLE role (
	//query = `INSERT INTO role(field) values("SDE")`
	//	ID int AUTO_INCREMENT PRIMARY Key,
	//	NAME varchar(20) NOT NULL
	//)`

	res, err := db.Query(query)
	if err != nil {
		log.Printf("Error %s when creating product table", err)
		//return err
	}
	//var mp = make(map[string]int)
	for res.Next() {

		var s Employee

		if err :=res.Scan(&s.Id,&s.Name,&s.Age,&s.Gender,&s.NewRole); err != nil {
			panic(err)
		}
		fmt.Println(s)

		//var (
		//	id int
		//	role string
		//)
		//res.Scan(&id,&role)
		//fmt.Printf("Name:%v Role:%v\n",id,role)

	}
}

func roleDB(){
	db, err := sql.Open("mysql", connectionString())
	// updating Employee using role
	if err!=nil{
		fmt.Errorf("error Occured during opening connection")
	}
	query := `CREATE TABLE products(
		ID int AUTO_INCREMENT PRIMARY Key,
		NAME varchar(20) NOT NULL,
        Brand_ID int,
        FOREIGN key(Brand_ID) references brands(ID) ON update cascade
    
	)`
	query = `CREATE TABLE brands(
		ID int AUTO_INCREMENT PRIMARY Key,
		NAME varchar(20) NOT NULL
	)`

	//query = `select role , count(role ) as count from Employee group by role;`
	//query = `INSERT INTO products values(2,"Television",5)`
	query = `select * from products`
	//query = `update role set field = "SOFTWARE DEVELOPER" where id = 1 `
	ans, er := db.Query(query)
	if er!=nil{
		log.Printf("Error")
	}
	defer ans.Close()
	for ans.Next(){
		var (
			id int
			name string
			bID int
		)
		if e:= ans.Scan(&id,&name,&bID); e!=nil{
			log.Println("Error")
		}
		fmt.Printf("ID:%v name:%v %v\n",id,name,bID)
	}
}

 */

func (db *sqlConnection)getAllData(w http.ResponseWriter, r *http.Request){

	w.Header().Set("Content-Type", "application/json")
	//db,err := sql.Open("mysql",connectionString())
	//if err!=nil{
	//	http.Error(w,"Error Opening dbms connection",500)
	//}
	res,er:= db.Db.Query("select * from employee;")

	if er!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		panic("error res")
	}
	var posts[] Employee
	var s Employee
	for res.Next(){
		if e:= res.Scan(&s.Id,&s.Name,&s.Age,&s.Gender,&s.NewRole); e!=nil{
			panic("Error reading" )
		}
		fmt.Println(s)
		posts = append(posts,s)
	}
	json.NewEncoder(w).Encode(posts)
}

func (db *sqlConnection)GetData(w http.ResponseWriter, r * http.Request){

	w.Header().Set("Content-Type", "application/json")
	//db,err := sql.Open("mysql",connectionString())
	//if err!=nil{
	//	http.Error(w,"Error Opening dbms connection",500)
	//}
	params := mux.Vars(r)
	id,_:= strconv.Atoi(params["id"])
	res,er:= db.Db.Query("select id,name,age,gender,role from employee where id=?",id)
	if er!=nil{
		http.Error(w,"error in query",http.StatusBadRequest)
		return
	}
	//var posts[] Employee
	var s Employee
	for res.Next(){
		if e:= res.Scan(&s.Id,&s.Name,&s.Age,&s.Gender,&s.NewRole); e!=nil{
			panic("Error reading" )
		}
		fmt.Println(s)
		//posts = append(posts,s)
	}
	if s.Id==0{
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(s)
}

func (db *sqlConnection)insertData(w http.ResponseWriter, r *http.Request){

	w.Header().Set("Content-Type", "application/json")
	//db,err := sql.Open("mysql",connectionString())
	//defer db.Close()
	//if err!=nil{
	//	http.Error(w,"Error Opening dbms connection",500)
	//}
	var p Employee
	//val,er:= json.Marshal(r.Body)
	er := json.NewDecoder(r.Body).Decode(&p)
	if er != nil {
		http.Error(w, er.Error(), http.StatusBadRequest)
		return
	}
	var id int64
	query := fmt.Sprintf("insert into employee(Name,Age,Gender,Role) values('%v',%v,'%v',%v);",p.Name,p.Age,p.Gender,p.NewRole)
	res,err:= db.Db.Exec(query)
	if err!=nil{
		http.Error(w,"Error while executing query ",http.StatusBadRequest)
		return
	}
	id ,_= res.LastInsertId()
	fmt.Println(res.LastInsertId())
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w,"ID:%v Record Inserted Succesfully! ",id)
}

func (db *sqlConnection)deleteRecord(w http.ResponseWriter, r *http.Request){

	w.Header().Set("Content-Type", "application/json")
	//db,err := sql.Open("mysql",connectionString())
	//if err!=nil{
	//	http.Error(w,"Error Opening dbms connection",500)
	//}
	params := mux.Vars(r)
	query := fmt.Sprintf("delete from employee where id=%v",params["id"])
	res,er:= db.Db.Exec(query)
	if er!=nil{
		http.Error(w,"Invalid ID",http.StatusBadRequest)
		return
	}
	numRow,_ := res.RowsAffected()
	if numRow==0{
		http.Error(w,"Invalid ID",http.StatusBadRequest)
		return
	}

	//fmt.Fprintf(w,"Record Deleted Successfully! ID:%v",params["id"])
	w.WriteHeader(204)
}

func (db *sqlConnection) updateData(w http.ResponseWriter, r *http.Request){

	w.Header().Set("Content-Type", "application/json")
	//db,err := sql.Open("mysql",connectionString())
	//if err!=nil{
	//	http.Error(w,"Error Opening dbms connection",500)
	//	return
	//}
	params := mux.Vars(r)
	var p Employee
	e := json.NewDecoder(r.Body).Decode(&p)
	if e!= nil{
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	query := fmt.Sprintf("update employee set Name='%v', Age=%v, Gender='%v',Role='%v' where id=%v",p.Name,p.Age,p.Gender,p.NewRole,params["id"])
	res,er:= db.Db.Exec(query)
	if er!=nil{
		http.Error(w,"Error executing query", http.StatusInternalServerError) //500
	}
	numRow,_ := res.RowsAffected()
	if numRow==0{
		http.Error(w,"Invalid ID",http.StatusBadRequest) //400
		return
	}
	//w.WriteHeader(204)

	fmt.Fprintf(w,"ID:%v Record Updated Successfully!",params["id"])
}

/*
func getParticularData(w http.ResponseWriter, r *http.Request){

}

 */


func dbTest(){
	fmt.Println("Server Running...")
	db, err := sql.Open("mysql", connectionString())
	if err!=nil{
		fmt.Println("Error")
	}
	defer db.Close()
	d := sqlConnection{Db: db}
	router := mux.NewRouter()
	router.HandleFunc("/employee",d.getAllData).Methods("GET") //done
	//router.HandleFunc()
	router.HandleFunc("/employee/{id}",d.GetData).Methods("GET") //done
	router.HandleFunc("/employee",d.insertData).Methods("POST") //done
	router.HandleFunc("/employee/{id}",d.deleteRecord).Methods("DELETE") //dome
	router.HandleFunc("/employee/{id}",d.updateData).Methods("PUT") //done
	http.ListenAndServe(":8080",router)
}
func main() {
	//dbms()
	//roleDB()
	dbTest()
}