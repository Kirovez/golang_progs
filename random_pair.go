package main


/*
this app randomly select paired sequences and write fasta file.
It takes two fastq files and number (percentage of the reads required).
*/
import (
	"os"
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}


func randomReads(file_path string, file2_path string, percent int){
	//file one opening
	file, err := os.Open(file_path)
	check(err)
	defer file.Close()

	// file2 opening
	file2, err := os.Open(file2_path)
	check(err)
	defer file2.Close()

	// create output file
	out_file, err := os.Create(file_path + "random.fasta")
	writer_obj := bufio.NewWriter(out_file)
	check(err)

	// run the app
	current_id := ""
	current_seq := ""
	line_number := 0
	cnt := 0
	cnt_total := 0
	dic_colllecte := make(map[string]string)

	// parse file one and collect sequences into map
	scanner_type := bufio.NewScanner(file)
	for scanner_type.Scan() {

		if strings.HasPrefix(scanner_type.Text(),"@") {
			cnt_total += 1
			line_number = 0
			if current_id != "" && 	rand.Intn(100) < percent {
				dic_colllecte[current_id] = current_seq
				//fmt.Println(current_id)
				//writer_obj.WriteString(strings.Join([]string{current_id, current_seq, ""}, "\n") )
			}
			current_id = scanner_type.Text()
			current_id = strings.Split(current_id, " ")[0]
			current_seq = ""

		} else if line_number == 1{
			current_seq = strings.Join([]string{current_seq, scanner_type.Text()},"")
		}

		line_number += 1
	}

	if err := scanner_type.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Number of sequences selected from first file", len(dic_colllecte))


	// parse file 2 and write sequences
	scanner_type2 := bufio.NewScanner(file2)
	for scanner_type2.Scan() {
		if strings.HasPrefix(scanner_type2.Text(),"@") {
			line_number = 0
			cnt_total += 1
			if val, ok := dic_colllecte[current_id]; ok 	 {
				cnt += 1
				//fmt.Println(strings.Join([]string{current_id + "/1", val, current_id + "/2", current_seq, ""}, "\n"))
				_, err = writer_obj.WriteString(strings.Join([]string{current_id + "/1", val, current_id + "/2", current_seq, ""}, "\n"))
				check(err)
			}
			current_id = scanner_type2.Text()
			current_id = strings.Split(current_id, " ")[0]
			current_seq = ""


		} else if line_number == 1{
			current_seq = strings.Join([]string{current_seq, scanner_type2.Text()},"")
		}
		line_number += 1
	}

	fmt.Println("Number of sequences selected from second file", cnt)


	fmt.Println("Number of sequence pairs in files is ", cnt_total/2)
	fmt.Println("Number of sequences selected is ", cnt)
}

func main(){
	st := time.Now()
	percent, _ := strconv.Atoi(os.Args[3])
	randomReads(os.Args[1], os.Args[2], percent)
	fmt.Println("Time spent: ", time.Since(st))
}
