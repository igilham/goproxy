package main

const (
	proxyHost    = "www-cache.mycompany.com"
	ftpProxyHost = "ftp-gw.mycompany.com"
	proxyPort    = "80"
	onNetwork    = "MyCompany On Network"
	offNetwork   = "MyCompany Off Network"
)

var (
	noProxyHosts = []string{
		"localhost",
		"127.0.0.1",
		"*.local",
		"sandbox",
	}
)
