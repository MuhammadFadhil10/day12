package main

import (
	"fmt"
	"mvcweb/connection"
	"mvcweb/controller"
	"mvcweb/middleware"
	"net/http"

	"github.com/gorilla/mux"
)




func main() {
	router := mux.NewRouter()

	directory := http.Dir("./public")
	fileServer := http.FileServer(directory)
	uploadDir := http.Dir("./uploads")
	uploadFileServer := http.FileServer(uploadDir)

    router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer));
    router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads", uploadFileServer));

	// router
	// get project
	router.HandleFunc("/", controller.GetHome).Methods("GET")
	router.HandleFunc("/form-add-project", controller.GetAddProject).Methods("GET")
	router.HandleFunc("/form-edit-project/{index}", controller.GetEditProject).Methods("GET")
	router.HandleFunc("/contact-me", controller.GetContactMe).Methods("GET")
	router.HandleFunc("/project/{projectId}", controller.GetProjectDetail).Methods("GET")
	// get auth
	router.HandleFunc("/form-register", controller.GetRegisterForm).Methods("GET")
	router.HandleFunc("/form-login", controller.GetLoginForm).Methods("GET")
	// post project
	router.HandleFunc("/add-project", middleware.UploadFile(controller.PostAddProject)).Methods("POST")
	router.HandleFunc("/update-project/{index}", controller.UpdateProject).Methods("POST")
	router.HandleFunc("/delete-project/{projectId}", controller.DeleteProject).Methods("POST")
	// post auth
	router.HandleFunc("/auth/register", controller.Register).Methods("POST")
	router.HandleFunc("/auth/login", controller.Login).Methods("POST")
	router.HandleFunc("/auth/logout", controller.Logout).Methods("GET")
	
	
	connection.DatabaseConnect(func() {
		fmt.Println("running on port 5000");
		http.ListenAndServe("localhost:5000", router);
	})

}












