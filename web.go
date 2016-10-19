package main
import(
  "fmt"
  "net/http"
  "net/url"
  "log"
  "time"
)

func root_handler(w http.ResponseWriter, r *http.Request) {
  time.Sleep(time.Duration(10) * time.Millisecond)
  log.Println("visit url /")
  fmt.Fprintf(w, "I'm fine.")
}

func click_handler(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  log.Printf("visit url /clicks with params %v", r.Form)
  redirect_to, err := fetch_redirect_url(r.Form)
  if err != nil {
    log.Print(err)
    fmt.Fprintf(w,"%s", err)
    return
  }
  log.Printf("redirect to %s", redirect_to)
  http.Redirect(w, r, redirect_to, http.StatusFound)
}

func fetch_redirect_url(form url.Values) (string, error) {
  redirect_to_values := form["redirect_to"]
  var redirect_to string

  if len(redirect_to_values) > 0 {
    redirect_to = redirect_to_values[0]
  }

  if redirect_to == "" {
    return "", fmt.Errorf("parameter 'redirect_to' not found")
  }

  redirect_url, err := url.Parse(redirect_to)
  if err != nil {
    return "", err
  }
  if redirect_url.Scheme == "" {
    redirect_url.Scheme = "http"
  }
  return redirect_url.String(), nil
}

func main() {
  http.HandleFunc("/", root_handler)
  http.HandleFunc("/clicks", click_handler)
  log.Print("server started")
  http.ListenAndServe(":8080", nil)
}
