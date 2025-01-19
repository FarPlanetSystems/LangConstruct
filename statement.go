package main

import (
	"strconv"
	"strings"
)

type Statement struct {
	rule_name   string
	conclusion  string
	params      []string
	premises    []string
	line int
}

func create_statement(rule_name string, concusion string, params []string, premises []string, line int, project *LC_project) Statement {
	res := Statement{
		rule_name:  rule_name,
		conclusion: concusion,
		params:     params,
		premises:   premises,
		line: line,
	}
	project.all_statements = append(project.all_statements, res)
	return res
}

func verify_statement(statement Statement, project *LC_project) bool{
	var applied_rule Rule
	for i := 0; i < len(project.all_rules); i++{
		if project.all_rules[i].name == statement.rule_name{
			applied_rule = deep_copy_rule(project.all_rules[i])
		}
	}
	if applied_rule.name == ""{
		message("no rule " + statement.rule_name + " was found. Line " + strconv.Itoa(statement.line), project)
		return false
	}
	if len(applied_rule.params) != len(statement.params) {
		message("derriving a statement, there must be as many parameters as there defined in the applied rule. Line "  + strconv.Itoa(statement.line), project)
		return false
	}
	if len(applied_rule.premises) != len(statement.premises) {
		message("derriving a statement, there must be as many premises as there defined in the applied rule. Line "  + strconv.Itoa(statement.line), project)
		return false
	}
	if !check_rule_applicability(statement, applied_rule, project){
		//message("the rule is unapplicable. Line " + strconv.Itoa(statement.line), project)
		return false
	}
	if !are_premises_verified(statement.premises, *project){
		message("not all premises are verified. Line " + strconv.Itoa(statement.line), project)
		return false
	}
	project.all_legal_expressions = append(project.all_legal_expressions, statement.conclusion)
	return true
}

func substitude_rule_with_params(statement Statement, rule Rule) Rule {
	substituted_rule := rule
	for i := 0; i<len(substituted_rule.params); i++{
		consequence := "[" + substituted_rule.params[i] + "]"
		// replacing params signs in premises of the rule (rule.params) with expressions in statements as arguments (statement.params)
		for j := 0; j < len(substituted_rule.premises); j++{
			substituted_rule.premises[j] = strings.Replace(substituted_rule.premises[j], consequence, statement.params[i], -1)
		}
		for j := 0; j < len(substituted_rule.conclusions); j++{
			substituted_rule.conclusions[j] = strings.Replace(substituted_rule.conclusions[j], consequence, statement.params[i], -1)
		}
	}

	return substituted_rule
}

func check_rule_applicability(statement Statement, rule Rule, project *LC_project) bool{
	substituted_rule := substitude_rule_with_params(statement, rule)
	// checking if there is correspondece with the statement's conclusion with one of the rule's conclusion
	correspondece_found := false

	for i := 0; i < len(substituted_rule.conclusions); i++{

		if substituted_rule.conclusions[i] == statement.conclusion{
			correspondece_found = true
		}
	}
	if !correspondece_found{
		msg_line := "conclusion "+ statement.conclusion + " does not correspond to any conclusion of the rule " + substituted_rule.name + ". Line " + strconv.Itoa(statement.line) + "\n See:"
		message(msg_line, project)
		for i := 0; i<len(substituted_rule.conclusions); i++{
			message(substituted_rule.conclusions[i], project)
		}
		return false
	}
	// checking the correspondence among premises
	for i := 0; i < len(substituted_rule.premises); i++{
		if substituted_rule.premises[i] != statement.premises[i]{
			msg_line := "a premise "+ statement.conclusion + " does not correspond to the required one " + substituted_rule.name + ". Line " + strconv.Itoa(statement.line) + "\n See:"
			message(msg_line, project)
			message(substituted_rule.premises[i] + " was expected, but " + statement.premises[i] + " was found", project)
			return false
		}
	}
	return true
}

func are_premises_verified(premises []string, project LC_project) bool{
	for i := 0; i < len(premises); i++{
		is_premise_found := false
		for j:=0; j < len(project.all_legal_expressions); j++{
			if project.all_legal_expressions[j] == premises[i]{
				is_premise_found = true
				break
			}
		}
		if !is_premise_found{
			return false
		}
	}
	return true
}