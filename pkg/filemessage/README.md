# FileMessage

`filemessage` package provides methods to encode and decode messages from Go structs to bytes and the other way around.

The messages handled by `filemessage` are meant to carry a file as a stream of bytes and a identifier for the message itsef.

## Encode the message

```go
file, _ := os.ReadFile("../../fixtures/sample.jpeg")
msg := FileMessage{
    ID:   1,
    File: file,
}

encodedMsg, err := msg.Bytes()
if err != nil {
    // handle error
}
```

## Decode the message

```go
decodedMsg, err := encodedMsg.Read()
if err != nil {
    // handle error
}
```

## Note on `FileMessage.BytesJSON`

Although easy to use, `FileMessage.Bytes` requires that the code decoding and using the message is also written in Go in order to call `Read(encoded)`.

To answer this caveat, `FileMessage.BytesJSON` offers a way for the decoder to simply decode to the universal JSON format.

However, in our benchmarks, it is about 4 times slower than `FileMessage.Bytes`, which is trivially explained by the extra work needed to encode the `File []byte` field to a base64 string before marshalling it to JSON. This is accentuated by the fact that the `File` field is always larger than the `ID` integer, which is ironically (is our case) well supported by the JSON format.

> Within this project scope we have full control over both the API server and the worker who would use this library. Thus there is no rational reason to use `FileMessage.BytesJSON`, as we can settle on using Go for writing the worker.
