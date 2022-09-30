package controller

import (
	"context"
	"fmt"
	"html/template"
	"mvcweb/connection"
	"mvcweb/helper"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type ProjectData struct {
	Id int
	AuthorId, AuthorName interface{}
	Name,Description,Image,Duration, StringStartDate, StringEndDate string
	StartDate,EndDate time.Time
	Technologies[]string	
	
}

type SessionStruct struct {
	Id int
	Email,Name,Flash string
	IsLogin bool
}


func GetHome(w http.ResponseWriter, r *http.Request) {
	data, err := connection.Conn.Query(context.Background(), `
		SELECT tb_projects.id,tb_projects.name,tb_projects.start_date,tb_projects.end_date,tb_projects.description,
		tb_projects.technologies,tb_projects.image, tb_user.name as authorName, tb_user.id as authorId 
		FROM public.tb_projects LEFT JOIN tb_user ON tb_projects.author_id = tb_user.id ORDER BY tb_projects.posted_at DESC
	`)

	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	var dataResult []ProjectData

	for data.Next() {
		var project = ProjectData{}

		var err = data.Scan(
			&project.Id,&project.Name, &project.StartDate, &project.EndDate, &project.Description,
			&project.Technologies, &project.Image, &project.AuthorName, &project.AuthorId,
		);
		if err != nil {
			panic(err.Error())
		}
		
		project.Duration = helper.GetDuration(project.StartDate.Format("2006-01-02"), project.EndDate.Format("2006-01-02"))
		project.Description = helper.CutString(project.Description, 30)
		project.Name = helper.CutString(project.Name, 20)

		dataResult = append(dataResult, project)
	}

	var view, templErr = template.ParseFiles("views/index.html", "views/layout/navigation.html");
	if templErr != nil {
		fmt.Println(templErr.Error())
		return
	}	
	store := sessions.NewCookieStore([]byte("MY_SESSION_KEY"))
	session, _ := store.Get(r,"MY_SESSION_KEY");
	if session.Values["Id"] == nil {
		session.Values["Id"] = 0
	} 
	if session.Values["IsLogin"] == nil {
		session.Values["IsLogin"] = false
	} 
	if session.Values["Name"] == nil {
		session.Values["Name"] = ""
	} 
	if session.Values["Email"] == nil{
		session.Values["Email"] = ""
	}
	

	flashMessage := session.Flashes("login message")

	var flashArr []string

	if len(flashMessage) > 0 {
		session.Save(r,w)
		for _, f := range flashMessage {
			flashArr = append(flashArr, f.(string))
		}
	}
	
	sessionData := SessionStruct{
		Id: session.Values["Id"].(int),
		Name: session.Values["Name"].(string),
		Email: session.Values["Email"].(string),
		IsLogin: session.Values["IsLogin"].(bool),
	}
	sessionData.Flash = strings.Join(flashArr,"")
	responseData := map[string]interface{} {
		"blogData": dataResult,
		"sessionData": sessionData,
	}
	view.Execute(w, responseData)
}

func GetAddProject(w http.ResponseWriter, r *http.Request) {
	var view, err = template.ParseFiles("views/project.html", "views/layout/navigation.html")	
	if err != nil {
		panic(err.Error())
	}
	
	store := sessions.NewCookieStore([]byte("MY_SESSION_KEY"))
	session, _ := store.Get(r,"MY_SESSION_KEY");
	if session.Values["IsLogin"] == nil {
		session.Values["IsLogin"] = false
	} 
	if session.Values["Name"] == nil {
		session.Values["Name"] = ""
	} 
	if session.Values["Email"] == nil{
		session.Values["Email"] = ""
	}
	sessionData := SessionStruct{
		Name: session.Values["Name"].(string),
		Email: session.Values["Email"].(string),
		IsLogin: session.Values["IsLogin"].(bool),
	}
	responseData := map[string]interface{} {
		"sessionData": sessionData,
	}
	view.Execute(w, responseData)
}

func PostAddProject(w http.ResponseWriter, r *http.Request) {
	parseErr := r.ParseForm()
	if parseErr != nil {
		fmt.Println(parseErr)
	}

	uploadContext := r.Context().Value("dataFile")
	image := uploadContext.(string)

	name := r.PostForm.Get("name")
	description := r.PostForm.Get("description")
	startDate := r.PostForm.Get("start-date")
	endDate := r.PostForm.Get("end-date")
	techlist := r.PostForm["checkbox"]

	store := sessions.NewCookieStore([]byte("MY_SESSION_KEY"))
	session, _ := store.Get(r,"MY_SESSION_KEY");
	sessionData := SessionStruct{
		Name: session.Values["Name"].(string),
		Id: session.Values["Id"].(int),
		Email: session.Values["Email"].(string),
		IsLogin: session.Values["IsLogin"].(bool),
	}

	query := `INSERT INTO public.tb_projects(
		name, start_date, end_date, description, technologies, image, author_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`

	_ , err := connection.Conn.Exec(context.Background(),query, name,startDate,endDate,description,techlist,image,sessionData.Id)

	if err != nil {
		panic(err.Error())
	}
	http.Redirect(w,r,"/form-add-project", http.StatusFound)
}

func GetEditProject(w http.ResponseWriter, r *http.Request) {
	projectId,idErr := strconv.Atoi(mux.Vars(r)["index"])
	if idErr != nil {
		panic(idErr.Error())
	}

	queryString := `
	SELECT id,name,start_date,end_date,description,technologies,image FROM public.tb_projects WHERE id = ($1)`

	data, err := connection.Conn.Query(context.Background(), queryString, projectId);

	if err != nil {
		panic(err.Error())
	}

	var project = ProjectData{};
	
	for data.Next() {
		err := data.Scan(&project.Id,&project.Name,&project.StartDate,&project.EndDate,&project.Description, &project.Technologies,&project.Image);
		if err != nil {
			fmt.Println(err.Error())
		}
		project.StringStartDate = project.StartDate.Format("2006-01-02")
		project.StringEndDate = project.EndDate.Format("2006-01-02")
	} 
	
	var view, viewErr = template.ParseFiles("views/edit-project.html", "views/layout/navigation.html")

	if viewErr != nil {
		panic(viewErr.Error())
	}
	
	store := sessions.NewCookieStore([]byte("MY_SESSION_KEY"))
	session, _ := store.Get(r,"MY_SESSION_KEY");
	if session.Values["IsLogin"] == nil {
		session.Values["IsLogin"] = false
	} 
	if session.Values["Name"] == nil {
		session.Values["Name"] = ""
	} 
	if session.Values["Email"] == nil{
		session.Values["Email"] = ""
	}
	sessionData := SessionStruct{
		Name: session.Values["Name"].(string),
		Email: session.Values["Email"].(string),
		IsLogin: session.Values["IsLogin"].(bool),
	}
	responseData := map[string]interface{} {
		"project": project,
		"sessionData": sessionData,
	}
	view.Execute(w, responseData)
}

func UpdateProject(w http.ResponseWriter, r *http.Request) {
	parseErr := r.ParseForm();
	if parseErr != nil {
		panic(parseErr.Error())
	}
	projectId,idErr := strconv.Atoi(mux.Vars(r)["index"]);
	if idErr != nil {
		panic(idErr.Error())
	}

	name := r.PostForm.Get("name")
	startDate := r.PostForm.Get("start-date")
	endDate := r.PostForm.Get("end-date")
	description := r.PostForm.Get("description")
	checkbox := r.PostForm["checkbox"]
	
	queryString2 := `
		UPDATE public.tb_projects
		SET name=$1, start_date=$2, end_date=$3, description=$4, technologies=array_cat(technologies, $5)
		WHERE id = ($6)
	`

	_, queryErr := connection.Conn.Exec(context.Background(),queryString2,name,startDate,endDate,description,checkbox ,projectId)
	
	if queryErr != nil {
		panic(queryErr.Error())
	}

	
	
	http.Redirect(w,r,"/",http.StatusFound)
}

func GetProjectDetail(w http.ResponseWriter, r *http.Request) {
	projectId, indexError := strconv.Atoi(mux.Vars(r)["projectId"]);
	if indexError != nil {
		panic(indexError.Error())
	}

	queryString := `
	SELECT tb_projects.id,tb_projects.name,tb_projects.start_date,tb_projects.end_date,tb_projects.description,
	tb_projects.technologies,tb_projects.image, tb_user.name as authorName
	FROM public.tb_projects LEFT JOIN tb_user ON tb_projects.author_id = tb_user.id WHERE tb_projects.id = $1
	`
	data, err := connection.Conn.Query(context.Background(),queryString, projectId )

	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	var project = ProjectData{}
	for data.Next() {
		// var project = ProjectData{}

		var scanErr = data.Scan(&project.Id,&project.Name, &project.StartDate, &project.EndDate, &project.Description, &project.Technologies, &project.Image, &project.AuthorName)
		if scanErr != nil {
			fmt.Println(scanErr.Error())
			return
		}
		project.Duration = helper.GetDuration(project.StartDate.Format("2006-01-02"), project.EndDate.Format("2006-01-02"))
		project.StringStartDate = project.StartDate.Format("January 02, 2006")
		project.StringEndDate = project.EndDate.Format("January 02, 2006")
	}

	var view,viewErr = template.ParseFiles("views/projectDetail.html", "views/layout/navigation.html")
	if viewErr != nil {
		panic(viewErr.Error())
	}
	store := sessions.NewCookieStore([]byte("MY_SESSION_KEY"))
	session, _ := store.Get(r,"MY_SESSION_KEY");
	if session.Values["IsLogin"] == nil {
		session.Values["IsLogin"] = false
	} 
	if session.Values["Name"] == nil {
		session.Values["Name"] = ""
	} 
	if session.Values["Email"] == nil{
		session.Values["Email"] = ""
	}
	sessionData := SessionStruct{
		Name: session.Values["Name"].(string),
		Email: session.Values["Email"].(string),
		IsLogin: session.Values["IsLogin"].(bool),
	}
	responseData := map[string]interface{} {
		"project": project,
		"sessionData": sessionData,
	}
	view.Execute(w, responseData)

}

func DeleteProject(w http.ResponseWriter, r *http.Request) {
	projectId, idErr := strconv.Atoi(mux.Vars(r)["projectId"]);

	if idErr != nil {
		panic(idErr.Error())
	}

	queryString := `
		DELETE FROM public.tb_projects WHERE id = $1
	`

	_, queryErr := connection.Conn.Exec(context.Background(), queryString, projectId)

	if queryErr != nil {
		panic(queryErr.Error())
	}

	http.Redirect(w,r,"/",http.StatusFound)
}