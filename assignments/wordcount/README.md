# NAME

wordcount

# SYNOPSIS

wordcount file1 file2 ...

# DESCRIPTION

The program reads the contents of the files passed as arguments and displays the words found in the files, each with corresponding frequency (list them in alphabetical order). All letters in each word must be transformed to lower case. Define words as consecutive sequence of letters or digits. Other symbols such as punctuations should be ignored. Do the word counting concurrently using goroutines, channels, and sync.Mutex.  The files must be opened and processed concurrently.

# SAMPLE OUTPUT

```
alligator 1
alpha 15
zebra 3
```

# ERROR HANDLING

Always report errors encountered (e.g., cannot open file).

