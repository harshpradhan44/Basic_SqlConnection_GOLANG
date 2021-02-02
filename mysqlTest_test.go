package main

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"io/ioutil"
	_ "io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func Test_GetData(t *testing.T) {
	d,mock,err := sqlmock.New()
	if err!=nil{
		t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
	}
	defer d.Close()
	db := &sqlConnection{d}

	testcases := []Employee{
		{6,"Bala",23,"M",3},
		{7,"Anshu",23,"F",2},

	}
	for i,tc := range testcases {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/employee/%v",tc.Id), nil)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		rows := sqlmock.NewRows([]string{"Id", "Name", "Age","Gender","Role"}).
			AddRow(tc.Id,tc.Name,tc.Age,tc.Gender,tc.NewRole)

		req = mux.SetURLVars(req, map[string]string{
			"id": strconv.Itoa(tc.Id),
		})
		query := "select id,name,age,gender,role from employee where id=?"
		mock.ExpectQuery(query).WithArgs(tc.Id).WillReturnRows(rows)
		db.GetData(w, req)
		res := w.Result()
		resByte, err := ioutil.ReadAll(res.Body)
		//resByte := w.Body.Bytes()
		if err != nil {
			t.Fail()
		}
		fmt.Printf("resbyte:%v tc: %v\n",string(resByte), tc)
		if  w.Code!=200{
			t.Errorf("test case failed at %v actual: %v expected: %v",i+1,string(resByte),w.Code)
		}
		//if resByte != tc{
		//	t.Errorf("test case failed at %v actual: %v expected: %v",i+1,string(resByte),expected)
		//}
	}
}

/*

func Test_InsertData(t *testing.T) {
	testcases := []struct{
		body string
		expected string
		status int
	}{
		//{"{      \n        \"name\": \"Anshu\",\n        \"age\": 23,\n        \"gender\": \"F\",\n        \"role\": 2\n}\n", "ID:12 Record Inserted Succesfully! ",202},
		{"{      \n        \"name\": \"Anshu\",\n        \"age\": 23,\n        \"gender\": \"F\"\n}","Error while executing query \n",400},
	}
	for i,tc := range testcases {
		//var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
		req, err := http.NewRequest("GET", fmt.Sprintf("/insertEmployee"),  bytes.NewBuffer([]byte(tc.body)))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		insertData(w, req)
		res := w.Result()
		resByte, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fail()
		}
		if  string(resByte)!= tc.expected{
			t.Errorf("test case failed at %v actual: %v expected: %v",i+1,string(resByte),tc.expected)
		}
		if res.StatusCode != tc.status {
			t.Errorf("test case failed at %v actual: %v expected: %v",i+1,res.StatusCode,tc.status)
		}
	}
}

 */
