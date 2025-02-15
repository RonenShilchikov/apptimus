package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
)

type Post struct {
	ID     int
	Title  string
	Body   string
	Author string
}

func main() {
	// Database connection
	db, err := sql.Open("mysql", "apptimus:1q2w3e@tcp(mysql_db:3306)/apptimus_db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Fetch posts
	rows, err := db.Query("SELECT id, title, body, author FROM posts")
	if err != nil {
		log.Fatalf("Failed to fetch posts: %v", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Body, &post.Author); err != nil {
			log.Fatalf("Failed to scan post: %v", err)
		}
		posts = append(posts, post)
	}

	// Create output directory
	outputDir := "/app/out"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Create funcMap before parsing templates
	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	// HTML Template
	tmpl, err := template.New("post").Funcs(funcMap).Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .Title }} - Apptimus Blog</title>
    <link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
</head>
<body class="bg-gray-100">
    <div class="container mx-auto px-4 py-8">
        <article class="max-w-2xl mx-auto bg-white p-6 rounded-lg shadow-md">
            <h1 class="text-3xl font-bold mb-4">{{ .Title }}</h1>
            <p class="text-gray-600 mb-4">By {{ .Author }}</p>
            <div class="prose">
                {{ .Body | safeHTML }}
            </div>
        </article>
    </div>
</body>
</html>
`)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	// Generate static pages
	for _, post := range posts {
		filename := filepath.Join(outputDir, fmt.Sprintf("post-%d.html", post.ID))
		f, err := os.Create(filename)
		if err != nil {
			log.Fatalf("Failed to create file %s: %v", filename, err)
		}
		defer f.Close()

		if err := tmpl.Execute(f, post); err != nil {
			log.Fatalf("Failed to write template: %v", err)
		}
	}

	// Generate index page
	indexTmpl, err := template.New("index").Funcs(funcMap).Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Apptimus Blog - All Posts</title>
    <link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
</head>
<body class="bg-gray-100">
    <div class="container mx-auto px-4 py-8">
        <h1 class="text-4xl font-bold mb-8 text-center">Apptimus Blog</h1>
        <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            {{range .}}
            <div class="bg-white p-6 rounded-lg shadow-md">
                <h2 class="text-2xl font-semibold mb-4">
                    <a href="post-{{.ID}}.html" class="text-blue-600 hover:underline">
                        {{.Title}}
                    </a>
                </h2>
                <p class="text-gray-600">By {{.Author}}</p>
            </div>
            {{end}}
        </div>
    </div>
</body>
</html>
`)
	if err != nil {
		log.Fatalf("Failed to parse index template: %v", err)
	}

	indexFile, err := os.Create(filepath.Join(outputDir, "index.html"))
	if err != nil {
		log.Fatalf("Failed to create index file: %v", err)
	}
	defer indexFile.Close()

	if err := indexTmpl.Execute(indexFile, posts); err != nil {
		log.Fatalf("Failed to write index template: %v", err)
	}

	fmt.Println("Static site generated successfully!")
}