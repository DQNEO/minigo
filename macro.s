// Macros
.macro FUNC_PROLOGUE
  push %rbp
  mov %rsp, %rbp
.endm

.macro POP_TO_ARG_0
  pop %rdi
.endm

.macro POP_TO_ARG_1
  pop %rsi
.endm

.macro POP_TO_ARG_2
  pop %rdx
.endm

.macro POP_TO_ARG_3
  pop %rcx
.endm

.macro POP_TO_ARG_4
  pop %r8
.endm

.macro POP_TO_ARG_5
  pop %r9
.endm

.macro POP_TO_ARG_6
  pop %r10
.endm

.macro POP_TO_ARG_7
  pop %r11
.endm

.macro POP_TO_ARG_8
  pop %r12
.endm

.macro POP_TO_ARG_9
  pop %r13
.endm

.macro POP_TO_ARG_10
  pop %r14
.endm

.macro POP_TO_ARG_11
  pop %r15
.endm

.macro PUSH_ARG_0
  push %rdi
.endm

.macro PUSH_ARG_1
  push %rsi
.endm

.macro PUSH_ARG_2
  push %rdx
.endm

.macro PUSH_ARG_3
  push %rcx
.endm

.macro PUSH_ARG_4
  push %r8
.endm

.macro PUSH_ARG_5
  push %r9
.endm

.macro PUSH_ARG_6
  push %r10
.endm

.macro PUSH_ARG_7
  push %r11
.endm

.macro PUSH_ARG_8
  push %r12
.endm

.macro PUSH_ARG_9
  push %r13
.endm

.macro PUSH_ARG_10
  push %r14
.endm

.macro PUSH_ARG_11
  push %r15
.endm

.macro PUSH_8
  push %rax # primitive
.endm

.macro PUSH_24
  push %rax # 1st
  push %rbx # 2nd
  push %rcx # 3rd
.endm

.macro PUSH_SLICE
  push %rax # slice.ptr
  push %rbx # slice.len
  push %rcx # slice.cap
.endm

.macro PUSH_INTERFACE
  push %rax # ifc.1st
  push %rbx # ifc.2nd
  push %rcx # ifc.3rd
.endm

.macro POP_8
  pop %rax # primitive
.endm

.macro POP_24
  pop %rcx # 3rd
  pop %rbx # 2nd
  pop %rax # 1st
.endm

.macro POP_SLICE
  pop %rcx # slice.cap
  pop %rbx # slice.len
  pop %rax # slice.ptr
.endm

.macro POP_MAP
  pop %rcx # map.cap
  pop %rbx # map.len
  pop %rax # map.ptr
.endm

.macro POP_INTERFACE
  pop %rcx # ifc.3rd
  pop %rbx # ifc.2nd
  pop %rax # ifc.1st
.endm

.macro LOAD_EMPTY_24
  mov $0, %rax
  mov $0, %rbx
  mov $0, %rcx
.endm

.macro LOAD_EMPTY_SLICE
  mov $0, %rax
  mov $0, %rbx
  mov $0, %rcx
.endm

.macro LOAD_EMPTY_MAP
  mov $0, %rax
  mov $0, %rbx
  mov $0, %rcx
.endm

.macro LOAD_EMPTY_INTERFACE
  mov $0, %rax
  mov $0, %rbx
  mov $0, %rcx
.endm

.macro LOAD_STRING_LITERAL slabel
  lea \slabel(%rip), %rax
.endm

.macro LOAD_NUMBER n
  mov $\n, %rax
.endm

.macro STORE_1_TO_LOCAL offset
  mov %al, \offset(%rbp)
.endm

.macro STORE_2_TO_LOCAL offset
  mov %ax, \offset(%rbp)
.endm

.macro STORE_8_TO_LOCAL offset
  mov %rax, \offset(%rbp)
.endm

.macro LOAD_GLOBAL_ADDR varname, offset
  lea \varname+\offset(%rip), %rax
.endm

.macro LOAD_LOCAL_ADDR offset
  lea \offset(%rbp), %rax
.endm

.macro LOAD_1_FROM_LOCAL_CAST offset
  movsbq \offset(%rbp), %rax
.endm

.macro LOAD_2_FROM_LOCAL_CAST offset
  movswq \offset(%rbp), %rax
.endm

.macro LOAD_1_FROM_LOCAL offset
  mov \offset(%rbp), %al
.endm

.macro LOAD_2_FROM_LOCAL offset
  mov \offset(%rbp), %ax
.endm

.macro LOAD_8_FROM_LOCAL offset
  mov \offset(%rbp), %rax
.endm

.macro STORE_1_TO_GLOBAL varname, offset
  mov %al, \varname+\offset(%rip)
.endm

.macro STORE_2_TO_GLOBAL varname, offset
  mov %ax, \varname+\offset(%rip)
.endm

.macro STORE_8_TO_GLOBAL varname, offset
  mov %rax, \varname+\offset(%rip)
.endm

.macro LOAD_1_FROM_GLOBAL_CAST varname, offset=0
  movsbq \varname+\offset(%rip), %rax
.endm

.macro LOAD_2_FROM_GLOBAL_CAST varname, offset=0
  movswq \varname+\offset(%rip), %rax
.endm

.macro LOAD_1_FROM_GLOBAL varname, offset=0
  mov \varname+\offset(%rip), %al
.endm

.macro LOAD_2_FROM_GLOBAL varname, offset=0
  mov \varname+\offset(%rip), %ax
.endm

.macro LOAD_8_FROM_GLOBAL varname, offset=0
  mov \varname+\offset(%rip), %rax
.endm

.macro LOAD_24_BY_DEREF
  mov 16(%rax), %rcx
  mov 8(%rax), %rbx
  mov 0(%rax), %rax
.endm

.macro LOAD_8_BY_DEREF
  mov (%rax), %rax
.endm

.macro LOAD_1_BY_DEREF
  movsbq (%rax), %rax
.endm

.macro LOAD_2_BY_DEREF
  movswq (%rax), %rax
.endm

.macro LOAD_24_FROM_GLOBAL varname
  mov \varname+0(%rip), %rax # 1st
  mov \varname+8(%rip), %rbx # 2nd
  mov \varname+16(%rip), %rcx # 3rd
.endm

.macro LOAD_24_FROM_LOCAL offset
  mov \offset+0(%rbp), %rax # 1st
  mov \offset+8(%rbp), %rbx # 2nd
  mov \offset+16(%rbp), %rcx # 3rd
.endm

.macro CAST_UINT8_TO_INT
  movzbq %al, %rax
.endm

.macro CAST_UINT16_TO_INT
  movzwq %al, %rax
.endm

.macro CMP_EQ_ZERO
  cmp $0, %rax
  sete %al
  movzb %al, %eax
.endm

.macro CMP_NE_ZERO
  cmp $0, %rax
  setne %al
  movzb %al, %eax
.endm

.macro CMP_FROM_STACK inst
  pop %rax # right
  pop %rcx # left
  cmp %rax, %rcx
  \inst %al
  movzb %al, %eax
.endm

.macro SUM_FROM_STACK
  pop %rcx
  pop %rax
  add %rcx , %rax
.endm

.macro SUB_FROM_STACK
  pop %rcx
  pop %rax
  sub %rcx , %rax
.endm

.macro IMUL_FROM_STACK
  pop %rcx
  pop %rax
  imul %rcx , %rax
.endm

.macro IMUL_NUMBER n
  imul $\n , %rax
.endm

.macro STORE_1_INDIRECT_FROM_STACK
  pop %rax # where
  pop %rcx # what
  mov %cl, (%rax)
.endm

.macro STORE_2_INDIRECT_FROM_STACK
  pop %rax # where
  pop %rcx # what
  mov %cx, (%rax)
.endm

.macro STORE_8_INDIRECT_FROM_STACK
  pop %rax # where
  pop %rcx # what
  mov %rcx, (%rax)
.endm

.macro STORE_24_INDIRECT_FROM_STACK
  pop %rax # target addr
  pop %rcx # load RHS value(c)
  mov %rcx, 16(%rax)
  pop %rcx # load RHS value(b)
  mov %rcx, 8(%rax)
  pop %rcx # load RHS value(a)
  mov %rcx, 0(%rax)
.endm

.macro ADD_NUMBER n
  add $\n , %rax
.endm

.macro SUB_NUMBER n
  sub $\n , %rax
.endm

.macro FUNCALL fname
  mov $0, %rax
  mov $0, %rbx
  call \fname
.endm

.macro TEST_IT
  test %rax, %rax
.endm

.macro LEAVE_AND_RET
  leave
  ret
.endm
