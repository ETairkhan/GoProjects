package main

//#include <stdio.h>
//void callC(){
//	printf("Calling C code\n"); 
//}

import (
	"C"
	"fmt"
) 


func main(){
	fmt.Println("A go statement")
	C.callC()
	fmt.Println("another go statement ")
}