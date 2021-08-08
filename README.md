# calc
Just trying to get a little hand at coding go with some homework for /r/programming_funny

## usage 
calc "1 +2 +3" or echo "1 + 2" | calc

## flags:
- \-s preceedes the outcome with the input and =

## operators supported:
- \+ plus
- \- minus
- / divide
- \* multiply
- % modulus
- V (uppercase v) square Root
- ^ power
- parentheses: (,)

example: V(2 ^ 8) or -2 + 2 + (5 % 3) * 12.7

## notes:
- stdin input streams are supported, intresting if you want to calculate multiple things, line per line
- negative numbers are supported obviously but will probably be interpreted as an unknown flag, to bypass echo the value to the program via a pipe, e.g. echo "-1 + 2" | calc
