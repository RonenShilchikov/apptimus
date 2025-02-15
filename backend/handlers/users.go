// handlers/users.go
package handlers
import (
	"database/sql"
	"encoding/json"
	"net/http"
	"log" 
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token,omitempty"`
}

func generateToken(userID int) string {
	return "mocked-jwt-token"
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		var credentials struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
			return
		}

		// ✅ Retrieve user from database
		var user User
		err = db.QueryRow("SELECT id, username, email FROM users WHERE email = ? AND password_hash = ?", credentials.Email, credentials.Password).
			Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			http.Error(w, `{"error": "Invalid email or password"}`, http.StatusUnauthorized)
			return
		}

		// ✅ Generate a JWT token
		user.Token = generateToken(user.ID)

		// ✅ Store the token in the database
		_, err = db.Exec("UPDATE users SET token = ? WHERE id = ?", user.Token, user.ID)
		if err != nil {
			http.Error(w, `{"error": "Failed to update user token"}`, http.StatusInternalServerError)
			return
		}

		// ✅ Return user with token
		json.NewEncoder(w).Encode(user)
	}
}


func SignupHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            log.Printf("Invalid method: %s", r.Method)
            http.Error(w, `{"error": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
            return
        }

        w.Header().Set("Content-Type", "application/json")

        var newUser struct {
            Name     string `json:"name"`
            Email    string `json:"email"`
            Password string `json:"password"`
        }
        
        err := json.NewDecoder(r.Body).Decode(&newUser)
        if err != nil {
            log.Printf("Failed to decode request body: %v", err)
            http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
            return
        }

        log.Printf("Attempting to create user: %s with email: %s", newUser.Name, newUser.Email)

        // Check if user already exists
        var exists bool
        err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", newUser.Email).Scan(&exists)
        if err != nil {
            log.Printf("Database error checking existence: %v", err)
            http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
            return
        }

        if exists {
            log.Printf("User with email %s already exists", newUser.Email)
            http.Error(w, `{"error": "Email already registered"}`, http.StatusConflict)
            return
        }

        // Insert new user
        result, err := db.Exec(
            "INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)", 
            newUser.Name, newUser.Email, newUser.Password,
        )
        if err != nil {
            log.Printf("Failed to insert user: %v", err)
            http.Error(w, `{"error": "Failed to create user"}`, http.StatusInternalServerError)
            return
        }

        userID, _ := result.LastInsertId()
        log.Printf("Successfully created user with ID: %d", userID)

        user := User{
            ID: int(userID),
            Name: newUser.Name,
            Email: newUser.Email,
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(user)
    }
}

func GetUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		rows, err := db.Query("SELECT id, username, email FROM users")
		if err != nil {
			http.Error(w, `{"error": "Error fetching users"}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Name, &user.Email)
			if err != nil {
				http.Error(w, `{"error": "Error reading users"}`, http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}

		json.NewEncoder(w).Encode(users)
	}
}

func DeleteUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, `{"error": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		userID := r.URL.Query().Get("id")
		if userID == "" {
			http.Error(w, `{"error": "User ID is required"}`, http.StatusBadRequest)
			return
		}

		result, err := db.Exec("DELETE FROM users WHERE id = ?", userID)
		if err != nil {
			http.Error(w, `{"error": "Error deleting user"}`, http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
	}
}