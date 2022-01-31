# ASSIGNMENT #3: SOLVER OF 3 UNKNOWNS
---
# NAME

solve

# SYNOPSIS

http://localhost:8080/solve

# DESCRIPTION

The webapp solves a system of 3 unknowns (x, y, z).
```
a11 x + a12 y + a13 z = a14
a21 x + a22 y + a23 z = a24
a31 x + a32 y + a33 z = a34
```
The coefficients of the equations are posted to the webapp via query string and then the webapp prints the system of equations and the solution if there are any.  If the system is inconsistent (i.e, no solution), the webapp must print 'inconsistent - no solution'. If the system is dependent (i.e., with multiple solutions), the webapp must print 'dependent - with multiple solutions' The coefficients must be inputted in row major.

# SAMPLE INTERACTION

```
http://localhost:8080/solve?coef=4,5,6,7,2,3,1,2,1,2,3,2

system:
4x + 5y + 6z = 7
2x + 3y + 1z = 2
1x + 2y + 3z = 2

solution:
x = 1.89, y = -0.78, z = 0.56
```
```
http://localhost:8080/solve?coef=1,-3,1,4,-1,2,-5,3,5,-13,13,8

system:
1x + -3y + 1z = 4
-1x + 2y + -5z = 3
5x + -13y + 13z = 8

inconsistent - no solution
```

# ERROR HANDLING

Always report errors encountered (e.g., incorrect url, incorrect number of coefficients).
