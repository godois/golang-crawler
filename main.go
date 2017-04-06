package main

import (
	"net/http"
	"log"
	"bufio"
	"fmt"
	"crypto/tls"

	"regexp"
	"strings"
)

func main(){

	loadWithScanner()

}



func loadWithScanner(){
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	url:= "https://www.confaz.fazenda.gov.br/legislacao/protocolos/2008/pt041_08"

	res, err := client.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	// scan the page
	scanner := bufio.NewScanner(res.Body)

	defer res.Body.Close()
	// Set the split function for the scanning operation.
	//scanner.Split(bufio.ScanWords)
	scanner.Split(bufio.ScanLines)


	// Create slice to hold counts

	// Loop over the words
	for scanner.Scan() {

		linhaOriginal:= scanner.Text()

		replacer := strings.NewReplacer(
			"<p","",
			"class=","",
			"&", "",
			"<", "",
			">", "",
			`"`, "",
			"'", "",
			"Tabelajustificadoverde","",
			"/p", "",
			" align=center ","",
			"justify","",
			"RedaoAnt style=text-align:", "")

		linha := replacer.Replace(linhaOriginal)


		var percdouble = regexp.MustCompile(`[0-9][0-9],[0-9][0-9]%`)
		var percsingle = regexp.MustCompile(`[0-9],[0-9]%`)
		var ncm4digits = regexp.MustCompile(`\d\d[.]\d\d\s`)
		var ncm4digits2 = regexp.MustCompile(`\d\d\d\d[.]`)
		var data = regexp.MustCompile(`\d\d[.]\d\d[.]\d\d\s`)

		if percdouble.MatchString(linha){
			fmt.Println(" *** Aliquota XX,XX", linha)
		} else if percsingle.MatchString(linha){
			fmt.Println(" *** Aliquota X,X",linha)
		} else if ncm4digits.MatchString(linha){
			if !data.MatchString(linha){
				fmt.Println(" *** NCM XX.XX",linha)
			}

		} else if ncm4digits2.MatchString(linha){
			fmt.Println(" *** NCM XXXX",linha)
		}

	}
}