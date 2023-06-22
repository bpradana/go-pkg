# Go Pkg
A collection of useful Go packages.

## Why?
My [friend](https://github.com/akhidnukhlis) and I are working on a project that requires us to write a lot of helper functions. We decided to create a package to store all of our helper functions so that we can reuse them in the future. Today (2022-06-22), he told me that he wants to make the packages public. So, here it is.

## Available Packages
### qbuildr: A package for building SQL queries.
```go
structParams := struct {
  ID        string `json:"id,omitempty"`
  UserName  string `json:"user_name,omitempty"`
  Age       int    `json:"age,omitempty"`
  StartDate string `json:"start_date,omitempty"`
  EndDate   string `json:"end_date,omitempty"`
}{
  ID:        "1",
  UserName:  "John",
  Age:       20,
  StartDate: "2020-01-01",
  EndDate:   "2020-01-31",
}

query := qbuildr.NewQueryBuilder("transactions").Select("*").WhereStruct(structParams).Build()
fmt.Println(query) // SELECT * FROM transactions WHERE user_name = 'John' AND age = 20 AND created_at >= '2020-01-01' AND created_at <= '2020-01-31' AND id = '1'
```
