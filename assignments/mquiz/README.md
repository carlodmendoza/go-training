# NAME

mquiz

# SYNOPSIS

mquiz -csv [file] -n [count]

# DESCRIPTION

The program randomly selects 10 (or n as specified in the command line) questions from the database and asks the user to input his/her answer via standard input. The program then counts the number of correct answers and reports it at the end of the program. The database is a csv file named problems.csv (or a file specified in the command line). The database format is [question],[answer].

## Options

-csv file

: change the default database to file

-n count

: change the default number of quiz questions

# SAMPLE INTERACTION

```
Q: 10+3 = 13
Q: 4*5 = 20
Q: 1+1 = 2
Q: 5*10 = 50
Q: 7+6 = 13
Q: 4*5 = 20
Q: 11+1 = 12
Q: 50*10 = 500
Q: 7+16 = 23
Q: 3*7 = 20
You answered 9 out of 10 questions correctly.
```

# ERROR HANDLING

Always report errors encountered (e.g., incorrect database format, insufficient questions).
