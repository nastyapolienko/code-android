package main
import(
	"fmt"
	"database/sql" 
	"encoding/json"
	"log"
	"net/http"
	_"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io/ioutil"
)
const(
	CONN_PORT = "8080"
	DRIVER_NAME = "mysql"
	DATA_SOURCE_NAME = "root:1111@/library"
)
var db *sql.DB
var connectionError error
func init(){
	db, connectionError = sql.Open(DRIVER_NAME, DATA_SOURCE_NAME)
	if connectionError != nil{
		log.Fatal("error connecting to database :: ", connectionError)
	}
}
type Book struct{
Id int `json:"bid"`
Name string `json:"bookname"`
}
func getBooks(w http.ResponseWriter, r *http.Request){
	log.Print("reading records from database")
	rows, err := db.Query("SELECT * FROM books")
	if err != nil{
		log.Print("error occurred while executing select query :: ",err)
		return
	}
	books := []Book{}
	for rows.Next(){
		var bid int
		var bookname string
		err = rows.Scan(&bid, &bookname)
		book := Book{Id: bid, Name: bookname}
		books = append(books, book)
	}
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	log.Print("reading a record from database")
	result, err := db.Query("SELECT Id, bookname FROM books WHERE Id = ?", params["Id"])
	if err != nil {
	  panic(err.Error())
	}
	defer result.Close()
	var book Book
	for result.Next() {
	  err := result.Scan(&book.Id, &book.Name)
	  if err != nil {
		panic(err.Error())
	  }
	}
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE books SET bookname = ? WHERE Id = ?")
	if err != nil {
	  panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
	  panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newName := keyVal["bookname"]
	_, err = stmt.Exec(newName, params["Id"])
	if err != nil {
	  panic(err.Error())
	}
	fmt.Fprintf(w, "Book with Id = %s was updated", params["Id"])
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stmt, err := db.Prepare("INSERT INTO books(bookname) VALUES(?)")
	if err != nil {
	  panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
	  panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	bookname := keyVal["bookname"]
	_, err = stmt.Exec(bookname)
	if err != nil {
	  panic(err.Error())
	}
	fmt.Fprintf(w, "New post was created")
  }

  func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM books WHERE Id = ?")
	if err != nil {
	  panic(err.Error())
	}
	_, err = stmt.Exec(params["Id"])
   if err != nil {
	  panic(err.Error())
	}
  fmt.Fprintf(w, "Book with Id = %s was deleted", params["Id"])
  }

func main(){
	router := mux.NewRouter()
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{Id}", getBook).Methods("GET")
	router.HandleFunc("/books", createBook).Methods("POST")
	router.HandleFunc("/books/{Id}", updateBook).Methods("PUT")
	router.HandleFunc("/books/{Id}", deleteBook).Methods("DELETE")
	defer db.Close()
	err := http.ListenAndServe(":"+CONN_PORT, router)
	if err != nil{
		log.Fatal("error starting http server :: ", err)
		return
	}
} 