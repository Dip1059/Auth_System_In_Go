package Middlewares

import (
	G "Auth_System_In_Go/Globals"
	Mod "Auth_System_In_Go/Models"
	R "Auth_System_In_Go/Repositories"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"net/http"
)


func IsGuest(c *gin.Context, store *sessions.CookieStore) (Mod.User, bool) {
	var user Mod.User
	session, _ := store.Get(c.Request, "login_token")
	email := session.Values["userEmail"]
	rememberToken := session.Values["remember_token"]
	var success bool
	if email != nil && rememberToken != nil {
		user.Email = session.Values["userEmail"].(string)
		user, success = R.ReadWithEmail(user)
		if !success {
			c.Redirect(http.StatusFound, "/logout")
			return user, false
		}
		if rememberToken.(string) != user.RememberToken.String {
			c.Redirect(http.StatusFound, "/logout")
			return user, false
		}
		if user.ActiveStatus == 2 {
			c.Redirect(http.StatusFound, "/logout")
		} else if user.ActiveStatus == 1 && user.Role == 1 {
			c.Redirect(http.StatusFound, "/dashboard")
		} else if user.ActiveStatus == 1 && user.Role == 2 {
			c.Redirect(http.StatusFound, "/home")
		}
		return user, false
	}
	return user, true
}


func IsAuthUser(c *gin.Context, store *sessions.CookieStore) (Mod.User, bool) {
	var user Mod.User
	session, _ := store.Get(c.Request, "login_token")
	email := session.Values["userEmail"]
	rememberToken := session.Values["remember_token"]
	var success bool
	if email != nil && rememberToken != nil {
		user.Email = session.Values["userEmail"].(string)
		user, success = R.ReadWithEmail(user)
		if !success {
			G.Msg.Fail = "User Doesn't Exist Anymore."
			c.Redirect(http.StatusFound, "/logout")
			return user, false
		}
		if rememberToken.(string) != user.RememberToken.String {
			//G.Msg.Fail = "Someone Stole Your Cookie From Your Browser. Please Be Cautious."
			c.Redirect(http.StatusFound, "/logout")
			return user, false
		}
		if user.ActiveStatus == 2 {
			G.Msg.Fail = "You Are Suspended. Contact With The Authority Quickly."
			c.Redirect(http.StatusFound, "/logout")
		} else if user.ActiveStatus == 1 && user.Role == 1 {
			c.Redirect(http.StatusFound, "/dashboard")
		} else if user.ActiveStatus == 1 && user.Role == 2 {
			return user, true
		}
		return user, false
	}
	c.Redirect(http.StatusFound, "/login")
	return user, false
}


func IsAuthAdminUser(c *gin.Context, store *sessions.CookieStore) (Mod.User, bool) {
	var user Mod.User
	session, _ := store.Get(c.Request, "login_token")
	email := session.Values["userEmail"]
	rememberToken := session.Values["remember_token"]
	var success bool

	if email != nil && rememberToken != nil {
		user.Email = session.Values["userEmail"].(string)
		user, success = R.ReadWithEmail(user)
		if !success {
			G.Msg.Fail = "User Doesn't Exist Anymore."
			c.Redirect(http.StatusFound, "/logout")
			return user, false
		}
		if rememberToken.(string) != user.RememberToken.String {
			//G.Msg.Fail = "Someone Stole Your Cookie From Your Browser. Please Be Cautious."
			c.Redirect(http.StatusFound, "/logout")
			return user, false
		}
		if user.ActiveStatus == 2 {
			G.Msg.Fail = "You Are Suspended. Contact With The Authority Quickly."
			c.Redirect(http.StatusFound, "/logout")
		} else if user.ActiveStatus == 1 && user.Role == 2 {
			c.Redirect(http.StatusFound, "/home")
		} else if user.ActiveStatus == 1 && user.Role == 1 {
			return user, true
		}
		return user, false
	}
	c.Redirect(http.StatusFound, "/login")
	return user, false
}
