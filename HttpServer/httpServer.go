package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"unicode"
)

func main() {
	dbConnect(dbUser, dbPass, dbName)
	http.HandleFunc("/", welcomeHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
}

type userInfo struct {
	Username  string
	password  string
	password2 string
	Email     string
}

type responseMsg struct {
	success bool
	Message string
	Heading string
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("welcomeform.html")
		check(err)
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		myUser := userInfo{}
		myUser.Username = r.Form.Get("username")
		myUser.password = r.Form.Get("password")
		myUser.password2 = r.Form.Get("password2")
		myUser.Email = r.Form.Get("email")
		fmt.Println(myUser)
		response := checkUserInfo(myUser)
		if response.success {
			response.Heading = "Weclome to Discord chat rooms"
			t, err := template.ParseFiles("welcomeresponse.html")
			check(err)
			t.Execute(w, response)
			//Add the user in the database
			addUser(myUser)
		} else {
			response.Heading = "Something went wrong.Please go back to fix it"
			t, err := template.ParseFiles("welcomeresponse.html")
			check(err)
			t.Execute(w, response)
		}
	}
}

//If s contains an uppercase letter returns true
func upperExists(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) && unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

//If s contains a lowercase letter returns true
func lowerExists(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) && unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

//If s contains a special char returns true
func hasSpecialChar(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

//Check if password is strong enough
func checkPassword(s string) bool {
	a := upperExists(s)
	b := lowerExists(s)
	c := hasSpecialChar(s)
	if a && b && c {
		return true
	}
	return false

}

//Checks if user's info are correct
func checkUserInfo(info userInfo) responseMsg {
	r := responseMsg{}
	if info.password != info.password2 {
		r.success = false
		r.Message = "Passwords do not match"
		return r
	}

	if checkPassword(info.password) == false {
		r.success = false
		r.Message = "Password is not strong enough.It should contain at least one capital letter and one special character"
		return r
	}
	//Check if a user with this username already exists
	if getUser(info.Username).Username != "" {
		r.success = false
		r.Message = "A user with this username already exists!"
		return r
	}
	r.success = true
	r.Message = "Open your terminal and join a room!"
	return r
}
