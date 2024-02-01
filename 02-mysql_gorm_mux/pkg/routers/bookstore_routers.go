package routers

import (
	"github.com/gorilla/mux"
	"mysql_book_management_system/pkg/controllers"
)

var RegisterBookStoreRoutes = func(r *mux.Router) {
	r.HandleFunc("/book", controllers.GetBook).Methods("GET")
	r.HandleFunc("/book", controllers.CreateBook).Methods("POST")
	r.HandleFunc("/book/{bookId}", controllers.UpdateBook).Methods("PUT")
	r.HandleFunc("/book/{bookId}", controllers.DeleteBook).Methods("DELETE")
	r.HandleFunc("/book/{bookId}", controllers.GetBookById).Methods("GET")
}
