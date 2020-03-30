"".sum0 STEXT nosplit size=25 args=0x18 locals=0x0
	0x0000 00000 (sum.go:4)	TEXT	"".sum0(SB), NOSPLIT|ABIInternal, $0-24
	0x0000 00000 (sum.go:4)	PCDATA	$0, $-2
	0x0000 00000 (sum.go:4)	PCDATA	$1, $-2
	0x0000 00000 (sum.go:4)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (sum.go:4)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (sum.go:4)	FUNCDATA	$2, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (sum.go:4)	PCDATA	$0, $0
	0x0000 00000 (sum.go:4)	PCDATA	$1, $0
	0x0000 00000 (sum.go:4)	MOVQ	$0, "".~r2+24(SP)
	0x0009 00009 (sum.go:5)	MOVQ	"".a+8(SP), AX
	0x000e 00014 (sum.go:5)	ADDQ	"".b+16(SP), AX
	0x0013 00019 (sum.go:5)	MOVQ	AX, "".~r2+24(SP)
	0x0018 00024 (sum.go:5)	RET
	0x0000 48 c7 44 24 18 00 00 00 00 48 8b 44 24 08 48 03  H.D$.....H.D$.H.
	0x0010 44 24 10 48 89 44 24 18 c3                       D$.H.D$..
"".sum1 STEXT nosplit size=52 args=0x18 locals=0x10
	0x0000 00000 (sum.go:8)	TEXT	"".sum1(SB), NOSPLIT|ABIInternal, $16-24
	0x0000 00000 (sum.go:8)	SUBQ	$16, SP
	0x0004 00004 (sum.go:8)	MOVQ	BP, 8(SP)
	0x0009 00009 (sum.go:8)	LEAQ	8(SP), BP
	0x000e 00014 (sum.go:8)	PCDATA	$0, $-2
	0x000e 00014 (sum.go:8)	PCDATA	$1, $-2
	0x000e 00014 (sum.go:8)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x000e 00014 (sum.go:8)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x000e 00014 (sum.go:8)	FUNCDATA	$2, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x000e 00014 (sum.go:8)	PCDATA	$0, $0
	0x000e 00014 (sum.go:8)	PCDATA	$1, $0
	0x000e 00014 (sum.go:8)	MOVQ	$0, "".~r2+40(SP)
	0x0017 00023 (sum.go:9)	MOVQ	"".a+24(SP), AX
	0x001c 00028 (sum.go:9)	ADDQ	"".b+32(SP), AX
	0x0021 00033 (sum.go:9)	MOVQ	AX, "".c(SP)
	0x0025 00037 (sum.go:10)	MOVQ	AX, "".~r2+40(SP)
	0x002a 00042 (sum.go:10)	MOVQ	8(SP), BP
	0x002f 00047 (sum.go:10)	ADDQ	$16, SP
	0x0033 00051 (sum.go:10)	RET
	0x0000 48 83 ec 10 48 89 6c 24 08 48 8d 6c 24 08 48 c7  H...H.l$.H.l$.H.
	0x0010 44 24 28 00 00 00 00 48 8b 44 24 18 48 03 44 24  D$(....H.D$.H.D$
	0x0020 20 48 89 04 24 48 89 44 24 28 48 8b 6c 24 08 48   H..$H.D$(H.l$.H
	0x0030 83 c4 10 c3                                      ....
"".assignstring STEXT nosplit size=61 args=0x20 locals=0x18
	0x0000 00000 (sum.go:13)	TEXT	"".assignstring(SB), NOSPLIT|ABIInternal, $24-32
	0x0000 00000 (sum.go:13)	SUBQ	$24, SP
	0x0004 00004 (sum.go:13)	MOVQ	BP, 16(SP)
	0x0009 00009 (sum.go:13)	LEAQ	16(SP), BP
	0x000e 00014 (sum.go:13)	PCDATA	$0, $-2
	0x000e 00014 (sum.go:13)	PCDATA	$1, $-2
	0x000e 00014 (sum.go:13)	FUNCDATA	$0, gclocals·9fad110d66c97cf0b58d28cccea80b12(SB)
	0x000e 00014 (sum.go:13)	FUNCDATA	$1, gclocals·d8b28f51bb91e05d264803f0f600a200(SB)
	0x000e 00014 (sum.go:13)	FUNCDATA	$2, gclocals·9fb7f0986f647f17cb53dda1484e0f7a(SB)
	0x000e 00014 (sum.go:13)	PCDATA	$0, $0
	0x000e 00014 (sum.go:13)	PCDATA	$1, $0
	0x000e 00014 (sum.go:13)	XORPS	X0, X0
	0x0011 00017 (sum.go:13)	MOVUPS	X0, "".~r1+48(SP)
	0x0016 00022 (sum.go:14)	PCDATA	$0, $1
	0x0016 00022 (sum.go:14)	MOVQ	"".a+32(SP), AX
	0x001b 00027 (sum.go:14)	PCDATA	$1, $1
	0x001b 00027 (sum.go:14)	MOVQ	"".a+40(SP), CX
	0x0020 00032 (sum.go:14)	MOVQ	AX, "".b(SP)
	0x0024 00036 (sum.go:14)	MOVQ	CX, "".b+8(SP)
	0x0029 00041 (sum.go:15)	PCDATA	$0, $0
	0x0029 00041 (sum.go:15)	PCDATA	$1, $2
	0x0029 00041 (sum.go:15)	MOVQ	AX, "".~r1+48(SP)
	0x002e 00046 (sum.go:15)	MOVQ	CX, "".~r1+56(SP)
	0x0033 00051 (sum.go:15)	MOVQ	16(SP), BP
	0x0038 00056 (sum.go:15)	ADDQ	$24, SP
	0x003c 00060 (sum.go:15)	RET
	0x0000 48 83 ec 18 48 89 6c 24 10 48 8d 6c 24 10 0f 57  H...H.l$.H.l$..W
	0x0010 c0 0f 11 44 24 30 48 8b 44 24 20 48 8b 4c 24 28  ...D$0H.D$ H.L$(
	0x0020 48 89 04 24 48 89 4c 24 08 48 89 44 24 30 48 89  H..$H.L$.H.D$0H.
	0x0030 4c 24 38 48 8b 6c 24 10 48 83 c4 18 c3           L$8H.l$.H....
"".concatestring STEXT size=127 args=0x30 locals=0x40
	0x0000 00000 (sum.go:18)	TEXT	"".concatestring(SB), ABIInternal, $64-48
	0x0000 00000 (sum.go:18)	MOVQ	(TLS), CX
	0x0009 00009 (sum.go:18)	CMPQ	SP, 16(CX)
	0x000d 00013 (sum.go:18)	PCDATA	$0, $-2
	0x000d 00013 (sum.go:18)	JLS	120
	0x000f 00015 (sum.go:18)	PCDATA	$0, $-1
	0x000f 00015 (sum.go:18)	SUBQ	$64, SP
	0x0013 00019 (sum.go:18)	MOVQ	BP, 56(SP)
	0x0018 00024 (sum.go:18)	LEAQ	56(SP), BP
	0x001d 00029 (sum.go:18)	PCDATA	$0, $-2
	0x001d 00029 (sum.go:18)	PCDATA	$1, $-2
	0x001d 00029 (sum.go:18)	FUNCDATA	$0, gclocals·2625d1fdbbaf79a2e52296235cb6527c(SB)
	0x001d 00029 (sum.go:18)	FUNCDATA	$1, gclocals·f6bd6b3389b872033d462029172c8612(SB)
	0x001d 00029 (sum.go:18)	FUNCDATA	$2, gclocals·1cf923758aae2e428391d1783fe59973(SB)
	0x001d 00029 (sum.go:18)	PCDATA	$0, $0
	0x001d 00029 (sum.go:18)	PCDATA	$1, $0
	0x001d 00029 (sum.go:18)	XORPS	X0, X0
	0x0020 00032 (sum.go:18)	MOVUPS	X0, "".~r2+104(SP)
	0x0025 00037 (sum.go:19)	MOVQ	$0, (SP)
	0x002d 00045 (sum.go:19)	PCDATA	$0, $1
	0x002d 00045 (sum.go:19)	MOVQ	"".a+72(SP), AX
	0x0032 00050 (sum.go:19)	PCDATA	$1, $1
	0x0032 00050 (sum.go:19)	MOVQ	"".a+80(SP), CX
	0x0037 00055 (sum.go:19)	PCDATA	$0, $0
	0x0037 00055 (sum.go:19)	MOVQ	AX, 8(SP)
	0x003c 00060 (sum.go:19)	MOVQ	CX, 16(SP)
	0x0041 00065 (sum.go:19)	PCDATA	$0, $1
	0x0041 00065 (sum.go:19)	MOVQ	"".b+88(SP), AX
	0x0046 00070 (sum.go:19)	PCDATA	$1, $2
	0x0046 00070 (sum.go:19)	MOVQ	"".b+96(SP), CX
	0x004b 00075 (sum.go:19)	PCDATA	$0, $0
	0x004b 00075 (sum.go:19)	MOVQ	AX, 24(SP)
	0x0050 00080 (sum.go:19)	MOVQ	CX, 32(SP)
	0x0055 00085 (sum.go:19)	CALL	runtime.concatstring2(SB)
	0x005a 00090 (sum.go:19)	MOVQ	48(SP), AX
	0x005f 00095 (sum.go:19)	PCDATA	$0, $2
	0x005f 00095 (sum.go:19)	MOVQ	40(SP), CX
	0x0064 00100 (sum.go:19)	PCDATA	$0, $0
	0x0064 00100 (sum.go:19)	PCDATA	$1, $3
	0x0064 00100 (sum.go:19)	MOVQ	CX, "".~r2+104(SP)
	0x0069 00105 (sum.go:19)	MOVQ	AX, "".~r2+112(SP)
	0x006e 00110 (sum.go:19)	MOVQ	56(SP), BP
	0x0073 00115 (sum.go:19)	ADDQ	$64, SP
	0x0077 00119 (sum.go:19)	RET
	0x0078 00120 (sum.go:19)	NOP
	0x0078 00120 (sum.go:18)	PCDATA	$1, $-1
	0x0078 00120 (sum.go:18)	PCDATA	$0, $-2
	0x0078 00120 (sum.go:18)	CALL	runtime.morestack_noctxt(SB)
	0x007d 00125 (sum.go:18)	PCDATA	$0, $-1
	0x007d 00125 (sum.go:18)	JMP	0
	0x0000 64 48 8b 0c 25 00 00 00 00 48 3b 61 10 76 69 48  dH..%....H;a.viH
	0x0010 83 ec 40 48 89 6c 24 38 48 8d 6c 24 38 0f 57 c0  ..@H.l$8H.l$8.W.
	0x0020 0f 11 44 24 68 48 c7 04 24 00 00 00 00 48 8b 44  ..D$hH..$....H.D
	0x0030 24 48 48 8b 4c 24 50 48 89 44 24 08 48 89 4c 24  $HH.L$PH.D$.H.L$
	0x0040 10 48 8b 44 24 58 48 8b 4c 24 60 48 89 44 24 18  .H.D$XH.L$`H.D$.
	0x0050 48 89 4c 24 20 e8 00 00 00 00 48 8b 44 24 30 48  H.L$ .....H.D$0H
	0x0060 8b 4c 24 28 48 89 4c 24 68 48 89 44 24 70 48 8b  .L$(H.L$hH.D$pH.
	0x0070 6c 24 38 48 83 c4 40 c3 e8 00 00 00 00 eb 81     l$8H..@........
	rel 5+4 t=17 TLS+0
	rel 86+4 t=8 runtime.concatstring2+0
	rel 121+4 t=8 runtime.morestack_noctxt+0
"".main STEXT size=578 args=0x0 locals=0x128
	0x0000 00000 (sum.go:22)	TEXT	"".main(SB), ABIInternal, $296-0
	0x0000 00000 (sum.go:22)	MOVQ	(TLS), CX
	0x0009 00009 (sum.go:22)	LEAQ	-168(SP), AX
	0x0011 00017 (sum.go:22)	CMPQ	AX, 16(CX)
	0x0015 00021 (sum.go:22)	PCDATA	$0, $-2
	0x0015 00021 (sum.go:22)	JLS	568
	0x001b 00027 (sum.go:22)	PCDATA	$0, $-1
	0x001b 00027 (sum.go:22)	SUBQ	$296, SP
	0x0022 00034 (sum.go:22)	MOVQ	BP, 288(SP)
	0x002a 00042 (sum.go:22)	LEAQ	288(SP), BP
	0x0032 00050 (sum.go:22)	PCDATA	$0, $-2
	0x0032 00050 (sum.go:22)	PCDATA	$1, $-2
	0x0032 00050 (sum.go:22)	FUNCDATA	$0, gclocals·f14a5bc6d08bc46424827f54d2e3f8ed(SB)
	0x0032 00050 (sum.go:22)	FUNCDATA	$1, gclocals·cde884c6f8ebff321c52f642fdb453a8(SB)
	0x0032 00050 (sum.go:22)	FUNCDATA	$2, gclocals·1cf923758aae2e428391d1783fe59973(SB)
	0x0032 00050 (sum.go:23)	PCDATA	$0, $0
	0x0032 00050 (sum.go:23)	PCDATA	$1, $0
	0x0032 00050 (sum.go:23)	MOVQ	$0, "".i+72(SP)
	0x003b 00059 (sum.go:24)	MOVQ	$2, "".a+112(SP)
	0x0044 00068 (sum.go:24)	MOVQ	$3, "".b+88(SP)
	0x004d 00077 (sum.go:24)	MOVQ	$0, "".~r2+64(SP)
	0x0056 00086 (<unknown line number>)	NOP
	0x0056 00086 (sum.go:5)	MOVQ	"".a+112(SP), AX
	0x005b 00091 (sum.go:5)	ADDQ	"".b+88(SP), AX
	0x0060 00096 (sum.go:24)	MOVQ	AX, ""..autotmp_15+120(SP)
	0x0065 00101 (sum.go:24)	MOVQ	AX, "".~r2+64(SP)
	0x006a 00106 (sum.go:24)	JMP	108
	0x006c 00108 (sum.go:24)	MOVQ	AX, "".i+72(SP)
	0x0071 00113 (sum.go:25)	CALL	runtime.printlock(SB)
	0x0076 00118 (sum.go:25)	MOVQ	"".i+72(SP), AX
	0x007b 00123 (sum.go:25)	MOVQ	AX, (SP)
	0x007f 00127 (sum.go:25)	CALL	runtime.printint(SB)
	0x0084 00132 (sum.go:25)	CALL	runtime.printnl(SB)
	0x0089 00137 (sum.go:25)	CALL	runtime.printunlock(SB)
	0x008e 00142 (sum.go:27)	MOVQ	$2, "".a+104(SP)
	0x0097 00151 (sum.go:27)	MOVQ	$3, "".b+96(SP)
	0x00a0 00160 (sum.go:27)	MOVQ	$0, "".~r2+56(SP)
	0x00a9 00169 (<unknown line number>)	NOP
	0x00a9 00169 (sum.go:9)	MOVQ	"".a+104(SP), AX
	0x00ae 00174 (sum.go:9)	ADDQ	"".b+96(SP), AX
	0x00b3 00179 (sum.go:9)	MOVQ	AX, "".c+80(SP)
	0x00b8 00184 (sum.go:27)	MOVQ	AX, "".~r2+56(SP)
	0x00bd 00189 (sum.go:27)	JMP	191
	0x00bf 00191 (sum.go:27)	MOVQ	AX, "".i+72(SP)
	0x00c4 00196 (sum.go:28)	CALL	runtime.printlock(SB)
	0x00c9 00201 (sum.go:28)	MOVQ	"".i+72(SP), AX
	0x00ce 00206 (sum.go:28)	MOVQ	AX, (SP)
	0x00d2 00210 (sum.go:28)	CALL	runtime.printint(SB)
	0x00d7 00215 (sum.go:28)	CALL	runtime.printnl(SB)
	0x00dc 00220 (sum.go:28)	CALL	runtime.printunlock(SB)
	0x00e1 00225 (sum.go:30)	PCDATA	$0, $1
	0x00e1 00225 (sum.go:30)	PCDATA	$1, $1
	0x00e1 00225 (sum.go:30)	LEAQ	go.string."hello"(SB), AX
	0x00e8 00232 (sum.go:30)	PCDATA	$0, $0
	0x00e8 00232 (sum.go:30)	MOVQ	AX, "".a+256(SP)
	0x00f0 00240 (sum.go:30)	MOVQ	$5, "".a+264(SP)
	0x00fc 00252 (sum.go:30)	XORPS	X0, X0
	0x00ff 00255 (sum.go:30)	MOVUPS	X0, "".~r1+176(SP)
	0x0107 00263 (<unknown line number>)	NOP
	0x0107 00263 (sum.go:14)	MOVQ	"".a+264(SP), AX
	0x010f 00271 (sum.go:14)	PCDATA	$0, $2
	0x010f 00271 (sum.go:14)	PCDATA	$1, $0
	0x010f 00271 (sum.go:14)	MOVQ	"".a+256(SP), CX
	0x0117 00279 (sum.go:14)	MOVQ	CX, "".b+224(SP)
	0x011f 00287 (sum.go:14)	MOVQ	AX, "".b+232(SP)
	0x0127 00295 (sum.go:30)	PCDATA	$0, $0
	0x0127 00295 (sum.go:30)	MOVQ	CX, "".~r1+176(SP)
	0x012f 00303 (sum.go:30)	MOVQ	AX, "".~r1+184(SP)
	0x0137 00311 (sum.go:30)	JMP	313
	0x0139 00313 (sum.go:32)	PCDATA	$0, $1
	0x0139 00313 (sum.go:32)	PCDATA	$1, $2
	0x0139 00313 (sum.go:32)	LEAQ	go.string."foo"(SB), AX
	0x0140 00320 (sum.go:32)	PCDATA	$0, $0
	0x0140 00320 (sum.go:32)	MOVQ	AX, "".a+240(SP)
	0x0148 00328 (sum.go:32)	MOVQ	$3, "".a+248(SP)
	0x0154 00340 (sum.go:32)	PCDATA	$0, $1
	0x0154 00340 (sum.go:32)	PCDATA	$1, $3
	0x0154 00340 (sum.go:32)	LEAQ	go.string."bar"(SB), AX
	0x015b 00347 (sum.go:32)	PCDATA	$0, $0
	0x015b 00347 (sum.go:32)	MOVQ	AX, "".b+208(SP)
	0x0163 00355 (sum.go:32)	MOVQ	$3, "".b+216(SP)
	0x016f 00367 (sum.go:32)	XORPS	X0, X0
	0x0172 00370 (sum.go:32)	MOVUPS	X0, "".~r2+160(SP)
	0x017a 00378 (<unknown line number>)	NOP
	0x017a 00378 (sum.go:19)	PCDATA	$0, $1
	0x017a 00378 (sum.go:19)	LEAQ	""..autotmp_17+128(SP), AX
	0x0182 00386 (sum.go:19)	PCDATA	$0, $0
	0x0182 00386 (sum.go:19)	MOVQ	AX, (SP)
	0x0186 00390 (sum.go:19)	MOVQ	"".a+248(SP), AX
	0x018e 00398 (sum.go:19)	PCDATA	$0, $2
	0x018e 00398 (sum.go:19)	PCDATA	$1, $4
	0x018e 00398 (sum.go:19)	MOVQ	"".a+240(SP), CX
	0x0196 00406 (sum.go:19)	PCDATA	$0, $0
	0x0196 00406 (sum.go:19)	MOVQ	CX, 8(SP)
	0x019b 00411 (sum.go:19)	MOVQ	AX, 16(SP)
	0x01a0 00416 (sum.go:19)	PCDATA	$0, $1
	0x01a0 00416 (sum.go:19)	MOVQ	"".b+208(SP), AX
	0x01a8 00424 (sum.go:19)	PCDATA	$1, $0
	0x01a8 00424 (sum.go:19)	MOVQ	"".b+216(SP), CX
	0x01b0 00432 (sum.go:19)	PCDATA	$0, $0
	0x01b0 00432 (sum.go:19)	MOVQ	AX, 24(SP)
	0x01b5 00437 (sum.go:19)	MOVQ	CX, 32(SP)
	0x01ba 00442 (sum.go:19)	CALL	runtime.concatstring2(SB)
	0x01bf 00447 (sum.go:19)	PCDATA	$0, $1
	0x01bf 00447 (sum.go:19)	MOVQ	40(SP), AX
	0x01c4 00452 (sum.go:19)	MOVQ	48(SP), CX
	0x01c9 00457 (sum.go:32)	MOVQ	AX, ""..autotmp_16+272(SP)
	0x01d1 00465 (sum.go:32)	MOVQ	CX, ""..autotmp_16+280(SP)
	0x01d9 00473 (sum.go:32)	MOVQ	AX, "".~r2+160(SP)
	0x01e1 00481 (sum.go:32)	MOVQ	CX, "".~r2+168(SP)
	0x01e9 00489 (sum.go:32)	JMP	491
	0x01eb 00491 (sum.go:32)	PCDATA	$0, $0
	0x01eb 00491 (sum.go:32)	PCDATA	$1, $5
	0x01eb 00491 (sum.go:32)	MOVQ	AX, "".s+192(SP)
	0x01f3 00499 (sum.go:32)	MOVQ	CX, "".s+200(SP)
	0x01fb 00507 (sum.go:33)	CALL	runtime.printlock(SB)
	0x0200 00512 (sum.go:33)	PCDATA	$0, $1
	0x0200 00512 (sum.go:33)	MOVQ	"".s+192(SP), AX
	0x0208 00520 (sum.go:33)	PCDATA	$1, $0
	0x0208 00520 (sum.go:33)	MOVQ	"".s+200(SP), CX
	0x0210 00528 (sum.go:33)	PCDATA	$0, $0
	0x0210 00528 (sum.go:33)	MOVQ	AX, (SP)
	0x0214 00532 (sum.go:33)	MOVQ	CX, 8(SP)
	0x0219 00537 (sum.go:33)	CALL	runtime.printstring(SB)
	0x021e 00542 (sum.go:33)	CALL	runtime.printnl(SB)
	0x0223 00547 (sum.go:33)	CALL	runtime.printunlock(SB)
	0x0228 00552 (sum.go:34)	MOVQ	288(SP), BP
	0x0230 00560 (sum.go:34)	ADDQ	$296, SP
	0x0237 00567 (sum.go:34)	RET
	0x0238 00568 (sum.go:34)	NOP
	0x0238 00568 (sum.go:22)	PCDATA	$1, $-1
	0x0238 00568 (sum.go:22)	PCDATA	$0, $-2
	0x0238 00568 (sum.go:22)	CALL	runtime.morestack_noctxt(SB)
	0x023d 00573 (sum.go:22)	PCDATA	$0, $-1
	0x023d 00573 (sum.go:22)	JMP	0
	0x0000 64 48 8b 0c 25 00 00 00 00 48 8d 84 24 58 ff ff  dH..%....H..$X..
	0x0010 ff 48 3b 41 10 0f 86 1d 02 00 00 48 81 ec 28 01  .H;A.......H..(.
	0x0020 00 00 48 89 ac 24 20 01 00 00 48 8d ac 24 20 01  ..H..$ ...H..$ .
	0x0030 00 00 48 c7 44 24 48 00 00 00 00 48 c7 44 24 70  ..H.D$H....H.D$p
	0x0040 02 00 00 00 48 c7 44 24 58 03 00 00 00 48 c7 44  ....H.D$X....H.D
	0x0050 24 40 00 00 00 00 48 8b 44 24 70 48 03 44 24 58  $@....H.D$pH.D$X
	0x0060 48 89 44 24 78 48 89 44 24 40 eb 00 48 89 44 24  H.D$xH.D$@..H.D$
	0x0070 48 e8 00 00 00 00 48 8b 44 24 48 48 89 04 24 e8  H.....H.D$HH..$.
	0x0080 00 00 00 00 e8 00 00 00 00 e8 00 00 00 00 48 c7  ..............H.
	0x0090 44 24 68 02 00 00 00 48 c7 44 24 60 03 00 00 00  D$h....H.D$`....
	0x00a0 48 c7 44 24 38 00 00 00 00 48 8b 44 24 68 48 03  H.D$8....H.D$hH.
	0x00b0 44 24 60 48 89 44 24 50 48 89 44 24 38 eb 00 48  D$`H.D$PH.D$8..H
	0x00c0 89 44 24 48 e8 00 00 00 00 48 8b 44 24 48 48 89  .D$H.....H.D$HH.
	0x00d0 04 24 e8 00 00 00 00 e8 00 00 00 00 e8 00 00 00  .$..............
	0x00e0 00 48 8d 05 00 00 00 00 48 89 84 24 00 01 00 00  .H......H..$....
	0x00f0 48 c7 84 24 08 01 00 00 05 00 00 00 0f 57 c0 0f  H..$.........W..
	0x0100 11 84 24 b0 00 00 00 48 8b 84 24 08 01 00 00 48  ..$....H..$....H
	0x0110 8b 8c 24 00 01 00 00 48 89 8c 24 e0 00 00 00 48  ..$....H..$....H
	0x0120 89 84 24 e8 00 00 00 48 89 8c 24 b0 00 00 00 48  ..$....H..$....H
	0x0130 89 84 24 b8 00 00 00 eb 00 48 8d 05 00 00 00 00  ..$......H......
	0x0140 48 89 84 24 f0 00 00 00 48 c7 84 24 f8 00 00 00  H..$....H..$....
	0x0150 03 00 00 00 48 8d 05 00 00 00 00 48 89 84 24 d0  ....H......H..$.
	0x0160 00 00 00 48 c7 84 24 d8 00 00 00 03 00 00 00 0f  ...H..$.........
	0x0170 57 c0 0f 11 84 24 a0 00 00 00 48 8d 84 24 80 00  W....$....H..$..
	0x0180 00 00 48 89 04 24 48 8b 84 24 f8 00 00 00 48 8b  ..H..$H..$....H.
	0x0190 8c 24 f0 00 00 00 48 89 4c 24 08 48 89 44 24 10  .$....H.L$.H.D$.
	0x01a0 48 8b 84 24 d0 00 00 00 48 8b 8c 24 d8 00 00 00  H..$....H..$....
	0x01b0 48 89 44 24 18 48 89 4c 24 20 e8 00 00 00 00 48  H.D$.H.L$ .....H
	0x01c0 8b 44 24 28 48 8b 4c 24 30 48 89 84 24 10 01 00  .D$(H.L$0H..$...
	0x01d0 00 48 89 8c 24 18 01 00 00 48 89 84 24 a0 00 00  .H..$....H..$...
	0x01e0 00 48 89 8c 24 a8 00 00 00 eb 00 48 89 84 24 c0  .H..$......H..$.
	0x01f0 00 00 00 48 89 8c 24 c8 00 00 00 e8 00 00 00 00  ...H..$.........
	0x0200 48 8b 84 24 c0 00 00 00 48 8b 8c 24 c8 00 00 00  H..$....H..$....
	0x0210 48 89 04 24 48 89 4c 24 08 e8 00 00 00 00 e8 00  H..$H.L$........
	0x0220 00 00 00 e8 00 00 00 00 48 8b ac 24 20 01 00 00  ........H..$ ...
	0x0230 48 81 c4 28 01 00 00 c3 e8 00 00 00 00 e9 be fd  H..(............
	0x0240 ff ff                                            ..
	rel 5+4 t=17 TLS+0
	rel 114+4 t=8 runtime.printlock+0
	rel 128+4 t=8 runtime.printint+0
	rel 133+4 t=8 runtime.printnl+0
	rel 138+4 t=8 runtime.printunlock+0
	rel 197+4 t=8 runtime.printlock+0
	rel 211+4 t=8 runtime.printint+0
	rel 216+4 t=8 runtime.printnl+0
	rel 221+4 t=8 runtime.printunlock+0
	rel 228+4 t=16 go.string."hello"+0
	rel 316+4 t=16 go.string."foo"+0
	rel 343+4 t=16 go.string."bar"+0
	rel 443+4 t=8 runtime.concatstring2+0
	rel 508+4 t=8 runtime.printlock+0
	rel 538+4 t=8 runtime.printstring+0
	rel 543+4 t=8 runtime.printnl+0
	rel 548+4 t=8 runtime.printunlock+0
	rel 569+4 t=8 runtime.morestack_noctxt+0
go.cuinfo.packagename. SDWARFINFO dupok size=0
	0x0000 6d 61 69 6e                                      main
go.info."".sum0$abstract SDWARFINFO dupok size=26
	0x0000 04 2e 73 75 6d 30 00 01 01 11 61 00 00 00 00 00  ..sum0....a.....
	0x0010 00 11 62 00 00 00 00 00 00 00                    ..b.......
	rel 13+4 t=29 go.info.int+0
	rel 21+4 t=29 go.info.int+0
go.info."".sum1$abstract SDWARFINFO dupok size=34
	0x0000 04 2e 73 75 6d 31 00 01 01 11 61 00 00 00 00 00  ..sum1....a.....
	0x0010 00 11 62 00 00 00 00 00 00 0c 63 00 09 00 00 00  ..b.......c.....
	0x0020 00 00                                            ..
	rel 13+4 t=29 go.info.int+0
	rel 21+4 t=29 go.info.int+0
	rel 29+4 t=29 go.info.int+0
go.info."".assignstring$abstract SDWARFINFO dupok size=34
	0x0000 04 2e 61 73 73 69 67 6e 73 74 72 69 6e 67 00 01  ..assignstring..
	0x0010 01 11 61 00 00 00 00 00 00 0c 62 00 0e 00 00 00  ..a.......b.....
	0x0020 00 00                                            ..
	rel 21+4 t=29 go.info.string+0
	rel 29+4 t=29 go.info.string+0
go.info."".concatestring$abstract SDWARFINFO dupok size=35
	0x0000 04 2e 63 6f 6e 63 61 74 65 73 74 72 69 6e 67 00  ..concatestring.
	0x0010 01 01 11 61 00 00 00 00 00 00 11 62 00 00 00 00  ...a.......b....
	0x0020 00 00 00                                         ...
	rel 22+4 t=29 go.info.string+0
	rel 30+4 t=29 go.info.string+0
go.loc."".sum0 SDWARFLOC size=0
go.info."".sum0 SDWARFINFO size=53
	0x0000 05 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0010 00 00 00 00 00 01 9c 12 00 00 00 00 01 9c 12 00  ................
	0x0020 00 00 00 02 91 08 0f 7e 72 32 00 01 04 00 00 00  .......~r2......
	0x0030 00 02 91 10 00                                   .....
	rel 0+0 t=24 type.int+0
	rel 1+4 t=29 go.info."".sum0$abstract+0
	rel 5+8 t=1 "".sum0+0
	rel 13+8 t=1 "".sum0+25
	rel 24+4 t=29 go.info."".sum0$abstract+9
	rel 31+4 t=29 go.info."".sum0$abstract+17
	rel 45+4 t=29 go.info.int+0
go.range."".sum0 SDWARFRANGE size=0
go.debuglines."".sum0 SDWARFMISC size=12
	0x0000 04 02 12 6a 06 41 04 01 03 7c 06 01              ...j.A...|..
go.loc."".sum1 SDWARFLOC size=0
go.info."".sum1 SDWARFINFO size=61
	0x0000 05 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0010 00 00 00 00 00 01 9c 12 00 00 00 00 01 9c 12 00  ................
	0x0020 00 00 00 02 91 08 0f 7e 72 32 00 01 08 00 00 00  .......~r2......
	0x0030 00 02 91 10 0d 00 00 00 00 02 91 68 00           ...........h.
	rel 0+0 t=24 type.int+0
	rel 1+4 t=29 go.info."".sum1$abstract+0
	rel 5+8 t=1 "".sum1+0
	rel 13+8 t=1 "".sum1+52
	rel 24+4 t=29 go.info."".sum1$abstract+9
	rel 31+4 t=29 go.info."".sum1$abstract+17
	rel 45+4 t=29 go.info.int+0
	rel 53+4 t=29 go.info."".sum1$abstract+25
go.range."".sum1 SDWARFRANGE size=0
go.debuglines."".sum1 SDWARFMISC size=19
	0x0000 04 02 0a 03 02 14 f6 06 41 06 6a 06 41 04 01 03  ........A.j.A...
	0x0010 77 06 01                                         w..
go.loc."".assignstring SDWARFLOC size=0
go.info."".assignstring SDWARFINFO size=53
	0x0000 05 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0010 00 00 00 00 00 01 9c 12 00 00 00 00 01 9c 0f 7e  ...............~
	0x0020 72 31 00 01 0d 00 00 00 00 02 91 10 0d 00 00 00  r1..............
	0x0030 00 02 91 60 00                                   ...`.
	rel 0+0 t=24 type.string+0
	rel 1+4 t=29 go.info."".assignstring$abstract+0
	rel 5+8 t=1 "".assignstring+0
	rel 13+8 t=1 "".assignstring+61
	rel 24+4 t=29 go.info."".assignstring$abstract+17
	rel 37+4 t=29 go.info.string+0
	rel 45+4 t=29 go.info."".assignstring$abstract+25
go.range."".assignstring SDWARFRANGE size=0
go.debuglines."".assignstring SDWARFMISC size=22
	0x0000 04 02 0a 03 07 14 06 b9 06 42 06 41 06 9c 06 41  .........B.A...A
	0x0010 04 01 03 72 06 01                                ...r..
go.loc."".concatestring SDWARFLOC size=0
go.info."".concatestring SDWARFINFO size=53
	0x0000 05 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0010 00 00 00 00 00 01 9c 12 00 00 00 00 01 9c 12 00  ................
	0x0020 00 00 00 02 91 10 0f 7e 72 32 00 01 12 00 00 00  .......~r2......
	0x0030 00 02 91 20 00                                   ... .
	rel 0+0 t=24 type.string+0
	rel 1+4 t=29 go.info."".concatestring$abstract+0
	rel 5+8 t=1 "".concatestring+0
	rel 13+8 t=1 "".concatestring+127
	rel 24+4 t=29 go.info."".concatestring$abstract+18
	rel 31+4 t=29 go.info."".concatestring$abstract+26
	rel 45+4 t=29 go.info.string+0
go.range."".concatestring SDWARFRANGE size=0
go.debuglines."".concatestring SDWARFMISC size=26
	0x0000 04 02 03 0c 14 0a a5 06 b9 06 42 06 5f 06 08 af  ..........B._...
	0x0010 06 41 06 08 4a 04 01 03 6f 01                    .A..J...o.
go.string."hello" SRODATA dupok size=5
	0x0000 68 65 6c 6c 6f                                   hello
go.string."foo" SRODATA dupok size=3
	0x0000 66 6f 6f                                         foo
go.string."bar" SRODATA dupok size=3
	0x0000 62 61 72                                         bar
go.loc."".main SDWARFLOC size=0
go.info."".main SDWARFINFO size=244
	0x0000 03 22 22 2e 6d 61 69 6e 00 00 00 00 00 00 00 00  ."".main........
	0x0010 00 00 00 00 00 00 00 00 00 01 9c 00 00 00 00 01  ................
	0x0020 0a 69 00 17 00 00 00 00 03 91 98 7e 0a 73 00 20  .i.........~.s. 
	0x0030 00 00 00 00 03 91 90 7f 06 00 00 00 00 00 00 00  ................
	0x0040 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0050 00 18 12 00 00 00 00 03 91 c0 7e 12 00 00 00 00  ..........~.....
	0x0060 03 91 a8 7e 00 06 00 00 00 00 00 00 00 00 00 00  ...~............
	0x0070 00 00 00 00 00 00 00 00 00 00 00 00 00 00 1b 12  ................
	0x0080 00 00 00 00 03 91 b8 7e 12 00 00 00 00 03 91 b0  .......~........
	0x0090 7e 0d 00 00 00 00 03 91 a0 7e 00 06 00 00 00 00  ~........~......
	0x00a0 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x00b0 00 00 00 00 1e 12 00 00 00 00 02 91 50 0d 00 00  ............P...
	0x00c0 00 00 03 91 b0 7f 00 06 00 00 00 00 00 00 00 00  ................
	0x00d0 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x00e0 20 12 00 00 00 00 02 91 40 12 00 00 00 00 03 91   .......@.......
	0x00f0 a0 7f 00 00                                      ....
	rel 0+0 t=24 type.[32]uint8+0
	rel 0+0 t=24 type.int+0
	rel 0+0 t=24 type.string+0
	rel 9+8 t=1 "".main+0
	rel 17+8 t=1 "".main+578
	rel 27+4 t=30 gofile../root/go/src/github.com/DQNEO/minigo/samples/gosamples/sum/sum.go+0
	rel 36+4 t=29 go.info.int+0
	rel 48+4 t=29 go.info.string+0
	rel 57+4 t=29 go.info."".sum0$abstract+0
	rel 61+8 t=1 "".main+86
	rel 69+8 t=1 "".main+96
	rel 77+4 t=30 gofile../root/go/src/github.com/DQNEO/minigo/samples/gosamples/sum/sum.go+0
	rel 83+4 t=29 go.info."".sum0$abstract+9
	rel 92+4 t=29 go.info."".sum0$abstract+17
	rel 102+4 t=29 go.info."".sum1$abstract+0
	rel 106+8 t=1 "".main+169
	rel 114+8 t=1 "".main+184
	rel 122+4 t=30 gofile../root/go/src/github.com/DQNEO/minigo/samples/gosamples/sum/sum.go+0
	rel 128+4 t=29 go.info."".sum1$abstract+9
	rel 137+4 t=29 go.info."".sum1$abstract+17
	rel 146+4 t=29 go.info."".sum1$abstract+25
	rel 156+4 t=29 go.info."".assignstring$abstract+0
	rel 160+8 t=1 "".main+263
	rel 168+8 t=1 "".main+295
	rel 176+4 t=30 gofile../root/go/src/github.com/DQNEO/minigo/samples/gosamples/sum/sum.go+0
	rel 182+4 t=29 go.info."".assignstring$abstract+17
	rel 190+4 t=29 go.info."".assignstring$abstract+25
	rel 200+4 t=29 go.info."".concatestring$abstract+0
	rel 204+8 t=1 "".main+378
	rel 212+8 t=1 "".main+457
	rel 220+4 t=30 gofile../root/go/src/github.com/DQNEO/minigo/samples/gosamples/sum/sum.go+0
	rel 226+4 t=29 go.info."".concatestring$abstract+18
	rel 234+4 t=29 go.info."".concatestring$abstract+26
go.range."".main SDWARFRANGE size=0
go.debuglines."".main SDWARFMISC size=107
	0x0000 04 02 03 10 14 0a 08 2d f6 6a 06 69 06 03 71 bf  .......-.j.i..q.
	0x0010 06 41 06 03 0e 46 06 41 06 88 06 41 06 08 11 06  .A...F.A...A....
	0x0020 69 06 03 72 bf 06 41 06 03 0d 78 06 41 06 56 06  i..r..A...x.A.V.
	0x0030 41 06 08 11 06 55 06 08 03 74 51 06 5f 06 03 10  A....U...tQ._...
	0x0040 ff 06 5f 06 75 06 55 06 02 22 03 77 fb 06 5f 06  .._.u.U..".w.._.
	0x0050 02 20 ff 06 41 06 03 08 78 06 5f 06 08 c4 06 41  . ..A...x._....A
	0x0060 06 08 b0 03 78 ab 04 01 03 6b 01                 ....x....k.
""..inittask SNOPTRDATA size=24
	0x0000 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0010 00 00 00 00 00 00 00 00                          ........
runtime.memequal64·f SRODATA dupok size=8
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=1 runtime.memequal64+0
runtime.gcbits.01 SRODATA dupok size=1
	0x0000 01                                               .
type..namedata.*[]uint8- SRODATA dupok size=11
	0x0000 00 00 08 2a 5b 5d 75 69 6e 74 38                 ...*[]uint8
type.*[]uint8 SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 a5 8e d0 69 08 08 08 36 00 00 00 00 00 00 00 00  ...i...6........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=1 runtime.memequal64·f+0
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+4 t=5 type..namedata.*[]uint8-+0
	rel 48+8 t=1 type.[]uint8+0
type.[]uint8 SRODATA dupok size=56
	0x0000 18 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 df 7e 2e 38 02 08 08 17 00 00 00 00 00 00 00 00  .~.8............
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+4 t=5 type..namedata.*[]uint8-+0
	rel 44+4 t=6 type.*[]uint8+0
	rel 48+8 t=1 type.uint8+0
type..eqfunc32 SRODATA dupok size=16
	0x0000 00 00 00 00 00 00 00 00 20 00 00 00 00 00 00 00  ........ .......
	rel 0+8 t=1 runtime.memequal_varlen+0
type..namedata.*[32]uint8- SRODATA dupok size=13
	0x0000 00 00 0a 2a 5b 33 32 5d 75 69 6e 74 38           ...*[32]uint8
type.*[32]uint8 SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 f4 c7 79 15 08 08 08 36 00 00 00 00 00 00 00 00  ..y....6........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=1 runtime.memequal64·f+0
	rel 32+8 t=1 runtime.gcbits.01+0
	rel 40+4 t=5 type..namedata.*[32]uint8-+0
	rel 48+8 t=1 type.[32]uint8+0
runtime.gcbits. SRODATA dupok size=0
type.[32]uint8 SRODATA dupok size=72
	0x0000 20 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00   ...............
	0x0010 9c 59 ff a8 0a 01 01 11 00 00 00 00 00 00 00 00  .Y..............
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0040 20 00 00 00 00 00 00 00                           .......
	rel 24+8 t=1 type..eqfunc32+0
	rel 32+8 t=1 runtime.gcbits.+0
	rel 40+4 t=5 type..namedata.*[32]uint8-+0
	rel 44+4 t=6 type.*[32]uint8+0
	rel 48+8 t=1 type.uint8+0
	rel 56+8 t=1 type.[]uint8+0
gclocals·33cdeccccebe80329f1fdbee7f5874cb SRODATA dupok size=8
	0x0000 01 00 00 00 00 00 00 00                          ........
gclocals·9fad110d66c97cf0b58d28cccea80b12 SRODATA dupok size=11
	0x0000 03 00 00 00 03 00 00 00 01 00 04                 ...........
gclocals·d8b28f51bb91e05d264803f0f600a200 SRODATA dupok size=11
	0x0000 03 00 00 00 02 00 00 00 00 00 00                 ...........
gclocals·9fb7f0986f647f17cb53dda1484e0f7a SRODATA dupok size=10
	0x0000 02 00 00 00 01 00 00 00 00 01                    ..........
gclocals·2625d1fdbbaf79a2e52296235cb6527c SRODATA dupok size=12
	0x0000 04 00 00 00 05 00 00 00 05 04 00 10              ............
gclocals·f6bd6b3389b872033d462029172c8612 SRODATA dupok size=8
	0x0000 04 00 00 00 00 00 00 00                          ........
gclocals·1cf923758aae2e428391d1783fe59973 SRODATA dupok size=11
	0x0000 03 00 00 00 02 00 00 00 00 01 02                 ...........
gclocals·f14a5bc6d08bc46424827f54d2e3f8ed SRODATA dupok size=8
	0x0000 06 00 00 00 00 00 00 00                          ........
gclocals·cde884c6f8ebff321c52f642fdb453a8 SRODATA dupok size=20
	0x0000 06 00 00 00 10 00 00 00 00 00 00 10 00 04 40 04  ..............@.
	0x0010 40 00 10 00                                      @...
