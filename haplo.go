package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func fileImport() {

}

const alphabet = 2

func main() {

	// inputFilePref := "hap_gen_10000_20220.txt"
	// inputFilePref := "test2.txt"
	inputFilePref := "test.txt"

	var inFile = flag.String("in", inputFilePref, "input file relative path as a string")
	flag.Parse()

	// Open input file and defer closure
	file, err := os.Open(*inFile)

	if err != nil {
		return
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string
	rows := 0
	for scanner.Scan() {
		rows++
		lines = append(lines, scanner.Text())
	}

	file.Close()

	columns := len(lines[0])

	fmt.Println("First Column ") // remembver to -48
	for i := 0; i < len(lines); i++ {
		fmt.Print(lines[i][4] - 48)
	}

	// // read file in a byte array, ram==file maybe needed, to check during tests.
	// bs := make([]byte, stat.Size())
	// _, err = file.Read(bs)
	// if err != nil {
	// 	fmt.Println("err2", err)

	// 	return
	// }
	// maybe read first row insted of using the for cycle -- @TODO check goLang Impl
	// var column int
	// for i := 0; i < int(stat.Size()); i++ {
	// 	if bs[i] == 10 { // check LineFeed==10 vs CarriageReturn
	// 		column = i
	// 		break
	// 	}
	// }

	// rows := (int(stat.Size()) / column)

	fmt.Println("Input size size is ", rows, " rows (samples)  x ", columns, " columns (SNPs)")
	fmt.Println("Assuming: alphabet size ==", alphabet) // in further implementation the program could recognize itself input type
	fmt.Println("Assuming: wildcard are present")
	// fmt.Println(lines)
	// arrays a_k and d_k will be built dinamically, blocks will be reported during the execution, MPHB will be computed right after

	// init starting arrays

	ak0 := make([]int, 0, rows)
	for i := 0; i < rows; i++ {
		ak0 = append(ak0, i)
	}

	dk0 := make([]int, 0, rows)
	for i := 0; i < rows; i++ {
		dk0 = append(dk0, 0)
	}

	pivot := 0
	fmt.Println(ak0, dk0)

	ak0, dk0 = computeNextArrays(ak0, dk0, pivot, lines)
	pivot++

	fmt.Println(ak0, dk0)

	ak0, dk0 = computeNextArrays(ak0, dk0, pivot, lines)
	pivot++

	fmt.Println(ak0, dk0)

	ak0, dk0 = computeNextArrays(ak0, dk0, pivot, lines)
	pivot++

	fmt.Println(ak0, dk0)
}

func computeNextArrays(ak, dk []int, k int, matrix []string) ([]int, []int) {
	dim := len(ak)
	// allocing dim size, that's not really memory wise to be honest
	// go slices got me covered ehehehe
	a := make([][]int, alphabet)
	for i := range a {
		a[i] = make([]int, 0, dim)
	}

	d := make([][]int, alphabet)
	for i := range d {
		d[i] = make([]int, 0, dim)
	}

	p := [alphabet]int{}
	u := [alphabet]int{}

	for i := 0; i < alphabet; i++ {
		u[i] = 0
		p[i] = k + 1
	}

	var allele int

	for i := 0; i < dim; i++ {
		allele = int(matrix[ak[i]][k] - 48)
		for l := 0; l < alphabet; l++ {
			if dk[i] > p[l] {
				p[l] = dk[i]
			}
		}
		if allele == -6 {
			for m := 0; m < alphabet; m++ {
				a[m] = append(a[m], ak[i])
				d[m] = append(d[m], p[m])
				p[m] = 0
				u[m] = u[m] + 1
			}
		} else {
			a[allele] = append(a[allele], ak[i])
			d[allele] = append(d[allele], p[allele])
			p[allele] = 0
			u[allele] = u[allele] + 1
		}

	}

	newdim := 0
	for i := 0; i < alphabet; i++ {
		newdim += len(a[i])
	}
	var akk []int
	for i := 0; i < alphabet; i++ {
		akk = append(akk, a[i]...)
	}

	var dkk []int
	for i := 0; i < alphabet; i++ {
		for j := range d[i] {
			dkk = append(dkk, d[i][j])
		}
	}
	return akk, dkk

}