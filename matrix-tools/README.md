# Scripts to manage binary matrices

## Subset a binary matrix

Extract a portion of a given binary matrix.

``` sh
python subset_matrix.py x y dx dy < IN-mat > OUT-mat
```

where:

* `x` first column (0-based)
* `y` first line (0-based)
* `dx` no. of columns
* `dy` no. of lines


## Add wildcards

Add wildcards to a binary matrix.

``` sh
python addwildcards.py IN-mat wr
```

where:

* `IN-mat` is the input binary matrix
* `wr` is the wildcard rate (between 0 and 1)

---

Tools provided by Yuri Pirola @yp