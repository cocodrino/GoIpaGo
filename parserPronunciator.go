package main

import (
	"fmt"
	. "github.com/yhirose/go-peg"
)

type ElementType int

const (
	World ElementType = iota
	Symbol
	Space
)

type Element struct {
	elementType ElementType
	value string
}

type ParsePronunciator struct{
	parser *Parser
}


func NewParsePronunciator() (myparser *ParsePronunciator,err error){
	parser, err := NewParser(`
		# Grammar for simple calculator...
		DOCUMENT <- TEXT*
		TEXT <- WORD / SPACE / SYMBOL
		SYMBOL <- [.,?!#$%&/()='¿¡{};:] 
		WORD <- LETTER+
		SPACE <- ' '
		LETTER <- [a-zA-Z0-9]
	`)

	if err != nil{
		fmt.Println(err)
		return
	}

	g := parser.Grammar

	g["DOCUMENT"].Action = func(v *Values, d Any) (any Any, e error) {
		fmt.Println("running document")
		return v.Vs,nil
	}

	g["SYMBOL"].Action = func(v *Values, d Any) (any Any, e error) {
		return &Element{elementType:Symbol,value:v.S},nil
	}
	g["WORD"].Action = func(v *Values, d Any) (any Any, e error) {
		return &Element{elementType:World, value: v.S},nil
	}
	g["SPACE"].Action = func(v *Values, d Any) (any Any, e error) {
		return &Element{elementType:Space,value:" "},nil
	}

	myparser = &ParsePronunciator{parser: parser}

	/*
	input := "hola como estas"

	values, _ := parser.ParseAndGetValue(input, nil)

	var results []Element
	for _,v := range values.([]Any){
		elemento := (v).(*Element)
		results = append(results,*elemento)
	}

	fmt.Printf("%v", results)
	for _,v:=range results{
		fmt.Println(v.elementType)
	}

	*/

return
}

func (prs *ParsePronunciator) getElements(txt string)[]Element{
	values, _ := prs.parser.ParseAndGetValue(txt,nil)

	var results []Element
	for _,v := range values.([]Any){
		elemento := (v).(*Element)
		results = append(results,*elemento)
	}

	return results

}
