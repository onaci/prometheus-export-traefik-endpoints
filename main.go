package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr       = flag.String("listen-address", "0.0.0.0:8080", "The address to listen on for HTTP requests.")
	traefikAPI = flag.String("traefik-api", "https://traefik.yoga260.alho.st/api", "The Traefik server's API URL")

	traefikEndpoints = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "traefik_endpoints",
			Help: "Traefik https hostname endpoints and certificate info",
		},
		[]string{"host", "ip", "subject", "dns", "isurer", "notbefore", "notafter"})
)

func main() {
	go func() {
		updateEndpoints()
		time.Sleep(5 * time.Minute)
	}()

	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	fmt.Printf("exporting metrics to http://%s/metrics\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func updateEndpoints() {
	response, err := http.Get(*traefikAPI)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		//fmt.Println(string(data))
		var config map[string]interface{}
		err := json.Unmarshal(data, &config)
		if err != nil {
			fmt.Println(err)
		} else {
			for _, val := range config {
				provider := val.(map[string]interface{})
				if provider["frontends"] == nil {
					continue
				}
				frontends := provider["frontends"].(map[string]interface{})
				for _, f := range frontends {
					entrypoints := f.(map[string]interface{})["entryPoints"].([]interface{})
					hasHTTPS := false
					for _, v := range entrypoints {
						if v.(string) == "https" {
							hasHTTPS = true
							break
						}
					}
					if !hasHTTPS {
						continue
					}
					routes := f.(map[string]interface{})["routes"].(map[string]interface{})
					for _, r := range routes {
						//fmt.Printf("%#v: %v\n", entrypoints, r.(map[string]interface{})["rule"])
						rule := r.(map[string]interface{})["rule"].(string)
						rules := strings.Split(rule, ";")
						for _, r := range rules {
							if strings.HasPrefix(r, "Host:") || strings.HasPrefix(r, "HostRegexp:") {
								//Host: and HostRegexp:
								hoststring := strings.TrimPrefix(r, "Host:")
								hoststring = strings.TrimPrefix(hoststring, "HostRegexp:")
								hosts := strings.Split(hoststring, ",")
								for _, host := range hosts {
									cert, ip, err := serverCertfunc(host, "443")
									if err != nil {
										fmt.Printf("https://%s/ (ERROR: %s)\n", host, err)
										continue
									}
									// fmt.Printf("https://%s/ (%s)\n", host, ip)
									// fmt.Printf("%s\n", cert[0].Issuer.CommonName)
									// fmt.Printf("%s\n", cert[0].Subject.CommonName)
									// fmt.Printf("%s\n", cert[0].DNSNames)
									// fmt.Printf("%s\n", cert[0].NotBefore.In(time.Local).String())
									// fmt.Printf("%s\n", cert[0].NotAfter.In(time.Local).String())

									traefikEndpoints.With(prometheus.Labels{
										"host":      host,
										"ip":        ip,
										"subject":   cert[0].Subject.CommonName,
										"dns":       strings.Join(cert[0].DNSNames, ","),
										"isurer":    cert[0].Issuer.CommonName,
										"notbefore": cert[0].NotBefore.String(),
										"notafter":  cert[0].NotAfter.String(),
									}).Set(float64(cert[0].NotAfter.Unix()))

								}
							}
						}
					}
				}
			}
		}
	}
}

var TimeoutSeconds = 3

func serverCertfunc(host, port string) ([]*x509.Certificate, string, error) {
	d := &net.Dialer{
		Timeout: time.Duration(TimeoutSeconds) * time.Second,
	}

	// cs, err := cipherSuite()
	// if err != nil {
	// 	return []*x509.Certificate{&x509.Certificate{}}, "", err
	// }

	conn, err := tls.DialWithDialer(d, "tcp", host+":"+port, &tls.Config{
		InsecureSkipVerify: true,
		//CipherSuites:       cs,
		//MaxVersion: tlsVersion(),
	})
	if err != nil {
		return []*x509.Certificate{&x509.Certificate{}}, "", err
	}
	defer conn.Close()

	addr := conn.RemoteAddr()
	ip, _, _ := net.SplitHostPort(addr.String())
	cert := conn.ConnectionState().PeerCertificates

	return cert, ip, nil
}
