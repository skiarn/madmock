//Copyright (C) 2016  Andreas Westberg

//This program is free software: you can redistribute it and/or modify
//it under the terms of the GNU Lesser General Public License as published by
//the Free Software Foundation, version 3 of the License.

//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU Lesser General Public License for more details.

//You should have received a copy of the GNU Lesser General Public License
//along with this program.  If not, see <http://www.gnu.org/licenses/lgpl-3.0.txt>.

package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/skiarn/madmock/handler"
	"github.com/skiarn/madmock/setting"
	"github.com/skiarn/madmock/ws"

	"golang.org/x/net/websocket"
)

var logger *log.Logger
var errorlog *os.File

func main() {
	//errorlog, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//if err != nil {
	//	fmt.Printf("error opening file: %v", err)
	//	os.Exit(1)
	//}
	//defer errorlog.Close()
	logger = log.New(os.Stdout, "applog: ", log.Lshortfile|log.LstdFlags)

	settings, err := setting.Create()
	if err != nil {
		log.Fatal(err)
	}
	err = settings.CreateDir()
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	pagehandler := handler.NewPageHandler(settings.DataDirPath, settings.TargetURL)
	curdhandler := handler.NewMockCURDHandler(settings.DataDirPath)
	viewDataHandler := handler.NewViewDataHandler(settings.DataDirPath)

	wsHandler := ws.NewHandler(logger)
	mockhandler := handler.NewMockhandler(settings.TargetURL, settings.DataDirPath, wsHandler)
	mux.Handle("/mock/wsmockinfo", websocket.Handler(wsHandler.WSMockInfoServer))
	go wsHandler.Run()

	mux.Handle("/mock", &pagehandler)
	mux.Handle("/mock/api/mock/", &curdhandler)
	mux.Handle("/mock/api/mock/data/", &viewDataHandler)

	mux.HandleFunc("/mock/edit/", pagehandler.EditHandler)
	mux.HandleFunc("/mock/new/", pagehandler.New)
	mux.Handle("/", &mockhandler)
	logger.Printf("Server to listen on a port: %v \n", settings.Port)
	protocol := "http"
	if settings.TLS {
		protocol = "https"
	}
	openBrowser(protocol + "://localhost:" + strconv.Itoa(settings.Port) + "/mock")

	if settings.TLS {
		err := generateSelfSignedCertificate()
		if err != nil {
			log.Fatal("Error generating self-signed certificate:", err)
			os.Exit(1)
		}
		err = http.ListenAndServeTLS(":"+strconv.Itoa(settings.Port), "server.crt", "server.key", mux)
		if err != nil {
			log.Fatal("Error starting HTTPS server:", err)
		}
		return
	}
	logger.Fatal(http.ListenAndServe(":"+strconv.Itoa(settings.Port), mux))
}

func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("cmd", "/c", "start", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = os.ErrInvalid
	}

	return err
}

func generateSelfSignedCertificate() error {
	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return err
	}

	certTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Self-Signed Certificate"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	certFile, err := os.Create("server.crt")
	if err != nil {
		return err
	}
	defer certFile.Close()
	err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	if err != nil {
		return err
	}

	keyFile, err := os.Create("server.key")
	if err != nil {
		return err
	}
	defer keyFile.Close()
	keyBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return err
	}
	err = pem.Encode(keyFile, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})
	if err != nil {
		return err
	}

	return nil
}
