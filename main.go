package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	baseJavaOpts = os.Getenv("BASE_JAVA_OPTS")

	proxyEnvKeys = []string{
		"http_proxy",
		"HTTP_PROXY",
		"https_proxy",
		"HTTPS_PROXY",
	}
)

func which(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func printenv(key string, value string) {
	fmt.Printf("export %s=\"%s\"\n", key, value)
}

func sh(args []string) {
	s := strings.Join(args, " ") + " &>/dev/null"
	fmt.Println(s)
}

func getLocation() string {
	cmd := exec.Command("networksetup", "-getcurrentlocation")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(out.String())
}

func configureCommon() {
	platformJavaOpts := os.Getenv("PLATFORM_JAVA_OPTS")
	envNoProxyHosts := strings.Join(noProxyHosts, ",")
	printenv("no_proxy", envNoProxyHosts)
	printenv("MAVEN_OPTS", "-Xms256m -Xmx512m "+platformJavaOpts)
}

func configureProxyOn() {
	hostPort := proxyHost + ":" + proxyPort
	httpProxy := "http://" + proxyHost

	for _, key := range proxyEnvKeys {
		printenv(key, httpProxy)
	}

	printenv("FTP_PROXY", ftpProxyHost)

	javaNoProxyHosts := strings.Join(noProxyHosts, "|")
	platformJavaOpts := fmt.Sprintf("%s -Dhttp.proxyHost=%s -Dhttp.proxyPort=%s -Dhttps.proxyHost=%s -Dhttps.proxyPort=%s -Dhttp.nonProxyHosts='%s'", baseJavaOpts, proxyHost, proxyPort, proxyHost, proxyPort, javaNoProxyHosts)
	printenv("PLATFORM_JAVA_OPTS", platformJavaOpts)

	if which("git") {
		sh([]string{"git", "config", "--global", "http.proxy", hostPort})
		sh([]string{"git", "config", "--global", "https.proxy", hostPort})
	}

	if which("npm") {
		sh([]string{"npm", "config", "-g", "set", "proxy", httpProxy})
		sh([]string{"npm", "config", "-g", "set", "https-proxy", httpProxy})
	}

	configureCommon()
}

func configureProxyOff() {
	for _, key := range proxyEnvKeys {
		printenv(key, "")
	}

	printenv("FTP_PROXY", "")
	printenv("PLATFORM_JAVA_OPTS", baseJavaOpts)

	if which("git") {
		sh([]string{"git", "config", "--global", "--remove-section", "http"})
		sh([]string{"git", "config", "--global", "--remove-section", "https"})
	}

	if which("npm") {
		sh([]string{"npm", "config", "-g", "delete", "proxy"})
		sh([]string{"npm", "config", "-g", "delete", "https-proxy"})
	}

	configureCommon()
}

func statusCommand() {
	location := getLocation()
	fmt.Printf("Network location: %s", location)
}

func onCommand() {
	cmd := exec.Command("networksetup", "-switchtolocation", onNetwork)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	configureProxyOn()
}

func offCommand() {
	cmd := exec.Command("networksetup", "-switchtolocation", offNetwork)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	configureProxyOff()
}

func resetCommand() {
	if location := getLocation(); location == onNetwork {
		configureProxyOn()
	} else {
		configureProxyOff()
	}
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatal("not enough arguments")
	}
	if flag.NArg() > 1 {
		log.Fatal("too many arguments")
	}

	cmd := flag.Arg(0)
	switch cmd {
	case "status":
		statusCommand()
	case "on":
		onCommand()
	case "off":
		offCommand()
	case "reset":
		resetCommand()
	default:
		log.Fatal("invalid command")
	}
}
