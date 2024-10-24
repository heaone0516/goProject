package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type Post struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	CreatedAt string `json:"created_at"`
}

// 게시물 목록 API 핸들러
func ListPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		searchQuery := r.URL.Query().Get("search")
		var rows *sql.Rows
		var err error

		if searchQuery != "" {
			query := "%" + searchQuery + "%"
			rows, err = db.Query("SELECT id, title, content, author, created_at FROM posts WHERE title LIKE ? OR content LIKE ?", query, query)
		} else {
			rows, err = db.Query("SELECT id, title, content, author, created_at FROM posts")
		}

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

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts) // JSON으로 응답
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
