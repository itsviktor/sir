.PHONY: dev grammar

dev:
	air

grammar:
	python grammar/transformGrammar.py
	java -jar tools/antlr/antlr-4.13.2-complete.jar -Dlanguage=Go -package parser -visitor -Xexact-output-dir -o src/internal/parser grammar/postgres/*.g4