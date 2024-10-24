package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

// 회원가입 핸들러 (API)
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// POST 요청이 아닌 경우 오류 처리
		if r.Method != http.MethodPost {
			http.Error(w, "잘못된 요청 방법", http.StatusMethodNotAllowed)
			return
		}

		var reqBody struct {
			UserID   string `json:"userid"`
			Password string `json:"password"`
		}

		// JSON 데이터 파싱
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			log.Println("JSON 파싱 실패:", err) // 오류 로그 추가
			http.Error(w, "잘못된 요청 본문", http.StatusBadRequest)
			return
		}

		// 사용자 ID가 이미 존재하는지 확인
		var existingUserID string
		err := db.QueryRow("SELECT userid FROM users WHERE userid = ?", reqBody.UserID).Scan(&existingUserID)
		if err != nil && err != sql.ErrNoRows {
			log.Println("데이터베이스 조회 실패:", err) // 오류 로그 추가
			http.Error(w, "데이터베이스 오류", http.StatusInternalServerError)
			return
		}

		if existingUserID != "" {
			http.Error(w, "이미 존재하는 사용자 ID입니다", http.StatusConflict)
			return
		}

		// 비밀번호 해싱
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("비밀번호 해싱 실패:", err) // 오류 로그 추가
			http.Error(w, "비밀번호 해싱 실패", http.StatusInternalServerError)
			return
		}

		// 사용자 등록 (userid와 password만 사용)
		_, err = db.Exec("INSERT INTO users(userid, password) VALUES(?, ?)", reqBody.UserID, string(hashedPassword))
		if err != nil {
			log.Println("사용자 등록 실패:", err) // 오류 로그 추가
			http.Error(w, "사용자 등록 실패", http.StatusInternalServerError)
			return
		}

		// 성공 응답
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "회원가입 성공", "userid": reqBody.UserID})
	}
}

// 로그인 핸들러 (API)
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			session, _ := store.Get(r, "session")

			var reqBody struct {
				UserID   string `json:"userid"`
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
				UserID   string
				Password string
			}

			// 데이터베이스에서 사용자 정보를 조회
			err = db.QueryRow("SELECT id, userid, password FROM users WHERE userid = ?", reqBody.UserID).Scan(&user.ID, &user.UserID, &user.Password)
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

			// 세션에 사용자 ID 저장 (DB의 ID 대신 userid 저장)
			session.Values["user_id"] = user.UserID
			session.Save(r, w)

			// 로그인 성공 JSON 응답
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "로그인 성공",
				"user_id": user.UserID,
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

// 로그인한 사용자 정보를 반환하는 핸들러 (API)
func GetUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		userID, ok := session.Values["user_id"]
		if !ok {
			http.Error(w, "로그인되지 않은 사용자입니다.", http.StatusUnauthorized)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id": userID,
		})
	}
}
