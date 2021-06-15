package main

import (
	"fmt"
	"lexer/parser"
	"strconv"
	"strings"

	// "lexer/ast"
	"encoding/json"
	"lexer/generator"
)

const JsonQuote = '"'
const FALSE_LEN = len("false")
const TRUE_LEN = len("true")
const NULL_LEN = len("null")

func JsonWhitespace() []rune {
	return []rune{' ', '\t', '\b', '\n', '\r'}
}
func JsonSyntax() []rune {
	return []rune{',', ':', '{', '}', '[', ']'}
}

func lex_number(inputstr string) (interface{}, string) {
	json_number := ""
	number_characters := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', 'e', '.'}
	input := strings.Split(inputstr, "")
	for _, c := range input {
		if Search(number_characters, c) != -1 {
			json_number += c
		} else {
			break
		}
	}
	if len(json_number) == 0 {
		return nil, inputstr
	}
	inputstr = inputstr[len(json_number):]
	i, err := strconv.Atoi(json_number)
	if err == nil {
		return i, inputstr
	} else {
		i, err := strconv.ParseFloat(json_number, 64)
		if err == nil {
			return i, inputstr
		}
	}
	// if '.' in json_number{
	//     return float(json_number), rest
	// }
	return nil, inputstr
}

func lex_bool(inputstr string) (interface{}, string) {
	string_len := len(inputstr)
	if string_len >= TRUE_LEN && inputstr[:TRUE_LEN] == "true" {
		return true, inputstr[TRUE_LEN:]
	} else if string_len >= FALSE_LEN && inputstr[:FALSE_LEN] == "false" {
		return false, inputstr[FALSE_LEN:]
	}
	return nil, inputstr
}

func lex_null(inputstr string) (*string, string) {
	string_len := len(inputstr)
	if string_len >= TRUE_LEN && inputstr[:TRUE_LEN] == "null" {
		s := "null"
		return &s, inputstr[TRUE_LEN:]
	}
	return nil, inputstr
}

// foo" test" test "
// "{"foo": [1, 2, {"bar": 2}]}"
func lex_string(inputstr string) (string, string) {
	var json_string string
	if inputstr[0] == JsonQuote {
		inputstr = inputstr[1:]
	} else {
		return "", inputstr
	}

	for _, c := range inputstr {
		if c == JsonQuote {
			return json_string, inputstr[len(json_string)+1:]
		} else {
			quoted := strconv.QuoteRune(c)
			json_string += quoted[1 : len(quoted)-1]
		}
	}
	return "", ""
}

func lex(inputstr string) []interface{} {
	var tokens []interface{}
	for len(inputstr) > 0 {
		//Check if value is string
		var json_string string
		json_string, inputstr = lex_string(inputstr)
		if json_string != "" {
			tokens = append(tokens, json_string)
			continue
		}

		//Check if value is num
		var num interface{}
		num, inputstr = lex_number(inputstr)
		if num != nil {
			switch i := num.(type) {
			case int:
				tokens = append(tokens, i)
			case float64:
				tokens = append(tokens, i)
			}
			continue
		}

		//Check if value is bool
		var json_bool interface{}
		json_bool, inputstr = lex_bool(inputstr)
		if json_bool != nil {
			tokens = append(tokens, json_bool)
			continue
		}
		var json_null *string
		json_null, inputstr = lex_null(inputstr)
		if json_null != nil {
			tokens = append(tokens, []interface{}{nil}...)
			continue
		}
		c := strings.Split(inputstr, "")[0]

		//Ignore whitespace
		if Search(JsonWhitespace(), c) != -1 {
			inputstr = inputstr[1:]
		} else if Search(JsonSyntax(), c) != -1 { //Append syntax
			tokens = append(tokens, c)
			inputstr = inputstr[1:]
		} else { //Skip invalid characters
			inputstr = inputstr[1:]
			break
		}
	}
	tokens = append(tokens, []interface{}{"EOF"}...)
	return tokens
}

func Search(text []rune, what string) int {
	whatRunes := []rune(what)

	for i := range text {
		found := true
		for j := range whatRunes {
			if text[i+j] != whatRunes[j] {
				found = false
				break
			}
		}
		if found {
			return i
		}
	}
	return -1
}

func main() {
	tokens := lex(`{"resr":"1","ac":1.1,"abc":null,"isTrue":true,"integer":2, "extra" : { "field1" : "val1", "field2" : {"f1": "v2" }}, "array" : [12.3, 12]}`)

	fmt.Println("lexer output", tokens)
	fmt.Println("_______")
	p := parser.New(tokens)
	tree, err := p.ParseJSON()

	if err != nil {
		fmt.Println(err)
	}
	treeJSON, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("MarshalIndent funnction output %s\n", string(treeJSON))

	generator.JsonStruct(tree)
	// walkTree(tree)
}

// func walkTree(node ast.RootNode){
//     fmt.Println("---------------")
//     switch node.Type{
//     case ast.ObjectRoot:
//         // fmt.Printf("%#v",node.RootValue)
//         walkObj(node.RootValue.Content.(ast.Object))
//     case ast.ArrayRoot:
//         fmt.Println("Your root JSON type is an array.")
//     }
// }
// func walkObj(obj ast.Object,a ...interface{}){
//     if a!=nil{
//         fmt.Println("---CHILDREN OF",a)
//     }
//     for _,childObj := range obj.Children{
//         // fmt.Println("[Key: ",childObj.Key.Value," Value: ",childObj.Value,"]\n")
//         switch childObj.Value.(type) {
//         case ast.Value:
//             result := childObj.Value.(ast.Value)
//             fmt.Println("KEY:",childObj.Key.Value)
//             switch result.Content.(type) {
//             case ast.Literal:
//                 res := result.Content.(ast.Literal)
//                 fmt.Printf("\n VALUE %+v",res)
//             case ast.Object:
//                 res := result.Content.(ast.Object)
//                 walkObj(res,childObj.Key.Value)
//             case ast.Array:
//                 res := result.Content.(ast.Array)
//                 walkArr(res,childObj.Key.Value)
//             }
//         }
//         fmt.Println()
//     }
// }
// func walkArr(arr ast.Array,a ...interface{}){
//     if a!=nil{
//         fmt.Println("---CHILDREN OF",a)
//     }
//     for _,childObj := range arr.Children{
//         // fmt.Println("[Key: ",childObj.Key.Value," Value: ",childObj.Value,"]\n")
//         switch childObj.Value.(type) {
//             case ast.Literal:
//                 res := childObj.Value.(ast.Literal)
//                 fmt.Printf("\n VALUE %+v",res)
//             case ast.Object:
//                 res := childObj.Value.(ast.Object)
//                 walkObj(res)
//             case ast.Array:
//                 res := childObj.Value.(ast.Array)
//                 fmt.Println("Array",res)
//         }
//         fmt.Println()
//     }
// }
