package main

import (
	"embed"
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"html/template"
	"image/png"
	"log"
	"net/http"
	"strings"
)

//go:embed templates
var templates embed.FS

//go:embed assets
var assets embed.FS

var (
	port        = 8080
	templateDir = "templates"
	assetsDir   = "assets"
	baseURL     = "http://localhost:8080"
)

type TemplateData struct {
	Data        []WakeUp
	ShowMessage bool
	Message     string
	Severity    string
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	loadData()
	message := ""
	severity := "sucess"
	showMessage := false
	if r.RequestURI == "" || r.RequestURI == "/" {
		c, err := r.Cookie("message")
		if err == nil {
			s := strings.Split(c.Value, "|")
			if len(s) == 2 {
				severity = s[0]
				message = s[1]
				showMessage = true
				http.SetCookie(w, &http.Cookie{Name: "message", Value: "", Path: "/"})
			}
		}
		if r.Method == http.MethodPost {
			w, odevice, scope, err := wakeUpFromRequest(r)
			if err != nil {
				// add error message to page
				message = fmt.Sprintf("No changes made: %v", err)
				severity = "danger"
				showMessage = true
			} else {
				err = insertOrUpdateData(w, odevice, scope)
				if err != nil {
					message = fmt.Sprintf("No changes made: %v", err)
					severity = "danger"
					showMessage = true
				} else {
					saveData()
					message = fmt.Sprintf("Changes written")
					severity = "success"
					showMessage = true
				}
			}
		}
		td := TemplateData{
			Data:        data,
			ShowMessage: showMessage,
			Message:     message,
			Severity:    severity,
		}
		renderTemplateIndex(w, td)
		//t.Execute(w, td)
		accessLog(r, http.StatusOK, "")
	} else if strings.HasSuffix(r.RequestURI, "/delete") {
		handlerDelete(w, r)
		accessLog(r, http.StatusOK, "")
	} else if strings.HasSuffix(r.RequestURI, "/clone") {
		handlerClone(w, r)
		accessLog(r, http.StatusOK, "")
	} else if strings.HasSuffix(r.RequestURI, "/qrcode") {
		handlerQrCode(w, r)
		accessLog(r, http.StatusOK, "")
	} else if strings.HasSuffix(r.RequestURI, "/wakeup") {
		handlerWakeup(w, r)
		accessLog(r, http.StatusOK, "")
	} else if r.RequestURI == "/index.html" {
		http.Redirect(w, r, "./", http.StatusMovedPermanently)
		accessLog(r, http.StatusMovedPermanently, "Redirect to '/'")
	} else if f, err := assets.Open(assetsDir + r.RequestURI); err == nil {
		f.Close()
		handlerStaticFiles(w, r)
	} else {
		renderNotFound(w, r)
		accessLog(r, http.StatusNotFound, "")
	}
}

func handlerStaticFiles(w http.ResponseWriter, r *http.Request) {
	data, err := assets.ReadFile(assetsDir + r.URL.Path)
	if err != nil {
		accessLog(r, http.StatusInternalServerError, err.Error())
		renderServerError(w, r)
		return
	}
	accessLog(r, 200, "")
	lc := strings.ToLower(r.RequestURI)
	switch {
	case strings.HasSuffix(lc, ".css"):
		w.Header().Add("Content-Type", "text/css")
	case strings.HasSuffix(lc, ".jpg"):
		w.Header().Add("Content-Type", "image/jpeg")
	case strings.HasSuffix(lc, ".jpeg"):
		w.Header().Add("Content-Type", "image/jpeg")
	case strings.HasSuffix(lc, ".png"):
		w.Header().Add("Content-Type", "image/png")
	case strings.HasSuffix(lc, ".gif"):
		w.Header().Add("Content-Type", "image/gif")
	case strings.HasSuffix(lc, ".ico"):
		w.Header().Add("Content-Type", "image/x-icon")
	case strings.HasSuffix(lc, ".html"):
		w.Header().Add("Content-Type", "text/html")
	case strings.HasSuffix(lc, ".js"):
		w.Header().Add("Content-Type", "application/javascript")
	case strings.HasSuffix(lc, ".map"):
		w.Header().Add("Content-Type", "application/json")
	case strings.HasSuffix(lc, ".svg"):
		w.Header().Add("Content-Type", "image/svg+xml")
	case strings.HasSuffix(lc, ".woff2"):
		w.Header().Add("Content-Type", "font/woff2")
	case strings.HasSuffix(lc, ".woff"):
		w.Header().Add("Content-Type", "application/font-woff")
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GET /{device}/delete
func handlerDelete(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	if len(s) != 3 || s[2] != "delete" || !deviceExists(s[1]) {
		http.SetCookie(w, &http.Cookie{
			Name:  "message",
			Value: "danger|Bad request.",
			Path:  "/",
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	log.Printf("DELETE %s", s[1])
	err := deleteItem(s[1])
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:  "message",
			Value: "danger|Device was not deleted.",
			Path:  "/",
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else {
		saveData()
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "message",
		Value: "success|Device deleted.",
		Path:  "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// GET /{device}/clone
func handlerClone(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	if len(s) != 3 || s[2] != "clone" || !deviceExists(s[1]) {
		http.SetCookie(w, &http.Cookie{
			Name:  "message",
			Value: "danger|Bad request.",
			Path:  "/",
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	log.Printf("CLONE %s", s[1])
	err := cloneItem(s[1])
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:  "message",
			Value: "danger|Device was not cloned.",
			Path:  "/",
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else {
		saveData()
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "message",
		Value: "success|Device cloned.",
		Path:  "/",
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// GET /{device}/qrcode
func handlerQrCode(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	if len(s) != 3 || s[2] != "qrcode" || !deviceExists(s[1]) {
		renderNotFound(w, r)
		return
	}
	qrCode, err := qr.Encode(fmt.Sprintf("%s/%s/wakeup", baseURL, s[1]), qr.M, qr.Auto)
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
		return
	}
	qrCode, _ = barcode.Scale(qrCode, 256, 256)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, qrCode)
}

// GET /{device}/qrcode
func handlerWakeup(w http.ResponseWriter, r *http.Request) {
	s := strings.Split(r.URL.Path, "/")
	if len(s) != 3 || s[2] != "wakeup" || !deviceExists(s[1]) {
		http.SetCookie(w, &http.Cookie{
			Name:  "message",
			Value: "danger|Bad request.",
			Path:  "/",
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	wu, ok := wakeupData(s[1])
	if ok {
		wolUdp(wu.Ip, wu.Mac, nil)
		http.SetCookie(w, &http.Cookie{
			Name:  "message",
			Value: "success|Wake up signal sent.",
			Path:  "/",
		})
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func accessLog(r *http.Request, httpCode int, payload string) {
	log.Printf("%s %s, %d, %s", r.Method, r.RequestURI, httpCode, payload)
}

func renderNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	renderWithMessage(w, "Page not found!", "danger")
}

func renderServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Internal Server Error: %s.", r.RequestURI)
}

func wakeUpFromRequest(r *http.Request) (WakeUp, string, string, error) {
	w := WakeUp{}
	err := r.ParseForm()
	if err != nil {
		return w, "", "", fmt.Errorf("error parsing form")
	}
	device, ok := r.Form["device"]
	if !ok {
		return w, "", "", fmt.Errorf("device must be set")
	}
	odevice, ok := r.Form["odevice"]
	if !ok {
		odevice = device
	}
	mac, ok := r.Form["mac"]
	if !ok {
		return w, "", "", fmt.Errorf("mac must be set")
	}
	ip, ok := r.Form["ip"]
	if !ok {
		return w, "", "", fmt.Errorf("ip must be set")
	}
	scope, ok := r.Form["scope"]
	if !ok {
		return w, "", "", fmt.Errorf("scope must be set")
	}
	w.Device = device[0]
	w.Mac = mac[0]
	w.Ip = ip[0]
	log.Printf("Got POST request with data (%s, %s, %s, %s)", w.Device, w.Mac, w.Ip, scope[0])
	return w, odevice[0], scope[0], nil
}
func renderTemplateIndex(w http.ResponseWriter, td TemplateData) {
	t, _ := template.ParseFS(templates, templateDir+"/index.html")
	t.Execute(w, td)
}

func renderWithMessage(w http.ResponseWriter, message string, severity string) {
	td := TemplateData{
		Data:        data,
		ShowMessage: true,
		Message:     message,
		Severity:    severity,
	}
	renderTemplateIndex(w, td)
}
