package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
)

// JWT 비밀 키
var jwtKey = []byte("super-secret-key")

// JWT claims 구조체 정의
type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// JWT 토큰 생성 함수
func GenerateJWT(userID string) (string, error) {
	// 토큰 만료 시간 (예: 72시간 후)
	expirationTime := time.Now().Add(72 * time.Hour)

	// Claims 설정
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// 토큰 생성 및 서명
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// JWT 검증 함수
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return claims, nil
}

// 로그인 핸들러 (API)
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
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

			// JWT 토큰 생성
			token, err := GenerateJWT(user.UserID)
			if err != nil {
				http.Error(w, "토큰 생성 실패", http.StatusInternalServerError)
				return
			}

			// 로그인 성공 시 JWT 응답
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "로그인 성공",
				"token":   token,
			})
		} else {
			http.Error(w, "잘못된 요청 방법", http.StatusMethodNotAllowed)
		}
	}
}

// 로그아웃은 JWT에서는 보통 클라이언트 측에서 토큰을 제거하는 방식으로 처리되므로 서버 측 로그아웃 함수는 필요하지 않음

// 인증 미들웨어 (JWT 검증)
func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authorization 헤더에서 토큰 가져오기
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "토큰이 없습니다", http.StatusUnauthorized)
			return
		}

		// Bearer 제거
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// JWT 토큰 검증
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "유효하지 않은 토큰입니다", http.StatusUnauthorized)
			return
		}

		// 인증이 성공하면 요청을 계속 처리
		r.Header.Set("user_id", claims.UserID) // 사용자 ID를 요청 헤더에 저장 (선택사항)
		next.ServeHTTP(w, r)
	}
}

// 로그인 여부 확인 핸들러 (API)
func IsLoggedInHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			json.NewEncoder(w).Encode(map[string]bool{"is_logged_in": false})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ValidateJWT(tokenString)
		if err != nil || claims == nil {
			json.NewEncoder(w).Encode(map[string]bool{"is_logged_in": false})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"is_logged_in": true,
			"user_id":      claims.UserID,
		})
	}
}

// 로그인한 사용자 정보를 반환하는 핸들러 (API)
func GetUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "토큰이 없습니다", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "유효하지 않은 토큰입니다", http.StatusUnauthorized)
			return
		}

		// 성공적으로 토큰이 검증되면 사용자 ID 반환
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id": claims.UserID,
		})
	}
}

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
			// JSON으로 데이터베이스 오류 반환
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error_code": "DB_ERROR", "message": "데이터베이스 오류"})
			return
		}

		if existingUserID != "" {
			// JSON으로 사용자 이미 존재하는 오류 반환
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error_code": "USER_EXISTS", "message": "이미 존재하는 사용자 ID입니다"})
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
