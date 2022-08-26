# mapq 语法

仅做参考，与实际不一定一致（懒得改了）  
bool_exp: BE->C ((AND|OR) BE)?
compare_exp: C->(B) ((EQ|NEQ|LG|SM|LEQ|SEQ) (B))*
boolean: B->TRUE|FALSE|(LP BE RP)|NOT B
added_factor: AF->F((ADD|MIN)F)*
factor: F->S|S((MUL|DIV|PS)S)*
symbol: S->N|((ADD|MIN) N)
number: N->n|VB|str|LP BE RP
varblock: VB->var (DOT var)*