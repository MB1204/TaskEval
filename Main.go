package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "html/template"
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
    // Test the API endpoint
    testApiEndpoint()

    // Create a new CORS handler
    corsHandler := cors.New(cors.Options{
        AllowedOrigins:   []string{"https://your-production-frontend-url.com"}, // Replace with your actual frontend URL
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
    http.ListenAndServe(":8080", corsHandler.Handler(http.DefaultServeMux))
}

func testApiEndpoint() {
    apiUrl := "https://taskeval-production.up.railway.app/Main" // Your API endpoint
    apiToken := os.Getenv("REPLICATE_API_TOKEN")
    if apiToken == "" {
        fmt.Println("API token is not set in the environment variables")
        return
    }

    requestBody, err := json.Marshal(map[string]string{"query": "Test query"})
    if err != nil {
        fmt.Println("Error marshaling request body:", err)
        return
    }

    req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(requestBody))
    if err != nil {
        fmt.Println("Error creating request:", err)
        return
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+apiToken)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error making request:", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        fmt.Printf("API request failed with status: %s\n", resp.Status)
        return
    }

    var apiResponse ApiResponse
    if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
        fmt.Println("Error decoding response:", err)
        return
    }

    fmt.Println("API response:", apiResponse)
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
        <form id="evaluation-form"
      hx-post="/submit" hx-target="#suggestions-container" hx-swap="innerHTML">
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

    <script>
        const BACKEND_URL = "https://taskeval-production.up.railway.app/Main"; // Replace with your actual backend URL
    </script>
</body>
</html>
        `
        t, err := template.New("form").Parse(tmpl)
        if err != nil {
            http.Error(w, "Failed to parse template ", http.StatusInternalServerError)
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
        // Parse the form data
        err := r.ParseForm()
        if err != nil {
            fmt.Println("Error parsing form:", err)
            http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
            return
        }

        data := FormData{
            Task1:     r.FormValue("task1"),
            Task2:     r.FormValue("task2"),
            Tools1:    r.FormValue("tools1"),
            Tracking1: r.FormValue("tracking1"),
            Pain1:     r.FormValue("pain1"),
            Pain2:     r.FormValue("pain2"),
            Goals1:    r.FormValue("goals1"),
        }

        // Send the data to an external API
        suggestions, err := getAutomationSuggestions(data)
        if err != nil {
            http.Error(w, "Failed to get suggestions: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // Create an HTML response to be injected into the suggestions container
        responseHTML := "<h2>Form submitted successfully!</h2>"
        responseHTML += "<p><strong>Task 1:</strong> " + data.Task1 + "</p>"
        responseHTML += "<p><strong>Task 2:</strong> " + data.Task2 + "</p>"
        responseHTML += "<p><strong>Tools:</strong> " + data.Tools1 + "</p>"
        responseHTML += "<p><strong>Tracking Method:</strong> " + data.Tracking1 + "</p>"
        responseHTML += "<p><strong>Pain Points:</strong> " + data.Pain1 + "</p>"
        responseHTML += "<p><strong>Repetitive Tasks:</strong> " + data.Pain2 + "</p>"
        responseHTML += "<p><strong>Goals:</strong> " + data.Goals1 + "</p>"
        
        // Add suggestions to the response
        responseHTML += "<h3>Automation Suggestions:</h3><ul>"
        for _, suggestion := range suggestions {
            responseHTML += "<li>" + suggestion + "</li>"
        }
        responseHTML += "</ul>"

        // Set the content type to HTML
        w.Header().Set("Content-Type", "text/html")
        fmt.Fprint(w, responseHTML)
    } else {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func getAutomationSuggestions(data FormData) ([]string, error) {
    apiUrl := "https://taskeval-production.up.railway.app/Main" // Corrected API endpoint
    apiToken := os.Getenv("REPLICATE_API_TOKEN")
    if apiToken == "" { 
        return nil, fmt.Errorf("API token is not set in the environment variables")
    }

    query := fmt.Sprintf("Give the best automation suggestions based on the answers: %s, %s, %s, %s, %s, %s, %s",
        data.Task1, data.Task2, data.Tools1, data.Tracking1, data.Pain1, data.Pain2, data.Goals1)

    requestBody, err := json.Marshal(map[string]string{"query": query})
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(requestBody))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+apiToken) // Add the API token to the request header

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
    }

    var apiResponse ApiResponse
    if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
        return nil, err
    }

    return apiResponse.Suggestions, nil
}


