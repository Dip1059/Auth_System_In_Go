package Repositories

import (
	M "Auth_System_In_Go/Models"
	G "Auth_System_In_Go/Globals"
	"database/sql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type DB_ENV struct {
	Host, Port, Dialect, Username, Password, DBname string
}


func init() {
	godotenv.Load()
	G.DbEnv = G.DB_ENV{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Dialect:  os.Getenv("DB_DIALECT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBname:   os.Getenv("DB_NAME"),
	}
}

func DBConnect() (*sql.DB, error) {
	db, _ := sql.Open(G.DbEnv.Dialect, G.DbEnv.Username+":"+G.DbEnv.Password+"@tcp("+G.DbEnv.Host+":"+G.DbEnv.Port+")/"+G.DbEnv.DBname+"?parseTime=true")
	return db, nil
}

func ReadWithEmail(user M.User) (M.User, bool) {
	db, _ := DBConnect()
	var results *sql.Rows
	var err error
	results, err = db.Query("SELECT * FROM users WHERE email=?;", user.Email)
	if results.Next() {
		err = results.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Password, &user.ActiveStatus, &user.Role, &user.EmailVf, &user.RememberToken, &user.Deleted_At, &user.Created_At, &user.Updated_At)
		if err != nil {
			log.Println("AuthRepo.go Log1", err.Error())
		}
		return user, true
	} else {
		return user, false
	}

	defer db.Close()
	defer results.Close()
	return user, true
}

func Register(user M.User) (M.User, bool) {
	db, _ := DBConnect()
	var results *sql.Rows
	var err error
	var success bool

	_, err = db.Query("INSERT INTO users(full_name, email, password, role_id, email_verification) VALUES(?, ?, ?, ?, ?);", user.Name, user.Email, user.Password, user.Role, user.EmailVf)
	if err != nil {
		log.Println("AuthRepo.go Log2", err.Error())
		return user, false
	}
	user, success = ReadWithEmail(user)
	if success {
		return user, true
	} else {
		return user, false
	}

	log.Println("AuthRepo.go Log3 Data Inserterd Successfully.\n")
	defer db.Close()
	defer results.Close()
	return user, true
}

func Login(user M.User) (M.User, bool) {
	var success bool
	user, success = ReadWithEmail(user)
	if success {
		return user, true
	} else {
		return user, false
	}
}


func SetRememberToken(user M.User) bool {
	db, _ := DBConnect()

	results, err := db.Query("UPDATE users set remember_token=? where email=?;", user.RememberToken, user.Email)
	if err != nil {
		log.Println("AuthRepo.go Log4", err.Error())
		return false
	}

	defer db.Close()
	defer results.Close()
	return true
}


func ActivateAccount(user M.User) (M.User, bool) {
	db, _ := DBConnect()
	var success bool

	results, err := db.Query("SELECT * FROM users WHERE email=? and email_verification=?;", user.Email, user.EmailVf.String)

	if results.Next() {
		results, err = db.Query("UPDATE users SET activestatus=1, email_verification=NULL WHERE email=? and email_verification=?;", user.Email, user.EmailVf.String)

		if err != nil {
			log.Println("AuthRepo.go Log5", err.Error())
			return user, false
		}

		user, success = ReadWithEmail(user)
		if success {
			return user, true
		} else {
			return user, false
		}
	} else {
		return user, false
	}

	defer db.Close()
	defer results.Close()
	return user, true
}


func SendPasswordResetLink(ps M.PasswordReset) bool {
	db, _ := DBConnect()

		results, err := db.Query("INSERT INTO password_resets(email,token) VALUES(?, ?);", ps.Email, ps.Token)
		if err != nil {
			log.Println("AuthRepo.go Log6", err.Error())
			return false
		}

	defer db.Close()
	defer results.Close()
	return true
}


func ResetPasswordGet(ps M.PasswordReset) bool {
	db, _ := DBConnect()

	results, err := db.Query("SELECT * from password_resets where email=? and token=? and status=0;", ps.Email, ps.Token)
	if err != nil {
		log.Println("AuthRepo.go Log7", err.Error())
		return false
	}
	if !results.Next() {
		return false
	}

	defer db.Close()
	defer results.Close()
	return true
}


func ResetPasswordPost(user M.User, ps M.PasswordReset) bool {
	db, _ := DBConnect()

	results, err := db.Query("UPDATE users SET password=? where email=?;", user.Password, user.Email)
	if err != nil {
		log.Println("AuthRepo.go Log8", err.Error())
		return false
	}

	results, err = db.Query("UPDATE password_resets SET status=1 where email=? and token=?;", ps.Email, ps.Token)
	if err != nil {
		log.Println("AuthRepo.go Log9", err.Error())
		return false
	}

	results, err = db.Query("UPDATE password_resets SET token=NULL where email=?;", ps.Email)
	if err != nil {
		log.Println("AuthRepo.go Log10", err.Error())
		return false
	}

	defer db.Close()
	defer results.Close()
	return true
}


func Logout(user M.User) {
	db, _ := DBConnect()

	results, err := db.Query("UPDATE users set remember_token=NULL where email=?;", user.Email)
	if err != nil {
		log.Println("AuthRepo.go Log11", err.Error())
		return
	}

	defer db.Close()
	defer results.Close()
	return
}
