package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

// 회원가입 핸들러 (API)
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "잘못된 요청 방법", http.StatusMethodNotAllowed)
			return
		}

		var reqBody struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		// JSON 데이터 파싱
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "잘못된 요청 본문", http.StatusBadRequest)
			return
		}

		// 사용자명이 이미 존재하는지 확인
		var existingUser string
		err := db.QueryRow("SELECT username FROM users WHERE username = ?", reqBody.Username).Scan(&existingUser)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, "데이터베이스 오류", http.StatusInternalServerError)
			return
		}

		if existingUser != "" {
			http.Error(w, "이미 존재하는 사용자명입니다", http.StatusConflict)
			return
		}

		// 비밀번호 해싱
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "비밀번호 해싱 실패", http.StatusInternalServerError)
			return
		}

		// 사용자 등록
		_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", reqBody.Username, string(hashedPassword))
		if err != nil {
			http.Error(w, "사용자 등록 실패", http.StatusInternalServerError)
			return
		}

		// JSON 응답
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "회원가입 성공", "username": reqBody.Username})
	}
}

// 로그인 핸들러 (API)
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			session, _ := store.Get(r, "session")

			var reqBody struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}

			// JSON 데이터 파싱
			err := json.NewDecoder(r.Body).Decode(&reqBody)
			if err != nil {
				http.Error(w, "잘못된 요청 본문", http.StatusBadRequest)
				return
			}

			var user struct {
				ID       int
				Username string
				Password string
			}

			err = db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", reqBody.Username).Scan(&user.ID, &user.Username, &user.Password)
			if err != nil {
				http.Error(w, "사용자를 찾을 수 없습니다", http.StatusUnauthorized)
				return
			}

			// 비밀번호 비교
			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password))
			if err != nil {
				http.Error(w, "비밀번호가 틀립니다", http.StatusUnauthorized)
				return
			}

			// 세션에 사용자 정보 저장
			session.Values["user_id"] = user.ID
			session.Save(r, w)

			// 로그인 성공 JSON 응답
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "로그인 성공",
				"user_id": user.ID,
			})
		} else {
			http.Error(w, "잘못된 요청 방법", http.StatusMethodNotAllowed)
		}
	}
}

// 로그아웃 핸들러 (API)
func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")

		// 세션에서 사용자 정보 삭제
		delete(session.Values, "user_id")
		session.Save(r, w)

		// 로그아웃 성공 JSON 응답
		json.NewEncoder(w).Encode(map[string]string{"message": "로그아웃 성공"})
	}
}

// 로그인 여부 확인 핸들러 (API)
func IsLoggedInHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		userID, ok := session.Values["user_id"]
		if ok {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"is_logged_in": true,
				"user_id":      userID,
			})
		} else {
			json.NewEncoder(w).Encode(map[string]bool{"is_logged_in": false})
		}
	}
}
