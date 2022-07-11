## fix length
描述：
客户端在发送数据包的时候，每个包都固定长度，比如1024个字节大小，如果客户端发送的数据长度不足1024个字节，则通过补充空格的方式补全到指定长度.
应用：
使用固定长度的封包和拆包的场景，例如FixedLengthFrameDecoder。解码一器会每次读取固定长度的消息，如果当前读取到的消息不足指定长度，那么就会等待下一个消息到达后进行补足。

## delimiter based
#### 描述
客户端在每个包的末尾使用固定的分隔符，例如\r\n，如果一个包被拆分了，则等待下一个包发送过来之后找到其中的\r\n，然后对其拆分后的头部部分与前一个包的剩余部分进行合并，这样就得到了一个完整的包。
#### 应用
对于通过分隔符进行粘包和拆包问题的处理，Netty提供了两个编解码的类，LineBasedFrameDecoder和DelimiterBasedFrameDecoder。这里LineBasedFrameDecoder的作用主要是通过换行符，即\n或者\r\n对数据进行处理；而DelimiterBasedFrameDecoder的作用则是通过用户指定的分隔符对数据进行粘包和拆包处理。
## length field based frame decoder
#### 描述
将消息分为头部和消息体，在头部中保存有当前整个消息的长度，只有在读取到足够长度的消息之后才算是读到了一个完整的消息。
#### 应用
goim的解码器就是使用package length来表示整个包的长度，只有读到一个
