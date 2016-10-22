package main

import(
  "fmt"
  "net/http"
  "net/url"
  "log"
  "time"
  "flag"
  "os"
  "syscall"

  "github.com/fvbock/endless"
  "github.com/gorilla/mux"
)

func root_handler(w http.ResponseWriter, r *http.Request) {
  time.Sleep(time.Duration(10) * time.Second)
  log.Println("visit url /")
  fmt.Fprintf(w, "You're fine.")
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

func write_pid() {
  if *pid_path == "" {
    return
  }
  f, err := os.Create(*pid_path)
  if err != nil {
    log.Fatal(err)
  }
  _, err = f.WriteString(fmt.Sprintf("%d", syscall.Getpid()))
  f.Close()
  if err != nil {
    log.Fatal(err)
  }
}

func is_file_exist() bool {
  if *pid_path == "" {
    return false
  }
  f, err := os.Open(*pid_path)
  if err == nil {
    f.Close()
  }
  return err == nil
}

func clear_pid() {
  if *pid_path != "" {
    os.Remove(*pid_path)
  }
}

var pid_path = flag.String("p", "", "pid file")

func main() {
  flag.Parse()

  // if is_file_exist() {
    // log.Fatal("already running")
  // }
  mux := mux.NewRouter()
  mux.HandleFunc("/", root_handler).Methods("GET")
  mux.HandleFunc("/clicks", click_handler).Methods("GET")

  server := endless.NewServer(":8080", mux)
  server.BeforeBegin = func(add string) {
    write_pid()
    log.Println("pid changed")
  }
  err := server.ListenAndServe()
  if err != nil {
    log.Println(err)
  }
  clear_pid()
  log.Println("Server on 4242 stopped")
}
