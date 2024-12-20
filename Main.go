package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "html/template"
    "io/ioutil"
    "net/http"
    "os"

    "github.com/rs/cors"
)

type FormData struct {
    Task1     string `json:"task1"`
    Task2     string `json:"task2"`
    Tools1    string `json:"tools1"`
    Tracking1 string `json:"tracking1"`
    Pain1     string `json:"pain1"`
    Pain2     string `json:"pain2"`
    Goals1    string `json:"goals1"`
}

type ApiResponse struct {
    Suggestions []string `json:"suggestions"`
}

func main() {
    // Create a new CORS handler
    corsHandler := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"}, // Allow all origins for testing
        AllowCredentials: true,
    })

    http.HandleFunc("/", formHandler)
    http.HandleFunc("/submit", submitHandler)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Fallback to 8080 if PORT is not set
    }
    fmt.Printf("Server started at :%s\n", port)
    // Use the CORS handler
    http.ListenAndServe(":"+port, corsHandler.Handler(http.DefaultServeMux))
}

func formHandler(w http.ResponseWriter, r *http.Request) {
    // Serve the HTML form
    if r.Method == http.MethodGet {
        tmpl := `
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <link rel="stylesheet" href="TAEForm.css">
            <title>Task Automation Evaluation Form</title>
            <script src="https://unpkg.com/htmx.org@1.6.1"></script>
        </head>
        <body>
            <div class="background"></div>
            <div class="form-container">
                <h1>Task Automation Evaluation Form</h1>
                <form id="evaluation-form" hx-post="/submit" hx-target="#suggestions-container" hx-swap="innerHTML">
                    <fieldset>
                        <legend>Task Identification</legend>
                        <label for="task1">What specific tasks or processes are you currently performing that you believe could be automated? Please list them.</label>
                        <textarea id="task1" name="task1" required></textarea>
                    
                        <label for="task2">Which tasks in your workflow consume the most time? Please describe them.</label>
                        <textarea id="task2" name="task2" required></textarea>
                    </fieldset>

                    <fieldset>
                        <legend>Current Tools and Processes</legend>
                        <label for="tools1">What tools or software are you currently using to manage these tasks? Please specify.</label>
                        <textarea id="tools1" name="tools1" required></textarea>
                    
                        <label for="tracking1">How do you currently track the progress of these tasks? Describe your method.</label>
                        <textarea id="tracking1" name="tracking1" required></textarea>
                    </fieldset>

                    <fieldset>
                        <legend>Pain Points</legend>
                        <label for="pain1">What challenges or frustrations do you encounter with your current task management process? Please elaborate.</label>
                        <textarea id="pain1" name="pain1" required></textarea>
                    
                        <label for="pain2">Are there any repetitive tasks that you find particularly tedious or prone to errors? If so, please describe.</label>
                        <textarea id="pain2" name="pain2" required></textarea>
                    </fieldset>

                    <fieldset>
                        <legend>Goals and Outcomes</legend>
                        <label for="goals1">What are your primary goals for automating these tasks? (e.g., saving time, reducing errors, improving efficiency)</label>
                        <textarea id="goals1" name="goals1" required></textarea>
                    </fieldset>

                    <button type="submit">Submit Feedback</button>
                </form>
                <div id="suggestions-container" class="suggestions-output"></div>
            </div>
        </body>
        </html>
        `
        t, err := template.New("form").Parse(tmpl)
        if err != nil {
            http.Error(w, "Failed to parse template", http.StatusInternalServerError)
            return
        }
        err = t.Execute(w, nil)
        if err != nil {
 http.Error(w, "Failed to execute template", http.StatusInternalServerError)
            return
        }
    }
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        var formData FormData
        err := json.NewDecoder(r.Body).Decode(&formData)
        if err != nil {
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }

        // Here you would typically process the form data and generate suggestions
        suggestions := []string{
            "Consider automating task 1 using Tool A.",
            "Task 2 can be streamlined with Tool B.",
            "Implementing a tracking system could help with task management.",
        }

        response := ApiResponse{Suggestions: suggestions}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    } else {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}