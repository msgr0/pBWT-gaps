package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"
)

const alphabet = 2         // alphabet cardinality - assuming {0, 1, ..., t-1} plus * symbol and t = 5
const wildcard bool = true // wildcard presence assumed true
const min_block_width = 10 // almost x SNPs
const min_block_rows = 2   // almost y rows

type block struct {
	i int   // begin column
	j int   // end column
	k []int //row indexes
}

func main() {

	// inputFilePref := "bigFile.txt"
	// inputFilePref := "test.txt"
	// inputFilePref := "alph3.txt"
	inputFilePref := "test2.txt"
	// inputFilePref := "test.txt"

	var inFile = flag.String("in", inputFilePref, "input file relative path as a string")
	flag.Parse()

	// Open input file and defer closure
	file, err := os.Open(*inFile)

	if err != nil {
		return
	}

	// big buffer
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*5000)
	scanner.Buffer(buf, 1024*1024)

	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	file.Close()

	columns := len(lines[0])
	rows := len(lines)

	// DEBUG
	fmt.Println("First Column ") // remembver to -48
	for i := 0; i < len(lines); i++ {
		fmt.Print(lines[i][0] - 48)
	}

	// INPUT INFO

	fmt.Println("Input matrix is", rows, "rows (samples)  x", columns, "columns (SNPs)")
	fmt.Println("Assuming: alphabet size ==", alphabet) // in further implementation the program could recognize itself input type
	fmt.Println("Assuming: wildcard presence ==", wildcard)
	fmt.Println("Assuming min block size ==", min_block_rows, "rows x", min_block_width, "columns")

	// INIT Ak Dk Arraya

	ak0 := make([]int, 0, rows)
	for i := 0; i < rows; i++ {
		ak0 = append(ak0, i)
	}

	dk0 := make([]int, 0, rows)
	for i := 0; i < rows; i++ {
		dk0 = append(dk0, 0)
	}

	v := make([][]int8, alphabet)

	stopper := len(lines[0])
	pivot := 0

	// fmt.Println(ak0, dk0)
	var since time.Duration
	var start time.Time = time.Now()

	// main loop che scorre tutta la matrice
	for pivot < stopper {

		ak0, dk0 = computeNextArrays(ak0, dk0, pivot, lines)
		// fmt.Println("Currently printing k = ", i+1)
		// fmt.Println("A_k[0]:", ak0[0], "\tD_k[0]:", dk0[0])
		// fmt.Println("Collapsing ...")
		ak0, dk0 = collapse(ak0, dk0)
		// fmt.Println("A_k[0]:", ak0[0], "\tnD_k[0]:", dk0[0])
		// fmt.Println("----####----")

		blocks := availableBlocks(ak0, dk0)

		v = computeBitVectors(ak0, dk0, pivot, lines)

		if pivot == stopper-2 {
			fmt.Println(blocks)
		}

		pivot++

	}
	fmt.Println("TOTAL COLUMNS : ", pivot, "Last size arrays: ak= ", len(ak0), " and dk= ", len(dk0))
	// print last (j) column
	for i := 0; i < len(ak0); i++ {
		fmt.Println("--> ", ak0[i], " and -> ", dk0[i], " and ->? ", v[0][i])
	}

	fmt.Println("Print Available blocks at end")

	fmt.Println("Print last bit vector", v)

	since = time.Since(start)

	fmt.Println("Started at : ", start, "\nRAN in ss: ", since)
}

func computeBitVectors(ak, dk []int, k int, matrix []string) [][]int8 {
	dim := len(ak)
	v := make([][]int8, alphabet)
	for i := range v {
		v[i] = make([]int8, 1, dim)
	}
	// init first position
	for t := 0; t < alphabet; t++ {
		v[t][0] = 1
	}

	if k == len(matrix[0])-1 {
		for t := 0; t < alphabet; t++ {
			for i := 1; i < dim; i++ {
				v[t] = append(v[t], 1)
			}
		}
		return v
	}

	var allele int
	var prec int
	prec = int(matrix[ak[0]][k+1] - 48) //
	fmt.Println("Prec ", prec)
	for i := 1; i < dim; i++ {
		allele = int(matrix[ak[i]][k+1] - 48)
		// if allele == prec {
		// 	for t := range v {
		// 		v[t] = append(v[t], 0)
		// 	}

		// }

		for t := 0; t < alphabet; t++ {
			if allele == prec {
				v[t] = append(v[t], 0)
			} else if allele > 9 {
				if t == prec {
					v[t] = append(v[t], 0)
				} else {
					v[t] = append(v[t], 1)
				}
			} else {
				v[t] = append(v[t], 1)
			}
		}
		prec = allele
	}
	return v
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

		if (allele < 10) && (allele >= 0) {

			a[allele] = append(a[allele], ak[i])
			d[allele] = append(d[allele], p[allele])
			p[allele] = 0
			u[allele] = u[allele] + 1

		} else {
			for m := 0; m < alphabet; m++ {
				a[m] = append(a[m], ak[i])
				d[m] = append(d[m], p[m])
				p[m] = 0
				u[m] = u[m] + 1
			}
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
		dkk = append(dkk, d[i]...)
	}

	return akk, dkk

}

func collapse(a, d []int) ([]int, []int) {
	ac := make([]int, 0, len(a))
	dc := make([]int, 0, len(d))

	j := 0
	pivot := 0
	for j < len(a)-1 {
		if a[j] == a[j+1] {
			j++
		} else {
			ac = append(ac, a[pivot])
			dc = append(dc, d[pivot])
			j++
			pivot = j
		}

	}

	ac = append(ac, a[pivot])
	dc = append(dc, d[pivot])

	return ac, dc
}
func availableBlocks(a, d []int) []block {
	available_blocks := make([]block, 0, len(d))
	for i := 1; i < len(d)-1; i++ {
		x := i - 1
		y := i + 1
		if d[i] <= d[0]-min_block_width {

			if d[x] <= d[i] {
				x--
			}

			if d[y] <= d[i] {
				y++
			}

			// block trim
			newblock := makeblock(d[i], d[0], a[x:y])
			fmt.Println("appending block: ", newblock)
			available_blocks = append(available_blocks, newblock)

		}

	}
	return available_blocks

	//TODO
	//TRIM to 1 instance of block
	// check maximality with BItVector

}

// func availableBlocks(a, d []int) []block {
// 	available_blocks := make([]block, 0, len(a))
// 	curr_indexes := make([]int, 0, len(a))
// 	// scorri d e trova le c >= k - min_x

// 	best_divergence := 0

// 	for i := 1; i < len(a); i++ {

// 		if d[i] <= d[0]-min_block_rows {
// 			curr_indexes = append(curr_indexes, a[i-1])
// 			if best_divergence <= d[i] {
// 				best_divergence = d[i]
// 			} else { //skip ??/
// 			}
// 		} else {
// 			curr_indexes = append(curr_indexes, a[i-1])

// 			if len(curr_indexes) >= min_block_rows {
// 				available_blocks = append(available_blocks, makeblock(best_divergence, d[0], curr_indexes))
// 			}

// 			curr_indexes = make([]int, 0, len(a))
// 			best_divergence = 0

// 		}
// 	}

// 	return available_blocks
// }

// func allblocks(a, d []int, k int)

// func findblocks(a, d []int, min_divergence int) []block {

// }

func makeblock(i, j int, ind []int) block {
	indexes := ind
	var output block
	output.i = i
	output.j = j

	// any item can be sorted
	fmt.Println("beg ", output.i, "\nin", indexes)
	sort.Ints(indexes)
	x := 0
	for y := 1; y < len(indexes); y++ {
		if indexes[x] != indexes[y] {
			x++
			// preserve the original data
			// in[i], in[j] = in[j], in[i]
			// only set what is required
			indexes[x] = indexes[y]
		}
	}
	output.k = indexes[:x+1]
	fmt.Println("out ", output.k, " BIgOut ", indexes)
	// output.k = indexes
	return output
}

func printMatrixForAk(matrix []string, ak []int) {

}

func printMatrixAtK(matrix []string, k int) {

}
