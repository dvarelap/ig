package main

func checkIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
