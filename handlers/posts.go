package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

type Post struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	CreatedAt string `json:"created_at"`
}

// 게시판 리스트 가져오는 핸들러
func ListPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 쿼리 파라미터에서 page, limit, search 값 가져오기
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")
		searchQuery := r.URL.Query().Get("search")

		// 기본값 설정 (page = 1, limit = 10)
		page := 1
		limit := 10
		var err error

		// page와 limit이 쿼리 파라미터로 제공된 경우 변환
		if pageStr != "" {
			page, err = strconv.Atoi(pageStr)
			if err != nil || page < 1 {
				page = 1 // 잘못된 값이면 1페이지로 설정
			}
		}

		if limitStr != "" {
			limit, err = strconv.Atoi(limitStr)
			if err != nil || limit < 1 {
				limit = 10 // 잘못된 값이면 기본 limit 사용
			}
		}

		// 검색 쿼리 설정
		var rows *sql.Rows
		offset := (page - 1) * limit
		query := "SELECT id, title, content, author, created_at FROM posts"
		countQuery := "SELECT COUNT(*) FROM posts"
		params := []interface{}{}

		if searchQuery != "" {
			searchTerm := "%" + searchQuery + "%"
			query += " WHERE title LIKE ? OR content LIKE ?"
			countQuery += " WHERE title LIKE ? OR content LIKE ?"
			params = append(params, searchTerm, searchTerm)
		}

		// 전체 게시물 수 조회 (페이지네이션 계산을 위해)
		var totalRecords int
		err = db.QueryRow(countQuery, params...).Scan(&totalRecords)
		if err != nil {
			http.Error(w, "게시물 총 개수 조회 실패", http.StatusInternalServerError)
			return
		}

		// 페이지네이션을 위한 SQL 쿼리 추가
		query += " LIMIT ? OFFSET ?"
		params = append(params, limit, offset)

		// 게시물 조회
		rows, err = db.Query(query, params...)
		if err != nil {
			http.Error(w, "데이터 조회 실패", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts []Post
		for rows.Next() {
			var post Post
			err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.CreatedAt)
			if err != nil {
				http.Error(w, "데이터 파싱 실패", http.StatusInternalServerError)
				return
			}
			posts = append(posts, post)
		}

		// 총 페이지 수 계산
		totalPages := (totalRecords + limit - 1) / limit

		// 응답 데이터 생성
		response := map[string]interface{}{
			"posts":      posts,
			"totalPages": totalPages,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response) // JSON으로 응답
	}
}

// GetPostHandler - 특정 게시물 조회 핸들러
func GetPostHandler(db *sql.DB, postID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var post Post
		err := db.QueryRow("SELECT id, title, content, author, created_at FROM posts WHERE id = ?", postID).Scan(
			&post.ID, &post.Title, &post.Content, &post.Author, &post.CreatedAt,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Post not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post) // 게시물 정보를 JSON으로 반환
	}
}

// 새 글 저장 API 핸들러
func CreatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "잘못된 요청 방법입니다", http.StatusMethodNotAllowed)
			return
		}

		var post Post
		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			http.Error(w, "잘못된 요청 본문입니다", http.StatusBadRequest)
			return
		}

		if post.Title == "" || post.Content == "" || post.Author == "" {
			http.Error(w, "모든 필드를 입력해주세요", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("INSERT INTO posts (title, content, author, created_at) VALUES (?, ?, ?, NOW())", post.Title, post.Content, post.Author)
		if err != nil {
			http.Error(w, "게시물 저장 실패: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// 게시글 수정 페이지 API 핸들러
func EditPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		var post Post
		err := db.QueryRow("SELECT id, title, content, author FROM posts WHERE id = ?", id).Scan(&post.ID, &post.Title, &post.Content, &post.Author)
		if err != nil {
			http.Error(w, "게시물을 찾을 수 없습니다", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post) // JSON으로 응답
	}
}

// 게시글 수정 처리 API 핸들러
func UpdatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		var post Post
		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			http.Error(w, "잘못된 요청 본문입니다", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("UPDATE posts SET title=?, content=? WHERE id=?", post.Title, post.Content, id)
		if err != nil {
			http.Error(w, "게시물 수정 실패", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// 게시글 삭제 API 핸들러
func DeletePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")

		_, err := db.Exec("DELETE FROM posts WHERE id=?", id)
		if err != nil {
			http.Error(w, "게시물 삭제 실패", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
