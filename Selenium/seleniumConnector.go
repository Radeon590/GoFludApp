package main

import(
	"fmt"
	"os"
	"time"

	"github.com/tebeka/selenium"
)

func main(){
	ConnectToSite("http://library.ru")
}

func ConnectToSite(url string){
	fmt.Println(url)
	const(
		seleniumPath = "./vendor/selenium-server.jar"
		geckoDriverPath = "./vendor/geckodriver"
		chromeDriverPath = "./vendor/chromedriver"
		port = 8080
	)
	opts := []selenium.ServiceOption{
		selenium.GeckoDriver(geckoDriverPath),
		selenium.Output(os.Stderr),
		selenium.ChromeDriver(chromeDriverPath),
	}
	selenium.SetDebug(true)
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil{
		panic(err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil{
		panic(err)
	}
	defer wd.Quit()

	newErr := wd.Get(url)
	if newErr != nil {
		panic(newErr)
	}
	time.Sleep(time.Duration(15) * time.Second)

	elem, err := wd.FindElement(selenium.ByClassName, "hdr")
	if err != nil{
		panic(err)
	}
	_ = elem.Click()
	fmt.Println("success")
}

func checkRequiredPath(path string) (bool){
	_, err := os.Stat(path)
	if err != nil{
		fmt.Println(err)
		return false
	}
	return true
}
