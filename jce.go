package gojce

/*
0: INT8/BYTE 单字节, 如果内容为0x00就变为ZERO_TAG
1: INT16
2: INT32
3: INT64
4: FLOAT32
5: FLOAT64
6: STRING1 长度小于256的字符串(用uint8标识长度)
7: STRING4 长度大于256的字符串(用uint64标识长度)
8: MAP 用一个id为0的INT8/16/32/64类型标识项数, 之后先k再v, k的id为0, v的id为1
9: LIST 用一个id为0的INT32类型标识项数, 之后每一项的id都为0
10: STRUCT_START 标识一个jceStruct的开始..?
11: STRUCT_END 标识一个jceStruct的结束..?
12: ZERO_TAG 该项为0.jpg
13: BYTES 现时一个内容为0x00的UNIT8(不变为ZERO_TAG), 再用一个id为0的INT32类型标识长度, 然后直接跟上内容
*/

/*
id可以被理解为该项在第几个(从0开始)...吧
*/
