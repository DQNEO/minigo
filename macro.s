// Macros
.macro FUNC_PROLOGUE
  pushq %rbp
  movq %rsp, %rbp
.endm

.macro POP_TO_ARG_0
  popq %rdi
.endm

.macro POP_TO_ARG_1
  popq %rsi
.endm

.macro POP_TO_ARG_2
  popq %rdx
.endm

.macro POP_TO_ARG_3
  popq %rcx
.endm

.macro POP_TO_ARG_4
  popq %r8
.endm

.macro POP_TO_ARG_5
  popq %r9
.endm

.macro POP_TO_ARG_6
  popq %r10
.endm

.macro POP_TO_ARG_7
  popq %r11
.endm

.macro POP_TO_ARG_8
  popq %r12
.endm

.macro POP_TO_ARG_9
  popq %r13
.endm

.macro POP_TO_ARG_10
  popq %r14
.endm

.macro POP_TO_ARG_11
  popq %r15
.endm

.macro PUSH_ARG_0
  pushq %rdi
.endm

.macro PUSH_ARG_1
  pushq %rsi
.endm

.macro PUSH_ARG_2
  pushq %rdx
.endm

.macro PUSH_ARG_3
  pushq %rcx
.endm

.macro PUSH_ARG_4
  pushq %r8
.endm

.macro PUSH_ARG_5
  pushq %r9
.endm

.macro PUSH_ARG_6
  pushq %r10
.endm

.macro PUSH_ARG_7
  pushq %r11
.endm

.macro PUSH_ARG_8
  pushq %r12
.endm

.macro PUSH_ARG_9
  pushq %r13
.endm

.macro PUSH_ARG_10
  pushq %r14
.endm

.macro PUSH_ARG_11
  pushq %r15
.endm

.macro PUSH_8
  pushq %rax # primitive
.endm

.macro PUSH_24
  pushq %rax # 1st
  pushq %rbx # 2nd
  pushq %rcx # 3rd
.endm

.macro PUSH_SLICE
  pushq %rax # slice.ptr
  pushq %rbx # slice.len
  pushq %rcx # slice.cap
.endm

.macro PUSH_INTERFACE
  pushq %rax # ifc.1st
  pushq %rbx # ifc.2nd
  pushq %rcx # ifc.3rd
.endm

.macro POP_8
  popq %rax # primitive
.endm

.macro POP_24
  popq %rcx # 3rd
  popq %rbx # 2nd
  popq %rax # 1st
.endm

.macro POP_SLICE
  popq %rcx # slice.cap
  popq %rbx # slice.len
  popq %rax # slice.ptr
.endm

.macro POP_MAP
  popq %rcx # map.cap
  popq %rbx # map.len
  popq %rax # map.ptr
.endm

.macro POP_INTERFACE
  popq %rcx # ifc.3rd
  popq %rbx # ifc.2nd
  popq %rax # ifc.1st
.endm

.macro LOAD_EMPTY_24
  movq $0, %rax
  movq $0, %rbx
  movq $0, %rcx
.endm

.macro LOAD_EMPTY_SLICE
  movq $0, %rax
  movq $0, %rbx
  movq $0, %rcx
.endm

.macro LOAD_EMPTY_MAP
  movq $0, %rax
  movq $0, %rbx
  movq $0, %rcx
.endm

.macro LOAD_EMPTY_INTERFACE
  movq $0, %rax
  movq $0, %rbx
  movq $0, %rcx
.endm

.macro LOAD_STRING_LITERAL slabel
  leaq \slabel(%rip), %rax
.endm

.macro LOAD_NUMBER n
  movq $\n, %rax
.endm

.macro STORE_1_TO_LOCAL offset
  movb %al, \offset(%rbp)
.endm

.macro STORE_2_TO_LOCAL offset
  movw %ax, \offset(%rbp)
.endm

.macro STORE_8_TO_LOCAL offset
  movq %rax, \offset(%rbp)
.endm

.macro LOAD_GLOBAL_ADDR varname, offset
  leaq \varname+\offset(%rip), %rax
.endm

.macro LOAD_LOCAL_ADDR offset
  leaq \offset(%rbp), %rax
.endm

.macro LOAD_1_FROM_LOCAL_CAST offset
  movsbq \offset(%rbp), %rax
.endm

.macro LOAD_2_FROM_LOCAL_CAST offset
  movswq \offset(%rbp), %rax
.endm

.macro LOAD_1_FROM_LOCAL offset
  movb \offset(%rbp), %al
.endm

.macro LOAD_2_FROM_LOCAL offset
  movw \offset(%rbp), %ax
.endm

.macro LOAD_8_FROM_LOCAL offset
  movq \offset(%rbp), %rax
.endm

.macro STORE_1_TO_GLOBAL varname, offset
  movb %al, \varname+\offset(%rip)
.endm

.macro STORE_2_TO_GLOBAL varname, offset
  movw %ax, \varname+\offset(%rip)
.endm

.macro STORE_8_TO_GLOBAL varname, offset
  movq %rax, \varname+\offset(%rip)
.endm

.macro LOAD_1_FROM_GLOBAL_CAST varname, offset=0
  movsbq \varname+\offset(%rip), %rax
.endm

.macro LOAD_2_FROM_GLOBAL_CAST varname, offset=0
  movswq \varname+\offset(%rip), %rax
.endm

.macro LOAD_1_FROM_GLOBAL varname, offset=0
  movb \varname+\offset(%rip), %al
.endm

.macro LOAD_2_FROM_GLOBAL varname, offset=0
  movw \varname+\offset(%rip), %ax
.endm

.macro LOAD_8_FROM_GLOBAL varname, offset=0
  movq \varname+\offset(%rip), %rax
.endm

.macro LOAD_24_BY_DEREF
  movq 16(%rax), %rcx
  movq 8(%rax), %rbx
  movq 0(%rax), %rax
.endm

.macro LOAD_8_BY_DEREF
  movq (%rax), %rax
.endm

.macro LOAD_1_BY_DEREF
  movsbq (%rax), %rax
.endm

.macro LOAD_2_BY_DEREF
  movswq (%rax), %rax
.endm

.macro LOAD_24_FROM_GLOBAL varname
  movq \varname+0(%rip), %rax # 1st
  movq \varname+8(%rip), %rbx # 2nd
  movq \varname+16(%rip), %rcx # 3rd
.endm

.macro LOAD_24_FROM_LOCAL offset
  movq \offset+0(%rbp), %rax # 1st
  movq \offset+8(%rbp), %rbx # 2nd
  movq \offset+16(%rbp), %rcx # 3rd
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
  popq %rax # right
  popq %rcx # left
  cmpq %rax, %rcx
  \inst %al
  movzb %al, %eax
.endm

.macro SUM_FROM_STACK
  popq %rcx
  popq %rax
  addq %rcx , %rax
.endm

.macro SUB_FROM_STACK
  popq %rcx
  popq %rax
  subq %rcx , %rax
.endm

.macro IMUL_FROM_STACK
  popq %rcx
  popq %rax
  imulq %rcx , %rax
.endm

.macro IMUL_NUMBER n
  imulq $\n , %rax
.endm

.macro STORE_1_INDIRECT_FROM_STACK
  popq %rax # where
  popq %rcx # what
  movb %cl, (%rax)
.endm

.macro STORE_2_INDIRECT_FROM_STACK
  popq %rax # where
  popq %rcx # what
  movw %cx, (%rax)
.endm

.macro STORE_8_INDIRECT_FROM_STACK
  popq %rax # where
  popq %rcx # what
  movq %rcx, (%rax)
.endm

.macro STORE_24_INDIRECT_FROM_STACK
  popq %rax # target addr
  popq %rcx # load RHS value(c)
  movq %rcx, 16(%rax)
  popq %rcx # load RHS value(b)
  movq %rcx, 8(%rax)
  popq %rcx # load RHS value(a)
  movq %rcx, 0(%rax)
.endm

.macro ADD_NUMBER n
  addq $\n , %rax
.endm

.macro SUB_NUMBER n
  subq $\n , %rax
.endm

.macro FUNCALL fname
  movq $0, %rax
  movq $0, %rbx
  callq \fname
.endm

.macro TEST_IT
  test %rax, %rax
.endm

.macro LEAVE_AND_RET
  leave
  ret
.endm
