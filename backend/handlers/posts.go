package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"fmt"
)

// ✅ Updated Post Struct (Stores Author Name Instead of `UserID`)
type Post struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Author string `json:"author"`
}

// ✅ Get All Posts (Sorted by Latest)
func GetPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		rows, err := db.Query("SELECT id, title, body, author FROM posts ORDER BY id DESC")
		if err != nil {
			http.Error(w, `{"error": "Error fetching posts"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts []Post
		for rows.Next() {
			var post Post
			err := rows.Scan(&post.ID, &post.Title, &post.Body, &post.Author)
			if err != nil {
				http.Error(w, `{"error": "Error reading posts"}`, http.StatusInternalServerError)
				return
			}
			posts = append(posts, post)
		}

		// ✅ Ensure correct response structure
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(posts)
	}
}


// ✅ Create a New Post (Only Logged-in Users)
func CreatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			http.Error(w, `{"error": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		// ✅ Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		fmt.Println("🔹 Received Authorization Header:", authHeader)

		// ✅ Check if token exists
		if authHeader == "" {
			fmt.Println("❌ No Authorization Header Found")
			http.Error(w, `{"error": "Unauthorized. No token provided."}`, http.StatusUnauthorized)
			return
		}

		// ✅ Ensure correct token format (Bearer token)
		var userToken string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			userToken = authHeader[7:]
		} else {
			fmt.Println("❌ Invalid Token Format:", authHeader)
			http.Error(w, `{"error": "Invalid token format"}`, http.StatusUnauthorized)
			return
		}

		fmt.Println("🔹 Extracted Token:", userToken)

		// ✅ Verify user token in the database
		var userID int
		var author string
		err := db.QueryRow("SELECT id, username FROM users WHERE token = ?", userToken).Scan(&userID, &author)
		if err != nil {
			fmt.Println("❌ Invalid Token / User Not Found:", err)
			http.Error(w, `{"error": "Invalid token. Please log in again."}`, http.StatusUnauthorized)
			return
		}

		fmt.Println("✅ Authenticated User:", userID, author)

		// ✅ Decode request body
		var post Post
		err = json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			fmt.Println("❌ Invalid Request Format:", err)
			http.Error(w, `{"error": "Invalid request format"}`, http.StatusBadRequest)
			return
		}

		// ✅ Ensure fields are not empty
		if post.Title == "" || post.Body == "" {
			fmt.Println("❌ Missing Fields:", post)
			http.Error(w, `{"error": "Title and body are required"}`, http.StatusBadRequest)
			return
		}

		// ✅ Insert post into database
		_, err = db.Exec("INSERT INTO posts (title, body, author, author_id) VALUES (?, ?, ?, ?)", post.Title, post.Body, author, userID)
		if err != nil {
			fmt.Println("❌ Database Insert Error:", err)
			http.Error(w, `{"error": "Error creating post"}`, http.StatusInternalServerError)
			return
		}

		// ✅ Return success response
		fmt.Println("✅ Post Created Successfully!")
		json.NewEncoder(w).Encode(map[string]string{"message": "Post created successfully"})
	}
}




// ✅ Edit Post (Only Author Can Edit)
func EditPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPut {
			http.Error(w, `{"error": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		// ✅ Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "Unauthorized. Please log in."}`, http.StatusUnauthorized)
			return
		}

		// ✅ Extract the token from "Bearer <token>" format
		var token string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			http.Error(w, `{"error": "Invalid token format"}`, http.StatusUnauthorized)
			return
		}

		// ✅ Get user ID from token
		var userID int
		err := db.QueryRow("SELECT id FROM users WHERE token = ?", token).Scan(&userID)
		if err != nil {
			http.Error(w, `{"error": "Invalid token. Please log in again."}`, http.StatusUnauthorized)
			return
		}

		// ✅ Get Post ID from query parameters
		postID := r.URL.Query().Get("id")
		if postID == "" {
			http.Error(w, `{"error": "Post ID is required"}`, http.StatusBadRequest)
			return
		}

		// ✅ Decode request body
		var post struct {
			Title string `json:"title"`
			Body  string `json:"body"`
		}
		err = json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			http.Error(w, `{"error": "Invalid request format"}`, http.StatusBadRequest)
			return
		}

		// ✅ Validate fields
		if post.Title == "" || post.Body == "" {
			http.Error(w, `{"error": "Title and body are required"}`, http.StatusBadRequest)
			return
		}

		// ✅ Check if the post exists and belongs to the user
		var authorID int
		err = db.QueryRow("SELECT author_id FROM posts WHERE id = ?", postID).Scan(&authorID)
		if err != nil {
			http.Error(w, `{"error": "Post not found"}`, http.StatusNotFound)
			return
		}

		// ✅ Ensure logged-in user is the author
		if authorID != userID {
			http.Error(w, `{"error": "Unauthorized. You can only edit your own posts."}`, http.StatusForbidden)
			return
		}

		// ✅ Update the post
		_, err = db.Exec("UPDATE posts SET title = ?, body = ? WHERE id = ?", post.Title, post.Body, postID)
		if err != nil {
			http.Error(w, `{"error": "Error updating post"}`, http.StatusInternalServerError)
			return
		}

		// ✅ Success Response
		json.NewEncoder(w).Encode(map[string]string{"message": "Post updated successfully"})
	}
}





// **Delete a Post Handler**
func DeletePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodDelete {
			http.Error(w, `{"error": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		// ✅ Extract post ID from query parameters
		postID := r.URL.Query().Get("id")
		if postID == "" {
			fmt.Println("❌ Missing Post ID")
			http.Error(w, `{"error": "Post ID is required"}`, http.StatusBadRequest)
			return
		}

		fmt.Println("🔹 Attempting to delete post with ID:", postID)

		// ✅ Check if the post exists before deleting
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postID).Scan(&exists)
		if err != nil {
			fmt.Println("❌ Database Error:", err)
			http.Error(w, `{"error": "Error checking post existence"}`, http.StatusInternalServerError)
			return
		}
		if !exists {
			fmt.Println("❌ Post Not Found:", postID)
			http.Error(w, `{"error": "Post not found"}`, http.StatusNotFound)
			return
		}

		// ✅ Delete the post
		result, err := db.Exec("DELETE FROM posts WHERE id = ?", postID)
		if err != nil {
			fmt.Println("❌ Error Deleting Post:", err)
			http.Error(w, `{"error": "Error deleting post"}`, http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			fmt.Println("❌ No Rows Affected. Post might not exist.")
			http.Error(w, `{"error": "Post not found"}`, http.StatusNotFound)
			return
		}

		fmt.Println("✅ Post Deleted Successfully:", postID)
		json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted successfully"})
	}
}


