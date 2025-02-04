package main

import (
	"goProject1/config"   // DB 설정을 위한 config 패키지 import
	"goProject1/handlers" // handlers 패키지 import
	"html/template"
	"log"
	"net/http"
	"strings"
)

// 템플릿 캐싱
var tmpl = template.Must(template.ParseFiles(
	"templates/login.html",
	"templates/index.html",
	"templates/register.html",
	"templates/refresh-posts.html",
))

// 로그인 페이지 핸들러
func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// 로그인 페이지 렌더링
		tmpl.ExecuteTemplate(w, "login.html", nil)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// 인덱스 페이지 핸들러
func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// 인덱스 페이지 렌더링
		tmpl.ExecuteTemplate(w, "index.html", nil)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// 회원가입 페이지 핸들러
func registerPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// 회원가입 페이지 렌더링
		tmpl.ExecuteTemplate(w, "register.html", nil)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// 새로고침 게시판 핸들러
func refreshPostsPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// 회원가입 페이지 렌더링
		tmpl.ExecuteTemplate(w, "refresh-posts.html", nil)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	// DB 연결
	db, err := config.ConnectDB() // config 패키지에서 DB 연결 설정 호출
	if err != nil {
		log.Fatal("DB 연결 실패:", err)
	}
	defer db.Close()

	// 정적 파일(css, js, images 등) 제공 설정
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 기본 경로에서 index.html 파일 템플릿 서빙
	http.HandleFunc("/", indexPageHandler)

	// 로그인 페이지 핸들러 추가
	http.HandleFunc("/login", loginPageHandler)

	// 회원가입 페이지 핸들러 추가
	http.HandleFunc("/register", registerPageHandler)

	//새로고침 게시판
	http.HandleFunc("/refresh-posts", refreshPostsPageHandler)

	// 핸들러 설정 (ServeMux 사용하여 경로와 메서드 처리)
	http.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.ListPostsHandler(db)(w, r) // 게시물 목록 조회
		case http.MethodPost:
			handlers.CreatePostHandler(db)(w, r) // 게시물 생성
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/posts/edit", handlers.EditPostHandler(db))     // 게시물 수정 조회
	http.HandleFunc("/api/posts/update", handlers.UpdatePostHandler(db)) // 게시물 수정 처리
	http.HandleFunc("/api/posts/delete", handlers.DeletePostHandler(db)) // 게시물 삭제 처리

	// 특정 게시물 조회 핸들러 추가 (GET 요청 처리)
	http.HandleFunc("/api/posts/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// 게시물 ID 추출 (/api/posts/{id}에서 id 추출)
			id := strings.TrimPrefix(r.URL.Path, "/api/posts/")
			if id == "" {
				http.Error(w, "Missing post ID", http.StatusBadRequest)
				return
			}
			handlers.GetPostHandler(db, id)(w, r) // 특정 게시물 조회
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/register", handlers.RegisterHandler(db))     // 회원가입
	http.HandleFunc("/api/login", handlers.LoginHandler(db))           // 로그인
	http.HandleFunc("/api/is_logged_in", handlers.IsLoggedInHandler()) // 로그인 여부 확인
	http.HandleFunc("/api/get_user", handlers.GetUserHandler())        // 로그인 계정 정보 조회

	// 서버 실행
	log.Println("Server started at :1000")
	log.Fatal(http.ListenAndServe(":1000", nil))
}
