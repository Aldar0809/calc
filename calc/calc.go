package calc

import (
	"fmt"
	"strconv"
)

// Calc вычисляет арифметическое выражение
func Calc(expression string) (float64, error) {
	tokens := tokenize(expression)
	if len(tokens) == 0 {
		return 0, fmt.Errorf("пустое выражение")
	}
	postfix, err := shuntingYard(tokens)
	if err != nil {
		return 0, err
	}
	result, err := evaluatePostfix(postfix)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func tokenize(expression string) []string {
	var tokens []string
	var currentToken string
	for _, char := range expression {
		if char >= '0' && char <= '9' || char == '.' {
			currentToken += string(char)
		} else if char == '+' || char == '-' || char == '*' || char == '/' || char == '(' || char == ')' {
			if currentToken != "" {
				tokens = append(tokens, currentToken)
				currentToken = ""
			}
			tokens = append(tokens, string(char))
		} else if char != ' ' {
			return nil
		}
	}
	if currentToken != "" {
		tokens = append(tokens, currentToken)
	}
	return tokens
}

func shuntingYard(tokens []string) ([]string, error) {
	var output []string
	var operatorStack []string
	precedence := map[string]int{"+": 1, "-": 1, "*": 2, "/": 2}

	for _, token := range tokens {
		if _, err := strconv.ParseFloat(token, 64); err == nil {
			output = append(output, token)
		} else if token == "(" {
			operatorStack = append(operatorStack, token)
		} else if token == ")" {
			for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1] != "(" {
				output = append(output, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			if len(operatorStack) == 0 {
				return nil, fmt.Errorf("некорректное выражение: непарные скобки")
			}
			operatorStack = operatorStack[:len(operatorStack)-1]
		} else if token == "+" || token == "-" || token == "*" || token == "/" {
			for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1] != "(" && precedence[token] <= precedence[operatorStack[len(operatorStack)-1]] {
				output = append(output, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			operatorStack = append(operatorStack, token)
		} else {
			return nil, fmt.Errorf("недопустимый токен: %s", token)
		}
	}

	for len(operatorStack) > 0 {
		if operatorStack[len(operatorStack)-1] == "(" {
			return nil, fmt.Errorf("некорректное выражение: непарные скобки")
		}
		output = append(output, operatorStack[len(operatorStack)-1])
		operatorStack = operatorStack[:len(operatorStack)-1]
	}

	return output, nil
}

func evaluatePostfix(tokens []string) (float64, error) {
	var stack []float64
	for _, token := range tokens {
		if val, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, val)
		} else {
			if len(stack) < 2 {
				return 0, fmt.Errorf("некорректное выражение")
			}
			operand2 := stack[len(stack)-1]
			operand1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			switch token {
			case "+":
				stack = append(stack, operand1+operand2)
			case "-":
				stack = append(stack, operand1-operand2)
			case "*":
				stack = append(stack, operand1*operand2)
			case "/":
				if operand2 == 0 {
					return 0, fmt.Errorf("деление на ноль")
				}
				stack = append(stack, operand1/operand2)
			default:
				return 0, fmt.Errorf("недопустимый оператор: %s", token)
			}
		}
	}
	if len(stack) != 1 {
		return 0, fmt.Errorf("некорректное выражение")
	}
	return stack[0], nil
}
