# mapq 语法


bool_exp: BE->C ((AND|OR) BE)?
compare_exp: C->(B) ((EQ|NEQ|LG|SM|LEQ|SEQ) (B))*
boolean: B->TRUE|FALSE|(LP BE RP)|NOT B
added_factor: AF->F((ADD|MIN)F)*
factor: F->S|S((MUL|DIV|PS)S)*
symbol: S->N|((ADD|MIN) N)
number: N->n|var|str|LP BE RP